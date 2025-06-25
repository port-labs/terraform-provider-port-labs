package webhook_test

import (
	"fmt"
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

func testAccCreateBlueprintConfigWithRelations(identifier string, authorIdentifier string) string {
	return fmt.Sprintf(`
	resource "port_blueprint" "author" {
		title = "Author"
		icon = "User"
		identifier = "%s"
		properties = {
			string_props = {
				"name" = {
					type = "string"
					title = "Name"
				}
			}
		}
	}
	
	resource "port_blueprint" "microservice" {
		title = "TF test microservice"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_props = {
				"url" = {
					type = "string"
					title = "URL"
				}
			}
		}
		relations = {
			"author" = {
				title = "Author"
				target = port_blueprint.author.identifier
			}
		}
	}
	`, authorIdentifier, identifier)
}

func testAccCreateBlueprintConfigWithMultipleRelations(identifier string, authorIdentifier string, teamIdentifier string) string {
	return fmt.Sprintf(`
	resource "port_blueprint" "author" {
		title = "Author"
		icon = "User"
		identifier = "%s"
		properties = {
			string_props = {
				"name" = {
					type = "string"
					title = "Name"
				}
			}
		}
	}
	
	resource "port_blueprint" "team" {
		title = "Team"
		icon = "Team"
		identifier = "%s"
		properties = {
			string_props = {
				"name" = {
					type = "string"
					title = "Team Name"
				}
			}
		}
	}
	
	resource "port_blueprint" "microservice" {
		title = "TF test microservice"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_props = {
				"url" = {
					type = "string"
					title = "URL"
				}
			}
		}
		relations = {
			"author" = {
				title = "Author"
				target = port_blueprint.author.identifier
			}
			"team" = {
				title = "Team"
				target = port_blueprint.team.identifier
			}
		}
	}
	`, authorIdentifier, teamIdentifier, identifier)
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

func TestAccPortWebhookWithAllOperationOptions(t *testing.T) {
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
			"filter" = ".headers.\"X-GitHub-Event\" == \"pull_request\" and .body.pull_request.action == \"opened\"",
			"items_to_parse" = ".body.pull_request",
			"operation" = {
				"type" = "create"
			},
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
			"filter" = ".headers.\"X-GitHub-Event\" == \"pull_request\" and .body.pull_request.state == \"closed\"",
			"items_to_parse" = ".body.pull_request",
			"operation": {
				"type": "delete"	
			},
			"entity" = {
					"identifier" = ".body.pull_request.id | tostring"
				}
			},
			{
			"blueprint" = port_blueprint.microservice.identifier,
			"filter" = ".headers.\"X-GitHub-Event\" == \"pull_request\" and .body.pull_request.state == \"edited\"",
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
			},
						{
			"blueprint" = port_blueprint.microservice.identifier,
			"filter" = ".headers.\"X-GitHub-Event\" == \"pull_request\" and .body.pull_request.state == \"closed\"",
			"items_to_parse" = ".body.pull_request",
			"operation": {
				"type": "delete",
				"delete_dependents": true
			},
			"entity" = {
					"identifier" = ".body.pull_request.id | tostring"
				}
			},
			{
			"blueprint" = port_blueprint.microservice.identifier,
			"filter" = ".headers.\"X-GitHub-Event\" == \"pull_request\" and .body.pull_request.state == \"closed\"",
			"items_to_parse" = ".body.pull_request",
			"operation": {
				"type": "delete",
				"delete_dependents": false
			},
			"entity" = {
					"identifier" = ".body.pull_request.id | tostring"
				}
			},
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
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.filter", ".headers.\"X-GitHub-Event\" == \"pull_request\" and .body.pull_request.action == \"opened\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.operation.type", "create"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.items_to_parse", ".body.pull_request"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.icon", "\"Terraform\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.team", "\"port\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.author", ".body.pull_request.user.login"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.url", ".body.pull_request.html_url"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.filter", ".headers.\"X-GitHub-Event\" == \"pull_request\" and .body.pull_request.state == \"closed\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.operation.type", "delete"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.items_to_parse", ".body.pull_request"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.1.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.filter", ".headers.\"X-GitHub-Event\" == \"pull_request\" and .body.pull_request.state == \"edited\""),
					resource.TestCheckNoResourceAttr("port_webhook.create_pr", "mappings.2.operation"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.items_to_parse", ".body.pull_request"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.entity.icon", "\"Terraform\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.entity.team", "\"port\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.entity.properties.author", ".body.pull_request.user.login"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.2.entity.properties.url", ".body.pull_request.html_url"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.3.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.3.filter", ".headers.\"X-GitHub-Event\" == \"pull_request\" and .body.pull_request.state == \"closed\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.3.operation.type", "delete"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.3.operation.delete_dependents", "true"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.4.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.4.filter", ".headers.\"X-GitHub-Event\" == \"pull_request\" and .body.pull_request.state == \"closed\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.4.operation.type", "delete"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.4.operation.delete_dependents", "false"),
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

func TestAccPortWebhookCreateWithRelations(t *testing.T) {
	identifier := utils.GenID()
	authorIdentifier := utils.GenID()
	teamIdentifier := utils.GenID()
	webhookIdentifier := utils.GenID()

	// Test case 1: JSON relations with combinator/rules structure
	var testPortWebhookConfigJSON = testAccCreateBlueprintConfigWithRelations(identifier, authorIdentifier) + fmt.Sprintf(`
	resource "port_webhook" "create_pr" {
		identifier = "%s"
		title      = "Webhook with json relation"
		icon       = "Terraform"
  		enabled    = true
 		mappings = [
    		{
      			blueprint = port_blueprint.microservice.identifier
				operation = { "type" = "create" }
				filter    = ".headers.\"x-github-event\" == \"pull_request\""
				entity = {
					identifier = ".body.pull_request.id | tostring"
					title      = ".body.pull_request.title"
					properties = {
						url = ".body.pull_request.html_url"
					}
					relations = {
						author = jsonencode({
							combinator = "'and'",
							rules = [
								{
									property = "'$identifier'"
									operator = "'='"
									value    = ".body.pull_request.user.login | tostring"
								}
							]
						})
        			}
      			}
    		}
  		]
		depends_on = [
			port_blueprint.microservice,
			port_blueprint.author
		]
	}`, webhookIdentifier)

	// Test case 2: String relations (simple string values)
	var testPortWebhookConfigString = testAccCreateBlueprintConfigWithRelations(identifier, authorIdentifier) + fmt.Sprintf(`
	resource "port_webhook" "create_pr" {
		identifier = "%s"
		title      = "Webhook with string relation"
		icon       = "Terraform"
  		enabled    = true
 		mappings = [
    		{
      			blueprint = port_blueprint.microservice.identifier
				operation = { "type" = "create" }
				filter    = ".headers.\"x-github-event\" == \"pull_request\""
				entity = {
					identifier = ".body.pull_request.id | tostring"
					title      = ".body.pull_request.title"
					properties = {
						url = ".body.pull_request.html_url"
					}
					relations = {
						author = ".body.pull_request.user.login | tostring"
        			}
      			}
    		}
  		]
		depends_on = [
			port_blueprint.microservice,
			port_blueprint.author
		]
	}`, webhookIdentifier)

	// Test case 3: Mixed relations (both JSON and string in the same webhook)
	var testPortWebhookConfigMixed = testAccCreateBlueprintConfigWithMultipleRelations(identifier, authorIdentifier, teamIdentifier) + fmt.Sprintf(`
	resource "port_webhook" "create_pr" {
		identifier = "%s"
		title      = "Webhook with mixed relations"
		icon       = "Terraform"
  		enabled    = true
 		mappings = [
    		{
      			blueprint = port_blueprint.microservice.identifier
				operation = { "type" = "create" }
				filter    = ".headers.\"x-github-event\" == \"pull_request\""
				entity = {
					identifier = ".body.pull_request.id | tostring"
					title      = ".body.pull_request.title"
					properties = {
						url = ".body.pull_request.html_url"
					}
					relations = {
						author = jsonencode({
							combinator = "'and'",
							rules = [
								{
									property = "'$identifier'"
									operator = "'='"
									value    = ".body.pull_request.user.login | tostring"
								}
							]
						})
						team = ".body.repository.owner.login | tostring"
        			}
      			}
    		}
  		]
		depends_on = [
			port_blueprint.microservice,
			port_blueprint.author,
			port_blueprint.team
		]
	}`, webhookIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testPortWebhookConfigJSON,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_webhook.create_pr", "identifier", webhookIdentifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "title", "Webhook with json relation"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "enabled", "true"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.operation.type", "create"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.filter", ".headers.\"x-github-event\" == \"pull_request\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.url", ".body.pull_request.html_url"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.relations.author", `{"combinator":"'and'","rules":[{"operator":"'='","property":"'$identifier'","value":".body.pull_request.user.login | tostring"}]}`),
				),
			},
			{
				Config: acctest.ProviderConfig + testPortWebhookConfigString,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_webhook.create_pr", "identifier", webhookIdentifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "title", "Webhook with string relation"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "enabled", "true"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.operation.type", "create"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.filter", ".headers.\"x-github-event\" == \"pull_request\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.url", ".body.pull_request.html_url"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.relations.author", ".body.pull_request.user.login | tostring"),
				),
			},
			{
				Config: acctest.ProviderConfig + testPortWebhookConfigMixed,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_webhook.create_pr", "identifier", webhookIdentifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "title", "Webhook with mixed relations"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "enabled", "true"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.operation.type", "create"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.filter", ".headers.\"x-github-event\" == \"pull_request\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.url", ".body.pull_request.html_url"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.relations.author", `{"combinator":"'and'","rules":[{"operator":"'='","property":"'$identifier'","value":".body.pull_request.user.login | tostring"}]}`),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.relations.team", ".body.repository.owner.login | tostring"),
				),
			},
		},
	})
}

