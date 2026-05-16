//go:build !integration

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iancullinane/prisoner/internal/store/memory"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := memory.NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest("players", player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest("players", player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest("players", player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest("players", player))
		assertResponseStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertResponseStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{Name: "Pepper", Wins: 3},
		}
		assertLeague(t, got, want)
	})
}
