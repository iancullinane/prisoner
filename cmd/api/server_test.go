package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

// MARK: POST test

func TestPOSTPlayers(t *testing.T) {
	store := StubPlayerStore{
		scores:   map[string]int{},
		winCalls: nil,
	}
	server := NewPlayerServer(&store)

	tests := []struct {
		name string
		code int
	}{
		{
			"test_post",
			http.StatusAccepted,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			buf := bytes.NewBufferString("body")
			request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", buf)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertResponseStatus(t, response.Code, tc.code)

			if len(store.winCalls) != 1 {
				t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
			}

			if store.winCalls[0] != "Pepper" {
				t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], "Pepper")
			}

		})
	}
}

// MARK: GET test

func TestGETPlayers(t *testing.T) {

	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Steve":  10,
		},
		[]string{
			"Pepper",
			"Steve",
		},
	}

	server := NewPlayerServer(&store)

	tests := []struct {
		name       string
		url        string
		playerName string
		response   string
		code       int
	}{
		{
			"test_pepper",
			"players",
			"Pepper",
			"20",
			http.StatusOK,
		},
		{
			"test_steve",
			"players",
			"Steve",
			"10",
			http.StatusOK,
		},
		{
			"test_missing_player",
			"players",
			"Apollo",
			"0",
			http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request := newGetScoreRequest(tc.url, tc.playerName)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			got := response.Body.String()
			want := tc.response

			assertResponseStatus(t, response.Code, tc.code)
			assertResponseBody(t, got, want)
		})
	}
}

func TestLeague(t *testing.T) {
	store := StubPlayerStore{}
	server := NewPlayerServer(&store)

	t.Run("it returns 200 on /league", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusOK)
	})
}

// MARK: Request Builders

func newGetScoreRequest(route, name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", route, name), nil)
	return req
}

func newPostWinRequest(route, name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/%s/%s", route, name), nil)

	return req
}

// MARK: Assertions

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body wrong, got %s, want %s", got, want)
	}
}

func assertResponseStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response header wrong, got %d, want %d", got, want)
	}
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest("players", player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest("players", player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest("players", player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest("players", player))
	assertResponseStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}
