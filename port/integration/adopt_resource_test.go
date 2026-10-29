package integration

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

type integrationPlan struct {
	installationType          string
	appType                   string
	adoptExisting             *bool
	createPortResourcesOrigin string
}

func (p integrationPlan) values() map[string]tftypes.Value {
	values := map[string]tftypes.Value{
		"installation_id":       adopttest.String("my-gitlab"),
		"installation_type":     adopttest.String(p.installationType),
		"installation_app_type": adopttest.String(p.appType),
		"title":                 adopttest.String("GitLab"),
	}
	if p.adoptExisting != nil {
		values["adopt_existing"] = adopttest.Bool(*p.adoptExisting)
	}
	if p.createPortResourcesOrigin != "" {
		values["create_port_resources_origin"] = adopttest.String(p.createPortResourcesOrigin)
	}
	return values
}

// selfRegisteredGitlab is a self-hosted GitLab integration as Ocean registers it when it starts
func selfRegisteredGitlab() map[string]any {
	return map[string]any{
		"installationId":      "my-gitlab",
		"installationType":    consts.InstallationTypeOnPrem,
		"installationAppType": "gitlab",
		"title":               "Registered by Ocean",
	}
}

var adoptionCases = []struct {
	name        string
	exists      bool
	policy      cli.AdoptionPolicy
	plan        integrationPlan
	wantError   string
	wantWarning []string
}{
	{name: "missing integration is created", plan: integrationPlan{installationType: consts.InstallationTypeOnPrem, appType: "gitlab"}},
	{name: "existing integration without opt-in fails", exists: true, plan: integrationPlan{installationType: consts.InstallationTypeOnPrem, appType: "gitlab"}, wantError: "integration already exists"},
	{name: "existing integration with resource opt-in is taken over", exists: true, plan: integrationPlan{installationType: consts.InstallationTypeOnPrem, appType: "gitlab", adoptExisting: ptr(true)}, wantWarning: []string{"integration will be taken over"}},
	{name: "existing integration with provider policy is taken over", exists: true, policy: adopttest.Policy(consts.ResourceKindIntegration), plan: integrationPlan{installationType: consts.InstallationTypeOnPrem, appType: "gitlab"}, wantWarning: []string{"integration will be taken over"}},
	{name: "blueprint policy does not cover integrations", exists: true, policy: adopttest.Policy(consts.ResourceKindBlueprint, consts.ResourceKindBlueprintRelation), plan: integrationPlan{installationType: consts.InstallationTypeOnPrem, appType: "gitlab"}, wantError: "integration already exists"},
	{name: "resource opt-out overrides provider policy", exists: true, policy: adopttest.Policy(consts.ResourceKindIntegration), plan: integrationPlan{installationType: consts.InstallationTypeOnPrem, appType: "gitlab", adoptExisting: ptr(false)}, wantError: "integration already exists"},
	{name: "different app type cannot be taken over", exists: true, policy: adopttest.Policy(consts.ResourceKindIntegration), plan: integrationPlan{installationType: consts.InstallationTypeOnPrem, appType: "github-ocean"}, wantError: "cannot take over existing integration"},
	{name: "create_port_resources_origin is reported as ignored", exists: true, policy: adopttest.Policy(consts.ResourceKindIntegration), plan: integrationPlan{installationType: consts.InstallationTypeOnPrem, appType: "gitlab", createPortResourcesOrigin: consts.CreatePortResourcesOriginEmpty}, wantWarning: []string{"integration will be taken over", "create_port_resources_origin is ignored"}},
}

func TestIntegrationModifyPlanAdoption(t *testing.T) {
	for _, tc := range adoptionCases {
		t.Run(tc.name, func(t *testing.T) {
			fake := adopttest.NewFakePort(t)
			if tc.exists {
				fake.AddIntegration("my-gitlab", selfRegisteredGitlab())
			}
			r := &IntegrationResource{portClient: fake.Client(tc.policy)}
			s := adopttest.ResourceSchema(t, r)
			plan := adopttest.Plan(t, s, tc.plan.values())

			resp := &resource.ModifyPlanResponse{Plan: plan}
			r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan, Config: adopttest.Config(plan), State: adopttest.NullState(s)}, resp)

			if tc.wantError == "" {
				assert.Empty(t, adopttest.ErrorSummaries(resp.Diagnostics))
			} else {
				assert.Equal(t, []string{tc.wantError}, adopttest.ErrorSummaries(resp.Diagnostics))
			}
			assert.Equal(t, tc.wantWarning, adopttest.WarningSummaries(resp.Diagnostics))
			assert.Empty(t, fake.Writes(), "planning must never write to Port")
		})
	}
}

func TestIntegrationCreateAdoption(t *testing.T) {
	for _, tc := range adoptionCases {
		t.Run(tc.name, func(t *testing.T) {
			fake := adopttest.NewFakePort(t)
			if tc.exists {
				fake.AddIntegration("my-gitlab", selfRegisteredGitlab())
			}
			r := &IntegrationResource{portClient: fake.Client(tc.policy)}
			s := adopttest.ResourceSchema(t, r)
			plan := adopttest.Plan(t, s, tc.plan.values())

			resp := &resource.CreateResponse{State: adopttest.NullState(s)}
			r.Create(context.Background(), resource.CreateRequest{Plan: plan}, resp)

			switch {
			case tc.wantError != "":
				assert.Equal(t, []string{tc.wantError}, adopttest.ErrorSummaries(resp.Diagnostics))
				assert.Empty(t, fake.Writes(), "an integration that may not be taken over must not be written")
				assert.True(t, resp.State.Raw.IsNull())
			case tc.exists:
				assert.Empty(t, adopttest.ErrorSummaries(resp.Diagnostics))
				assert.Equal(t, []string{"integration was taken over"}, adopttest.WarningSummaries(resp.Diagnostics), "taking over during apply is reported")
				assert.Equal(t, []string{"PATCH /v1/integration/my-gitlab"}, fake.Writes())
				assert.False(t, resp.State.Raw.IsNull())
			default:
				assert.Empty(t, adopttest.ErrorSummaries(resp.Diagnostics))
				assert.Equal(t, []string{"POST /v1/integration"}, fake.Writes())
			}
		})
	}
}

func TestIntegrationModifyPlanInstallationTypeConflict(t *testing.T) {
	fake := adopttest.NewFakePort(t)
	portHosted := selfRegisteredGitlab()
	portHosted["installationType"] = consts.InstallationTypeSaas
	fake.AddIntegration("my-gitlab", portHosted)
	r := &IntegrationResource{portClient: fake.Client(adopttest.Policy(consts.ResourceKindIntegration))}
	s := adopttest.ResourceSchema(t, r)
	plan := adopttest.Plan(t, s, integrationPlan{installationType: consts.InstallationTypeOnPrem, appType: "gitlab"}.values())

	resp := &resource.ModifyPlanResponse{Plan: plan}
	r.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan, Config: adopttest.Config(plan), State: adopttest.NullState(s)}, resp)

	if assert.Len(t, resp.Diagnostics.Errors(), 1) {
		assert.Equal(t, "cannot take over existing integration", resp.Diagnostics.Errors()[0].Summary())
		assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "`installation_type` = \"Saas\", but the configuration sets \"OnPrem\"")
	}
}
