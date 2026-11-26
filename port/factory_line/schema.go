package factory_line

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var executorValues = []string{
	"claude_managed_agent",
	"devin_managed_agent",
	"factory_ai_managed_agent",
	"custom_workflow",
	"human",
	"none",
}

var factoryAIReasoningEffort = []string{
	"none", "dynamic", "off", "minimal", "low", "medium", "high", "xhigh", "max",
}

func (r *FactoryLineResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Port software factory line (`v1/factory/lines`). Each `stage` declares an `executor` " +
			"and a matching executor settings block. Use `factory_ai_managed_agent` for Factory.ai Droid stages " +
			"(requires exactly one of `environment_id` or `computer_id`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Line identifier (same as `identifier`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identifier": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique line identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"title": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Line title.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Line description.",
			},
			"team": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Team identifiers associated with the line.",
				ElementType:         types.StringType,
			},
			"max_running": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Maximum concurrent running work items (1–1000).",
				Validators:          []validator.Int64{int64validator.Between(1, 1000)},
			},
			"status": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Line status. One of `active`, `paused`.",
				Default:             stringdefault.StaticString("active"),
				Validators: []validator.String{
					stringvalidator.OneOf("active", "paused"),
				},
			},
			"status_reason": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Reason for the current status.",
			},
			"stage_order": schema.ListAttribute{
				Computed:            true,
				MarkdownDescription: "Ordered stage identifiers returned by the API.",
				ElementType:         types.StringType,
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Creation timestamp.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last update timestamp.",
			},
		},
		Blocks: map[string]schema.Block{
			"stage": stageBlock(),
		},
	}
}

func stageBlock() schema.Block {
	return schema.ListNestedBlock{
		MarkdownDescription: "Ordered stages on the line (max 30).",
		Validators: []validator.List{
			listvalidator.SizeAtMost(30),
		},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"identifier": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Stage identifier.",
				},
				"title": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Stage title.",
				},
				"description": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Stage description.",
				},
				"instructions": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Instructions passed to the stage executor.",
				},
				"team": schema.ListAttribute{
					Optional:            true,
					MarkdownDescription: "Team identifiers for this stage.",
					ElementType:         types.StringType,
				},
				"executor": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Stage executor type.",
					Validators: []validator.String{
						stringvalidator.OneOf(executorValues...),
					},
				},
				"max_tries": schema.Int64Attribute{
					Optional:            true,
					Computed:            true,
					MarkdownDescription: "Maximum attempts per stage (1–10).",
					Default:             int64default.StaticInt64(3),
					Validators:          []validator.Int64{int64validator.Between(1, 10)},
				},
				"allow_take_over": schema.BoolAttribute{
					Optional:            true,
					Computed:            true,
					MarkdownDescription: "Whether operators can take over a running stage.",
					Default:             booldefault.StaticBool(true),
				},
				"timeout": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Stage timeout duration (for example `2h`, `30m`, `1d`).",
				},
				"checkpoints": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Stage checkpoints as a JSON array string. Use [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode). Defaults to `[]`.",
				},
			},
			Blocks: map[string]schema.Block{
				"factory_ai_managed_agent": factoryAISettingsBlock(),
				"claude_managed_agent":     claudeSettingsBlock(),
				"devin_managed_agent":      devinSettingsBlock(),
				"custom_workflow":          customWorkflowSettingsBlock(),
				"human":                    humanSettingsBlock(),
			},
		},
	}
}

func factoryAISettingsBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "Settings for executor `factory_ai_managed_agent`. Set exactly one of `environment_id` or `computer_id`.",
		Attributes: map[string]schema.Attribute{
			"installation_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Factory.ai integration installation ID. Optional when the org has a single installation.",
			},
			"environment_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Factory environment where the Droid runs (new session).",
			},
			"computer_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Existing Factory computer where the Droid runs (new session).",
			},
			"model": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Model for a new Droid session.",
			},
			"reasoning_effort": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Reasoning effort for a new session.",
				Validators: []validator.String{
					stringvalidator.OneOf(factoryAIReasoningEffort...),
				},
			},
		},
	}
}

func claudeSettingsBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "Settings for executor `claude_managed_agent`.",
		Attributes: map[string]schema.Attribute{
			"agent_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Claude agent identifier.",
			},
			"environment_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Claude environment identifier.",
			},
			"installation_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Claude managed agents installation ID.",
			},
			"session": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Session options as a JSON object string.",
			},
			"git_access_token": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Organization secret name for git access.",
			},
			"execution_identity": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Execution identity as a JSON object (`type` = `service_account`).",
			},
		},
	}
}

func devinSettingsBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "Settings for executor `devin_managed_agent`.",
		Attributes: map[string]schema.Attribute{
			"installation_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Devin installation ID.",
			},
			"playbook_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Devin playbook identifier.",
			},
			"repos": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Repository identifiers.",
				ElementType:         types.StringType,
			},
			"session": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Session options as a JSON object string.",
			},
			"execution_identity": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Execution identity as a JSON object (`type` = `service_account`).",
			},
		},
	}
}

func customWorkflowSettingsBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "Settings for executor `custom_workflow`.",
		Attributes: map[string]schema.Attribute{
			"workflow": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Workflow identifier to run.",
			},
		},
	}
}

func humanSettingsBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "Settings for executor `human`.",
		Attributes: map[string]schema.Attribute{
			"assignees": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Assignees as a JSON object (`users`, `teams`, `roles`, `usersQuery`).",
			},
			"notify": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Notification options as a JSON object (`slack`, `email`).",
			},
		},
	}
}
