package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

const jsonContentType = "application/json"

var (
	ErrPlayerCannotPlaySelf = errors.New("player cannot play itself")
)

// var (
// 	luigi, _ = uuid.Parse("00000000-0000-0000-0000-000000000000")
// )

type Player = types.Player

type PlayerServer struct {
	logger       *slog.Logger
	playerStore  types.PlayerStore
	historyStore types.HistoryStore
	http.Handler
}

func NewPlayerServer(logger *slog.Logger, playerStore types.PlayerStore, historyStore types.HistoryStore) *PlayerServer {

	p := new(PlayerServer)
	p.logger = logger

	p.playerStore = playerStore
	p.historyStore = historyStore

	/*
		2026-06-05
		As of the player needs to exist, and be passed as a parameter into the url
	*/

	router := http.NewServeMux()
	router.Handle("GET /players", http.HandlerFunc(p.listPlayersHandler))
	// Get a player by id
	router.Handle("GET /players/{name}", http.HandlerFunc(p.playersHandler))
	// Create a player by name
	router.Handle("POST /players/{name}", http.HandlerFunc(p.playersHandler))
	router.Handle("GET /history", http.HandlerFunc(p.historyHandler))
	// History of a player
	router.Handle("GET /history/{id}", http.HandlerFunc(p.historyHandler))
	router.Handle("POST /play", http.HandlerFunc(p.playHandler))

	p.Handler = router

	return p
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
		http.Error(w, "could not find player a", http.StatusInternalServerError)
		return
	}

	interaction := types.NewInteraction(playerA.ID, playerB.ID, req.PlayerAMove, req.PlayerBMove)
	if err := p.historyStore.RecordInteraction(interaction); err != nil {
		fmt.Printf("record interaction: %v", err)
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
	pID := r.PathValue("id")
	if pID == "" {
		//do nothing
	}

	history, err := p.historyStore.GetHistory()
	if err != nil {
		http.Error(w, "could not load history", http.StatusInternalServerError)
		return
	}

	if rawID := r.PathValue("id"); rawID != "" {
		playerID, err := uuid.Parse(rawID)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not parse player id: %v", err), http.StatusBadRequest)
			return
		}
		history = filterHistoryByPlayer(history, playerID)
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(history)
}

func filterHistoryByPlayer(history types.History, playerID uuid.UUID) types.History {
	filtered := make(types.History, 0, len(history))
	for _, interaction := range history {
		if interaction.PlayerA == playerID || interaction.PlayerB == playerID {
			filtered = append(filtered, interaction)
		}
	}
	return filtered
}
