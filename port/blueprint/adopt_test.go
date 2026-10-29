package blueprint

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/adopttest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
)

func ptr(b bool) *bool { return &b }

func blueprintPlanValues(adoptExisting *bool) map[string]tftypes.Value {
	values := map[string]tftypes.Value{
		"identifier": adopttest.String("gitlabGroup"),
		"title":      adopttest.String("GitLab Group"),
	}
	if adoptExisting != nil {
		values["adopt_existing"] = adopttest.Bool(*adoptExisting)
	}
	return values
}

// existingIntegrationBlueprint is a blueprint as an integration's default resources create it, with a relation
// and an aggregation property that Terraform does not declare
func existingIntegrationBlueprint() map[string]any {
	return map[string]any{
		"identifier": "gitlabGroup",
		"title":      "Created by the integration",
		"schema":     map[string]any{"properties": map[string]any{}},
		"relations": map[string]any{
			"parent": map[string]any{"target": "gitlabGroup", "title": "Parent"},
		},
		"aggregationProperties": map[string]any{
			"projects_count": map[string]any{"title": "Projects", "target": "project", "calculationSpec": map[string]any{"func": "count", "calculationBy": "entities"}},
		},
	}
}

var adoptionCases = []struct {
	name          string
	exists        bool
	policy        cli.AdoptionPolicy
	adoptExisting *bool
	wantError     string
	wantWarning   string
}{
	{name: "missing blueprint is created", exists: false},
	{name: "missing blueprint with adoption allowed is created", exists: false, policy: adopttest.Policy(consts.ResourceKindBlueprint)},
	{name: "existing blueprint without opt-in fails", exists: true, wantError: "blueprint already exists"},
	{name: "existing blueprint with resource opt-in is taken over", exists: true, adoptExisting: ptr(true), wantWarning: "blueprint will be taken over"},
	{name: "existing blueprint with provider policy is taken over", exists: true, policy: adopttest.Policy(consts.ResourceKindBlueprint), wantWarning: "blueprint will be taken over"},
	{name: "provider policy for other kinds does not apply", exists: true, policy: adopttest.Policy(consts.ResourceKindPage, consts.ResourceKindBlueprintRelation), wantError: "blueprint already exists"},
	{name: "resource opt-out overrides provider policy", exists: true, policy: adopttest.Policy(consts.ResourceKindBlueprint), adoptExisting: ptr(false), wantError: "blueprint already exists"},
}

func TestBlueprintModifyPlanAdoption(t *testing.T) {
	for _, tc := range adoptionCases {
		t.Run(tc.name, func(t *testing.T) {
			fake := adopttest.NewFakePort(t)
			if tc.exists {
				fake.AddBlueprint("gitlabGroup", existingIntegrationBlueprint())
			}
			r := &BlueprintResource{portClient: fake.Client(tc.policy)}
			s := adopttest.ResourceSchema(t, r)
			plan := adopttest.Plan(t, s, blueprintPlanValues(tc.adoptExisting))

			resp := &resource.ModifyPlanResponse{Plan: plan}
			r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan, State: adopttest.NullState(s)}, resp)

			assertSummaries(t, tc.wantError, adopttest.ErrorSummaries(resp.Diagnostics))
			assertSummaries(t, tc.wantWarning, adopttest.WarningSummaries(resp.Diagnostics))
			assert.Empty(t, fake.Writes(), "planning must never write to Port")
		})
	}
}

