package file

import (
	"os"
	"testing"
)

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

	removeFile := func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}

	return tempfile, removeFile
}

// func newStoreFromLeagueJSON(t testing.TB, leagueJSON string) *FilesSystemPlayerStore {
// 	t.Helper()
// 	return NewFilesSystemPlayerStore(bytes.NewReader([]byte(leagueJSON)))
// }

func assertScoreEquals(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
