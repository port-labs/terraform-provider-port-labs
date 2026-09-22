package scorecard_group_test

import (
	"context"
	"encoding/json"
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
	scorecardPropKey := "tf_scorecard_" + strings.ReplaceAll(utils.GenID(), "-", "")
	scorecardRelationKey := "tf_scorecard_rel_" + strings.ReplaceAll(utils.GenID(), "-", "")

	// `_scorecard_group` rejects schema changes (even PATCH) with protected_blueprint_violation.
	// Cover scorecard_* live against `_scorecard`; group_* mapping is covered by unit tests.
	config := testAccCreateBlueprintConfig("microservice", blueprintIdentifier) + fmt.Sprintf(`
	resource "port_team" "platform" {
		name = "%s"
	}

	resource "port_scorecard_group" "test" {
		identifier = "%s"
		title      = "Scorecard Group Props Relations"
		blueprints = [port_blueprint.microservice.identifier]
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
		scorecardPropKey,
		scorecardRelationKey,
		testAccHasAuthorRuleHCL(),
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheckScorecardGroups(t)
			testAccExtendScorecardBlueprint(t, scorecardPropKey, scorecardRelationKey)
		},
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard_group.test", "identifier", groupIdentifier),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "title", "Scorecard Group Props Relations"),
					resource.TestMatchResourceAttr("port_scorecard_group.test", "scorecard_properties", regexp.MustCompile(regexp.QuoteMeta(scorecardPropKey))),
					resource.TestMatchResourceAttr("port_scorecard_group.test", "scorecard_properties", regexp.MustCompile(`"platform-team"`)),
					resource.TestMatchResourceAttr("port_scorecard_group.test", "scorecard_relations", regexp.MustCompile(regexp.QuoteMeta(scorecardRelationKey))),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "blueprints.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard_group.test", "rules.#", "1"),
				),
			},
		},
	})
}

func testAccExtendScorecardBlueprint(t *testing.T, propKey, relationKey string) {
	t.Helper()

	client, ctx := testAccPortClient(t)
	teamTarget := "_team"
	title := "TF Acc Test"
	many := false
	required := false

	if _, err := client.PatchBlueprintRelation(ctx, "_scorecard", relationKey, &cli.Relation{
		Title:    &title,
		Target:   &teamTarget,
		Many:     &many,
		Required: &required,
	}); err != nil {
		t.Fatalf("failed to add relation %q on _scorecard: %v", relationKey, err)
	}

	if err := testAccPatchBlueprintProperty(ctx, client, "_scorecard", propKey, title); err != nil {
		t.Fatalf("failed to add property %q on _scorecard: %v", propKey, err)
	}

	t.Cleanup(func() {
		_ = testAccPatchBlueprintProperty(ctx, client, "_scorecard", propKey, "")
		_ = testAccPatchBlueprintRelationDelete(ctx, client, "_scorecard", relationKey)
	})
}

func testAccPatchBlueprintProperty(ctx context.Context, client *cli.PortClient, blueprintID, propKey, title string) error {
	var property any
	if title == "" {
		property = nil
	} else {
		property = map[string]any{
			"type":  "string",
			"title": title,
		}
	}

	body := map[string]any{
		"schema": map[string]any{
			"properties": map[string]any{
				propKey: property,
			},
		},
	}

	resp, err := client.Client.R().
		SetContext(ctx).
		SetBody(body).
		SetPathParam("identifier", blueprintID).
		Patch("v1/blueprints/{identifier}")
	if err != nil {
		return err
	}

	var pb cli.PortBody
	if err := json.Unmarshal(resp.Body(), &pb); err != nil {
		return err
	}
	if !pb.OK {
		return fmt.Errorf("failed to patch property %q on blueprint %q, got: %s", propKey, blueprintID, resp.Body())
	}
	return nil
}

func testAccPatchBlueprintRelationDelete(ctx context.Context, client *cli.PortClient, blueprintID, relationKey string) error {
	body := map[string]any{
		"relations": map[string]any{
			relationKey: nil,
		},
	}
	resp, err := client.Client.R().
		SetContext(ctx).
		SetBody(body).
		SetPathParam("identifier", blueprintID).
		Patch("v1/blueprints/{identifier}")
	if err != nil {
		return err
	}
	var pb cli.PortBody
	if err := json.Unmarshal(resp.Body(), &pb); err != nil {
		return err
	}
	if !pb.OK {
		return fmt.Errorf("failed to delete relation %q on blueprint %q, got: %s", relationKey, blueprintID, resp.Body())
	}
	return nil
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
