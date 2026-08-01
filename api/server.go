package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/store"
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

const jsonContentType = "application/json"

var (
	ErrPlayerCannotPlaySelf = errors.New("player cannot play itself")
)

type PlayerServer struct {
	logger       *slog.Logger
	playerStore  types.PlayerStore
	historyStore types.HistoryStore
	http.Handler
}

type PlayerHistoryProvider interface {
	GetPrettyHistory(playerID *uuid.UUID) (types.PrettyHistory, error)
}

func NewPlayerServer(
	logger *slog.Logger,
	playerStore types.PlayerStore,
	historyStore types.HistoryStore) *PlayerServer {

	p := new(PlayerServer)
	p.logger = logger

	p.playerStore = playerStore
	p.historyStore = historyStore

	router := http.NewServeMux()
	router.Handle("GET /healthz", http.HandlerFunc(p.healthzHandler))

	// business routes
	router.Handle("GET /api/v1/players", http.HandlerFunc(p.listPlayersHandler))
	// Get a player by id
	router.Handle("GET /api/v1/players/{name}", http.HandlerFunc(p.playersHandler))
	// Create a player by name
	router.Handle("POST /api/v1/players/{name}", http.HandlerFunc(p.playersHandler))
	router.Handle("GET /api/v1/history", http.HandlerFunc(p.historyHandler))
	// History of a player
	router.Handle("GET /api/v1/history/{id}", http.HandlerFunc(p.historyHandler))
	router.Handle("POST /api/v1/play", http.HandlerFunc(p.playHandler))

	p.Handler = withCORS(router)

	return p
}

// withCORS allows a browser-based frontend running on a different origin
// (e.g. the Vite dev server) to call this API. The origin is reflected back
// rather than hardcoded so it works regardless of where the frontend runs.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// healthzHandler is a K8s liveness probe: 200 while the process serves.
func (p *PlayerServer) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

// MARK: Handlers
// ==========================================

func (p *PlayerServer) playHandler(w http.ResponseWriter, r *http.Request) {

	var req PlayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if req.PlayerA == req.PlayerB {
		http.Error(w, "player cannot play itself", http.StatusBadRequest)
		return
	}

	playerA, err := p.playerStore.GetPlayerByID(req.PlayerA)
	if err != nil {
		http.Error(w, "could not find player a", http.StatusInternalServerError)
		return
	}

	playerB, err := p.playerStore.GetPlayerByID(req.PlayerB)
	if err != nil {
		http.Error(w, "could not find player b", http.StatusInternalServerError)
		return
	}

	interaction := types.NewInteraction(playerA.ID, playerB.ID, req.PlayerAMove, req.PlayerBMove)
	if err := p.historyStore.RecordInteraction(interaction); err != nil {
		p.logger.Error("could not record interaction", "error", err)
		http.Error(w, "could not record interaction", http.StatusInternalServerError)
		return
	}

	playerAScore, playerBScore := prisoner.Play(req.PlayerAMove, req.PlayerBMove)
	response := PlayResponse{
		ID:           interaction.ID,
		PlayerAScore: playerAScore,
		PlayerBScore: playerBScore,
	}

	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (p *PlayerServer) listPlayersHandler(w http.ResponseWriter, r *http.Request) {

	players, err := p.playerStore.GetAllPlayers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(players)
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	playerName := r.PathValue("name")

	var player types.Player
	var statusCode int
	var err error

	switch r.Method {
	case http.MethodGet:
		player, err = p.playerStore.GetPlayerByName(playerName)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not get player: %v", err), http.StatusNotFound)
			return
		}
		statusCode = http.StatusOK
	case http.MethodPost:
		player, err = p.playerStore.GetOrCreatePlayer(playerName)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not create player: %v", err), http.StatusInternalServerError)
			return
		}
		statusCode = http.StatusOK
	}

	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(player)
}

func (p *PlayerServer) historyHandler(w http.ResponseWriter, r *http.Request) {
	var playerID *uuid.UUID
	if raw := r.PathValue("id"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not parse player id: %v", err), http.StatusBadRequest)
			return
		}
		playerID = &id
	}

	var ph types.PrettyHistory
	var err error
	if provider, ok := p.historyStore.(PlayerHistoryProvider); ok {
		ph, err = provider.GetPrettyHistory(playerID) // postgres: SQL join fast path
	} else {
		ph, err = store.GetPrettyHistoryFromStores(playerID, p.playerStore, p.historyStore)
	}
	if err != nil {
		p.logger.Error("could not load history", slog.Any("error", err))
		http.Error(w, "could not load history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(ph)
}
