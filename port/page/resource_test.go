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
        "blueprint" : port_blueprint.microservice.identifier,
        "dataset" : {
          "combinator" : "and",
          "rules" : [
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
        "blueprint" : port_blueprint.microservice.identifier,
        "dataset" : {
          "combinator" : "and",
          "rules" : [
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

func TestAccPortPageResourceWithoutFilters(t *testing.T) {
	pageIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortPageResourceWithoutFilters = fmt.Sprintf(`
resource "port_page" "page_without_filters" {
  identifier = "%s"
  title      = "Page Without Filters"
  icon       = "Dashboard"
  type       = "dashboard"
  widgets    = [
    jsonencode(
      {
        "id" = "dashboardWidget"
        "type" = "dashboard-widget"
        "layout" = [
          {
            "height" = 400
            "columns" = [
              {
                "id" = "widget1"
                "size" = 12
              }
            ]
          }
        ]
        "widgets" = [
          {
            "id" = "widget1"
            "type" = "markdown"
            "title" = "Test Widget"
            "markdown" = "# Test"
          }
        ]
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
				Config: acctest.ProviderConfig + testAccPortPageResourceWithoutFilters,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_page.page_without_filters", "identifier", pageIdentifier),
					resource.TestCheckResourceAttr("port_page.page_without_filters", "title", "Page Without Filters"),
					resource.TestCheckResourceAttr("port_page.page_without_filters", "type", "dashboard"),
					resource.TestCheckResourceAttr("port_page.page_without_filters", "widgets.#", "1"),
					resource.TestCheckResourceAttr("port_page.page_without_filters", "page_filters.#", "0"),
				),
			},
		},
	})
}

func TestAccPortPageResourceWithFilters(t *testing.T) {
	serviceBlueprintIdentifier := utils.GenID()
	clusterBlueprintIdentifier := utils.GenID()
	pageIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}

	var testAccPortPageResourceWithFilters = fmt.Sprintf(`
resource "port_blueprint" "service" {
  identifier  = "%s"
  title       = "Service Test"
  icon        = "Microservice"
  description = "Service entities for testing filters"

  properties = {
    string_props = {
      language = {
        type = "string"
      }
    }
  }
}

resource "port_blueprint" "cluster" {
  identifier  = "%s"
  title       = "Cluster Test"
  icon        = "Cluster"
  description = "Cluster entities for testing filters"

  properties = {
    string_props = {
      studio = {
        type = "string"
      }
    }
  }
}

resource "port_entity" "service_ruby" {
  title     = "Ruby Service"
  blueprint = port_blueprint.service.identifier

  properties = {
    string_props = {
      language = "Ruby"
    }
  }

  depends_on = [port_blueprint.service]
}

resource "port_entity" "service_python" {
  title     = "Python Service"
  blueprint = port_blueprint.service.identifier

  properties = {
    string_props = {
      language = "Python"
    }
  }

  depends_on = [port_blueprint.service]
}

resource "port_entity" "cluster_studio1" {
  title     = "Studio 1 Cluster"
  blueprint = port_blueprint.cluster.identifier

  properties = {
    string_props = {
      studio = "Studio1"
    }
  }

  depends_on = [port_blueprint.cluster]
}

resource "port_entity" "cluster_studio2" {
  title     = "Studio 2 Cluster"
  blueprint = port_blueprint.cluster.identifier

  properties = {
    string_props = {
      studio = "Studio2"
    }
  }

  depends_on = [port_blueprint.cluster]
}

resource "port_page" "page_with_filters" {
  identifier = "%s"
  title      = "Page With Filters"
  icon       = "Dashboard"
  type       = "dashboard"

  depends_on = [
    port_blueprint.service,
    port_blueprint.cluster,
    port_entity.service_ruby,
    port_entity.service_python,
    port_entity.cluster_studio1,
    port_entity.cluster_studio2
  ]

  page_filters = [
    jsonencode(
      {
        "identifier" = "filter-1"
        "title"      = "Service: language = Ruby"
        "query" = {
          "combinator" = "and"
          "rules" = [
            {
              "value"    = "Ruby"
              "property" = "language"
              "operator" = "="
            }
          ]
          "blueprint" = port_blueprint.service.identifier
        }
      }
    ),
    jsonencode(
      {
        "identifier" = "filter-2"
        "title"      = "Cluster: studio != Studio1"
        "query" = {
          "combinator" = "and"
          "rules" = [
            {
              "value"    = "Studio1"
              "property" = "studio"
              "operator" = "!="
            }
          ]
          "blueprint" = port_blueprint.cluster.identifier
        }
      }
    )
  ]

  widgets = [
    jsonencode(
      {
        "id" = "dashboardWidget"
        "type" = "dashboard-widget"
        "layout" = [
          {
            "height" = 400
            "columns" = [
              {
                "id" = "widget1"
                "size" = 12
              }
            ]
          }
        ]
        "widgets" = [
          {
            "id" = "widget1"
            "type" = "table-entities-explorer"
            "title" = "Services Table"
            "blueprint" = port_blueprint.service.identifier
            "dataset" = {
              "combinator" = "and"
              "rules" = []
            }
          }
        ]
      }
    )
  ]
}
`, serviceBlueprintIdentifier, clusterBlueprintIdentifier, pageIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccPortPageResourceWithFilters,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_page.page_with_filters", "identifier", pageIdentifier),
					resource.TestCheckResourceAttr("port_page.page_with_filters", "title", "Page With Filters"),
					resource.TestCheckResourceAttr("port_page.page_with_filters", "type", "dashboard"),
					resource.TestCheckResourceAttr("port_page.page_with_filters", "widgets.#", "1"),
					resource.TestCheckResourceAttr("port_page.page_with_filters", "page_filters.#", "2"),
				),
			},
		},
	})
}

func testAccCreateBlueprintWithStatusConfig(identifier string) string {
	return fmt.Sprintf(`
	resource "port_blueprint" "microservice_with_status" {
		title = "TF test microservice with status"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_props = {
				"status" = {
					type = "string"
					title = "Status"
					enum = ["Active", "Inactive", "Deprecated"]
				}
			}
		}
	}
	`, identifier)
}

func TestAccPortPageResourceWithDataViewModeTable(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	pageIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortPageResourceDataViewModeTable = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`

resource "port_page" "microservice_table_page" {
  identifier            = "%s"
  title                 = "Microservices Table View"
  icon                  = "Microservice"
  blueprint             = port_blueprint.microservice.identifier
  type                  = "blueprint-entities"
  widgets               = [
    jsonencode(
      {
        "id" : "microservice-table",
        "type" : "table-entities-explorer",
        "blueprint" : port_blueprint.microservice.identifier,
        "dataViewMode" : "table",
        "dataset" : {
          "combinator" : "and",
          "rules" : [
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
				Config: acctest.ProviderConfig + testAccPortPageResourceDataViewModeTable,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_page.microservice_table_page", "identifier", pageIdentifier),
					resource.TestCheckResourceAttr("port_page.microservice_table_page", "title", "Microservices Table View"),
					resource.TestCheckResourceAttr("port_page.microservice_table_page", "icon", "Microservice"),
					resource.TestCheckResourceAttr("port_page.microservice_table_page", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_page.microservice_table_page", "type", "blueprint-entities"),
					resource.TestCheckResourceAttr("port_page.microservice_table_page", "widgets.#", "1"),
				),
			},
		},
	})
}

func TestAccPortPageResourceWithDataViewModeBoard(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	pageIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortPageResourceDataViewModeBoard = testAccCreateBlueprintWithStatusConfig(blueprintIdentifier) + fmt.Sprintf(`

resource "port_page" "microservice_board_page" {
  identifier            = "%s"
  title                 = "Microservices Board View"
  icon                  = "Microservice"
  blueprint             = port_blueprint.microservice_with_status.identifier
  type                  = "blueprint-entities"
  widgets               = [
    jsonencode(
      {
        "id" : "microservice-board",
        "type" : "table-entities-explorer",
        "blueprint" : port_blueprint.microservice_with_status.identifier,
        "dataViewMode" : "board",
        "blueprintConfig" : {
          (port_blueprint.microservice_with_status.identifier) : {
            "groupSettings" : {
              "groupBy" : ["status"]
            }
          }
        },
        "dataset" : {
          "combinator" : "and",
          "rules" : [
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
				Config: acctest.ProviderConfig + testAccPortPageResourceDataViewModeBoard,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_page.microservice_board_page", "identifier", pageIdentifier),
					resource.TestCheckResourceAttr("port_page.microservice_board_page", "title", "Microservices Board View"),
					resource.TestCheckResourceAttr("port_page.microservice_board_page", "icon", "Microservice"),
					resource.TestCheckResourceAttr("port_page.microservice_board_page", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_page.microservice_board_page", "type", "blueprint-entities"),
					resource.TestCheckResourceAttr("port_page.microservice_board_page", "widgets.#", "1"),
				),
			},
		},
	})
}
