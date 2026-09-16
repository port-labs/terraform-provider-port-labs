package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
)

func TestValidatePlanRules_ConfigOnCreate(t *testing.T) {
	plan := &IntegrationModel{
		InstallationType: types.StringValue(consts.InstallationTypeSaas),
		Config:           types.StringValue(`{"resources":[]}`),
	}

	diags := diag.Diagnostics{}
	validatePlanRules(plan, nil, true, &diags)
	assert.True(t, diags.HasError())
}

func TestValidatePlanRules_ChangelogDestinationRemoval(t *testing.T) {
	state := &IntegrationModel{
		WebhookChangelogDestination: types.ObjectValueMust(webhookChangelogDestinationType, map[string]attr.Value{
			"url":   types.StringValue("https://example.com"),
			"agent": types.BoolValue(false),
		}),
	}
	plan := &IntegrationModel{
		WebhookChangelogDestination: types.ObjectNull(webhookChangelogDestinationType),
	}

	diags := diag.Diagnostics{}
	validatePlanRules(plan, state, false, &diags)
	assert.True(t, diags.HasError())
}
