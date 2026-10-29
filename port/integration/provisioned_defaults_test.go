package integration_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

// gitlabProvisionedDefaultsHCL installs a Port hosted GitLab integration, whose default resources Port provisions
// on creation (the gitlabGroup blueprint and the _team.gitlab_group relation among them), and declares those two in
// the same configuration, as a customer managing them with Terraform would
func gitlabProvisionedDefaultsHCL(providerConfig, installationID, secretName string) string {
	return providerConfig + fmt.Sprintf(`
resource "port_organization_secret" "gitlab_token" {
	secret_name  = "%s"
	secret_value = "glpat-acctest-dummy"
}

resource "port_integration" "gitlab" {
	depends_on = [port_organization_secret.gitlab_token]

	installation_id              = "%s"
	installation_app_type        = "gitlab-v2"
	installation_type            = "Saas"
	title                        = "GitLab acctest"
	create_port_resources_origin = "Port"

	spec = jsonencode({
		integrationSpec = {
			gitlabToken = port_organization_secret.gitlab_token.secret_name
		}
	})
}

resource "port_blueprint" "gitlab_group" {
	depends_on = [port_integration.gitlab]

	identifier = "gitlabGroup"
	title      = "Group"
	icon       = "GitLab"
}

resource "port_blueprint_relation" "team_gitlab_group" {
	blueprint  = "_team"
	identifier = "gitlab_group"
	target     = port_blueprint.gitlab_group.identifier
	title      = "GitLab Group"
}
`, secretName, installationID)
}

// cleanupGitlabProvisionedDefaults removes what Port provisioned for the integration and Terraform did not manage,
// retrying blueprints until the ones relating to others are gone
func cleanupGitlabProvisionedDefaults(t *testing.T, portClient *cli.PortClient, installationID, secretName string) {
	ctx := context.Background()
	_, _ = portClient.DeleteIntegration(ctx, installationID)
	_ = portClient.DeleteOrganizationSecret(ctx, secretName)
	_, _ = portClient.DeleteBlueprintRelation(ctx, "_team", "gitlab_group")
	for attempt := 0; attempt < 5; attempt++ {
		remaining := 0
		for _, identifier := range []string{"gitlabMergeRequest", "gitlabRepository", "gitlabGroup", "gitlabMember", "gitlabPipeline", "gitlabJob", "gitlabIssue"} {
			if _, statusCode, err := portClient.ReadBlueprint(ctx, identifier); err != nil && statusCode == 404 {
				continue
			}
			if err := portClient.DeleteBlueprint(ctx, identifier); err != nil {
				remaining++
			}
		}
		if remaining == 0 {
			return
		}
	}
	t.Logf("some gitlab blueprints provisioned by Port could not be removed")
}

func blueprintHasRelation(portClient *cli.PortClient, blueprint, relation string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		b, _, err := portClient.ReadBlueprint(context.Background(), blueprint)
		if err != nil {
			return err
		}
		if _, ok := b.Relations[relation]; !ok {
			return fmt.Errorf("blueprint %q has no relation %q", blueprint, relation)
		}
		return nil
	}
}

// TestAccPortIntegrationProvisionedDefaultsAreTakenOver is case 2 of the side effects Terraform could not adopt:
// installing the GitLab integration makes Port create the gitlabGroup blueprint and the _team.gitlab_group relation,
// and declaring them in the same configuration takes them over in one apply instead of failing.
//
// Port provisions defaults only for organizations with provisioned defaults enabled, so it runs only against those.
func TestAccPortIntegrationProvisionedDefaultsAreTakenOver(t *testing.T) {
	acctest.SkipUnlessAcceptance(t)
	enableIntegrationBetaFeatures(t)

	portClient := acctest.PortClient(t)
	enabled, err := portClient.HasFeatureFlags(context.Background(), cli.FeatureFlagProvisionedDefaults)
	if err != nil {
		t.Fatalf("failed to read organization feature flags: %v", err)
	}
	if !enabled {
		t.Skipf("organization does not have %s, so Port does not provision integration defaults", cli.FeatureFlagProvisionedDefaults)
	}
	installationID := utils.GenID()
	secretName := strings.ReplaceAll(utils.GenID(), "-", "_")
	t.Cleanup(func() { cleanupGitlabProvisionedDefaults(t, portClient, installationID, secretName) })
	config := gitlabProvisionedDefaultsHCL(
		acctest.ProviderConfigWithAdoption("port_blueprint", "port_blueprint_relation"), installationID, secretName)

	steps := []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("port_blueprint.gitlab_group", "identifier", "gitlabGroup"),
				resource.TestCheckResourceAttr("port_blueprint_relation.team_gitlab_group", "id", "_team:gitlab_group"),
				// relations Port provisioned on the blueprint and Terraform does not declare are kept
				blueprintHasRelation(portClient, "gitlabGroup", "gitlabMembers"),
			),
		},
		{
			Config:   config,
			PlanOnly: true,
		},
	}
	if acctest.IsLocal() {
		// deleting a Port hosted integration needs the Ocean SaaS operator, which local environments do not run,
		// so leave everything to the cleanup instead of destroying it
		steps = append(steps, resource.TestStep{Config: acctest.ProviderConfig + `
removed {
	from = port_blueprint_relation.team_gitlab_group
	lifecycle { destroy = false }
}
removed {
	from = port_blueprint.gitlab_group
	lifecycle { destroy = false }
}
removed {
	from = port_integration.gitlab
	lifecycle { destroy = false }
}
removed {
	from = port_organization_secret.gitlab_token
	lifecycle { destroy = false }
}
`})
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}
