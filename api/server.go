package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

const jsonContentType = "application/json"

var (
	ianID, _ = uuid.Parse("9999999999-2222-bbbb-3333-333333333333")
)

type Player = types.Player

type PlayerServer struct {
	playerStore  types.PlayerStore
	historyStore types.HistoryStore
	http.Handler
}

func NewPlayerServer(playerStore types.PlayerStore, historyStore types.HistoryStore) *PlayerServer {

	p := new(PlayerServer)

	p.playerStore = playerStore
	p.historyStore = historyStore

	router := http.NewServeMux()
	// TODO: Update the player route to work with newer versions of player
	//   Revamp player to search by name or ID, to post with the ID and Move. This would be usual `GET /players` which then gets a history of all interactions from that user.
	// labels: feature
	router.Handle("GET /players/{name}", http.HandlerFunc(p.playersHandler))
	router.Handle("POST /players/{name}", http.HandlerFunc(p.playersHandler))
	// router.Handle("POST /players/{id}", http.HandlerFunc(p.playersHandler))
	router.Handle("GET /history", http.HandlerFunc(p.historyHandler))
	router.Handle("GET /history/{id}", http.HandlerFunc(p.historyHandler))
	router.Handle("POST /play/{id}", http.HandlerFunc(p.playHandler))

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
	// TODO: Implement play against a CPU opponent on the server
	//  until real user backing exists (specifically UUID support) have the play response return a random move, but always set the opponent as uuid:uuid.Parse("00000000-0000-aaaa-2222-222222222222") for testing and development purposes.
	// labels:story
	devOpponent, err := types.NewPlayerFromID("00000000-0000-aaaa-2222-222222222222")
	if err != nil {
		fmt.Println("could not create opponent: %w", err)
		http.Error(w, "could not create opponent", http.StatusInternalServerError)
		return
	}

	player := r.PathValue("id")
	playerID, err := uuid.Parse(player)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing player id: %v", err), http.StatusBadRequest)
		return
	}

	protagonist, err := p.playerStore.GetPlayerByID(playerID)
	if err != nil {
		http.Error(w, "could not get protagonist", http.StatusInternalServerError)
		return
	}

	fmt.Println("protagonist", protagonist)

	protagonistMove, opponentMove := prisoner.Cooperate, prisoner.Cooperate

	interaction := types.NewInteraction(protagonist.ID, devOpponent.ID, protagonistMove, opponentMove)
	if err := p.historyStore.RecordInteraction(interaction); err != nil {
		fmt.Println("record interaction: %w", err)
		http.Error(w, "could not record interaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, interaction)
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
		http.Error(w, fmt.Sprintf("could not get player: %v", err), http.StatusInternalServerError)
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
		if interaction.Protagonist == playerID || interaction.Opponent == playerID {
			filtered = append(filtered, interaction)
		}
	}
	return filtered
}
