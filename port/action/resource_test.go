package action_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/acctest"
)

func genID() string {
	id, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("t-%s", id[:18])
}

func TestAccPortActionBasic(t *testing.T) {
	identifier := genID()
	actionIdentifier := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
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
	resource "port-labs_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port-labs_blueprint.microservice.id
		trigger = "DAY-2"
		kafka_method = {}
	}`, identifier, actionIdentifier)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "trigger", "DAY-2"),
				),
			},
		},
	})
}

func TestAccPortActionWebhookInvocation(t *testing.T) {
	identifier := genID()
	actionIdentifier := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF test microservice"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_prop = {
				"text" = {
					title = "text"
				}
			}
		}
	}
	resource "port-labs_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port-labs_blueprint.microservice.id
		trigger = "DAY-2"
		webhook_method = {
			url = "https://example.com"
			agent = true
		}
	}`, identifier, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "webhook_method.url", "https://example.com"),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "webhook_method.agent", "true"),
				),
			},
		},
	})
}

func TestAccPortActionAzureInvocation(t *testing.T) {
	identifier := genID()
	actionIdentifier := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF test microservice"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_prop = {
				"text" = {
					title = "text"
				}
			}
		}
	}
	resource "port-labs_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port-labs_blueprint.microservice.id
		trigger = "DAY-2"
		azure_method = {
			org = "port",
			webhook = "https://example.com"
		}
	}`, identifier, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "azure_method.org", "port"),
					resource.TestCheckResourceAttr("port-labs_action.create_microservice", "azure_method.webhook", "https://example.com"),
				),
			},
		},
	})
}
