package workflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func WorkflowSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the workflow",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the workflow",
			Required:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the workflow",
			Optional:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the workflow",
			Optional:            true,
		},
		"tags": schema.ListAttribute{
			MarkdownDescription: "Optional list of free-form tags for filtering and grouping workflows. Maximum 20 tags, each up to 64 characters.",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtMost(20),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthBetween(1, 64),
				),
			},
		},
		"allow_anyone_to_view_runs": schema.BoolAttribute{
			MarkdownDescription: "Whether members can view the runs of this workflow",
			Optional:            true,
		},
		"nodes": schema.StringAttribute{
			MarkdownDescription: "Workflow nodes in JSON format, encoded as a string. Use [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode). Supports all workflow node types including `CURSOR_AGENT`, `INTEGRATION_ACTION`, `WEBHOOK`, `SELF_SERVE_TRIGGER`, and `EVENT_TRIGGER`.",
			Required:            true,
		},
		"connections": schema.StringAttribute{
			MarkdownDescription: "Workflow connections in JSON format, encoded as a string. Use [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode). Each connection requires `sourceIdentifier` and `targetIdentifier`.",
			Required:            true,
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the workflow",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the workflow",
			Computed:            true,
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the workflow",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the workflow",
			Computed:            true,
		},
	}
}

func (r *WorkflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: ResourceMarkdownDescription,
		Attributes:          WorkflowSchema(),
	}
}

var ResourceMarkdownDescription = `
# Workflow resource

Docs for Port workflows can be found [here](https://docs.port.io/workflows/overview/).

## Example Usage with CURSOR_AGENT node

` + "```hcl" + `
resource "port_workflow" "review_pr" {
  identifier  = "review-pr"
  title       = "Review PR with Cursor Agent"
  description = "Launch a Cursor cloud agent to review and improve a pull request"
  tags        = ["lorekeeper", "cursor-agent"]

  nodes = jsonencode([
    {
      identifier = "trigger"
      title      = "PR Trigger"
      config = {
        type = "EVENT_TRIGGER"
        event = {
          type                 = "ENTITY_UPDATED"
          blueprintIdentifier  = "githubPullRequest"
        }
      }
    },
    {
      identifier = "decision"
      title      = "Review and Improve PR"
      config = {
        type   = "CURSOR_AGENT"
        apiKey = "{{ .secrets[\"CURSOR_API_KEY\"] }}"
        prompt = {
          text = "Review this PR for code quality issues and apply fixes. Add any missing error handling."
        }
        source = {
          prUrl = "{{ .outputs.trigger.diff.after.properties.link }}"
        }
      }
    }
  ])

  connections = jsonencode([
    {
      sourceIdentifier = "trigger"
      targetIdentifier = "decision"
    }
  ])
}
` + "\n```" + `

## Example Usage with INTEGRATION_ACTION node

` + "```hcl" + `
resource "port_workflow" "dispatch_github_workflow" {
  identifier = "dispatch-github-workflow"
  title      = "Dispatch GitHub Workflow"

  nodes = jsonencode([
    {
      identifier = "trigger"
      title      = "Start"
      config = {
        type = "SELF_SERVE_TRIGGER"
        userInputs = {
          properties = {
            pr_link = {
              type  = "string"
              title = "PR Link"
            }
          }
          required = ["pr_link"]
        }
      }
    },
    {
      identifier = "fetch_pr_context"
      title      = "Fetch PR Context"
      config = {
        type                      = "INTEGRATION_ACTION"
        installationId            = "your-github-installation-id"
        integrationProvider       = "github-ocean"
        integrationInvocationType = "dispatch_workflow"
        integrationActionExecutionProperties = {
          org      = "port-labs"
          repo     = "Port"
          workflow = "hackathon_lore_keepres.yml"
          workflowInputs = {
            node_type       = "fetch_pr_context"
            pr_link         = "{{ .outputs.trigger.pr_link }}"
            workflow_run_id = "{{ .workflowRun.identifier }}"
          }
          reportWorkflowStatus = true
        }
      }
    }
  ])

  connections = jsonencode([
    {
      sourceIdentifier = "trigger"
      targetIdentifier = "fetch_pr_context"
    }
  ])
}
` + "\n```" + `
`
