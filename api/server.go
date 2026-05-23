package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

const jsonContentType = "application/json"

type Player = types.Player
type PlayerStore = types.PlayerStore

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {

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
	leagueTable := p.store.GetLeague()
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
	score := p.store.GetPlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
