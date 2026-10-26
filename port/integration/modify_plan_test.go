package integration

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanSpec(t *testing.T) {
	state := types.StringValue(`{"integrationSpec":{"githubToken":"_OLD"},"appSpec":{"liveEventsEnabled":false}}`)

	tests := []struct {
		name   string
		config string
		want   string
	}{
		{
			name:   "inherits appSpec when only integrationSpec changes",
			config: `{"integrationSpec":{"githubToken":"_NEW"}}`,
			want:   `{"appSpec":{"liveEventsEnabled":false},"integrationSpec":{"githubToken":"_NEW"}}`,
		},
		{
			name:   "keeps state when the configuration is unchanged",
			config: `{"integrationSpec":{"githubToken":"_OLD"}}`,
			want:   state.ValueString(),
		},
		{
			name:   "configuration wins when it manages appSpec",
			config: `{"integrationSpec":{"githubToken":"_OLD"},"appSpec":{"liveEventsEnabled":true}}`,
			want:   `{"appSpec":{"liveEventsEnabled":true},"integrationSpec":{"githubToken":"_OLD"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := planSpec(types.StringValue(tt.config), state)
			assert.JSONEq(t, tt.want, got.ValueString())
		})
	}
}

func TestExtractAppSpec(t *testing.T) {
	got := extractAppSpec(types.StringValue(
		`{"integrationSpec":{"githubAppPrivateKey":"x"},"appSpec":{"scheduledResyncInterval":"12h"}}`,
	))
	assert.JSONEq(t, `{"appSpec":{"scheduledResyncInterval":"12h"}}`, got.ValueString())

	assert.True(t, extractAppSpec(types.StringValue(`{"integrationSpec":{"token":"x"}}`)).IsNull())
	assert.True(t, extractAppSpec(types.StringNull()).IsNull())
	assert.True(t, extractAppSpec(types.StringUnknown()).IsUnknown())
}

func integrationSchema() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":                             tftypes.String,
			"installation_id":                tftypes.String,
			"installation_app_type":          tftypes.String,
			"installation_type":              tftypes.String,
			"title":                          tftypes.String,
			"version":                        tftypes.String,
			"config":                         tftypes.String,
			"spec":                           tftypes.String,
			"create_port_resources_origin":   tftypes.String,
			"status":                         tftypes.String,
			"kafka_changelog_destination":    tftypes.Object{AttributeTypes: map[string]tftypes.Type{}},
			"webhook_changelog_destination":  tftypes.Object{AttributeTypes: map[string]tftypes.Type{"url": tftypes.String, "agent": tftypes.Bool}},
		},
	}
}

func makeState(installationType string) tfsdk.State {
	schema := integrationSchema()
	val := tftypes.NewValue(schema, map[string]tftypes.Value{
		"id":                             tftypes.NewValue(tftypes.String, "test-id"),
		"installation_id":                tftypes.NewValue(tftypes.String, "test-id"),
		"installation_app_type":          tftypes.NewValue(tftypes.String, "kafka"),
		"installation_type":              tftypes.NewValue(tftypes.String, installationType),
		"title":                          tftypes.NewValue(tftypes.String, "Test"),
		"version":                        tftypes.NewValue(tftypes.String, "1.0.0"),
		"config":                         tftypes.NewValue(tftypes.String, nil),
		"spec":                           tftypes.NewValue(tftypes.String, nil),
		"create_port_resources_origin":   tftypes.NewValue(tftypes.String, nil),
		"status":                         tftypes.NewValue(tftypes.String, nil),
		"kafka_changelog_destination":    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}, nil),
		"webhook_changelog_destination":  tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"url": tftypes.String, "agent": tftypes.Bool}}, nil),
	})
	return tfsdk.State{Raw: val, Schema: IntegrationSchemaForTest()}
}

func makePlan(installationType *string) tfsdk.Plan {
	schema := integrationSchema()
	var installationTypeVal tftypes.Value
	if installationType == nil {
		installationTypeVal = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	} else {
		installationTypeVal = tftypes.NewValue(tftypes.String, *installationType)
	}
	val := tftypes.NewValue(schema, map[string]tftypes.Value{
		"id":                             tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"installation_id":                tftypes.NewValue(tftypes.String, "test-id"),
		"installation_app_type":          tftypes.NewValue(tftypes.String, "kafka"),
		"installation_type":              installationTypeVal,
		"title":                          tftypes.NewValue(tftypes.String, "Test"),
		"version":                        tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"config":                         tftypes.NewValue(tftypes.String, nil),
		"spec":                           tftypes.NewValue(tftypes.String, nil),
		"create_port_resources_origin":   tftypes.NewValue(tftypes.String, nil),
		"status":                         tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"kafka_changelog_destination":    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}, nil),
		"webhook_changelog_destination":  tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"url": tftypes.String, "agent": tftypes.Bool}}, nil),
	})
	return tfsdk.Plan{Raw: val, Schema: IntegrationSchemaForTest()}
}

func makeConfig(installationType *string) tfsdk.Config {
	schema := integrationSchema()
	var installationTypeVal tftypes.Value
	if installationType == nil {
		installationTypeVal = tftypes.NewValue(tftypes.String, nil)
	} else {
		installationTypeVal = tftypes.NewValue(tftypes.String, *installationType)
	}
	val := tftypes.NewValue(schema, map[string]tftypes.Value{
		"id":                             tftypes.NewValue(tftypes.String, nil),
		"installation_id":                tftypes.NewValue(tftypes.String, "test-id"),
		"installation_app_type":          tftypes.NewValue(tftypes.String, "kafka"),
		"installation_type":              installationTypeVal,
		"title":                          tftypes.NewValue(tftypes.String, "Test"),
		"version":                        tftypes.NewValue(tftypes.String, nil),
		"config":                         tftypes.NewValue(tftypes.String, nil),
		"spec":                           tftypes.NewValue(tftypes.String, nil),
		"create_port_resources_origin":   tftypes.NewValue(tftypes.String, nil),
		"status":                         tftypes.NewValue(tftypes.String, nil),
		"kafka_changelog_destination":    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}, nil),
		"webhook_changelog_destination":  tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"url": tftypes.String, "agent": tftypes.Bool}}, nil),
	})
	return tfsdk.Config{Raw: val, Schema: IntegrationSchemaForTest()}
}

func nullState() tfsdk.State {
	return tfsdk.State{Raw: tftypes.NewValue(integrationSchema(), nil), Schema: IntegrationSchemaForTest()}
}

func TestModifyPlan_InstallationType_CreateDefaultsToOnPrem(t *testing.T) {
	r := &IntegrationResource{}
	ctx := context.Background()

	req := resource.ModifyPlanRequest{
		State:  nullState(),
		Plan:   makePlan(nil),
		Config: makeConfig(nil),
	}
	resp := &resource.ModifyPlanResponse{Plan: req.Plan}
	r.ModifyPlan(ctx, req, resp)

	require.False(t, resp.Diagnostics.HasError(), "unexpected error: %v", resp.Diagnostics)

	var result IntegrationModel
	resp.Plan.Get(ctx, &result)
	assert.Equal(t, consts.InstallationTypeOnPrem, result.InstallationType.ValueString())
}

func TestModifyPlan_InstallationType_CreateWithExplicitSaas(t *testing.T) {
	r := &IntegrationResource{}
	ctx := context.Background()

	saas := consts.InstallationTypeSaas
	req := resource.ModifyPlanRequest{
		State:  nullState(),
		Plan:   makePlan(&saas),
		Config: makeConfig(&saas),
	}
	resp := &resource.ModifyPlanResponse{Plan: req.Plan}
	r.ModifyPlan(ctx, req, resp)

	var result IntegrationModel
	resp.Plan.Get(ctx, &result)
	assert.Equal(t, consts.InstallationTypeSaas, result.InstallationType.ValueString())
}

func TestModifyPlan_InstallationType_UpdatePreservesState_OnPrem(t *testing.T) {
	r := &IntegrationResource{}
	ctx := context.Background()

	req := resource.ModifyPlanRequest{
		State:  makeState(consts.InstallationTypeOnPrem),
		Plan:   makePlan(nil),
		Config: makeConfig(nil),
	}
	resp := &resource.ModifyPlanResponse{Plan: req.Plan}
	r.ModifyPlan(ctx, req, resp)

	require.False(t, resp.Diagnostics.HasError(), "unexpected error: %v", resp.Diagnostics)

	var result IntegrationModel
	resp.Plan.Get(ctx, &result)
	assert.Equal(t, consts.InstallationTypeOnPrem, result.InstallationType.ValueString())
}

func makeStateWithSpec(installationType, spec string) tfsdk.State {
	schema := integrationSchema()
	val := tftypes.NewValue(schema, map[string]tftypes.Value{
		"id":                             tftypes.NewValue(tftypes.String, "test-id"),
		"installation_id":                tftypes.NewValue(tftypes.String, "test-id"),
		"installation_app_type":          tftypes.NewValue(tftypes.String, "linear"),
		"installation_type":              tftypes.NewValue(tftypes.String, installationType),
		"title":                          tftypes.NewValue(tftypes.String, "Test"),
		"version":                        tftypes.NewValue(tftypes.String, "1.0.0"),
		"config":                         tftypes.NewValue(tftypes.String, nil),
		"spec":                           tftypes.NewValue(tftypes.String, spec),
		"create_port_resources_origin":   tftypes.NewValue(tftypes.String, nil),
		"status":                         tftypes.NewValue(tftypes.String, nil),
		"kafka_changelog_destination":    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}, nil),
		"webhook_changelog_destination":  tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"url": tftypes.String, "agent": tftypes.Bool}}, nil),
	})
	return tfsdk.State{Raw: val, Schema: IntegrationSchemaForTest()}
}

func makePlanWithSpec(installationType *string, spec string) tfsdk.Plan {
	schema := integrationSchema()
	var installationTypeVal tftypes.Value
	if installationType == nil {
		installationTypeVal = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	} else {
		installationTypeVal = tftypes.NewValue(tftypes.String, *installationType)
	}
	val := tftypes.NewValue(schema, map[string]tftypes.Value{
		"id":                             tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"installation_id":                tftypes.NewValue(tftypes.String, "test-id"),
		"installation_app_type":          tftypes.NewValue(tftypes.String, "linear"),
		"installation_type":              installationTypeVal,
		"title":                          tftypes.NewValue(tftypes.String, "Test"),
		"version":                        tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"config":                         tftypes.NewValue(tftypes.String, nil),
		"spec":                           tftypes.NewValue(tftypes.String, spec),
		"create_port_resources_origin":   tftypes.NewValue(tftypes.String, nil),
		"status":                         tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"kafka_changelog_destination":    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}, nil),
		"webhook_changelog_destination":  tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"url": tftypes.String, "agent": tftypes.Bool}}, nil),
	})
	return tfsdk.Plan{Raw: val, Schema: IntegrationSchemaForTest()}
}

func makeConfigWithSpec(installationType *string, spec string) tfsdk.Config {
	schema := integrationSchema()
	var installationTypeVal tftypes.Value
	if installationType == nil {
		installationTypeVal = tftypes.NewValue(tftypes.String, nil)
	} else {
		installationTypeVal = tftypes.NewValue(tftypes.String, *installationType)
	}
	val := tftypes.NewValue(schema, map[string]tftypes.Value{
		"id":                             tftypes.NewValue(tftypes.String, nil),
		"installation_id":                tftypes.NewValue(tftypes.String, "test-id"),
		"installation_app_type":          tftypes.NewValue(tftypes.String, "linear"),
		"installation_type":              installationTypeVal,
		"title":                          tftypes.NewValue(tftypes.String, "Test"),
		"version":                        tftypes.NewValue(tftypes.String, nil),
		"config":                         tftypes.NewValue(tftypes.String, nil),
		"spec":                           tftypes.NewValue(tftypes.String, spec),
		"create_port_resources_origin":   tftypes.NewValue(tftypes.String, nil),
		"status":                         tftypes.NewValue(tftypes.String, nil),
		"kafka_changelog_destination":    tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}, nil),
		"webhook_changelog_destination":  tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"url": tftypes.String, "agent": tftypes.Bool}}, nil),
	})
	return tfsdk.Config{Raw: val, Schema: IntegrationSchemaForTest()}
}

func TestModifyPlan_InstallationType_UpdatePreservesState_Saas(t *testing.T) {
	r := &IntegrationResource{}
	ctx := context.Background()

	saasSpec := `{"integrationSpec":{"linearApiKey":"secret"}}`
	req := resource.ModifyPlanRequest{
		State:  makeStateWithSpec(consts.InstallationTypeSaas, saasSpec),
		Plan:   makePlanWithSpec(nil, saasSpec),
		Config: makeConfigWithSpec(nil, saasSpec),
	}
	resp := &resource.ModifyPlanResponse{Plan: req.Plan}
	r.ModifyPlan(ctx, req, resp)

	require.False(t, resp.Diagnostics.HasError(), "unexpected error: %v", resp.Diagnostics)

	var result IntegrationModel
	resp.Plan.Get(ctx, &result)
	assert.Equal(t, consts.InstallationTypeSaas, result.InstallationType.ValueString())
}

func TestModifyPlan_InstallationType_ExplicitChangeRejected(t *testing.T) {
	r := &IntegrationResource{}
	ctx := context.Background()

	onprem := consts.InstallationTypeOnPrem
	req := resource.ModifyPlanRequest{
		State:  makeState(consts.InstallationTypeSaas),
		Plan:   makePlan(&onprem),
		Config: makeConfig(&onprem),
	}
	resp := &resource.ModifyPlanResponse{Plan: req.Plan}
	r.ModifyPlan(ctx, req, resp)

	require.True(t, resp.Diagnostics.HasError())
	found := false
	for _, diag := range resp.Diagnostics.Errors() {
		if diag.Summary() == "cannot change installation_type" {
			found = true
		}
	}
	assert.True(t, found, "expected 'cannot change installation_type' error")
}
