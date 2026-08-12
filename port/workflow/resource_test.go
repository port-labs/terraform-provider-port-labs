package workflow_test

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
			}
		}
	}
	`, identifier)
}

func TestAccPortWorkflowBasic(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	workflowIdentifier := utils.GenID()

	var testAccWorkflowConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_workflow" "review_pr" {
		identifier                = "%s"
		title                     = "Review PR"
		category                  = "engineering"
		description               = "Reviews pull requests"
		allow_anyone_to_view_runs = false

		node {
			identifier  = "trigger"
			title       = "PR Trigger"
			icon        = "Terraform"
			description = "Runs when the service changes"
			verbose     = true
			links       = ["https://example.com/runs"]
			variables = {
				env = ".diff.after.properties.author"
			}
			event_trigger {
				type                 = "ENTITY_UPDATED"
				blueprint_identifier = port_blueprint.microservice.identifier
				condition {
					expressions = [".diff.after.properties.author != null"]
					combinator  = "and"
				}
			}
		}

		node {
			identifier = "notify"
			webhook {
				url        = "https://example.com/hook"
				method     = "POST"
				body       = jsonencode({ text = "a service changed" })
				on_failure = "continue"
			}
		}

		connections {
			source_identifier = "trigger"
			target_identifier = "notify"
			description       = "notify on change"
		}
	}`, workflowIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccWorkflowConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_workflow.review_pr", "identifier", workflowIdentifier),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "title", "Review PR"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "category", "engineering"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "allow_anyone_to_view_runs", "false"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.0.identifier", "trigger"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.0.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.0.verbose", "true"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.0.links.0", "https://example.com/runs"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.0.event_trigger.type", "ENTITY_UPDATED"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.0.event_trigger.blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.0.event_trigger.condition.combinator", "and"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.1.webhook.url", "https://example.com/hook"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.1.webhook.on_failure", "continue"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "connections.0.source_identifier", "trigger"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "connections.0.target_identifier", "notify"),
				),
			},
		},
	})
}

func TestAccPortWorkflowUpdate(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	workflowIdentifier := utils.GenID()

	var testAccWorkflowConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_workflow" "review_pr" {
		identifier = "%s"
		title      = "Review PR"
		category   = "engineering"

		node {
			identifier = "trigger"
			event_trigger {
				type                 = "ENTITY_UPDATED"
				blueprint_identifier = port_blueprint.microservice.identifier
			}
		}

		node {
			identifier = "notify"
			webhook {
				url = "https://example.com/hook"
			}
		}

		connections {
			source_identifier = "trigger"
			target_identifier = "notify"
		}
	}`, workflowIdentifier)

	var testAccWorkflowConfigUpdate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_workflow" "review_pr" {
		identifier = "%s"
		title      = "Review PR (updated)"
		category   = "platform"

		node {
			identifier = "trigger"
			event_trigger {
				type                 = "ENTITY_UPDATED"
				blueprint_identifier = port_blueprint.microservice.identifier
			}
		}

		node {
			identifier = "notify"
			webhook {
				url    = "https://example.com/hook-updated"
				method = "PUT"
			}
		}

		node {
			identifier = "record"
			upsert_entity {
				blueprint_identifier = port_blueprint.microservice.identifier
				mapping {
					identifier = "run-{{ .run.id }}"
					title      = "Run record"
				}
			}
		}

		connections {
			source_identifier = "trigger"
			target_identifier = "notify"
		}

		connections {
			source_identifier = "notify"
			target_identifier = "record"
		}
	}`, workflowIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccWorkflowConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_workflow.review_pr", "title", "Review PR"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "category", "engineering"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.#", "2"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.1.webhook.url", "https://example.com/hook"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.1.webhook.method", "POST"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccWorkflowConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_workflow.review_pr", "title", "Review PR (updated)"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "category", "platform"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.#", "3"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.1.webhook.url", "https://example.com/hook-updated"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.1.webhook.method", "PUT"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.2.identifier", "record"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.2.upsert_entity.blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "connections.#", "2"),
				),
			},
		},
	})
}

