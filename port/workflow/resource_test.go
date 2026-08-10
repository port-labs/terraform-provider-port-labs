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
		identifier = "%s"
		title      = "Review PR"
		category   = "engineering"

		node {
			identifier = "trigger"
			title      = "PR Trigger"
			event_trigger {
				type                 = "ENTITY_UPDATED"
				blueprint_identifier = port_blueprint.microservice.identifier
			}
		}

		node {
			identifier = "decision"
			title      = "Review"
			cursor_agent {
				api_key = "{{ .secrets[\"CURSOR_API_KEY\"] }}"
				prompt {
					text = "Review this PR."
				}
				source {
					pr_url = "{{ .outputs.trigger.diff.after.properties.link }}"
				}
			}
		}

		connections {
			source_identifier = "trigger"
			target_identifier = "decision"
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
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.0.identifier", "trigger"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.0.title", "PR Trigger"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.0.event_trigger.type", "ENTITY_UPDATED"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.0.event_trigger.blueprint_identifier", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.1.identifier", "decision"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.1.cursor_agent.prompt.text", "Review this PR."),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.1.cursor_agent.source.pr_url", "{{ .outputs.trigger.diff.after.properties.link }}"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "connections.0.source_identifier", "trigger"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "connections.0.target_identifier", "decision"),
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
			identifier = "decision"
			cursor_agent {
				api_key = "{{ .secrets[\"CURSOR_API_KEY\"] }}"
				prompt {
					text = "Review this PR."
				}
				source {
					pr_url = "{{ .outputs.trigger.diff.after.properties.link }}"
				}
			}
		}

		connections {
			source_identifier = "trigger"
			target_identifier = "decision"
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
			identifier = "decision"
			cursor_agent {
				api_key = "{{ .secrets[\"CURSOR_API_KEY\"] }}"
				prompt {
					text = "Review this PR carefully."
				}
				source {
					pr_url = "{{ .outputs.trigger.diff.after.properties.link }}"
				}
			}
		}

		node {
			identifier = "summarize"
			cursor_agent {
				api_key = "{{ .secrets[\"CURSOR_API_KEY\"] }}"
				prompt {
					text = "Summarize the review."
				}
				source {
					pr_url = "{{ .outputs.trigger.diff.after.properties.link }}"
				}
			}
		}

		connections {
			source_identifier = "trigger"
			target_identifier = "decision"
		}

		connections {
			source_identifier = "decision"
			target_identifier = "summarize"
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
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.1.identifier", "decision"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.1.cursor_agent.prompt.text", "Review this PR."),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccWorkflowConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_workflow.review_pr", "title", "Review PR (updated)"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "category", "platform"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.#", "3"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.1.cursor_agent.prompt.text", "Review this PR carefully."),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.2.identifier", "summarize"),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "node.2.cursor_agent.prompt.text", "Summarize the review."),
					resource.TestCheckResourceAttr("port_workflow.review_pr", "connections.#", "2"),
				),
			},
		},
	})
}

func TestAccPortWorkflowImport(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	workflowIdentifier := utils.GenID()

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
