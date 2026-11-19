package system_blueprint_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"regexp"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/port-labs/terraform-provider-port-labs/v2/version"
)

func TestAccPortSystemBlueprintBasic(t *testing.T) {
	identifier := "_user"

	var basicConfig = fmt.Sprintf(`
	resource "port_system_blueprint" "test" {
		identifier = "%s"
	}`, identifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + basicConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_system_blueprint.test", "identifier", identifier),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "id", identifier),
				),
			},
			{
				ResourceName:      "port_system_blueprint.test",
				ImportState:       true,
				ImportStateId:     identifier,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"properties",
					"relations",
					"mirror_properties",
					"calculation_properties",
				},
			},
		},
	})
}

func TestAccPortSystemBlueprintProperties(t *testing.T) {
	identifier := "_user"

	var configWithProperties = fmt.Sprintf(`
	resource "port_system_blueprint" "test" {
		identifier = "%s"
		properties = {
			string_props = {
				"environment" = {
					title = "Environment"
					description = "The environment this service runs in"
					enum = ["dev", "staging", "prod"]
					enum_colors = {
						"dev" = "blue"
						"staging" = "yellow"
						"prod" = "green"
					}
				}
			}
			number_props = {
				"version" = {
					title = "Version"
					description = "The version number"
					minimum = 1
					maximum = 10
				}
			}
		}
	}`, identifier)

	var configWithoutProperties = fmt.Sprintf(`
	resource "port_system_blueprint" "test" {
		identifier = "%s"
		properties = {
		}
	}`, identifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + configWithProperties,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_system_blueprint.test", "identifier", identifier),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "properties.string_props.environment.title", "Environment"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "properties.string_props.environment.description", "The environment this service runs in"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "properties.number_props.version.title", "Version"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "properties.number_props.version.minimum", "1"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "properties.number_props.version.maximum", "10"),
				),
			},
			{
				Config: acctest.ProviderConfig + configWithoutProperties,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_system_blueprint.test", "identifier", identifier),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "properties.number_props.version.title"),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "properties.number_props.version.minimum"),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "properties.number_props.version.maximum"),
				),
			},
		},
	})
}

func TestAccPortBlueprintChangePropertyType(t *testing.T) {
	type data struct{ PropType string }
	cleanConfig := `
		resource "port_system_blueprint" "user" { identifier = "_user" }
	`
	tmpl, err := template.New("resource").Parse(`
	resource "port_system_blueprint" "user" {
		identifier = "_user"
		properties = {
			{{.PropType}}_props = {
				myProperty = {
					title = "My Property"
					description = "This is a {{.PropType}} property"
				}
			}
		}
	}`)
	require.NoErrorf(t, err, "failed to parse test template")

	var propTypes = [...]string{"string", "number", "boolean", "array", "object"}

	// Shuffle the prop types to make sure we don't have an issue transitioning from one type to the next.
	rand.Shuffle(len(propTypes), func(i, j int) { propTypes[i], propTypes[j] = propTypes[j], propTypes[i] })

	steps := make([]resource.TestStep, 1, 1+len(propTypes))
	steps[0] = resource.TestStep{
		ResourceName:       "port_system_blueprint.user",
		ImportState:        true,
		ImportStateId:      "_user",
		ImportStatePersist: true,
		Config:             acctest.ProviderConfigNoPropertyTypeProtection + cleanConfig,
	}
	for _, propType := range propTypes {
		var txt strings.Builder
		err = tmpl.Execute(&txt, data{PropType: propType})
		require.NoErrorf(t, err, "failed to execute template for propType: %s", propType)
		steps = append(steps, resource.TestStep{
			Config: acctest.ProviderConfigNoPropertyTypeProtection + txt.String(),
			Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(
				"port_system_blueprint.user",
				fmt.Sprintf("properties.%s_props.myProperty.title", propType),
				"My Property",
			)),
		})
	}
	steps = append(steps, resource.TestStep{Config: acctest.ProviderConfigNoPropertyTypeProtection + cleanConfig})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

func TestAccPortBlueprintChangePropertyTypeProtection(t *testing.T) {
	type data struct{ PropType string }
	cleanConfig := `
		resource "port_system_blueprint" "user" { identifier = "_user" }
	`
	tmpl, err := template.New("resource").Parse(`
	resource "port_system_blueprint" "user" {
		identifier = "_user"
		properties = {
			{{.PropType}}_props = {
				myProperty = {
					title = "My Property"
					description = "This is a {{.PropType}} property"
				}
			}
		}
	}`)
	require.NoErrorf(t, err, "failed to parse test template")

	var step1Text, step2Text strings.Builder
	require.NoError(t, tmpl.Execute(&step1Text, data{PropType: "string"}))
	require.NoError(t, tmpl.Execute(&step2Text, data{PropType: "number"}))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				ResourceName:       "port_system_blueprint.user",
				ImportState:        true,
				ImportStateId:      "_user",
				ImportStatePersist: true,
				Config:             acctest.ProviderConfig + cleanConfig,
			},
			{
				Config: acctest.ProviderConfig + step1Text.String(),
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(
					"port_system_blueprint.user",
					"properties.string_props.myProperty.title",
					"My Property",
				)),
			},
			{
				Config:      acctest.ProviderConfig + step2Text.String(),
				ExpectError: regexp.MustCompile(`The type of property "myProperty" changed from "string" to "number"`),
			},
			{
				Config: acctest.ProviderConfig + cleanConfig,
			},
		},
	})
}

