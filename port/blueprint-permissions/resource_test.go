package blueprint_permissions_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func createBlueprint(identifier string) string {
	return fmt.Sprintf(`

resource "port_blueprint" "microservice" {
  identifier            = "%s"
  title                 = "TF Provider Test"
  icon                  = "Terraform"
  description			= ""
}
`, identifier)
}

func createBlueprintWithProperties(identifier string) string {
	return fmt.Sprintf(`

resource "port_blueprint" "microservice" {
  identifier            = "%s"
  title                 = "TF Provider Test"
  icon                  = "Terraform"
  description			= ""
  properties = {
  	string_props = {
  		myStringIdentifier = {
  			description = "This is a string property"
  			title = "text"
  			icon = "Terraform"
  			required = true
  			min_length = 1
  			max_length = 10
  			default = "default"
  			enum = ["default", "default2"]
  			pattern = "^[a-zA-Z0-9]*$"
  			format = "user"
  			enum_colors = {
  				default = "red"
  				default2 = "green"
  			}
  		}
  	}
  }
}
`, identifier)
}

func TestAccPortBlueprintPermissionsBasic(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortBlueprintResourceBasic = createBlueprint(blueprintIdentifier)

	var testAccBaseBlueprintPermissionsConfigUpdate = `

	resource "port_blueprint_permissions" "microservice_permissions" {
		blueprint_identifier = port_blueprint.microservice.identifier
		entities = {
			"register" = {
				"roles": [
					"Member",
				],
				"users": [],
				"teams": []
			},
			"unregister" = {
				"roles": [
					"Member",
				],
				"users": [],
				"teams": []
			},
			"update" = {
				"roles": [
					"Member",
				],
				"users": [],
				"teams": []
			},
			"update_metadata_properties" = {
				"icon" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
				"identifier" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
				"team" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
				"title" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
			}
		}
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPortBlueprintResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
				),
			},
			{
				Config: testAccPortBlueprintResourceBasic + testAccBaseBlueprintPermissionsConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "id", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.register.roles.#", "1"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.register.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.register.users.#", "0"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.register.teams.#", "0"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.update_metadata_properties.icon.roles.#", "1"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.update_metadata_properties.icon.roles.0", "Member"),
				),
			},
		},
	})
}

func TestAccPortBlueprintPermissionsWithProperties(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortBlueprintResourceBasic = createBlueprintWithProperties(blueprintIdentifier)

	var testAccBaseBlueprintPermissionsConfigUpdate = `

	resource "port_blueprint_permissions" "microservice_permissions" {
		blueprint_identifier = port_blueprint.microservice.identifier
		entities = {
			"register" = {
				"roles": [
					"Member",
				],
				"users": [],
				"teams": []
			},
			"unregister" = {
				"roles": [
					"Member",
				],
				"users": [],
				"teams": []
			},
			"update" = {
				"roles": [
					"Member",
				],
				"users": [],
				"teams": []
			},
			"update_properties" = {
				"myStringIdentifier" = {
					"teams" = [],
					"users" = [],
					"roles" = ["Member"]
				}
			},
			"update_metadata_properties" = {
				"icon" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
				"identifier" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
				"team" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
				"title" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
			}
		}
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPortBlueprintResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
				),
			},
			{
				Config: testAccPortBlueprintResourceBasic + testAccBaseBlueprintPermissionsConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "id", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.register.roles.#", "1"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.register.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.register.users.#", "0"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.register.teams.#", "0"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.update_properties.myStringIdentifier.roles.#", "1"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.update_properties.myStringIdentifier.roles.0", "Member"),
				)},
		},
	})
}

func TestAccPortBlueprintPermissionsWithRelations(t *testing.T) {
	blueprintMicroserviceIdentifier := utils.GenID()
	blueprintEnvIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortBlueprintResourceBasic = fmt.Sprintf(`
resource "port_blueprint" "environment" {
  title      = "Environment"
  icon       = "Environment"
  identifier = "%s"
  properties = {
    string_props = {
      "name" = {
        type  = "string"
        title = "name"
      }
      "docs-url" = {
        title  = "Docs URL"
        format = "url"
      }
    }
  }
}


resource "port_blueprint" "microservice" {
  identifier            = "%s"
  title                 = "TF Provider Test"
  icon                  = "Terraform"
  description			= ""
  relations = {
    "environment" = {
      title    = "Test Relation"
      required = "true"
      target   = port_blueprint.environment.identifier
    }
  }
}
`, blueprintEnvIdentifier, blueprintMicroserviceIdentifier)

	teamName := utils.GenID()
	var testAccBaseBlueprintPermissionsConfigUpdate = fmt.Sprintf(`

	resource "port_team" "team" {
		name = "%s"
		description = "Test description"
		users = []
	}

	resource "port_blueprint_permissions" "microservice_permissions" {
		blueprint_identifier = port_blueprint.microservice.identifier
		entities = {
			"register" = {
				"roles": [
					"Member",
				],
				"users": [],
				"teams": [port_team.team.name]
			},
			"unregister" = {
				"roles": [
					"Member",
				],
				"users": [],
				"teams": []
			},
			"update" = {
				"roles": [
					"Member",
				],
				"users": [],
				"teams": []
			},
			"update_metadata_properties" = {
				"icon" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
				"identifier" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
				"team" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
				"title" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
			},
			"update_relations" = {
				"environment" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				}
			}
		}
	}`, teamName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPortBlueprintResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.environment", "identifier", blueprintEnvIdentifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", blueprintMicroserviceIdentifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
				),
			},
			{
				Config: testAccPortBlueprintResourceBasic + testAccBaseBlueprintPermissionsConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", blueprintMicroserviceIdentifier),
					resource.TestCheckResourceAttr("port_blueprint.environment", "identifier", blueprintEnvIdentifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "blueprint_identifier", blueprintMicroserviceIdentifier),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "id", blueprintMicroserviceIdentifier),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.register.roles.#", "1"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.register.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.register.teams.#", "1"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.register.teams.0", teamName),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.register.users.#", "0"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.update_relations.environment.roles.#", "1"),
					resource.TestCheckResourceAttr("port_blueprint_permissions.microservice_permissions", "entities.update_relations.environment.roles.0", "Member"),
				)},
		},
	})
}

func TestAccPortBlueprintPermissionsWithInvalidProperties(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortBlueprintResourceBasic = createBlueprintWithProperties(blueprintIdentifier)

	teamName := utils.GenID()
	var testAccBaseBlueprintPermissionsConfigUpdate = fmt.Sprintf(`

	resource "port_team" "team" {
		name = "%s"
		description = "Test description"
		users = []
	}

	resource "port_blueprint_permissions" "microservice_permissions" {
		blueprint_identifier = port_blueprint.microservice.identifier
		entities = {
			"register" = {
				"roles": [
					"Member",
				],
				"users": [],
				"teams": []
			},
			"unregister" = {
				"roles": [
					"Member",
				],
				"users": [],
				"teams": []
			},
			"update" = {
				"roles": [
					"Member",
				],
				"users": [],
				"teams": []
			},
			"update_properties" = {
				"$bla" = {
					"teams" = [],
					"users" = [],
					"roles" = ["Member"]
				}
			},
			"update_metadata_properties" = {
				"icon" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
				"identifier" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
				"team" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
				"title" = {
					"roles": [
						"Member",
					],
					"users": [],
					"teams": []
				},
			}
		}
	}`, teamName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccPortBlueprintResourceBasic + testAccBaseBlueprintPermissionsConfigUpdate,
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
		},
	})
}
