package api

import (
	"encoding/json"
	"fmt"
	"log"
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
	router.Handle("GET /players/{name}", http.HandlerFunc(p.playersHandler))
	router.Handle("POST /players/{name}", http.HandlerFunc(p.playersHandler))
	router.Handle("GET /history/{id}", http.HandlerFunc(p.historyHandler))
	router.Handle("POST /play/{id}", http.HandlerFunc(p.playHandler))

	p.Handler = router

	return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	leagueTable, err := p.playerStore.GetLeague()
	if err != nil {
		log.Printf("league: %v", err)
		http.Error(w, "could not load league", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(leagueTable)
}

func (p *PlayerServer) playHandler(w http.ResponseWriter, r *http.Request) {
	protagonistMove, opponentMove := prisoner.Cooperate, prisoner.Cooperate
	protagonistResult, _ := prisoner.Play(protagonistMove, opponentMove)

	interaction := types.NewInteraction(uuid.Nil, uuid.Nil, protagonistMove, opponentMove)
	if err := p.historyStore.RecordInteraction(interaction); err != nil {
		log.Printf("record interaction: %v", err)
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
		log.Printf("get score for %q: %v", player, err)
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
		log.Printf("record win for %q: %v", player, err)
		http.Error(w, "could not record win", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer) historyHandler(w http.ResponseWriter, r *http.Request) {

	pID := r.PathValue("id")
	if pID == "" {
		//do nothing
	}

	history, err := p.historyStore.GetHistory()
	if err != nil {
		log.Printf("get history: %v", err)
		http.Error(w, "could not load history", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(history)
}
