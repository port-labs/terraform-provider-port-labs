package webhook_test

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
			string_props = {
				"author" = {
					type = "string"
					title = "text"
				}
				"url" = {
					type = "string"
					title = "text"
				}
			}
		}
	}
	`, identifier)
}

func TestAccPortWebhookBasic(t *testing.T) {
	identifier := utils.GenID()
	webhookIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_webhook" "create_pr" {
		identifier = "%s"
		title      = "Test"
		icon       = "Terraform"
		enabled    = true
	}`, webhookIdentifier)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_webhook.create_pr", "identifier", webhookIdentifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "title", "Test"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccPortWebhook(t *testing.T) {
	identifier := utils.GenID()
	webhookIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_webhook" "create_pr" {
		identifier = "%s"
		title      = "Test"
		icon       = "Terraform"
		enabled    = true
		security = {
			secret                  = "test"
			signature_header_name   = "X-Hub-Signature-256"
			signature_algorithm     = "sha256"
			signature_prefix        = "sha256="
			request_identifier_path = "body.repository.full_name"
		  }
		mappings = [
			{
			"blueprint" = "%s",
			"filter" = ".headers.\"X-GitHub-Event\" == \"pull_request\"",
			"entity" = {
					"identifier" = ".body.pull_request.id | tostring",
					"title" = ".body.pull_request.title",
					"icon" = "Terraform",
					"team" = "port",
					"properties" = {
						"author" = ".body.pull_request.user.login",
						"url" = ".body.pull_request.html_url"
					}
				}
			}
		]
		lifecycle {
			ignore_changes = [
			  security.secret
			]
		  }
		  depends_on = [
			port_blueprint.microservice
			]		  		
	}`, webhookIdentifier, identifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_webhook.create_pr", "identifier", webhookIdentifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "title", "Test"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "enabled", "true"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_header_name", "X-Hub-Signature-256"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_algorithm", "sha256"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_prefix", "sha256="),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.request_identifier_path", "body.repository.full_name"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.filter", ".headers.\"X-GitHub-Event\" == \"pull_request\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.team", "port"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.author", ".body.pull_request.user.login"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.url", ".body.pull_request.html_url"),
				),
			},
		},
	})
}

func TestAccPortWebhookImport(t *testing.T) {
	identifier := utils.GenID()
	webhookIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_webhook" "create_pr" {
		identifier = "%s"
		title      = "Test"
		icon       = "Terraform"
		enabled    = true
		security = {
			signature_header_name   = "X-Hub-Signature-256"
			signature_algorithm     = "sha256"
			signature_prefix        = "sha256="
			request_identifier_path = "body.repository.full_name"
		  }
		mappings = [
			{
			"blueprint" = "%s",
			"filter" = ".headers.\"X-GitHub-Event\" == \"pull_request\"",
			"entity" = {
					"identifier" = ".body.pull_request.id | tostring",
					"title" = ".body.pull_request.title",
					"icon" = "Terraform",
					"team" = "port",
					"properties" = {
						"author" = ".body.pull_request.user.login",
						"url" = ".body.pull_request.html_url"
					}
				}
			}
		]
		lifecycle {
			ignore_changes = [
			  security.secret
			]
		  }
		depends_on = [
		  port_blueprint.microservice
		  ]	
	}`, webhookIdentifier, identifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_webhook.create_pr", "identifier", webhookIdentifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "title", "Test"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "enabled", "true"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_header_name", "X-Hub-Signature-256"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_algorithm", "sha256"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_prefix", "sha256="),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.request_identifier_path", "body.repository.full_name"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.filter", ".headers.\"X-GitHub-Event\" == \"pull_request\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.team", "port"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.author", ".body.pull_request.user.login"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.url", ".body.pull_request.html_url"),
				),
			},
			{
				ResourceName:      "port_webhook.create_pr",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     webhookIdentifier,
			},
		},
	})
}