func TestAccPortWebhookCreateWithInvalidRelations(t *testing.T) {
	identifier := utils.GenID()
	authorIdentifier := utils.GenID()
	webhookIdentifier := utils.GenID()

	// Test case 1: Missing combinator field
	var testPortWebhookConfigMissingCombinator = testAccCreateBlueprintConfigWithRelations(identifier, authorIdentifier) + fmt.Sprintf(`
	resource "port_webhook" "create_pr" {
		identifier = "%s"
		title      = "Invalid Relations Test"
		icon       = "Terraform"
  		enabled    = true
 		mappings = [
    		{
      			blueprint = port_blueprint.microservice.identifier
				operation = { "type" = "create" }
				filter    = ".headers.\"x-github-event\" == \"pull_request\""
				entity = {
					identifier = ".body.pull_request.id | tostring"
					title      = ".body.pull_request.title"
					properties = {
						url = ".body.pull_request.html_url"
					}
					relations = {
						author = jsonencode({
							rules = [
								{
									property = "'$identifier'"
									operator = "'='"
									value    = ".body.pull_request.user.login | tostring"
								}
							]
						})
        			}
      			}
    		}
  		]
		depends_on = [
			port_blueprint.microservice,
			port_blueprint.author
		]
	}`, webhookIdentifier)

	// Test case 2: Missing rules field
	var testPortWebhookConfigMissingRules = testAccCreateBlueprintConfigWithRelations(identifier+"2", authorIdentifier+"2") + fmt.Sprintf(`
	resource "port_webhook" "create_pr2" {
		identifier = "%s"
		title      = "Invalid Relations Test"
		icon       = "Terraform"
  		enabled    = true
 		mappings = [
    		{
      			blueprint = port_blueprint.microservice.identifier
				operation = { "type" = "create" }
				filter    = ".headers.\"x-github-event\" == \"pull_request\""
				entity = {
					identifier = ".body.pull_request.id | tostring"
					title      = ".body.pull_request.title"
					properties = {
						url = ".body.pull_request.html_url"
					}
					relations = {
						author = jsonencode({
							combinator = "'and'"
						})
        			}
      			}
    		}
  		]
		depends_on = [
			port_blueprint.microservice,
			port_blueprint.author
		]
	}`, webhookIdentifier+"2")

	// Test case 3: Missing required field in rule (missing operator)
	var testPortWebhookConfigMissingOperator = testAccCreateBlueprintConfigWithRelations(identifier+"3", authorIdentifier+"3") + fmt.Sprintf(`
	resource "port_webhook" "create_pr3" {
		identifier = "%s"
		title      = "Invalid Relations Test"
		icon       = "Terraform"
  		enabled    = true
 		mappings = [
    		{
      			blueprint = port_blueprint.microservice.identifier
				operation = { "type" = "create" }
				filter    = ".headers.\"x-github-event\" == \"pull_request\""
				entity = {
					identifier = ".body.pull_request.id | tostring"
					title      = ".body.pull_request.title"
					properties = {
						url = ".body.pull_request.html_url"
					}
					relations = {
						author = jsonencode({
							combinator = "'and'",
							rules = [
								{
									property = "'$identifier'"
									value    = ".body.pull_request.user.login | tostring"
								}
							]
						})
        			}
      			}
    		}
  		]
		depends_on = [
			port_blueprint.microservice,
			port_blueprint.author
		]
	}`, webhookIdentifier+"3")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      acctest.ProviderConfig + testPortWebhookConfigMissingCombinator,
				ExpectError: regexp.MustCompile("relation.*missing required field.*combinator"),
			},
			{
				Config:      acctest.ProviderConfig + testPortWebhookConfigMissingRules,
				ExpectError: regexp.MustCompile("relation.*missing required field.*rules"),
			},
			{
				Config:      acctest.ProviderConfig + testPortWebhookConfigMissingOperator,
				ExpectError: regexp.MustCompile("relation.*missing required field.*operator"),
			},
		},
	})
}

