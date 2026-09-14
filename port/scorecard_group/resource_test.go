package scorecard_group_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func testAccCreateBlueprintConfig(resourceName, identifier string) string {
	return fmt.Sprintf(`
	resource "port_blueprint" "%s" {
		title      = "TF test %s"
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
	`, resourceName, resourceName, identifier)
}

func testAccHasAuthorRuleHCL() string {
	return `[{
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
		}]`
}

func TestAccPortScorecardGroupSharedRules(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	groupIdentifier := utils.GenID()
	config := testAccCreateBlueprintConfig("microservice", blueprintIdentifier) + fmt.Sprintf(`
	resource "port_scorecard_group" "test" {
		identifier = "%s"
		title      = "Scorecard Group 1"
		blueprints = [port_blueprint.microservice.identifier]
		rules = %s
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
	}`, groupIdentifier, testAccHasAuthorRuleHCL())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheckScorecardGroups(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard_group.test", "title", "Scorecard Group 1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "identifier", groupIdentifier),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "blueprints.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "rules.0.identifier", "has-author"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "filters.%", "1"),
					resource.TestCheckNoResourceAttr("port_scorecard_group.test", "scorecards"),
				),
			},
		},
	})
}

func TestAccPortScorecardGroupSharedRulesMultipleBlueprints(t *testing.T) {
	svcIdentifier := utils.GenID()
	dbIdentifier := utils.GenID()
	groupIdentifier := utils.GenID()
	config := testAccCreateBlueprintConfig("microservice", svcIdentifier) +
		testAccCreateBlueprintConfig("database", dbIdentifier) +
		fmt.Sprintf(`
	resource "port_scorecard_group" "test" {
		identifier = "%s"
		title      = "Scorecard Group Shared Multi"
		blueprints = [
			port_blueprint.microservice.identifier,
			port_blueprint.database.identifier,
		]
		rules = %s
		filters = {
			(port_blueprint.microservice.identifier) = {
				combinator = "and"
				conditions = [jsonencode({
					property = "author"
					operator = "isNotEmpty"
				})]
			}
			(port_blueprint.database.identifier) = {
				combinator = "and"
				conditions = [jsonencode({
					property = "author"
					operator = "isNotEmpty"
				})]
			}
		}
		depends_on = [
			port_blueprint.microservice,
			port_blueprint.database,
		]
	}`, groupIdentifier, testAccHasAuthorRuleHCL())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheckScorecardGroups(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard_group.test", "title", "Scorecard Group Shared Multi"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "identifier", groupIdentifier),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "blueprints.#", "2"),
					resource.TestCheckTypeSetElemAttr("port_scorecard_group.test", "blueprints.*", svcIdentifier),
					resource.TestCheckTypeSetElemAttr("port_scorecard_group.test", "blueprints.*", dbIdentifier),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "rules.0.identifier", "has-author"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "filters.%", "2"),
					resource.TestCheckNoResourceAttr("port_scorecard_group.test", "scorecards"),
				),
			},
		},
	})
}

func TestAccPortScorecardGroupPerBlueprint(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	groupIdentifier := utils.GenID()
	config := testAccCreateBlueprintConfig("microservice", blueprintIdentifier) + fmt.Sprintf(`
	resource "port_scorecard_group" "test" {
		identifier = "%s"
		title      = "Scorecard Group 2"
		scorecards = {
			(port_blueprint.microservice.identifier) = {
				rules = %s
			}
		}
		depends_on = [port_blueprint.microservice]
	}`, groupIdentifier, testAccHasAuthorRuleHCL())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheckScorecardGroups(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard_group.test", "title", "Scorecard Group 2"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "identifier", groupIdentifier),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "scorecards.%", "1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "scorecards."+blueprintIdentifier+".rules.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "scorecards."+blueprintIdentifier+".rules.0.identifier", "has-author"),
					resource.TestCheckNoResourceAttr("port_scorecard_group.test", "blueprints"),
					resource.TestCheckNoResourceAttr("port_scorecard_group.test", "rules"),
				),
			},
		},
	})
}

func TestAccPortScorecardGroupPerBlueprintMultiple(t *testing.T) {
	svcIdentifier := utils.GenID()
	dbIdentifier := utils.GenID()
	groupIdentifier := utils.GenID()
	config := testAccCreateBlueprintConfig("microservice", svcIdentifier) +
		testAccCreateBlueprintConfig("database", dbIdentifier) +
		fmt.Sprintf(`
	resource "port_scorecard_group" "test" {
		identifier = "%s"
		title      = "Scorecard Group Per Blueprint Multi"
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
			(port_blueprint.database.identifier) = {
				filter = {
					combinator = "and"
					conditions = [jsonencode({
						property = "author"
						operator = "isNotEmpty"
					})]
				}
				rules = [{
					identifier = "db-has-author"
					title      = "Database Has Author"
					level      = "Silver"
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
		depends_on = [
			port_blueprint.microservice,
			port_blueprint.database,
		]
	}`, groupIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheckScorecardGroups(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard_group.test", "title", "Scorecard Group Per Blueprint Multi"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "identifier", groupIdentifier),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "scorecards.%", "2"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "scorecards."+svcIdentifier+".rules.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "scorecards."+svcIdentifier+".rules.0.identifier", "has-author"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "scorecards."+dbIdentifier+".rules.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "scorecards."+dbIdentifier+".rules.0.identifier", "db-has-author"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "scorecards."+dbIdentifier+".filter.combinator", "and"),
					resource.TestCheckNoResourceAttr("port_scorecard_group.test", "blueprints"),
					resource.TestCheckNoResourceAttr("port_scorecard_group.test", "rules"),
				),
			},
		},
	})
}
