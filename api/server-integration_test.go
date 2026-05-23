package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/iancullinane/prisoner/internal/store/file"
	"github.com/iancullinane/prisoner/internal/store/memory"
	"github.com/iancullinane/prisoner/internal/store/testhelpers"
)

func TestRecordingWinsAndRetrievingThemFromFile(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := file.NewFileSystemPlayerStore(database)
	assertNoError(t, err)

	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostRequest("players", player))
	server.ServeHTTP(httptest.NewRecorder(), newPostRequest("players", player))
	server.ServeHTTP(httptest.NewRecorder(), newPostRequest("players", player))

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
		testhelpers.AssertLeague(t, got, want)
	})
}

func TestRecordingWinsAndRetrievingThemInMemory(t *testing.T) {
	store := memory.NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostRequest("players", player))
	server.ServeHTTP(httptest.NewRecorder(), newPostRequest("players", player))
	server.ServeHTTP(httptest.NewRecorder(), newPostRequest("players", player))

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
		testhelpers.AssertLeague(t, got, want)
	})
}

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tempfile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create temp file for file store test %v", err)
	}

	if _, err := tempfile.Write([]byte(initialData)); err != nil {
		tempfile.Close()
		os.Remove(tempfile.Name())
		t.Fatalf("could not write initial data to temp file %v", err)
	}

	return tempfile, func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got an error but didn't want one %v", err)
	}
}