func TestAccPortSystemBlueprintRelations(t *testing.T) {
	identifier := "_user"

	var configWithRelations = fmt.Sprintf(`
	resource "port_system_blueprint" "test" {
		identifier = "%s"
		relations = {
			"groups" = {
				target = "_team"
				title = "Teams"
				description = "The teams that owns this service"
				many = true
				required = false
			}
			"owner" = {
				target = "_team"
				title = "Owner"
				description = "The team that owns this service"
				many = false
				required = true
			}
		}
	}`, identifier)

	var configWithoutRelations = fmt.Sprintf(`
	resource "port_system_blueprint" "test" {
		identifier = "%s"
		relations = {
		}
	}`, identifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + configWithRelations,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_system_blueprint.test", "identifier", identifier),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.target", "_team"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.title", "Teams"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.description", "The teams that owns this service"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.many", "true"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.required", "false"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.owner.target", "_team"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.owner.title", "Owner"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.owner.description", "The team that owns this service"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.owner.many", "false"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.owner.required", "true"),
					testAccCheckSystemBlueprintRelationInAPI(identifier, "owner", "_team"),
					testAccCheckSystemBlueprintRelationInAPI(identifier, "groups", "_team"),
				),
			},
			{
				Config: acctest.ProviderConfig + configWithoutRelations,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_system_blueprint.test", "identifier", identifier),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "relations.owner.target"),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "relations.owner.title"),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "relations.owner.description"),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "relations.owner.many"),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "relations.owner.required"),
				),
			},
		},
	})
}

func TestAccPortSystemBlueprintResourceAllowsUpdate(t *testing.T) {
	identifier := "_user"
	relationKey := fmt.Sprintf("tf_rel_%s", strings.ReplaceAll(utils.GenID(), "-", ""))
	propKey := fmt.Sprintf("tf_prop_%s", strings.ReplaceAll(utils.GenID(), "-", ""))

	configWithSchema := fmt.Sprintf(`
	resource "port_system_blueprint" "test" {
		identifier = "%s"
		properties = {
			string_props = {
				"%s" = {
					title = "TF Create Schema Prop"
				}
			}
		}
		relations = {
			"%s" = {
				target = "_team"
				title = "TF Create Schema Relation"
				many = false
				required = false
			}
		}
	}`, identifier, propKey, relationKey)

	configWithoutSchema := fmt.Sprintf(`
	resource "port_system_blueprint" "test" {
		identifier = "%s"
		properties = {}
		relations = {}
	}`, identifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + configWithSchema,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_system_blueprint.test", "identifier", identifier),
					resource.TestCheckResourceAttr(
						"port_system_blueprint.test",
						fmt.Sprintf("properties.string_props.%s.title", propKey),
						"TF Create Schema Prop",
					),
					resource.TestCheckResourceAttr(
						"port_system_blueprint.test",
						fmt.Sprintf("relations.%s.target", relationKey),
						"_team",
					),
					testAccCheckSystemBlueprintPropertyInAPI(identifier, propKey),
					testAccCheckSystemBlueprintRelationInAPI(identifier, relationKey, "_team"),
				),
			},
			{
				Config: acctest.ProviderConfig + configWithoutSchema,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_system_blueprint.test", "identifier", identifier),
					testAccCheckSystemBlueprintPropertyAbsentInAPI(identifier, propKey),
					testAccCheckSystemBlueprintRelationAbsentInAPI(identifier, relationKey),
				),
			},
		},
	})
}

func testAccPortClient() (*cli.PortClient, context.Context, error) {
	baseURL := os.Getenv("PORT_BASE_URL")
	if baseURL == "" {
		baseURL = consts.DefaultBaseUrl
	}
	client, err := cli.New(baseURL, cli.WithHeader("User-Agent", version.ProviderVersion))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Port client: %w", err)
	}
	ctx := context.Background()
	if _, err := client.Authenticate(ctx, os.Getenv("PORT_CLIENT_ID"), os.Getenv("PORT_CLIENT_SECRET")); err != nil {
		return nil, nil, fmt.Errorf("failed to authenticate with Port: %w", err)
	}
	return client, ctx, nil
}

