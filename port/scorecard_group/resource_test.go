package scorecard_group_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func enableScorecardGroupBeta(t *testing.T) {
	t.Helper()
	if err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true"); err != nil {
		t.Fatal(err)
	}
}

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
	enableScorecardGroupBeta(t)
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
	enableScorecardGroupBeta(t)
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
	enableScorecardGroupBeta(t)
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
	enableScorecardGroupBeta(t)
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

func TestAccPortScorecardGroupBetaDisabled(t *testing.T) {
	if err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "false"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	})

	config := `
	resource "port_scorecard_group" "test" {
		identifier = "beta-disabled"
		title      = "Beta Disabled"
		blueprints = ["service"]
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
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      acctest.ProviderConfig + config,
				ExpectError: regexp.MustCompile("Beta features are not enabled"),
			},
		},
	})
}

func TestAccPortScorecardGroupPropertiesAndRelations(t *testing.T) {
	enableScorecardGroupBeta(t)

	blueprintIdentifier := utils.GenID()
	teamBlueprintIdentifier := utils.GenID()
	teamEntityIdentifier := utils.GenID()
	groupIdentifier := utils.GenID()
	groupPropKey := "tf_group_" + utils.GenID()
	scorecardPropKey := "tf_scorecard_" + utils.GenID()
	groupRelationKey := "tf_group_rel_" + utils.GenID()
	scorecardRelationKey := "tf_scorecard_rel_" + utils.GenID()

	config := testAccCreateBlueprintConfig("microservice", blueprintIdentifier) + fmt.Sprintf(`
	resource "port_blueprint" "team" {
		title      = "TF test team"
		icon       = "Team"
		identifier = "%s"
	}

	resource "port_entity" "platform_team" {
		identifier = "%s"
		title      = "Platform Team"
		blueprint  = port_blueprint.team.identifier
	}

	resource "port_system_blueprint" "scorecard_group" {
		identifier = "_scorecard_group"
		properties = {
			string_props = {
				"%s" = {
					title = "TF Group Category"
				}
			}
		}
		relations = {
			"%s" = {
				title  = "TF Owning Team"
				target = port_blueprint.team.identifier
			}
		}
	}

	resource "port_system_blueprint" "scorecard" {
		identifier = "_scorecard"
		properties = {
			string_props = {
				"%s" = {
					title = "TF Scorecard Owner"
				}
			}
		}
		relations = {
			"%s" = {
				title  = "TF Owner Team"
				target = port_blueprint.team.identifier
			}
		}
	}

	resource "port_scorecard_group" "test" {
		identifier = "%s"
		title      = "Scorecard Group Props Relations"
		blueprints = [port_blueprint.microservice.identifier]
		group_properties = jsonencode({
			%s = "governance"
		})
		group_relations = jsonencode({
			%s = port_entity.platform_team.identifier
		})
		scorecard_properties = jsonencode({
			%s = "platform-team"
		})
		scorecard_relations = jsonencode({
			%s = port_entity.platform_team.identifier
		})
		rules = %s
		depends_on = [
			port_blueprint.microservice,
			port_system_blueprint.scorecard_group,
			port_system_blueprint.scorecard,
			port_entity.platform_team,
		]
	}`,
		teamBlueprintIdentifier,
		teamEntityIdentifier,
		groupPropKey,
		groupRelationKey,
		scorecardPropKey,
		scorecardRelationKey,
		groupIdentifier,
		groupPropKey,
		groupRelationKey,
		scorecardPropKey,
		scorecardRelationKey,
		testAccHasAuthorRuleHCL(),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheckScorecardGroups(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard_group.test", "identifier", groupIdentifier),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "title", "Scorecard Group Props Relations"),
					resource.TestMatchResourceAttr("port_scorecard_group.test", "group_properties", regexp.MustCompile(regexp.QuoteMeta(groupPropKey))),
					resource.TestMatchResourceAttr("port_scorecard_group.test", "group_properties", regexp.MustCompile(`"governance"`)),
					resource.TestMatchResourceAttr("port_scorecard_group.test", "scorecard_properties", regexp.MustCompile(regexp.QuoteMeta(scorecardPropKey))),
					resource.TestMatchResourceAttr("port_scorecard_group.test", "scorecard_properties", regexp.MustCompile(`"platform-team"`)),
					resource.TestMatchResourceAttr("port_scorecard_group.test", "group_relations", regexp.MustCompile(regexp.QuoteMeta(groupRelationKey))),
					resource.TestMatchResourceAttr("port_scorecard_group.test", "group_relations", regexp.MustCompile(regexp.QuoteMeta(teamEntityIdentifier))),
					resource.TestMatchResourceAttr("port_scorecard_group.test", "scorecard_relations", regexp.MustCompile(regexp.QuoteMeta(scorecardRelationKey))),
					resource.TestMatchResourceAttr("port_scorecard_group.test", "scorecard_relations", regexp.MustCompile(regexp.QuoteMeta(teamEntityIdentifier))),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "blueprints.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "rules.#", "1"),
				),
			},
		},
	})
}
