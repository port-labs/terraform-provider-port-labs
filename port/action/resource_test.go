package action_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func testAccCreateBlueprintConfig(identifier string) string {
	return fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF test microservice"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_prop = {
			"text" = {
				type = "string"
				title = "text"
				}
			}
		}
	}
	`, identifier)
}
func TestAccPortActionBasic(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		kafka_method = {}
	}`, actionIdentifier)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
				),
			},
		},
	})
}

func TestAccPortActionWebhookInvocation(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		webhook_method = {
			url = "https://example.com"
			agent = true
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "webhook_method.url", "https://example.com"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "webhook_method.agent", "true"),
				),
			},
		},
	})
}

func TestAccPortActionAzureInvocation(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		azure_method = {
			org = "port",
			webhook = "https://example.com"
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "azure_method.org", "port"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "azure_method.webhook", "https://example.com"),
				),
			},
		},
	})
}

func TestAccPortActionGithubInvocation(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		github_method = {
			org = "port",
			repo = "terraform-provider-port",
			workflow = "main.yml"
			omit_payload = true
			omit_user_inputs = true
			report_workflow_status = false
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "github_method.org", "port"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "github_method.repo", "terraform-provider-port"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "github_method.workflow", "main.yml"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "github_method.omit_payload", "true"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "github_method.omit_user_inputs", "true"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "github_method.report_workflow_status", "false"),
				),
			},
		},
	})
}

func TestAccPortActionImport(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		webhook_method = {
			url = "https://example.com"
		}
		user_properties = {
			"string_prop" = {
				"myStringIdentifier" = {
					"title" = "My String Identifier"
					"required" = true
				}
			}
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "webhook_method.url", "https://example.com"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_prop.myStringIdentifier.title", "My String Identifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_prop.myStringIdentifier.required", "true"),
				),
			},
			{
				ResourceName:      "port_action.create_microservice",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s:%s", blueprintIdentifier, actionIdentifier),
			},
		},
	})
}

func TestAccPortActionUpdate(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		webhook_method = {
			url = "https://example.com"
		}
		user_properties = {
			"string_prop" = {
				"myStringIdentifier" = {
					"title" = "My String Identifier"
					"required" = true
				}
			}
		}
	}`, actionIdentifier)

	var testAccActionConfigUpdate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		webhook_method = {
			url = "https://example.com"
		}
		user_properties = {
			"string_prop" = {
				"myStringIdentifier2" = {
					"title" = "My String Identifier"
					"required" = false
				}
			}
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "webhook_method.url", "https://example.com"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_prop.myStringIdentifier.title", "My String Identifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_prop.myStringIdentifier.required", "true"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "webhook_method.url", "https://example.com"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_prop.myStringIdentifier2.title", "My String Identifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_prop.myStringIdentifier2.required", "false"),
				),
			},
		},
	})
}
