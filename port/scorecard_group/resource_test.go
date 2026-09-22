package scorecard_group_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/port-labs/terraform-provider-port-labs/v2/version"
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
	groupIdentifier := utils.GenID()
	teamName := "tf-sc-group-" + strings.ReplaceAll(utils.GenID(), "-", "")
	groupPropKey := "tf_group_" + strings.ReplaceAll(utils.GenID(), "-", "")
	scorecardPropKey := "tf_scorecard_" + strings.ReplaceAll(utils.GenID(), "-", "")
	groupRelationKey := "tf_group_rel_" + strings.ReplaceAll(utils.GenID(), "-", "")
	scorecardRelationKey := "tf_scorecard_rel_" + strings.ReplaceAll(utils.GenID(), "-", "")

	// Extend system blueprints via the API from the live schema. Using port_system_blueprint
	// merges from the structure schema and can violate protected scorecard-group properties.
	config := testAccCreateBlueprintConfig("microservice", blueprintIdentifier) + fmt.Sprintf(`
	resource "port_team" "platform" {
		name = "%s"
	}

	resource "port_scorecard_group" "test" {
		identifier = "%s"
		title      = "Scorecard Group Props Relations"
		blueprints = [port_blueprint.microservice.identifier]
		group_properties = jsonencode({
			%s = "governance"
		})
		group_relations = jsonencode({
			%s = port_team.platform.identifier
		})
		scorecard_properties = jsonencode({
			%s = "platform-team"
		})
		scorecard_relations = jsonencode({
			%s = port_team.platform.identifier
		})
		rules = %s
		depends_on = [
			port_blueprint.microservice,
			port_team.platform,
		]
	}`,
		teamName,
		groupIdentifier,
		groupPropKey,
		groupRelationKey,
		scorecardPropKey,
		scorecardRelationKey,
		testAccHasAuthorRuleHCL(),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheckScorecardGroups(t)
			testAccExtendScorecardSystemBlueprints(t, groupPropKey, scorecardPropKey, groupRelationKey, scorecardRelationKey)
		},
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
					resource.TestMatchResourceAttr("port_scorecard_group.test", "scorecard_relations", regexp.MustCompile(regexp.QuoteMeta(scorecardRelationKey))),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "blueprints.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "rules.#", "1"),
				),
			},
		},
	})
}

func testAccExtendScorecardSystemBlueprints(t *testing.T, groupPropKey, scorecardPropKey, groupRelationKey, scorecardRelationKey string) {
	t.Helper()

	client, ctx := testAccPortClient(t)
	teamTarget := "_team"
	title := "TF Acc Test"
	many := false
	required := false

	addExtensions := func(blueprintID, propKey, relationKey string) {
		t.Helper()
		bp, statusCode, err := client.ReadBlueprint(ctx, blueprintID)
		if err != nil {
			t.Fatalf("failed to read %s blueprint: %v (status %d)", blueprintID, err, statusCode)
		}
		if bp.Schema.Properties == nil {
			bp.Schema.Properties = map[string]cli.BlueprintProperty{}
		}
		bp.Schema.Properties[propKey] = cli.BlueprintProperty{
			Type:  "string",
			Title: &title,
		}
		if bp.Relations == nil {
			bp.Relations = map[string]cli.Relation{}
		}
		bp.Relations[relationKey] = cli.Relation{
			Title:    &title,
			Target:   &teamTarget,
			Many:     &many,
			Required: &required,
		}
		if _, err := client.UpdateBlueprint(ctx, bp, blueprintID); err != nil {
			t.Fatalf("failed to extend %s blueprint: %v", blueprintID, err)
		}
	}

	addExtensions("_scorecard_group", groupPropKey, groupRelationKey)
	addExtensions("_scorecard", scorecardPropKey, scorecardRelationKey)

	t.Cleanup(func() {
		removeExtensions := func(blueprintID, propKey, relationKey string) {
			bp, statusCode, err := client.ReadBlueprint(ctx, blueprintID)
			if err != nil {
				t.Logf("cleanup: failed to read %s blueprint: %v (status %d)", blueprintID, err, statusCode)
				return
			}
			delete(bp.Schema.Properties, propKey)
			delete(bp.Relations, relationKey)
			if _, err := client.UpdateBlueprint(ctx, bp, blueprintID); err != nil {
				t.Logf("cleanup: failed to remove extensions from %s: %v", blueprintID, err)
			}
		}
		removeExtensions("_scorecard_group", groupPropKey, groupRelationKey)
		removeExtensions("_scorecard", scorecardPropKey, scorecardRelationKey)
	})
}

func testAccPortClient(t *testing.T) (*cli.PortClient, context.Context) {
	t.Helper()

	baseURL := os.Getenv("PORT_BASE_URL")
	if baseURL == "" {
		baseURL = consts.DefaultBaseUrl
	}
	client, err := cli.New(baseURL, cli.WithHeader("User-Agent", version.ProviderVersion))
	if err != nil {
		t.Fatalf("failed to create Port client: %v", err)
	}
	ctx := context.Background()
	if _, err := client.Authenticate(ctx, os.Getenv("PORT_CLIENT_ID"), os.Getenv("PORT_CLIENT_SECRET")); err != nil {
		t.Fatalf("failed to authenticate with Port: %v", err)
	}
	return client, ctx
}