func TestBlueprintCreateAdoption(t *testing.T) {
	for _, tc := range adoptionCases {
		t.Run(tc.name, func(t *testing.T) {
			fake := adopttest.NewFakePort(t)
			if tc.exists {
				fake.AddBlueprint("gitlabGroup", existingIntegrationBlueprint())
			}
			r := &BlueprintResource{portClient: fake.Client(tc.policy)}

			if tc.wantError != "" {
				s := adopttest.ResourceSchema(t, r)
				plan := adopttest.Plan(t, s, blueprintPlanValues(tc.adoptExisting))
				resp := &resource.CreateResponse{State: adopttest.NullState(s)}
				r.Create(context.Background(), resource.CreateRequest{Plan: plan}, resp)

				assert.Equal(t, []string{tc.wantError}, adopttest.ErrorSummaries(resp.Diagnostics))
				assert.Empty(t, fake.Writes(), "a blueprint that may not be taken over must not be written")
				assert.Equal(t, "Created by the integration", fake.Blueprint("gitlabGroup")["title"])
				return
			}

			// Create also records managed relations in private state, which the framework does not let tests
			// construct, so the write itself is checked through createOrTakeOverBlueprint
			var diags diag.Diagnostics
			r.createOrTakeOverBlueprint(context.Background(), gitlabGroupBlueprint(nil), nil, boolValue(tc.adoptExisting), &diags)
			assert.Empty(t, adopttest.ErrorSummaries(diags))

			if tc.exists {
				assert.Equal(t, []string{"blueprint was taken over"}, adopttest.WarningSummaries(diags), "taking over during apply is reported")
				assert.Equal(t, []string{"PUT /v1/blueprints/gitlabGroup"}, fake.Writes())
				adopted := fake.Blueprint("gitlabGroup")
				assert.Equal(t, "GitLab Group", adopted["title"])
				assert.Contains(t, adopted["relations"], "parent", "relations Terraform does not declare are kept")
				assert.Contains(t, adopted["aggregationProperties"], "projects_count", "aggregation properties are kept")
			} else {
				assert.Equal(t, []string{"POST /v1/blueprints"}, fake.Writes())
			}
		})
	}
}

func TestBlueprintCreateAdoptionDeclaredRelationWins(t *testing.T) {
	fake := adopttest.NewFakePort(t)
	fake.AddBlueprint("gitlabGroup", existingIntegrationBlueprint())
	r := &BlueprintResource{portClient: fake.Client(adopttest.Policy(consts.ResourceKindBlueprint))}

	target, title := "gitlabGroup", "Declared parent"
	declared := gitlabGroupBlueprint(map[string]cli.Relation{"parent": {Target: &target, Title: &title}})
	var diags diag.Diagnostics
	r.createOrTakeOverBlueprint(context.Background(), declared, nil, types.BoolValue(true), &diags)
	assert.Empty(t, adopttest.ErrorSummaries(diags))

	parent := fake.Blueprint("gitlabGroup")["relations"].(map[string]any)["parent"].(map[string]any)
	assert.Equal(t, "Declared parent", parent["title"], "a relation declared in Terraform replaces the existing one")
}

func TestBlueprintCreateTakesOverBlueprintCreatedConcurrently(t *testing.T) {
	testCases := []struct {
		name       string
		policy     cli.AdoptionPolicy
		wantError  []string
		wantWrites []string
	}{
		{name: "without opt-in fails", wantError: []string{"blueprint already exists"}, wantWrites: []string{"POST /v1/blueprints"}},
		{name: "with provider policy is taken over", policy: adopttest.Policy(consts.ResourceKindBlueprint), wantWrites: []string{"POST /v1/blueprints", "PUT /v1/blueprints/gitlabGroup"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fake := adopttest.NewFakePort(t)
			fake.BlueprintAppearsOnCreate("gitlabGroup", existingIntegrationBlueprint())
			r := &BlueprintResource{portClient: fake.Client(tc.policy)}

			var diags diag.Diagnostics
			r.createOrTakeOverBlueprint(context.Background(), gitlabGroupBlueprint(nil), nil, types.BoolNull(), &diags)

			assert.Equal(t, tc.wantError, adopttest.ErrorSummaries(diags))
			assert.Equal(t, tc.wantWrites, fake.Writes())
		})
	}
}

func gitlabGroupBlueprint(relations map[string]cli.Relation) *cli.Blueprint {
	if relations == nil {
		relations = map[string]cli.Relation{}
	}
	return &cli.Blueprint{
		Identifier:            "gitlabGroup",
		Title:                 "GitLab Group",
		Schema:                cli.BlueprintSchema{Properties: map[string]cli.BlueprintProperty{}},
		CalculationProperties: map[string]cli.BlueprintCalculationProperty{},
		AggregationProperties: map[string]cli.BlueprintAggregationProperty{},
		MirrorProperties:      map[string]cli.BlueprintMirrorProperty{},
		Relations:             relations,
	}
}

func boolValue(b *bool) types.Bool {
	if b == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*b)
}

func assertSummaries(t *testing.T, want string, got []string) {
	t.Helper()
	if want == "" {
		assert.Empty(t, got)
		return
	}
	assert.Equal(t, []string{want}, got)
}
