package cli

import (
	"encoding/json"
	"testing"
)

func TestPortBodyScorecardGroupsUnmarshal(t *testing.T) {
	payload := `{
		"ok": true,
		"scorecardGroups": [
			{
				"identifier": "production-readiness",
				"title": "Production Readiness"
			}
		]
	}`

	var pb PortBody
	if err := json.Unmarshal([]byte(payload), &pb); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if !pb.OK {
		t.Fatal("expected ok=true")
	}
	if len(pb.ScorecardGroups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(pb.ScorecardGroups))
	}
	if pb.ScorecardGroups[0].Identifier != "production-readiness" {
		t.Fatalf("unexpected identifier: %q", pb.ScorecardGroups[0].Identifier)
	}
}
