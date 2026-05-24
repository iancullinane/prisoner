package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

const jsonContentType = "application/json"

type Player = types.Player

type PlayerServer struct {
	store types.PlayerStore
	http.Handler
}

func NewPlayerServer(store types.PlayerStore) *PlayerServer {

	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("GET /league", http.HandlerFunc(p.leagueHandler))
	router.Handle("GET /players/{name}", http.HandlerFunc(p.playersHandler))
	router.Handle("POST /players/{name}", http.HandlerFunc(p.playersHandler))
	router.Handle("/play/", http.HandlerFunc(p.playHandler))

	p.Handler = router

	return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	leagueTable, err := p.store.GetLeague()
	if err != nil {
		log.Printf("league: %v", err)
		http.Error(w, "could not load league", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(leagueTable)
}

func (p *PlayerServer) playHandler(w http.ResponseWriter, r *http.Request) {
	r1, _ := prisoner.Play('C', 'C')

	fmt.Fprint(w, r1)
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
	score, err := p.store.GetPlayerScore(player)
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
	if err := p.store.RecordWin(player); err != nil {
		log.Printf("record win for %q: %v", player, err)
		http.Error(w, "could not record win", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