func TestAccPortWebhookUpdateRelationType(t *testing.T) {
	identifier := utils.GenID()
	authorIdentifier := utils.GenID()
	webhookIdentifier := utils.GenID()

	// Initial config with string relation
	var testAccWebhookConfigStringRelation = testAccCreateBlueprintConfigWithRelations(identifier, authorIdentifier) + fmt.Sprintf(`
	resource "port_webhook" "create_pr" {
		identifier = "%s"
		title      = "Webhook relation type test"
		icon       = "Terraform"
  		enabled    = true
 		mappings = [
    		{
      			blueprint = port_blueprint.microservice.identifier
				operation = { "type" = "create" }
				filter    = ".headers.\"x-github-event\" == \"pull_request\""
				entity = {
					identifier = ".body.pull_request.id | tostring"
					title      = ".body.pull_request.title"
					properties = {
						url = ".body.pull_request.html_url"
					}
					relations = {
						author = ".body.pull_request.user.login | tostring"
        			}
      			}
    		}
  		]
		depends_on = [
			port_blueprint.microservice,
			port_blueprint.author
		]
	}`, webhookIdentifier)

	// Updated config with JSON object relation
	var testAccWebhookConfigJSONRelation = testAccCreateBlueprintConfigWithRelations(identifier, authorIdentifier) + fmt.Sprintf(`
	resource "port_webhook" "create_pr" {
		identifier = "%s"
		title      = "Webhook relation type test"
		icon       = "Terraform"
  		enabled    = true
 		mappings = [
    		{
      			blueprint = port_blueprint.microservice.identifier
				operation = { "type" = "create" }
				filter    = ".headers.\"x-github-event\" == \"pull_request\""
				entity = {
					identifier = ".body.pull_request.id | tostring"
					title      = ".body.pull_request.title"
					properties = {
						url = ".body.pull_request.html_url"
					}
					relations = {
						author = jsonencode({
							combinator = "'and'",
							rules = [
								{
									property = "'$identifier'"
									operator = "'='"
									value    = ".body.pull_request.user.login | tostring"
								}
							]
						})
        			}
      			}
    		}
  		]
		depends_on = [
			port_blueprint.microservice,
			port_blueprint.author
		]
	}`, webhookIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccWebhookConfigStringRelation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_webhook.create_pr", "identifier", webhookIdentifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "title", "Webhook relation type test"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "enabled", "true"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.operation.type", "create"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.filter", ".headers.\"x-github-event\" == \"pull_request\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.url", ".body.pull_request.html_url"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.relations.author", ".body.pull_request.user.login | tostring"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccWebhookConfigJSONRelation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_webhook.create_pr", "identifier", webhookIdentifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "title", "Webhook relation type test"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "enabled", "true"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.blueprint", identifier),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.operation.type", "create"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.filter", ".headers.\"x-github-event\" == \"pull_request\""),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.identifier", ".body.pull_request.id | tostring"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.title", ".body.pull_request.title"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.properties.url", ".body.pull_request.html_url"),
					resource.TestCheckResourceAttr("port_webhook.create_pr", "mappings.0.entity.relations.author", `{"combinator":"'and'","rules":[{"operator":"'='","property":"'$identifier'","value":".body.pull_request.user.login | tostring"}]}`),
				),
			},
		},
	})
}
