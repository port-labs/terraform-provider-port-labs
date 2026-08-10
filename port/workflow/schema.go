package workflow

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

func secondsValidators() []validator.String {
	return []validator.String{
		stringvalidator.Any(
			// A positive integer up to 86400 (1 day), or a dynamic Port expression.
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^(?:[1-9][0-9]*)$`),
				"must be a positive integer or a dynamic expression ({{ ... }})",
			),
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^[\n\r\s]*{{.*}}[\n\r\s]*$`),
				"must be a positive integer or a dynamic expression ({{ ... }})",
			),
		),
	}
}

func nodeExactlyOneOfValidator() []validator.Object {
	return []validator.Object{
		objectvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("event_trigger"),
			path.MatchRelative().AtParent().AtName("cursor_agent"),
			path.MatchRelative().AtParent().AtName("delay"),
			path.MatchRelative().AtParent().AtName("integration_action"),
		),
	}
}

func WorkflowSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the workflow (computed)",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the workflow",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the workflow",
			Optional:            true,
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
			MarkdownDescription: "Free-form tags for filtering and grouping workflows (max 20 tags, each 1-64 characters)",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtMost(20),
				listvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 64)),
			},
		},
	}
}

func eventTriggerBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "An event trigger node that starts the workflow when an entity event occurs",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				MarkdownDescription: "The event type that triggers the workflow",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						consts.EntityCreated,
						consts.EntityUpdated,
						consts.EntityDeleted,
						consts.AnyEntityChange,
						consts.TimerPropertyExpired,
					),
				},
			},
			"blueprint_identifier": schema.StringAttribute{
				MarkdownDescription: "The blueprint identifier the event relates to",
				Optional:            true,
			},
		},
		Validators: nodeExactlyOneOfValidator(),
	}
}

func cursorAgentBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "A Cursor agent node that runs a Cursor agent decision",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The Cursor API key (supports Port secret references)",
				Optional:            true,
				Sensitive:           true,
			},
		},
		Blocks: map[string]schema.Block{
			"prompt": schema.SingleNestedBlock{
				MarkdownDescription: "The prompt passed to the Cursor agent",
				Attributes: map[string]schema.Attribute{
					"text": schema.StringAttribute{
						MarkdownDescription: "The prompt text",
						Optional:            true,
					},
				},
			},
			"source": schema.SingleNestedBlock{
				MarkdownDescription: "The source the Cursor agent operates on",
				Attributes: map[string]schema.Attribute{
					"pr_url": schema.StringAttribute{
						MarkdownDescription: "The pull request URL the agent operates on",
						Optional:            true,
					},
				},
			},
		},
	}
}

func delayBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "A delay node that pauses the workflow for a number of seconds",
		Attributes: map[string]schema.Attribute{
			"seconds": schema.StringAttribute{
				MarkdownDescription: "The number of seconds to wait (1-86400) or a dynamic expression",
				Optional:            true,
				Validators:          secondsValidators(),
			},
		},
	}
}

func integrationActionBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "An integration action node that triggers an integration workflow",
		Attributes: map[string]schema.Attribute{
			"installation_id": schema.StringAttribute{
				MarkdownDescription: "The installation id of the integration",
				Optional:            true,
			},
			"integration_action_type": schema.StringAttribute{
				MarkdownDescription: "The integration action type",
				Optional:            true,
			},
			"org": schema.StringAttribute{
				MarkdownDescription: "The org the workflow belongs to",
				Optional:            true,
			},
			"repo": schema.StringAttribute{
				MarkdownDescription: "The repo the workflow belongs to",
				Optional:            true,
			},
			"workflow": schema.StringAttribute{
				MarkdownDescription: "The workflow to run",
				Optional:            true,
			},
			"workflow_inputs": schema.StringAttribute{
				MarkdownDescription: "The workflow inputs as a JSON encoded string",
				Optional:            true,
			},
			"report_workflow_status": schema.StringAttribute{
				MarkdownDescription: "Whether to report the workflow status back to Port",
				Optional:            true,
			},
		},
	}
}

func nodeBlock() schema.Block {
	return schema.ListNestedBlock{
		MarkdownDescription: "A node of the workflow graph",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"identifier": schema.StringAttribute{
					MarkdownDescription: "The identifier of the node",
					Required:            true,
				},
				"title": schema.StringAttribute{
					MarkdownDescription: "The title of the node",
					Optional:            true,
				},
			},
			Blocks: map[string]schema.Block{
				"event_trigger":      eventTriggerBlock(),
				"cursor_agent":       cursorAgentBlock(),
				"delay":              delayBlock(),
				"integration_action": integrationActionBlock(),
			},
		},
	}
}

func connectionBlock() schema.Block {
	return schema.ListNestedBlock{
		MarkdownDescription: "A directed connection between two nodes",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"source_identifier": schema.StringAttribute{
					MarkdownDescription: "The identifier of the source node",
					Required:            true,
				},
				"target_identifier": schema.StringAttribute{
					MarkdownDescription: "The identifier of the target node",
					Required:            true,
				},
			},
		},
	}
}

func (r *WorkflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Workflow resource for managing Port workflows and their node/connection graph",
		Attributes:          WorkflowSchema(),
		Blocks: map[string]schema.Block{
			"node": nodeBlock(),
			// "connection" is a reserved Terraform block name, so the plural
			// "connections" block is used for the directed edges of the graph.
			"connections": connectionBlock(),
		},
	}
}
