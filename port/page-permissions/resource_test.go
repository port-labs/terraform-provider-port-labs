package page_permissions_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func createPage(identifier string) string {
	return fmt.Sprintf(`

resource "port_page" "microservice_dashboard_page" {
  identifier            = "%s"
  title                 = "dashboards"
  icon                  = "GitHub"
  type                  = "dashboard"
  widgets               = [
    jsonencode(
      {
        "id" : "dashboardWidget",
        "layout" : [
          {
            "height" : 400,
            "columns" : [
              {
                "id" : "microserviceGuide",
                "size" : 12
              }
            ]
          }
        ],
        "type" : "dashboard-widget",
        "widgets" : [
          {
            "title" : "Microservices Guide",
            "icon" : "BlankPage",
            "markdown" : "# This is the new Microservice Dashboard",
            "type" : "markdown",
            "description" : "",
            "id" : "microserviceGuide"
          }
        ],
      }
    )
  ]
}
`, identifier)
}

func TestAccPortPagePermissionsBasic(t *testing.T) {
	pageIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortPageResourceBasic = createPage(pageIdentifier)

	var testAccBasePagePermissionsConfigUpdate = `

	resource "port_page_permissions" "microservice_permissions" {
		page_identifier = port_page.microservice_dashboard_page.identifier
		read = {
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
				Config: testAccPortPageResourceBasic + testAccBasePagePermissionsConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "identifier", pageIdentifier),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "title", "dashboards"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "icon", "GitHub"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "type", "dashboard"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "widgets.#", "1"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "page_identifier", pageIdentifier),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read.roles.#", "1"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read.users.#", "0"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read.teams.#", "0"),
				),
			},
		},
	})
}

func TestAccPortPagePermissionsUpdateWithUsers(t *testing.T) {
	pageIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortPageResourceBasic = createPage(pageIdentifier)

	teamName := utils.GenID()

	var testAccBasePagePermissionsConfigUpdate = fmt.Sprintf(`

	resource "port_team" "team" {
		name = "%s"
		description = "Test description"
		users = []
	}

	resource "port_page_permissions" "microservice_permissions" {
		page_identifier = port_page.microservice_dashboard_page.identifier
		read = {
			"roles": [
				"Member",
			],
		  "users": [],
		  "teams": [port_team.team.name],
			}
	}`, teamName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPortPageResourceBasic + testAccBasePagePermissionsConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "identifier", pageIdentifier),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "title", "dashboards"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "icon", "GitHub"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "type", "dashboard"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "widgets.#", "1"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "page_identifier", pageIdentifier),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read.roles.#", "1"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read.roles.0", "Member"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read.users.#", "0"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read.teams.#", "1"),
					resource.TestCheckResourceAttr("port_page_permissions.microservice_permissions", "read.teams.0", teamName),
				),
			},
		},
	})
}