func TestAccPortWorkflowSelfServeTrigger(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	workflowIdentifier := utils.GenID()

	var testAccWorkflowConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_workflow" "deploy" {
		identifier = "%s"
		title      = "Deploy"

		node {
			identifier = "trigger"
			self_serve_trigger {
				action_card_button_text    = "Deploy"
				execute_action_button_text = "Go"

				user_inputs {
					user_properties = {
						string_props = {
							"service" = {
								title    = "Service"
								required = true
							}
						}
						number_props = {
							"replicas" = {
								title = "Replicas"
							}
						}
					}
				}

				permissions {
					roles = ["Member"]
				}
			}
		}

		node {
			identifier = "notify"
			webhook {
				url = "https://example.com/deploy"
			}
		}

		connections {
			source_identifier = "trigger"
			target_identifier = "notify"
		}
	}`, workflowIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccWorkflowConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_workflow.deploy", "identifier", workflowIdentifier),
					resource.TestCheckResourceAttr("port_workflow.deploy", "node.0.self_serve_trigger.action_card_button_text", "Deploy"),
					resource.TestCheckResourceAttr("port_workflow.deploy", "node.0.self_serve_trigger.user_inputs.user_properties.string_props.service.title", "Service"),
					resource.TestCheckResourceAttr("port_workflow.deploy", "node.0.self_serve_trigger.user_inputs.user_properties.number_props.replicas.title", "Replicas"),
					resource.TestCheckResourceAttr("port_workflow.deploy", "node.0.self_serve_trigger.permissions.roles.0", "Member"),
				),
			},
		},
	})
}

func TestAccPortWorkflowConditionOutlets(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	workflowIdentifier := utils.GenID()

	var testAccWorkflowConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_workflow" "branching" {
		identifier = "%s"
		title      = "Branching"

		node {
			identifier = "trigger"
			event_trigger {
				type                 = "ENTITY_CREATED"
				blueprint_identifier = port_blueprint.microservice.identifier
			}
		}

		node {
			identifier = "branch"
			condition {
				outlets {
					identifier = "has_author"
					title      = "Has author"
					expression = ".outputs.trigger.entity.properties.author != null"
					status_label {
						text    = "Author found"
						variant = "success"
					}
				}
			}
		}

		node {
			identifier = "matched"
			webhook {
				url = "https://example.com/matched"
			}
		}

		node {
			identifier = "unmatched"
			webhook {
				url = "https://example.com/unmatched"
			}
		}

		connections {
			source_identifier = "trigger"
			target_identifier = "branch"
		}

		connections {
			source_identifier        = "branch"
			target_identifier        = "matched"
			source_outlet_identifier = "has_author"
		}

		connections {
			source_identifier = "branch"
			target_identifier = "unmatched"
			fallback          = true
		}
	}`, workflowIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccWorkflowConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_workflow.branching", "node.1.condition.outlets.0.identifier", "has_author"),
					resource.TestCheckResourceAttr("port_workflow.branching", "node.1.condition.outlets.0.status_label.text", "Author found"),
					resource.TestCheckResourceAttr("port_workflow.branching", "connections.1.source_outlet_identifier", "has_author"),
					resource.TestCheckResourceAttr("port_workflow.branching", "connections.2.fallback", "true"),
				),
			},
		},
	})
}

func TestAccPortWorkflowImport(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	workflowIdentifier := utils.GenID()

	// event_trigger has no secrets and no server-applied defaults, so every
	// attribute round-trips exactly under ImportStateVerify.
	var testAccWorkflowConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_workflow" "review_pr" {
		identifier = "%s"
		title      = "Review PR"

		node {
			identifier = "trigger"
			title      = "PR Trigger"
			event_trigger {
				type                 = "ENTITY_UPDATED"
				blueprint_identifier = port_blueprint.microservice.identifier
			}
		}
	}`, workflowIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccWorkflowConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_workflow.review_pr", "identifier", workflowIdentifier),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.#", "1"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.0.event_trigger.type", "ENTITY_UPDATED"),
				),
			},
			{
				ResourceName:      "port_workflow.review_pr",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     workflowIdentifier,
			},
		},
	})
}
