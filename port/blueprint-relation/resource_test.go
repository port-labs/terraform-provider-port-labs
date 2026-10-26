package blueprint_relation_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

// createBlueprintWithRelation creates two blueprints outside Terraform and a relation between them, as an
// integration's default resources would (e.g. _team.gitlab_group to gitlabGroup), and deletes them when the
// test ends. It returns the source and target blueprint identifiers.
func createBlueprintWithRelation(t *testing.T, relation string) (*cli.PortClient, string, string) {
	t.Helper()
	acctest.SkipUnlessAcceptance(t)
	ctx := context.Background()
	portClient := acctest.PortClient(t)
	source, target := utils.GenID(), utils.GenID()

	// cleanups run in reverse, so the source and its relation go before the target it points at
	for _, identifier := range []string{target, source} {
		icon := "Terraform"
		_, err := portClient.CreateBlueprint(ctx, &cli.Blueprint{
			Identifier:            identifier,
			Icon:                  &icon,
			Title:                 "Created outside Terraform",
			Schema:                cli.BlueprintSchema{Properties: map[string]cli.BlueprintProperty{}},
			CalculationProperties: map[string]cli.BlueprintCalculationProperty{},
			AggregationProperties: map[string]cli.BlueprintAggregationProperty{},
			MirrorProperties:      map[string]cli.BlueprintMirrorProperty{},
			Relations:             map[string]cli.Relation{},
		}, nil)
		if err != nil {
			t.Fatalf("failed to create blueprint: %s", err)
		}
		t.Cleanup(func() { _ = portClient.DeleteBlueprint(context.Background(), identifier) })
	}

	title := "Added outside Terraform"
	if _, err := portClient.PatchBlueprintRelation(ctx, source, relation, &cli.Relation{Target: &target, Title: &title}); err != nil {
		t.Fatalf("failed to create relation: %s", err)
	}
	return portClient, source, target
}

func relationConfig(blueprint, relation, target string, extra string) string {
	return fmt.Sprintf(`
	resource "port_blueprint_relation" "parent" {
		blueprint  = "%s"
		identifier = "%s"
		target     = "%s"
		title      = "Declared in Terraform"
		%s
	}
`, blueprint, relation, target, extra)
}

func relationTitleIs(portClient *cli.PortClient, blueprint, relation, want string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		b, _, err := portClient.ReadBlueprint(context.Background(), blueprint)
		if err != nil {
			return err
		}
		got := b.Relations[relation].Title
		if got == nil || *got != want {
			return fmt.Errorf("relation %q title: got %v, want %q", relation, got, want)
		}
		return nil
	}
}

func TestAccPortBlueprintRelationAdoptsExistingRelation(t *testing.T) {
	portClient, blueprint, target := createBlueprintWithRelation(t, "parent")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      acctest.ProviderConfig + relationConfig(blueprint, "parent", target, ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("relation already exists"),
			},
			{
				Config: acctest.ProviderConfig + relationConfig(blueprint, "parent", target, "adopt_existing = true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint_relation.parent", "id", blueprint+":parent"),
					relationTitleIs(portClient, blueprint, "parent", "Declared in Terraform"),
				),
			},
		},
	})
}

func TestAccPortBlueprintRelationAdoptsExistingRelationByProviderPolicy(t *testing.T) {
	portClient, blueprint, target := createBlueprintWithRelation(t, "parent")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// the resource's own opt-out wins over the provider policy
				Config:      acctest.ProviderConfigWithAdoption("port_blueprint_relation") + relationConfig(blueprint, "parent", target, "adopt_existing = false"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("relation already exists"),
			},
			{
				Config: acctest.ProviderConfigWithAdoption("port_blueprint_relation") + relationConfig(blueprint, "parent", target, ""),
				Check:  relationTitleIs(portClient, blueprint, "parent", "Declared in Terraform"),
			},
		},
	})
}
