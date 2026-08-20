package scorecard_group_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func testAccCreateBlueprintConfig(identifier string) string {
	return fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title      = "TF test microservice"
		icon       = "Terraform"
		identifier = "%s"
		properties = {
			string_props = {
				"author" = {
					title = "Author"
				}
			}
		}
	}
	`, identifier)
}

func TestAccPortScorecardGroupSharedRules(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	groupIdentifier := utils.GenID()
	config := testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_scorecard_group" "test" {
		identifier = "%s"
		title      = "Scorecard Group 1"
		blueprints = [port_blueprint.microservice.identifier]
		rules = [{
			identifier = "has-author"
			title      = "Has Author"
			level      = "Gold"
			query = {
				combinator = "and"
				conditions = [jsonencode({
					property = "author"
					operator = "isNotEmpty"
				})]
			}
		}]
		filters = {
			(port_blueprint.microservice.identifier) = {
				combinator = "and"
				conditions = [jsonencode({
					property = "author"
					operator = "isNotEmpty"
				})]
			}
		}
		depends_on = [port_blueprint.microservice]
	}`, groupIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard_group.test", "title", "Scorecard Group 1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "identifier", groupIdentifier),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "rules.0.identifier", "has-author"),
					resource.TestCheckResourceAttrSet("port_scorecard_group.test", "blueprints.#"),
				),
			},
		},
	})
}

func TestAccPortScorecardGroupPerBlueprint(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	groupIdentifier := utils.GenID()
	config := testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_scorecard_group" "test" {
		identifier = "%s"
		title      = "Scorecard Group 2"
		scorecards = {
			(port_blueprint.microservice.identifier) = {
				rules = [{
					identifier = "has-author"
					title      = "Has Author"
					level      = "Gold"
					query = {
						combinator = "and"
						conditions = [jsonencode({
							property = "author"
							operator = "isNotEmpty"
						})]
					}
				}]
			}
		}
		depends_on = [port_blueprint.microservice]
	}`, groupIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard_group.test", "title", "Scorecard Group 2"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "identifier", groupIdentifier),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "scorecards.%", "1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "scorecards."+blueprintIdentifier+".rules.#", "1"),
				),
			},
		},
	})
}
