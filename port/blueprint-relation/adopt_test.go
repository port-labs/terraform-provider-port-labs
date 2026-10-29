package blueprint_relation

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/adopttest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
)

func ptr(b bool) *bool { return &b }

func relationPlanValues(adoptExisting *bool) map[string]tftypes.Value {
	values := map[string]tftypes.Value{
		"blueprint":  adopttest.String("_team"),
		"identifier": adopttest.String("gitlab_group"),
		"target":     adopttest.String("gitlabGroup"),
		"title":      adopttest.String("GitLab Group"),
		"many":       adopttest.Bool(false),
		"required":   adopttest.Bool(false),
	}
	if adoptExisting != nil {
		values["adopt_existing"] = adopttest.Bool(*adoptExisting)
	}
	return values
}

// teamBlueprint is the _team blueprint, with the relation an integration's default resources added when
// withRelation is set
func teamBlueprint(withRelation bool) map[string]any {
	relations := map[string]any{}
	if withRelation {
		relations["gitlab_group"] = map[string]any{"target": "gitlabGroup", "title": "Added by the integration"}
	}
	return map[string]any{"identifier": "_team", "title": "Team", "relations": relations}
}

var adoptionCases = []struct {
	name           string
	relationExists bool
	policy         cli.AdoptionPolicy
	adoptExisting  *bool
	wantError      string
	wantWarning    string
}{
	{name: "missing relation is created"},
	{name: "missing relation with adoption allowed is created", adoptExisting: ptr(true)},
	{name: "existing relation without opt-in fails", relationExists: true, wantError: "relation already exists"},
	{name: "existing relation with resource opt-in is taken over", relationExists: true, adoptExisting: ptr(true), wantWarning: "relation will be taken over"},
	{name: "existing relation with provider policy is taken over", relationExists: true, policy: adopttest.Policy(consts.ResourceKindBlueprintRelation), wantWarning: "relation will be taken over"},
	{name: "blueprint policy does not cover relations", relationExists: true, policy: adopttest.Policy(consts.ResourceKindBlueprint), wantError: "relation already exists"},
	{name: "resource opt-out overrides provider policy", relationExists: true, policy: adopttest.Policy(consts.ResourceKindBlueprintRelation), adoptExisting: ptr(false), wantError: "relation already exists"},
}

func TestRelationModifyPlanAdoption(t *testing.T) {
	for _, tc := range adoptionCases {
		t.Run(tc.name, func(t *testing.T) {
			fake := adopttest.NewFakePort(t)
			fake.AddBlueprint("_team", teamBlueprint(tc.relationExists))
			r := &BlueprintRelationResource{portClient: fake.Client(tc.policy)}
			s := adopttest.ResourceSchema(t, r)
			plan := adopttest.Plan(t, s, relationPlanValues(tc.adoptExisting))

			resp := &resource.ModifyPlanResponse{Plan: plan}
			r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan, State: adopttest.NullState(s)}, resp)

			assertSummaries(t, tc.wantError, adopttest.ErrorSummaries(resp.Diagnostics))
			assertSummaries(t, tc.wantWarning, adopttest.WarningSummaries(resp.Diagnostics))
			assert.Empty(t, fake.Writes(), "planning must never write to Port")
		})
	}
}

func TestRelationCreateAdoption(t *testing.T) {
	for _, tc := range adoptionCases {
		t.Run(tc.name, func(t *testing.T) {
			fake := adopttest.NewFakePort(t)
			fake.AddBlueprint("_team", teamBlueprint(tc.relationExists))
			r := &BlueprintRelationResource{portClient: fake.Client(tc.policy)}
			s := adopttest.ResourceSchema(t, r)
			plan := adopttest.Plan(t, s, relationPlanValues(tc.adoptExisting))

			resp := &resource.CreateResponse{State: adopttest.NullState(s)}
			r.Create(context.Background(), resource.CreateRequest{Plan: plan}, resp)

			relation := fake.Blueprint("_team")["relations"].(map[string]any)["gitlab_group"]
			if tc.wantError != "" {
				assert.Equal(t, []string{tc.wantError}, adopttest.ErrorSummaries(resp.Diagnostics))
				assert.Empty(t, fake.Writes(), "a relation that may not be taken over must not be written")
				assert.Equal(t, "Added by the integration", relation.(map[string]any)["title"])
				assert.True(t, resp.State.Raw.IsNull())
				return
			}
			assert.Empty(t, adopttest.ErrorSummaries(resp.Diagnostics))
			if tc.relationExists {
				assert.Equal(t, []string{"relation was taken over"}, adopttest.WarningSummaries(resp.Diagnostics), "taking over during apply is reported")
			} else {
				assert.Empty(t, adopttest.WarningSummaries(resp.Diagnostics))
			}
			assert.Equal(t, []string{"PATCH /v1/blueprints/_team"}, fake.Writes())
			assert.Equal(t, "GitLab Group", relation.(map[string]any)["title"])
		})
	}
}

func TestRelationModifyPlanMissingBlueprint(t *testing.T) {
	fake := adopttest.NewFakePort(t)
	r := &BlueprintRelationResource{portClient: fake.Client(nil)}
	s := adopttest.ResourceSchema(t, r)
	plan := adopttest.Plan(t, s, relationPlanValues(nil))

	resp := &resource.ModifyPlanResponse{Plan: plan}
	r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan, State: adopttest.NullState(s)}, resp)

	assert.Empty(t, resp.Diagnostics, "a blueprint created in the same apply has no relations to conflict with")
}

func TestRelationExistsErrorNamesImportID(t *testing.T) {
	fake := adopttest.NewFakePort(t)
	fake.AddBlueprint("_team", teamBlueprint(true))
	r := &BlueprintRelationResource{portClient: fake.Client(nil)}
	s := adopttest.ResourceSchema(t, r)
	plan := adopttest.Plan(t, s, relationPlanValues(nil))

	resp := &resource.ModifyPlanResponse{Plan: plan}
	r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan, State: adopttest.NullState(s)}, resp)

	if assert.Len(t, resp.Diagnostics.Errors(), 1) {
		detail := resp.Diagnostics.Errors()[0].Detail()
		assert.Contains(t, detail, "`terraform import <address> \"_team:gitlab_group\"`")
		assert.Contains(t, detail, `add "port_blueprint_relation" to the provider's`)
	}
}

func assertSummaries(t *testing.T, want string, got []string) {
	t.Helper()
	if want == "" {
		assert.Empty(t, got)
		return
	}
	assert.Equal(t, []string{want}, got)
}
