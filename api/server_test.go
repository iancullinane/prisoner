package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

var (
	testID1, _   = uuid.Parse("00000000-0000-aaaa-2222-222222222222")
	testID2, _   = uuid.Parse("11111111-1111-bbbb-3333-333333333333")
	playerID1, _ = uuid.Parse("22222222-2222-bbbb-3333-333333333333")
	playerID2, _ = uuid.Parse("22222222-2222-bbbb-3333-333333333333")
)

// MARK: GET test

func TestPlayers_Get(t *testing.T) {

	playerStore := StubPlayerStore{
		players: []types.Player{
			{
				ID:   playerID1,
				Name: "Chris",
			},
		},
	}

	server := NewPlayerServer(&playerStore, nil)

	tests := []struct {
		name       string
		url        string
		playerName string
		response   string
		code       int
	}{
		{
			"test base",
			"players",
			"Chris",
			`{"ID":"22222222-2222-bbbb-3333-333333333333","Name":"Chris"}` + "\n",
			http.StatusOK,
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

// ==================
// MARK: POST test
// ==================

// This actually records wins, which we are removing
func TestCreatePlayers(t *testing.T) {
	playerStore := StubPlayerStore{
		players: []types.Player{
			{
				ID:   playerID1,
				Name: "Chris",
			},
		},
	}

	server := NewPlayerServer(&playerStore, nil)

	tests := []struct {
		name     string
		response string
		code     int
	}{
		{
			"test post",
			"Pepper",
			http.StatusCreated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			buf := bytes.NewBufferString("body")
			request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", buf)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			if len(playerStore.createPlayerCalls) != 1 {
				t.Errorf("got %d calls to CreatePlayer want %d", len(playerStore.createPlayerCalls), 2)
			}

			if playerStore.createPlayerCalls[0] != "Pepper" {
				t.Errorf("did not store correct winner got %q want %q", playerStore.createPlayerCalls[0], "Pepper")
			}

			assertResponseStatus(t, response.Code, tc.code)
		})
	}
}

// MARK: Play Test

func TestHistory_Get(t *testing.T) {
	tests := []struct {
		name    string
		pID     uuid.UUID
		history types.History
		code    int
	}{
		{
			name: "should play the game - both cooperate",
			pID:  testID1,
			history: types.History{
				types.Interaction{
					Protagonist:     testID1,
					Opponent:        testID2,
					ProtagonistMove: prisoner.Cooperate,
					OpponentMove:    prisoner.Cooperate,
				},
			},
			code: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			historyStore := StubHistoryStore{history: tc.history}
			server := NewPlayerServer(&StubPlayerStore{}, &historyStore)

			request := newHistoryRequest(tc.pID)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			assertResponseStatus(t, response.Code, tc.code)
			assertContentType(t, response, jsonContentType)
			assertHistory(t, getHistoryFromResponse(t, response.Body), tc.history)
		})
	}
}

// MARK: League Test
// ========================================

// MARK: Request Builders
// ====================================

func newHistoryRequest(id uuid.UUID) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/history/%s", id.String()), nil)
	return req
}

func newPlayRequest(id uuid.UUID) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/play/%s", id.String()), nil)
	return req
}

func newGetScoreRequest(route, name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", route, name), nil)
	return req
}

// MARK: Assertions
// ===============================
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

// server_test.go
func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

// MARK: League helper

func getHistoryFromResponse(t testing.TB, body io.Reader) types.History {
	t.Helper()
	h, err := types.NewHistory(body)
	if err != nil {
		t.Fatalf("unable to parse response into History: %v", err)
	}
	return h
}

func assertHistory(t testing.TB, got, want types.History) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("history mismatch\n  got:  %+v\n  want: %+v", got, want)
	}
}
