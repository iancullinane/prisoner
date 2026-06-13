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
		and it will always select `Luigi` as the opponent.
	*/

	router := http.NewServeMux()
	router.Handle("GET /players", http.HandlerFunc(p.listPlayersHandler))
	// Get a player by id
	router.Handle("GET /players/{name}", http.HandlerFunc(p.playersHandler))
	// Create a player by name
	router.Handle("POST /players/{name}", http.HandlerFunc(p.playersHandler))
	// router.Handle("POST /players/{id}", http.HandlerFunc(p.playersHandler))
	router.Handle("GET /history", http.HandlerFunc(p.historyHandler))
	router.Handle("GET /history/{id}", http.HandlerFunc(p.historyHandler))
	router.Handle("POST /play/{player_a}/{player_b}", http.HandlerFunc(p.playHandler))
	router.Handle("POST /play/{player_a}", http.HandlerFunc(p.playHandler))

	p.Handler = router

	return p
}

// MARK: Logging Setup
// ==========================================

// func newLogger() *slog.Logger {
// 	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
// 		Level: slog.LevelInfo,
// 	})
// 	return slog.New(h).With(
// 		slog.String("service", "prisoner"),
// 	)
// }

// MARK: Handlers
// ==========================================

func (p *PlayerServer) playHandler(w http.ResponseWriter, r *http.Request) {

	playerAParam := r.PathValue("player_a")
	playerBParam := r.PathValue("player_b")

	if playerAParam == playerBParam {
		http.Error(w, ErrPlayerCannotPlaySelf.Error(), http.StatusBadRequest)
		return
	}

	playerAID, err := uuid.Parse(playerAParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing player id: %v", err), http.StatusBadRequest)
		return
	}

	playerA, err := p.playerStore.GetPlayerByID(playerAID)
	if err != nil {
		http.Error(w, "could not find player a", http.StatusInternalServerError)
		return
	}

	var playerB types.Player
	if playerBParam == "" {
		playerB, err = p.playerStore.GetRandomPlayerExcept(playerAID)
	} else {
		playerBID, err := uuid.Parse(playerBParam)
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing player id: %v", err), http.StatusBadRequest)
			return
		}
		playerB, err = p.playerStore.GetPlayerByID(playerBID)
	}

	if err != nil {
		http.Error(w, "could not find player b", http.StatusInternalServerError)
		return
	}

	playeAMove, playerBmove := prisoner.Cooperate, prisoner.RandomMove()

	p.logger.Info("play")

	interaction := types.NewInteraction(playerA.ID, playerB.ID, playeAMove, playerBmove)
	if err := p.historyStore.RecordInteraction(interaction); err != nil {
		fmt.Printf("record interaction: %v", err)
		http.Error(w, "could not record interaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(interaction)
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
		statusCode = http.StatusOK
	case http.MethodPost:
		player, err = p.playerStore.GetOrCreatePlayer(playerName)
		statusCode = http.StatusOK
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("could not get player: %v", err), http.StatusNotFound)
		return
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