func testAccReadSystemBlueprint(blueprintID string) (*cli.Blueprint, error) {
	client, ctx, err := testAccPortClient()
	if err != nil {
		return nil, err
	}
	bp, _, err := client.ReadBlueprint(ctx, blueprintID)
	if err != nil {
		return nil, fmt.Errorf("failed to read blueprint %q from API: %w", blueprintID, err)
	}
	return bp, nil
}

func testAccCheckSystemBlueprintRelationInAPI(blueprintID, relationKey, target string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		bp, err := testAccReadSystemBlueprint(blueprintID)
		if err != nil {
			return err
		}
		relation, ok := bp.Relations[relationKey]
		if !ok {
			return fmt.Errorf("relation %q was not found on blueprint %q in Port API", relationKey, blueprintID)
		}
		if relation.Target == nil || *relation.Target != target {
			return fmt.Errorf("relation %q target = %v, want %q", relationKey, relation.Target, target)
		}
		return nil
	}
}

func testAccCheckSystemBlueprintRelationAbsentInAPI(blueprintID, relationKey string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		bp, err := testAccReadSystemBlueprint(blueprintID)
		if err != nil {
			return err
		}
		if _, ok := bp.Relations[relationKey]; ok {
			return fmt.Errorf("relation %q still exists on blueprint %q in Port API", relationKey, blueprintID)
		}
		return nil
	}
}

func testAccCheckSystemBlueprintPropertyInAPI(blueprintID, propKey string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		bp, err := testAccReadSystemBlueprint(blueprintID)
		if err != nil {
			return err
		}
		if _, ok := bp.Schema.Properties[propKey]; !ok {
			return fmt.Errorf("property %q was not found on blueprint %q in Port API", propKey, blueprintID)
		}
		return nil
	}
}

func testAccCheckSystemBlueprintPropertyAbsentInAPI(blueprintID, propKey string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		bp, err := testAccReadSystemBlueprint(blueprintID)
		if err != nil {
			return err
		}
		if _, ok := bp.Schema.Properties[propKey]; ok {
			return fmt.Errorf("property %q still exists on blueprint %q in Port API", propKey, blueprintID)
		}
		return nil
	}
}

func TestAccPortSystemBlueprintIncludeInGlobalSearch(t *testing.T) {
	identifier := "_user"

	var configTrue = fmt.Sprintf(`
	resource "port_system_blueprint" "test" {
		identifier = "%s"
		include_in_global_search = true
	}`, identifier)

	var configFalse = fmt.Sprintf(`
	resource "port_system_blueprint" "test" {
		identifier = "%s"
		include_in_global_search = false
	}`, identifier)

	var configUnset = fmt.Sprintf(`
	resource "port_system_blueprint" "test" {
		identifier = "%s"
	}`, identifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + configTrue,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_system_blueprint.test", "identifier", identifier),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "include_in_global_search", "true"),
				),
			},
			{
				Config: acctest.ProviderConfig + configFalse,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_system_blueprint.test", "include_in_global_search", "false"),
				),
			},
			{
				Config: acctest.ProviderConfig + configUnset,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "include_in_global_search"),
				),
			},
		},
	})
}

func TestAccPortSystemBlueprintMirrorProperties(t *testing.T) {
	identifier := "_user"

	var configWithMirrorProps = fmt.Sprintf(`
	resource "port_system_blueprint" "test" {
		identifier = "%s"
		relations = {
			"groups" = {
				target = "_team"
				title = "Teams"
				description = "The teams that owns this service"
				many = true
				required = false
			}
		}
		mirror_properties = {
			"team_size" = {
				path = "groups.size"
				title = "Team Size"
			}
			"team_name" = {
				path = "groups.name"
				title = "Team Name"
			}
		}
	}`, identifier)

	var configWithoutMirrorProps = fmt.Sprintf(`
	resource "port_system_blueprint" "test" {
		identifier = "%s"
		relations = {
		}
		mirror_properties = {
		}
	}`, identifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + configWithMirrorProps,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_system_blueprint.test", "identifier", identifier),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.target", "_team"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.title", "Teams"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.description", "The teams that owns this service"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.many", "true"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.required", "false"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "mirror_properties.team_size.path", "groups.size"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "mirror_properties.team_size.title", "Team Size"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "mirror_properties.team_name.path", "groups.name"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "mirror_properties.team_name.title", "Team Name"),
				),
			},
			{
				Config: acctest.ProviderConfig + configWithoutMirrorProps,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_system_blueprint.test", "identifier", identifier),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "mirror_properties.team_name.path"),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "mirror_properties.team_name.title"),
				),
			},
		},
	})
}
