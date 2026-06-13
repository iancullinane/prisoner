package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

// testLogger discards output; these tests assert on HTTP behaviour, not logs.
func testLogger() *slog.Logger { return slog.New(slog.DiscardHandler) }

var (
	// luigi        = types.Player{Name: "Luigi", ID: luigiID}
	luigiID, _   = uuid.Parse("00000000-0000-0000-0000-000000000000")
	testID1, _   = uuid.Parse("00000000-0000-aaaa-2222-222222222222")
	testID2, _   = uuid.Parse("11111111-1111-bbbb-3333-333333333333")
	playerID1, _ = uuid.Parse("11111111-1111-bbbb-3333-333333333333")
	playerID2, _ = uuid.Parse("22222222-2222-bbbb-3333-333333333333")

	playerStore = StubPlayerStore{
		players: []types.Player{
			{ID: playerID1, Name: "Chris"},
			{ID: luigiID, Name: "Luigi"},
			{ID: playerID2, Name: "Pepper"},
		},
	}
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

	server := NewPlayerServer(testLogger(), &playerStore, nil)

	tests := []struct {
		name     string
		url      string
		response string
		code     int
	}{
		{
			"test base",
			"/players/Chris",
			`{"ID":"11111111-1111-bbbb-3333-333333333333","Name":"Chris"}` + "\n",
			http.StatusOK,
		},
		{
			"test player not found",
			"/players/NotChris",
			"could not get player: player not found\n",
			http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, tc.url, nil)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			got := response.Body.String()
			want := tc.response

			assertResponseStatus(t, response.Code, tc.code)
			assertResponseBody(t, got, want)
		})
	}
}

// ===================================
// MARK: POST test
// ===================================

func TestCreatePlayers(t *testing.T) {

	server := NewPlayerServer(testLogger(), &playerStore, nil)

	tests := []struct {
		name          string
		response      string
		code          int
		expectedError error
	}{
		{
			name:          "test post",
			response:      playerID2.String(),
			code:          http.StatusOK,
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			buf := bytes.NewBufferString("body") // if there were a body on post...
			request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", buf)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			if len(playerStore.getOrCreatePlayerCalls) != 1 {
				t.Errorf("got %d calls to GetOrCreatePlayer want %d", len(playerStore.getOrCreatePlayerCalls), 1)
			}

			if playerStore.getOrCreatePlayerCalls[0] != "Pepper" {
				t.Errorf("did not store correct winner got %q want %q", playerStore.getOrCreatePlayerCalls[0], "Pepper")
			}

			assertResponseStatus(t, response.Code, tc.code)
		})
	}
}

// ========================================
// MARK: play test
// ========================================

func TestPlay(t *testing.T) {
	historyStore := StubHistoryStore{}
	server := NewPlayerServer(testLogger(), &playerStore, &historyStore)

	tests := []struct {
		name          string
		playerA       string
		playerB       string
		result        types.Interaction
		code          int
		expectedError error
	}{
		{
			name:    "test play",
			playerA: playerID1.String(),
			playerB: playerID2.String(),
			result: types.Interaction{
				PlayerA:     playerID1,
				PlayerB:     playerID2,
				PlayerAMove: prisoner.Cooperate,
				PlayerBMove: prisoner.Cooperate,
			},
			code:          http.StatusCreated,
			expectedError: nil,
		},
		{
			name:          "player should not be able to player themselves",
			playerA:       playerID1.String(),
			playerB:       playerID1.String(),
			result:        types.Interaction{},
			code:          http.StatusBadRequest,
			expectedError: ErrPlayerCannotPlaySelf,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBufferString("body")
			req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/play/%s/%s", tc.playerA, tc.playerB), buf)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, req)

			if tc.expectedError != nil {
				assertResponseBody(t, response.Body.String(), tc.expectedError.Error()+"\n")
				return
			}

			if len(historyStore.recordInteractionCalls) != 1 {
				t.Errorf("incorrect number of calls on record interaction")
			}

			assertResponseStatus(t, response.Code, tc.code)
			assertInteraction(t, getInterationFromResponse(t, response.Body), tc.result)
		})

	}
}

// ========================================
// MARK: History Test
// ========================================

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
					PlayerA:     testID1,
					PlayerB:     testID2,
					PlayerAMove: prisoner.Cooperate,
					PlayerBMove: prisoner.Cooperate,
				},
			},
			code: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			historyStore := StubHistoryStore{history: tc.history}
			server := NewPlayerServer(testLogger(), &StubPlayerStore{}, &historyStore)

			request := newHistoryRequest(tc.pID)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			assertResponseStatus(t, response.Code, tc.code)
			assertContentType(t, response, jsonContentType)
			assertHistory(t, getHistoryFromResponse(t, response.Body), tc.history)
		})
	}
}

// MARK: Request Builders
// ====================================

func newHistoryRequest(id uuid.UUID) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/history/%s", id.String()), nil)
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

func getInterationFromResponse(t testing.TB, body io.Reader) types.Interaction {
	t.Helper()
	i, err := types.NewInteractionFromJSON(body)
	if err != nil {
		t.Fatalf("%s", err)
	}

	return i
}

func assertInteraction(t testing.TB, got, want types.Interaction) {
	t.Helper()
	// the interaction ID is randomly generated so ignore in comparison
	got.ID = uuid.UUID{}
	want.ID = uuid.UUID{}
	if reflect.DeepEqual(got, want) {
		return
	}
	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		t.Fatalf("marshal got interaction: %v", err)
	}
	wantJSON, err := json.MarshalIndent(want, "", "  ")
	if err != nil {
		t.Fatalf("marshal want interaction: %v", err)
	}
	t.Errorf("interaction mismatch\n  got:  %s\n  want: %s", gotJSON, wantJSON)
}
