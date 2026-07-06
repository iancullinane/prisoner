package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iancullinane/prisoner/internal/types"
)

func TestCORS_GetCarriesOriginHeader(t *testing.T) {
	playerStore := StubPlayerStore{
		players: []types.Player{{Name: "Chris"}},
	}
	server := NewPlayerServer(testLogger(), &playerStore, nil)

	request, _ := http.NewRequest(http.MethodGet, "/players", nil)
	request.Header.Set("Origin", "http://localhost:5173")
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assertResponseStatus(t, response.Code, http.StatusOK)
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Errorf("Access-Control-Allow-Origin got %q want %q", got, "http://localhost:5173")
	}
}

func TestCORS_PreflightOptionsRequest(t *testing.T) {
	playerStore := StubPlayerStore{}
	server := NewPlayerServer(testLogger(), &playerStore, nil)

	request, _ := http.NewRequest(http.MethodOptions, "/players", nil)
	request.Header.Set("Origin", "http://localhost:5173")
	request.Header.Set("Access-Control-Request-Method", "POST")
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assertResponseStatus(t, response.Code, http.StatusNoContent)

	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Errorf("Access-Control-Allow-Origin got %q want %q", got, "http://localhost:5173")
	}
	if got := response.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Errorf("Access-Control-Allow-Methods should not be empty")
	}
	if got := response.Header().Get("Access-Control-Allow-Headers"); got == "" {
		t.Errorf("Access-Control-Allow-Headers should not be empty")
	}
	if len(playerStore.getOrCreatePlayerCalls) != 0 {
		t.Errorf("preflight request should not reach the store, got %d calls", len(playerStore.getOrCreatePlayerCalls))
	}
}
