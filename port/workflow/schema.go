package workflow

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

func nodeExactlyOneOfValidator() []validator.Object {
	return []validator.Object{
		objectvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("event_trigger"),
			path.MatchRelative().AtParent().AtName("cursor_agent"),
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
		"category": schema.StringAttribute{
			MarkdownDescription: "A free-form category used to group the workflow (max 40 characters)",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(40),
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
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						consts.EntityCreated,
						consts.EntityUpdated,
						consts.EntityDeleted,
						consts.AnyEntityChange,
						consts.WorkflowTimerExpired,
					),
				},
			},
			"blueprint_identifier": schema.StringAttribute{
				MarkdownDescription: "The blueprint identifier the event relates to",
				Optional:            true,
			},
			"property_identifier": schema.StringAttribute{
				MarkdownDescription: "The property identifier the timer event relates to (only for the TIMER_EXPIRED event type)",
				Optional:            true,
			},
		},
		Validators: nodeExactlyOneOfValidator(),
	}
}

func cursorAgentBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "A Cursor agent node that runs a Cursor agent task",
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
				MarkdownDescription: "The source the Cursor agent operates on. Provide either repository (with an optional ref) or pr_url.",
				Attributes: map[string]schema.Attribute{
					"repository": schema.StringAttribute{
						MarkdownDescription: "The GitHub repository URL the agent operates on",
						Optional:            true,
					},
					"ref": schema.StringAttribute{
						MarkdownDescription: "The git ref (branch, tag or commit) to use as the base",
						Optional:            true,
					},
					"pr_url": schema.StringAttribute{
						MarkdownDescription: "The pull request URL the agent operates on. When set, repository and ref are ignored.",
						Optional:            true,
					},
				},
			},
		},
	}
}

func integrationActionBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "An integration action node that invokes an integration",
		Attributes: map[string]schema.Attribute{
			"installation_id": schema.StringAttribute{
				MarkdownDescription: "The installation id of the integration. Omit when defer_integration_installation is true.",
				Optional:            true,
			},
			"integration_provider": schema.StringAttribute{
				MarkdownDescription: "The provider of the integration action",
				Optional:            true,
			},
			"integration_invocation_type": schema.StringAttribute{
				MarkdownDescription: "The invocation type of the integration action",
				Optional:            true,
			},
			"execution_properties": schema.StringAttribute{
				MarkdownDescription: "The integration action execution properties as a JSON encoded string",
				Optional:            true,
			},
			"defer_integration_installation": schema.BoolAttribute{
				MarkdownDescription: "When true, allows saving the workflow before the integration is installed",
				Optional:            true,
			},
			"on_failure": schema.StringAttribute{
				MarkdownDescription: "The action to take if the integration action fails",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("continue", "terminate"),
				},
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
			"node":        nodeBlock(),
			"connections": connectionBlock(),
		},
	}
}
