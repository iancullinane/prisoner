package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/store/memory"
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
)

// MARK: Healthz test

func TestHealthz(t *testing.T) {
	server := NewPlayerServer(testLogger(), &StubPlayerStore{}, &StubHistoryStore{})

	request, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assertResponseStatus(t, response.Code, http.StatusOK)
	assertResponseBody(t, response.Body.String(), "ok")

}

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
			"/api/v1/players/Chris",
			`{"id":"11111111-1111-bbbb-3333-333333333333","name":"Chris","CreatedAt":"0001-01-01T00:00:00Z"}` + "\n",
			http.StatusOK,
		},
		{
			"test player not found",
			"/api/v1/players/NotChris",
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
		{
			name:          "test post error",
			response:      "could not create player: some error\n",
			code:          http.StatusInternalServerError,
			expectedError: errors.New("some error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			playerStore := StubPlayerStore{
				players: []types.Player{
					{ID: playerID1, Name: "Chris"},
					{ID: luigiID, Name: "Luigi"},
					{ID: playerID2, Name: "Pepper"},
				},
				getOrCreatePlayerError: tc.expectedError,
			}
			server := NewPlayerServer(testLogger(), &playerStore, nil)

			buf := bytes.NewBufferString("body") // if there were a body on post...
			request, _ := http.NewRequest(http.MethodPost, "/api/v1/players/Pepper", buf)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			if len(playerStore.getOrCreatePlayerCalls) != 1 {
				t.Errorf("got %d calls to GetOrCreatePlayer want %d", len(playerStore.getOrCreatePlayerCalls), 1)
			}

			if playerStore.getOrCreatePlayerCalls[0] != "Pepper" {
				t.Errorf("did not store correct winner got %q want %q", playerStore.getOrCreatePlayerCalls[0], "Pepper")
			}

			assertResponseStatus(t, response.Code, tc.code)
			if tc.expectedError != nil {
				assertResponseBody(t, response.Body.String(), tc.response)
			}
		})
	}
}

// ========================================
// MARK: play test
// ========================================

func TestPlay_Refactor(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   PlayRequest
		response      PlayResponse
		code          int
		expectedError error
	}{
		{
			name: "test play with body",
			requestBody: PlayRequest{
				PlayerA:     playerID1,
				PlayerB:     playerID2,
				PlayerAMove: prisoner.Cooperate,
				PlayerBMove: prisoner.Betray,
			},
			response: PlayResponse{
				PlayerAScore: prisoner.Sucker,
				PlayerBScore: prisoner.Temptation,
			},
			code:          http.StatusCreated,
			expectedError: nil,
		},
		{
			name: "player should not be able to play themselves",
			requestBody: PlayRequest{
				PlayerA:     playerID1,
				PlayerB:     playerID1,
				PlayerAMove: prisoner.Cooperate,
				PlayerBMove: prisoner.Betray,
			},
			response:      PlayResponse{},
			code:          http.StatusBadRequest,
			expectedError: ErrPlayerCannotPlaySelf,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			playerStore := StubPlayerStore{
				players: []types.Player{
					{ID: playerID1, Name: "Chris"},
					{ID: luigiID, Name: "Luigi"},
					{ID: playerID2, Name: "Pepper"},
				},
			}
			historyStore := StubHistoryStore{}
			server := NewPlayerServer(testLogger(), &playerStore, &historyStore)

			body, _ := json.Marshal(tc.requestBody)
			req, err := http.NewRequest(http.MethodPost, "/api/v1/play", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			server.ServeHTTP(response, req)

			if tc.expectedError != nil {
				assertResponseBody(t, response.Body.String(), tc.expectedError.Error()+"\n")
				return
			}

			if err != nil {
				t.Errorf("err should be nil: %v", err)
				return
			}

			if len(historyStore.recordInteractionCalls) != 1 {
				t.Errorf("incorrect number of calls on record interaction got %d want 1", len(historyStore.recordInteractionCalls))
				return
			}

			assertPlayResponse(t, getPlayResponseFromResponse(t, response.Body), tc.response)

		})
	}
}

// ========================================
// MARK: History Test
// ========================================

