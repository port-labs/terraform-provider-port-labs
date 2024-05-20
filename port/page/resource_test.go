package page_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func testAccCreateBlueprintConfig(identifier string) string {
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
	`, identifier)
}

func TestAccPortPageResourceBasicBetaEnabled(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	pageIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortPageResourceBasic = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`

resource "port_page" "microservice_blueprint_page" {
  identifier            = "%s"
  title                 = "Microservices"
  icon                  = "Microservice"
  blueprint             = port_blueprint.microservice.identifier
  type                  = "blueprint-entities"
  widgets               = [
    jsonencode(
      {
        "id" : "blabla",
        "type" : "table-entities-explorer",
        "dataset" : {
          "combinator" : "and",
          "rules" : [
            {
              "operator" : "=",
              "property" : "$blueprint",
              "value" : "{{blueprint}}"
            }
          ]
        }
      }
    )
  ]
}
`, pageIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccPortPageResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_page.microservice_blueprint_page", "identifier", pageIdentifier),
					resource.TestCheckResourceAttr("port_page.microservice_blueprint_page", "title", "Microservices"),
					resource.TestCheckResourceAttr("port_page.microservice_blueprint_page", "icon", "Microservice"),
					resource.TestCheckResourceAttr("port_page.microservice_blueprint_page", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_page.microservice_blueprint_page", "type", "blueprint-entities"),
					resource.TestCheckResourceAttr("port_page.microservice_blueprint_page", "widgets.#", "1"),
				),
			},
		},
	})
}

func TestAccPortPageResourceBasicBetaDisabled(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	pageIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "false")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortPageResourceBasic = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`

resource "port_page" "microservice_blueprint_page" {
  identifier            = "%s"
  title                 = "Microservices"
  icon                  = "Microservice"
  blueprint             = port_blueprint.microservice.identifier
  type                  = "blueprint-entities"
  widgets               = [
    jsonencode(
      {
        "id" : "blabla",
        "type" : "table-entities-explorer",
        "dataset" : {
          "combinator" : "and",
          "rules" : [
            {
              "operator" : "=",
              "property" : "$blueprint",
              "value" : "{{blueprint}}"
            }
          ]
        }
      }
    )
  ]
}
`, pageIdentifier)

	// expect to fail on beta feature not enabled
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      acctest.ProviderConfig + testAccPortPageResourceBasic,
				ExpectError: regexp.MustCompile("Beta features are not enabled"),
			},
		},
	})
}

func TestAccPortPageResourceCreateDashboardPage(t *testing.T) {
	pageIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortPageResourceBasic = fmt.Sprintf(`

resource "port_page" "microservice_dashboard_page" {
  identifier            = "%s"
  title                 = "dashboards"
  icon                  = "GitHub"
  type                  = "dashboard"
  description           = "My Dashboard Page Description"
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
`, pageIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccPortPageResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "identifier", pageIdentifier),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "description", "My Dashboard Page Description"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "title", "dashboards"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "icon", "GitHub"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "type", "dashboard"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "widgets.#", "1"),
				),
			},
		},
	})
}

func TestAccPortPageResourceCreatePageAfterPage(t *testing.T) {
	pageIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortPageResourceBasic = fmt.Sprintf(`

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
`, pageIdentifier)

	pageIdentifier2 := utils.GenID()
	var testAccPortPageResourceBasic2 = fmt.Sprintf(`

resource "port_page" "microservice_dashboard_page_2" {
  identifier            = "%s"
  title                 = "Microservices_2"
  icon                  = "GitHub"
	after								 	= port_page.microservice_dashboard_page.identifier
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
`, pageIdentifier2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccPortPageResourceBasic + testAccPortPageResourceBasic2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "identifier", pageIdentifier),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "title", "dashboards"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "icon", "GitHub"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "type", "dashboard"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page", "widgets.#", "1"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page_2", "identifier", pageIdentifier2),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page_2", "title", "Microservices_2"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page_2", "after", pageIdentifier),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page_2", "icon", "GitHub"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page_2", "type", "dashboard"),
					resource.TestCheckResourceAttr("port_page.microservice_dashboard_page_2", "widgets.#", "1"),
				),
			},
		},
	})
}
