package page_permissions_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func generateBlueprintPluralIdentifier() string {
	// set the blueprint identifier to end with plural so the page identifier won't add plural to the page identifier
	return fmt.Sprintf("%s-tests", utils.GenID()[:10])
}

func testAccCreateBlueprintConfig(blueprintIdentifier string) string {
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
	}`, blueprintIdentifier)
}
func TestAccPortPagePermissionsBasic(t *testing.T) {
	blueprintIdentifier := generateBlueprintPluralIdentifier()
	var testAccBaseBlueprintConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier)

	var testAccBasePagePermissionsConfigUpdate = testAccCreateBlueprintConfig(blueprintIdentifier) + `

	resource "port_page_permissions" "microservice_permissions" {
		page_identifier = port_blueprint.microservice.identifier
		read_permissions = {
			"roles": [
				"Member",
			],
			"users": [],
			"teams": []
			}
	}`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBaseBlueprintConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", blueprintIdentifier),
				),
			},
			{
				Config: testAccBasePagePermissionsConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "page_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read_permissions.roles.#", "1"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read_permissions.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read_permissions.users.#", "0"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read_permissions.teams.#", "0"),
				),
			},
			{
				Config: testAccBaseBlueprintConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", blueprintIdentifier),
				),
			},
		},
	})
}

func TestAccPortPagePermissionsUpdateWithUsers(t *testing.T) {
	blueprintIdentifier := generateBlueprintPluralIdentifier()
	teamName := utils.GenID()

	var testAccBasePagePermissionsConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier)
	var testAccBasePagePermissionsConfigUpdate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	
	resource "port_team" "team" {
		name = "%s"
		description = "Test description"
		users = []
	}

	resource "port_page_permissions" "microservice_permissions" {
		page_identifier = port_blueprint.microservice.identifier
		read_permissions = {
			"roles": [
				"Member",
			],
		  "users": ["devops-port@port-test.io"],
		  "teams": [port_team.team.name],
			}
	}`, teamName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBasePagePermissionsConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", blueprintIdentifier),
				),
			},
			{
				Config: testAccBasePagePermissionsConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "page_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read_permissions.roles.#", "1"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read_permissions.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read_permissions.users.#", "1"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read_permissions.users.0", "devops-port@port-test.io"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read_permissions.teams.#", "1"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read_permissions.teams.0", teamName),
				),
			},
		},
	})
}
