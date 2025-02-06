package webhook_test

import (
	"fmt"
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
					resource.TestCheckResourceAttrWith("port_webhook.create_pr", "url",
						func(value string) error {
							if value == "" {
								return fmt.Errorf("value is empty")
							}
							return nil
						}),
					resource.TestCheckResourceAttrWith("port_webhook.create_pr", "webhook_key",
						func(value string) error {
							if value == "" {
								return fmt.Errorf("value is empty")
							}
							return nil
						}),
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
			//secret                  = "test"
			signature_header_name   = "X-Hub-Signature-256"
			signature_algorithm     = "sha256"
			signature_prefix        = "sha256="
			request_identifier_path = ".body.repository.full_name"
		  }
		mappings = [
			{
			"blueprint" = "%s",
			"filter" = ".headers.\"X-GitHub-Event\" == \"pull_request\"",
			"items_to_parse" = ".body.pull_request",
			"entity" = {
					"identifier" = ".body.pull_request.id | tostring",
					"title" = ".body.pull_request.title",
					"icon" = "\"Terraform\"",
					"team" = "\"port\"",
					"properties" = {
						"author" = ".body.pull_request.user.login",
						"url" = ".body.pull_request.html_url"
					}
				}
			}
		]
		lifecycle {
			ignore_changes = [
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
					resource.TestCheckResourceAttrWith("port_webhook.create_pr", "url",
						func(value string) error {
							if value == "" {
								return fmt.Errorf("value is empty")
							}
							return nil
						}),
					resource.TestCheckResourceAttrWith("port_webhook.create_pr", "webhook_key",
						func(value string) error {
							if value == "" {
								return fmt.Errorf("value is empty")
							}
							return nil
						}),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_header_name", "X-Hub-Signature-256"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_algorithm", "sha256"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_prefix", "sha256="),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.request_identifier_path", ".body.repository.full_name"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.filter", ".headers.\"X-GitHub-Event\" == \"pull_request\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.items_to_parse", ".body.pull_request"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.icon", "\"Terraform\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.team", "\"port\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.author", ".body.pull_request.user.login"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.url", ".body.pull_request.html_url"),
				),
			},
		},
	})
}

