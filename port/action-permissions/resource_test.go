package action_permissions_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func testAccCreateBlueprintAndActionConfig(blueprintIdentifier string, actionIdentifier string) string {
	return fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF test microservice"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_props = {
			"text" = {
				type = "string"
				title = "text"
				}
			}
		}
	}

	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.identifier
		trigger = "DAY-2"
		kafka_method = {}
	}`, blueprintIdentifier, actionIdentifier)
}
func TestAccPortActionPermissionsBasic(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionPermissionsConfigCreate = testAccCreateBlueprintAndActionConfig(blueprintIdentifier, actionIdentifier) + `
	resource "port_action_permissions" "create_microservice_permissions" {
	  action_identifier = port_action.create_microservice.identifier
	  blueprint_identifier = port_blueprint.microservice.identifier
	  permissions = {
		"execute": {
		  "roles": [
			"Member",
		  ],
		  "users": [],
		  "teams": [],
		  "owned_by_team": false
		},
		"approve": {
		  "roles": [],
		  "users": [],
		  "teams": []
		}
	  }
	}`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccActionPermissionsConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "action_identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.#", "1"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.teams.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.owned_by_team", "false"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.roles.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.teams.#", "0"),
				),
			},
		},
	})
}

func TestAccPortActionPermissionsUpdate(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	actionIdentifier := utils.GenID()
	teamName := utils.GenID()
	var testAccActionPermissionsConfigCreate = testAccCreateBlueprintAndActionConfig(blueprintIdentifier, actionIdentifier) + `
	resource "port_action_permissions" "create_microservice_permissions" {
	  action_identifier = port_action.create_microservice.identifier
	  blueprint_identifier = port_blueprint.microservice.identifier
	  permissions = {
		"execute": {
		  "roles": [
			"Member",
		  ],
		  "users": [],
		  "teams": [],
		  "owned_by_team": false
		},
		"approve": {
		  "roles": [],
		  "users": [],
		  "teams": []
		}
	  }
	}`
	var testAccActionPermissionsConfigUpdate = testAccCreateBlueprintAndActionConfig(blueprintIdentifier, actionIdentifier) + fmt.Sprintf(`
   	resource "port_team" "team" {
		name = "%s"
		description = "Test description"
		users = []
	}

	resource "port_action_permissions" "create_microservice_permissions" {
	  action_identifier = port_action.create_microservice.identifier
	  blueprint_identifier = port_blueprint.microservice.identifier
	  permissions = {
		"execute": {
		  "roles": [
			"Member",
		  ],
		  "users": ["devops-port@port-test.io"],
		  "teams": [port_team.team.name],
		  "owned_by_team": false
		},
		"approve": {
		  "roles": [
			"Member",
		  ],
		  "users": [],
		  "teams": []
		}
	  }
	}`, teamName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccActionPermissionsConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "action_identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.#", "1"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.teams.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.owned_by_team", "false"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.roles.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.teams.#", "0"),
				),
			},
			{
				Config: testAccActionPermissionsConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "action_identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.#", "1"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.users.#", "1"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.users.0", "devops-port@port-test.io"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.teams.#", "1"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.teams.0", teamName),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.owned_by_team", "false"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.roles.#", "1"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.teams.#", "0"),
				),
			},
		},
	})
}

func TestAccPortActionPermissionsWithPolicy(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionPermissionsConfigCreate = testAccCreateBlueprintAndActionConfig(blueprintIdentifier, actionIdentifier) + `
	resource "port_action_permissions" "create_microservice_permissions" {
	  action_identifier = port_action.create_microservice.identifier
	  blueprint_identifier = port_blueprint.microservice.identifier	
	  permissions = {
		"execute": {
		  "roles": [
			"Member",	
		  ],	
		  "users": [],
		  "teams": [],
		  "owned_by_team": false
		},	
		"approve": {	
		  "roles": [],
		  "users": [],
		  "teams": [],
		  "policy": jsonencode(
        {
          queries: {
            executingUser: {
              rules: [
                {
                  value: "user",
                  operator: "=",
                  property: "$blueprint"
                },
                {
                  value: "{{.trigger.user.email}}",
                  operator: "=",
                  property: "$identifier"
                },
                {
                    value: "true",
                    operator: "=",
                    property: "$owned_by_team"

                }
              ],
              combinator: "or"
            }
          },
          conditions: [
          "true"]
          })
		}
      }
	}`

	var testAccActionPermissionsConfigUpdate = testAccCreateBlueprintAndActionConfig(blueprintIdentifier, actionIdentifier) + `
	resource "port_action_permissions" "create_microservice_permissions" {
	  action_identifier = port_action.create_microservice.identifier
	  blueprint_identifier = port_blueprint.microservice.identifier	
	  permissions = {
		"execute": {
		  "roles": [
			"Member",	
		  ],	
		  "users": [],
		  "teams": [],
		  "owned_by_team": false
		},	
		"approve": {	
		  "roles": [],
		  "users": [],
		  "teams": []
		}
	  }
    }`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccActionPermissionsConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "action_identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.#", "1"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.teams.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.owned_by_team", "false"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.roles.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.teams.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.policy", "{\"conditions\":[\"true\"],\"queries\":{\"executingUser\":{\"combinator\":\"or\",\"rules\":[{\"operator\":\"=\",\"property\":\"$blueprint\",\"value\":\"user\"},{\"operator\":\"=\",\"property\":\"$identifier\",\"value\":\"{{.trigger.user.email}}\"},{\"operator\":\"=\",\"property\":\"$owned_by_team\",\"value\":\"true\"}]}}}"),
				),
			},
			{
				Config: testAccActionPermissionsConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "action_identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.#", "1"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.teams.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.owned_by_team", "false"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.roles.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.teams.#", "0"),
				),
			},
		},
	})
}

func TestAccPortActionPermissionsWithPolicyUpdate(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionPermissionsConfigCreate = testAccCreateBlueprintAndActionConfig(blueprintIdentifier, actionIdentifier) + `
	resource "port_action_permissions" "create_microservice_permissions" {
	  action_identifier = port_action.create_microservice.identifier
	  blueprint_identifier = port_blueprint.microservice.identifier	
	  permissions = {
		"execute": {
		  "roles": [
			"Member",	
		  ],	
		  "users": [],
		  "teams": [],
		  "owned_by_team": false
		},	
		"approve": {	
		  "roles": [],
		  "users": [],
		  "teams": [],
		  "policy": jsonencode(
		{
		  queries: {
			executingUser: {
			  rules: [
				{
				  value: "user",
				  operator: "=",
				  property: "$blueprint"
				},
				{
				  value: "{{.trigger.user.email}}",
				  operator: "=",
				  property: "$identifier"
				},
				{
					value: "true",
					operator: "=",
					property: "$owned_by_team"

				}
			  ],
			  combinator: "and"
			}
		  },
		  conditions: [
		  "true"]
		  })
		}
	  }
	}`

	var testAccActionPermissionsConfigUpdate = testAccCreateBlueprintAndActionConfig(blueprintIdentifier, actionIdentifier) + `
	resource "port_action_permissions" "create_microservice_permissions" {
	  action_identifier = port_action.create_microservice.identifier
	  blueprint_identifier = port_blueprint.microservice.identifier	
	  permissions = {
		"execute": {
		  "roles": [
			"Member",	
		  ],	
		  "users": [],
		  "teams": [],
		  "owned_by_team": false
		},	
		"approve": {	
		  "roles": [],
		  "users": [],
		  "teams": [],
		  "policy": jsonencode(
		{
		  queries: {
			executingUser: {
			  rules: [
				{
				  value: "user",
				  operator: "=",
				  property: "$blueprint"
				},
				{
				  value: "{{.trigger.user.email}}",
				  operator: "=",
				  property: "$identifier"
				},
				{
					value: "true",
					operator: "=",
					property: "$owned_by_team"

				}
			  ],
			  combinator: "or"
			}
		  },
		  conditions: [
		  "true"]	
		  })
		}	
	  }
	}`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccActionPermissionsConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "action_identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.#", "1"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.teams.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.owned_by_team", "false"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.roles.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.teams.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.policy", "{\"conditions\":[\"true\"],\"queries\":{\"executingUser\":{\"combinator\":\"and\",\"rules\":[{\"operator\":\"=\",\"property\":\"$blueprint\",\"value\":\"user\"},{\"operator\":\"=\",\"property\":\"$identifier\",\"value\":\"{{.trigger.user.email}}\"},{\"operator\":\"=\",\"property\":\"$owned_by_team\",\"value\":\"true\"}]}}}"),
				),
			},
			{
				Config: testAccActionPermissionsConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "action_identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.#", "1"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.teams.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.owned_by_team", "false"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.roles.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.teams.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.policy", "{\"conditions\":[\"true\"],\"queries\":{\"executingUser\":{\"combinator\":\"or\",\"rules\":[{\"operator\":\"=\",\"property\":\"$blueprint\",\"value\":\"user\"},{\"operator\":\"=\",\"property\":\"$identifier\",\"value\":\"{{.trigger.user.email}}\"},{\"operator\":\"=\",\"property\":\"$owned_by_team\",\"value\":\"true\"}]}}}"),
				),
			},
		},
	})
}

func TestAccPortActionPermissionsImportState(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionPermissionsConfigCreate = testAccCreateBlueprintAndActionConfig(blueprintIdentifier, actionIdentifier) + `
	resource "port_action_permissions" "create_microservice_permissions" {
	  action_identifier = port_action.create_microservice.identifier
	  blueprint_identifier = port_blueprint.microservice.identifier	
	  permissions = {
		"execute": {
		  "roles": [
			"Member",	
		  ],	
		  "users": [],
		  "teams": [],
		  "owned_by_team": false
		},	
		"approve": {	
		  "roles": [],
		  "users": [],
		  "teams": [],
		  "policy": jsonencode(
		{
		  queries: {
			executingUser: {
			  rules: [
				{
				  value: "user",
				  operator: "=",
				  property: "$blueprint"
				},
				{
				  value: "{{.trigger.user.email}}",
				  operator: "=",
				  property: "$identifier"
				},
				{
					value: "true",
					operator: "=",
					property: "$owned_by_team"
				}
			  ],
			  combinator: "and"
			}
		  },
		  conditions: [
		  "true"]
		  })
		}
	  }
	}`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccActionPermissionsConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "action_identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.#", "1"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.teams.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.owned_by_team", "false"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.roles.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.teams.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.policy", "{\"conditions\":[\"true\"],\"queries\":{\"executingUser\":{\"combinator\":\"and\",\"rules\":[{\"operator\":\"=\",\"property\":\"$blueprint\",\"value\":\"user\"},{\"operator\":\"=\",\"property\":\"$identifier\",\"value\":\"{{.trigger.user.email}}\"},{\"operator\":\"=\",\"property\":\"$owned_by_team\",\"value\":\"true\"}]}}}"),
				),
			},
			{
				ResourceName:      "port_action_permissions.create_microservice_permissions",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPortActionWithEmptyFieldsExpectDefaultsToApply(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionPermissionsConfigCreate = testAccCreateBlueprintAndActionConfig(blueprintIdentifier, actionIdentifier) + `
	resource "port_action_permissions" "create_microservice_permissions" {
	  action_identifier = port_action.create_microservice.identifier
	  blueprint_identifier = port_blueprint.microservice.identifier	
	  permissions = {
		"execute": {}
		"approve": {}
		}
	}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccActionPermissionsConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "action_identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.roles.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.teams.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.execute.owned_by_team", "true"),
					resource.TestCheckNoResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.policy"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.roles.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.users.#", "0"),
					resource.TestCheckResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.teams.#", "0"),
					resource.TestCheckNoResourceAttr("port_action_permissions.create_microservice_permissions", "permissions.approve.policy"),
				),
			},
		},
	})
}
