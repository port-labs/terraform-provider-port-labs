package system_blueprint_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
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
					resource.TestCheckResourceAttr("port_system_blueprint.test", "properties.string_props.environment.title", "Environment"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "properties.string_props.environment.description", "The environment this service runs in"),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "properties.number_props.version.title"),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "properties.number_props.version.minimum"),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "properties.number_props.version.maximum"),
				),
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
			"groups" = {
				target = "_team"
				title = "Teams"
				description = "The teams that owns this service"
				many = true
				required = false
			}
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
				),
			},
			{
				Config: acctest.ProviderConfig + configWithoutRelations,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_system_blueprint.test", "identifier", identifier),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.target", "_team"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.title", "Teams"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.description", "The teams that owns this service"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.many", "true"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.required", "false"),
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
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.target", "_team"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.title", "Teams"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.description", "The teams that owns this service"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.many", "true"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "relations.groups.required", "false"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "mirror_properties.team_size.path", "groups.size"),
					resource.TestCheckResourceAttr("port_system_blueprint.test", "mirror_properties.team_size.title", "Team Size"),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "mirror_properties.team_name.path"),
					resource.TestCheckNoResourceAttr("port_system_blueprint.test", "mirror_properties.team_name.title"),
				),
			},
		},
	})
} 