func TestHistory_Get(t *testing.T) {
	t.Run("resolves player names via the store fallback", func(t *testing.T) {
		playerStore := StubPlayerStore{players: []types.Player{
			{ID: testID1, Name: "Alice"},
			{ID: testID2, Name: "Bob"},
		}}
		historyStore := StubHistoryStore{history: types.History{
			types.Interaction{
				PlayerA:     testID1,
				PlayerB:     testID2,
				PlayerAMove: prisoner.Cooperate,
				PlayerBMove: prisoner.Cooperate,
			},
		}}
		server := NewPlayerServer(testLogger(), &playerStore, &historyStore)

		request := newHistoryRequest(testID1)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)

		want := types.PrettyHistory{{
			PlayerAName: "Alice",
			PlayerBName: "Bob",
			PlayerAMove: prisoner.Cooperate,
			PlayerBMove: prisoner.Cooperate,
		}}
		assertPrettyHistory(t, getPrettyHistoryFromResponse(t, response.Body), want)
	})

	t.Run("memory store history comes back pretty, not empty", func(t *testing.T) {
		playerStore := memory.NewInMemoryPlayerStore(testLogger())
		historyStore := memory.NewInMemoryHistoryStore(testLogger())

		alice, err := playerStore.GetOrCreatePlayer("Alice")
		if err != nil {
			t.Fatal(err)
		}
		bob, err := playerStore.GetOrCreatePlayer("Bob")
		if err != nil {
			t.Fatal(err)
		}
		interaction := types.NewInteraction(alice.ID, bob.ID, prisoner.Cooperate, prisoner.Betray)
		if err := historyStore.RecordInteraction(interaction); err != nil {
			t.Fatal(err)
		}

		server := NewPlayerServer(testLogger(), playerStore, historyStore)

		request, _ := http.NewRequest(http.MethodGet, "/api/v1/history", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusOK)
		got := getPrettyHistoryFromResponse(t, response.Body)
		if len(got) != 1 {
			t.Fatalf("got %d pretty interactions, want 1: %+v", len(got), got)
		}
		if got[0].PlayerAName != "Alice" || got[0].PlayerBName != "Bob" {
			t.Errorf("got names %q vs %q, want Alice vs Bob", got[0].PlayerAName, got[0].PlayerBName)
		}
	})

	t.Run("store failure returns 500, not an empty 200", func(t *testing.T) {
		// history references a player the player store doesn't know
		historyStore := StubHistoryStore{history: types.History{
			types.Interaction{PlayerA: testID1, PlayerB: testID2},
		}}
		server := NewPlayerServer(testLogger(), &StubPlayerStore{}, &historyStore)

		request, _ := http.NewRequest(http.MethodGet, "/api/v1/history", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertResponseStatus(t, response.Code, http.StatusInternalServerError)
	})
}

// MARK: Request Builders
// ====================================

func newHistoryRequest(id uuid.UUID) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/history/%s", id.String()), nil)
	return req
}

// MARK: Assertions
// ===============================
func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body wrong\ngot %s\nwant %s", got, want)
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

func getPrettyHistoryFromResponse(t testing.TB, body io.Reader) types.PrettyHistory {
	t.Helper()
	var ph types.PrettyHistory
	if err := json.NewDecoder(body).Decode(&ph); err != nil {
		t.Fatalf("unable to parse response into PrettyHistory: %v", err)
	}
	return ph
}

func assertPrettyHistory(t testing.TB, got, want types.PrettyHistory) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("pretty history mismatch\n  got:  %+v\n  want: %+v", got, want)
	}
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

func getPlayResponseFromResponse(t testing.TB, body io.Reader) PlayResponse {
	t.Helper()
	var resp PlayResponse
	err := json.NewDecoder(body).Decode(&resp)
	if err != nil {
		t.Fatalf("unable to parse response into PlayResponse: %v", err)
	}
	return resp
}

func assertPlayResponse(t testing.TB, got, want PlayResponse) {
	t.Helper()
	got.ID = uuid.UUID{}
	want.ID = uuid.UUID{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("play response mismatch\n  got:  %+v\n  want: %+v", got, want)
	}
}
