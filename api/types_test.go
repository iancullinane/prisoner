package api

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

// TestPlayRequestResponseJSONCasing locks PlayRequest/PlayResponse to the
// same camelCase convention already used by internal/types.Interaction, so
// the API's wire format is consistent across endpoints.
func TestPlayRequestResponseJSONCasing(t *testing.T) {
	req := PlayRequest{
		PlayerA:     uuid.New(),
		PlayerB:     uuid.New(),
		PlayerAMove: prisoner.Cooperate,
		PlayerBMove: prisoner.Betray,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal PlayRequest: %v", err)
	}

	var reqFields map[string]json.RawMessage
	if err := json.Unmarshal(reqBytes, &reqFields); err != nil {
		t.Fatalf("unmarshal PlayRequest into map: %v", err)
	}

	for _, key := range []string{"playerA", "playerB", "playerAMove", "playerBMove"} {
		if _, ok := reqFields[key]; !ok {
			t.Errorf("PlayRequest JSON missing camelCase key %q, got %v", key, reqFields)
		}
	}

	resp := PlayResponse{
		ID:           uuid.New(),
		PlayerAScore: prisoner.Reward,
		PlayerBScore: prisoner.Punish,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal PlayResponse: %v", err)
	}

	var respFields map[string]json.RawMessage
	if err := json.Unmarshal(respBytes, &respFields); err != nil {
		t.Fatalf("unmarshal PlayResponse into map: %v", err)
	}

	for _, key := range []string{"id", "playerAScore", "playerBScore"} {
		if _, ok := respFields[key]; !ok {
			t.Errorf("PlayResponse JSON missing camelCase key %q, got %v", key, respFields)
		}
	}
}
