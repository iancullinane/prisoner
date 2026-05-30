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
	router.Handle("GET /league", http.HandlerFunc(p.leagueHandler))
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

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	leagueTable, err := p.playerStore.GetLeague()
	if err != nil {
		http.Error(w, "could not load league", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(leagueTable)
}

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

	proUUID, _ := uuid.NewRandom()

	protagonist := types.Player{
		ID:   proUUID,
		Name: "Ian",
		Wins: 0,
	}

	protagonistMove, opponentMove := prisoner.Cooperate, prisoner.Cooperate
	protagonistResult, _ := prisoner.Play(protagonistMove, opponentMove)

	interaction := types.NewInteraction(protagonist.ID, devOpponent.ID, protagonistMove, opponentMove)
	if err := p.historyStore.RecordInteraction(interaction); err != nil {
		// TODO: Add createPlayer before recordInteraction
		//  right now when you attempt to record an interaction it fails because the protagonist and opponent don't exit on the players table they are fk'd to
		// labels: story
		fmt.Println("record interaction: %w", err)
		http.Error(w, "could not record interaction", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, protagonistResult)
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := r.PathValue("name")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score, err := p.playerStore.GetPlayerScore(player)
	if err != nil {
		http.Error(w, "could not load score", http.StatusInternalServerError)
		return
	}
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	if err := p.playerStore.RecordWin(player); err != nil {
		http.Error(w, "could not record win", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer) historyHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: Get history per player
	//   We need to be able to get the history of a player as well, since the data is flat
	//   and only tracks the actual interactions, we can likely reduce on the map and get
	//   what we need. This would operate differently depending on store types, ensure
	//   it is behind a swappable interface. Different stores will have totally different
	//   performance costs.
	//  labels: story
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
			http.Error(w, "invalid player id", http.StatusBadRequest)
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
