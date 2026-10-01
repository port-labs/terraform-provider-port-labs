package workflow

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

func TestValidateConfigNodeTypeExclusivity(t *testing.T) {
	for _, tc := range []struct {
		name    string
		node    configObject
		invalid bool
	}{
		{"missing", configObject{"identifier": "trigger"}, true},
		{"single", configObject{"identifier": "trigger", "schedule_trigger": configObject{"cron": "0 * * * *"}}, false},
		{"multiple", configObject{"identifier": "trigger", "schedule_trigger": configObject{"cron": "0 * * * *"}, "webhook": configObject{"url": "https://example.invalid"}}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			diags := validateWorkflowConfig(context.Background(), t, configObject{"identifier": "wf", "node": configList{tc.node}}, nil)
			if tc.invalid {
				assert.Contains(t, errorSummaries(diags), "Invalid node type combination")
			} else {
				assert.False(t, diags.HasError(), diags)
			}
		})
	}
}

func TestValidateConfigNodeTypeExclusivityDefersUnknownBlocks(t *testing.T) {
	for _, tc := range []struct {
		name string
		node configObject
	}{
		{"unknown_only", configObject{"identifier": "trigger"}},
		{"known_and_unknown", configObject{"identifier": "trigger", "schedule_trigger": configObject{"cron": "0 * * * *"}}},
		{"multiple_known_and_unknown", configObject{"identifier": "trigger", "schedule_trigger": configObject{"cron": "0 * * * *"}, "webhook": configObject{"url": "https://example.invalid"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			unknownAt := tftypes.NewAttributePath().WithAttributeName("node").WithElementKeyInt(0).WithAttributeName("condition")
			diags := validateWorkflowConfig(context.Background(), t, configObject{"identifier": "wf", "node": configList{tc.node}}, unknownAt)
			assert.False(t, diags.HasError(), diags)
		})
	}
}