func TestAccPortWebhookWithOperation(t *testing.T) {
	identifier := utils.GenID()
	webhookIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_webhook" "create_pr" {
		identifier = "%s"
		title      = "Test"
		icon       = "Terraform"
		enabled    = true
		mappings = [
			{
			"blueprint" = port_blueprint.microservice.identifier,
			"filter" = ".headers.\"X-GitHub-Event\" == \"pull_request\" && .body.pull_request.action == \"opened\"",
			"items_to_parse" = ".body.pull_request",
			"operation" = "create",
			"entity" = {
					"identifier" = ".body.pull_request.id | tostring",
					"title" = ".body.pull_request.title",
					"icon" = "\"Terraform\"",
					"team" = "\"port\"",
					"properties" = {
						"author" = ".body.pull_request.user.login",
						"url" = ".body.pull_request.html_url"
					}
				}
			},
			{
			"blueprint" = port_blueprint.microservice.identifier,
			"filter" = ".headers.\"X-GitHub-Event\" == \"pull_request\" && .body.pull_request.state == \"closed\"",
			"items_to_parse" = ".body.pull_request",
			"operation": "delete",
			"entity" = {
					"identifier" = ".body.pull_request.id | tostring",
					"title" = ".body.pull_request.title",
					"icon" = "\"Terraform\"",
					"team" = "\"port\"",
					"properties" = {
						"author" = ".body.pull_request.user.login",
						"url" = ".body.pull_request.html_url"
					}
				}
			},
			{
			"blueprint" = port_blueprint.microservice.identifier,
			"filter" = ".headers.\"X-GitHub-Event\" == \"pull_request\" && .body.pull_request.state == \"edited\"",
			"items_to_parse" = ".body.pull_request",
			"entity" = {
					"identifier" = ".body.pull_request.id | tostring",
					"title" = ".body.pull_request.title",
					"icon" = "\"Terraform\"",
					"team" = "\"port\"",
					"properties" = {
						"author" = ".body.pull_request.user.login",
						"url" = ".body.pull_request.html_url"
					}
				}
			}
		]
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
					resource.TestCheckResourceAttrWith("port_webhook.create_pr", "url",
						func(value string) error {
							if value == "" {
								return fmt.Errorf("value is empty")
							}
							return nil
						}),
					resource.TestCheckResourceAttrWith("port_webhook.create_pr", "webhook_key",
						func(value string) error {
							if value == "" {
								return fmt.Errorf("value is empty")
							}
							return nil
						}),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_header_name", "X-Hub-Signature-256"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_algorithm", "sha256"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_prefix", "sha256="),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.request_identifier_path", ".body.repository.full_name"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.filter", ".headers.\"X-GitHub-Event\" == \"pull_request\" && .body.pull_request.action == \"opened\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.operation", "create"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.items_to_parse", ".body.pull_request"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.icon", "\"Terraform\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.team", "\"port\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.author", ".body.pull_request.user.login"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.url", ".body.pull_request.html_url"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.filter", ".headers.\"X-GitHub-Event\" == \"pull_request\" && .body.pull_request.state == \"closed\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.operation", "delete"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.items_to_parse", ".body.pull_request"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.entity.icon", "\"Terraform\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.entity.team", "\"port\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.entity.properties.author", ".body.pull_request.user.login"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.entity.properties.url", ".body.pull_request.html_url"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.filter", ".headers.\"X-GitHub-Event\" == \"pull_request\" && .body.pull_request.state == \"edited\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.operation", "create"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.items_to_parse", ".body.pull_request"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.entity.icon", "\"Terraform\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.entity.team", "\"port\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.entity.properties.author", ".body.pull_request.user.login"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.entity.properties.url", ".body.pull_request.html_url"),
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
			request_identifier_path = ".body.repository.full_name"
		  }
		mappings = [
			{
			"blueprint" = "%s",
			"filter" = ".headers.\"X-GitHub-Event\" == \"pull_request\"",
			"entity" = {
					"identifier" = ".body.pull_request.id | tostring",
					"title" = ".body.pull_request.title",
					"icon" = "\"Terraform\"",
					"team" = "\"port\"",
					"properties" = {
						"author" = ".body.pull_request.user.login",
						"url" = ".body.pull_request.html_url"
					}
				}
			}
		]
		lifecycle {
			ignore_changes = [
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
					resource.TestCheckResourceAttrWith("port_webhook.create_pr", "url",
						func(value string) error {
							if value == "" {
								return fmt.Errorf("value is empty")
							}
							return nil
						}),
					resource.TestCheckResourceAttrWith("port_webhook.create_pr", "webhook_key",
						func(value string) error {
							if value == "" {
								return fmt.Errorf("value is empty")
							}
							return nil
						}),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_header_name", "X-Hub-Signature-256"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_algorithm", "sha256"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.signature_prefix", "sha256="),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "security.request_identifier_path", ".body.repository.full_name"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.filter", ".headers.\"X-GitHub-Event\" == \"pull_request\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.icon", "\"Terraform\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.team", "\"port\""),
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

func TestAccPortWebhookUpdateIdentifier(t *testing.T) {
	identifier := utils.GenID()
	webhookIdentifier := utils.GenID()
	webhookIdentifierUpdated := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_webhook" "create_pr" {
		identifier = "%s"
		title      = "Test"
		icon       = "Terraform"
		enabled    = true
	}`, webhookIdentifier)

	var testAccActionConfigUpdate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_webhook" "create_pr" {
		identifier = "%s"
		title      = "Test"
		icon       = "Terraform"
		enabled    = true
	}`, webhookIdentifierUpdated)

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
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_webhook.create_pr", "identifier", webhookIdentifierUpdated),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "title", "Test"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "enabled", "true"),
				),
			},
		},
	})
}
