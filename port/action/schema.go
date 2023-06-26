package action

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ActionSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"ID": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "Identifier",
			Required:            true,
		},
		"blueprint": schema.StringAttribute{
			MarkdownDescription: "The blueprint identifier the action relates to",
			Required:            true,
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "Title",
			Required:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "Icon",
			Optional:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "Description",
			Optional:            true,
		},
		"required_approval": schema.BoolAttribute{
			MarkdownDescription: "Require approval before invoking the action",
			Optional:            true,
		},
		"kafka_method": schema.SingleNestedAttribute{
			MarkdownDescription: "The invocation method of the action",
			Optional:            true,
			Validators: []validator.Object{
				objectvalidator.ConflictsWith(path.Expressions{
					path.MatchRoot("webhook_method"),
				}...),
			},
		},
		"webhook_method": schema.SingleNestedAttribute{
			MarkdownDescription: "The invocation method of the action",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"url": schema.StringAttribute{
					MarkdownDescription: "Required when selecting type WEBHOOK. The URL to invoke the action",
					Required:            true,
				},
				"agent": schema.BoolAttribute{
					MarkdownDescription: "Use the agent to invoke the action",
					Optional:            true,
				},
			},
		},
		"github_method": schema.SingleNestedAttribute{
			MarkdownDescription: "The invocation method of the action",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"org": schema.StringAttribute{
					MarkdownDescription: "Required when selecting type GITHUB. The GitHub org that the workflow belongs to",
					Required:            true,
				},
				"repo": schema.StringAttribute{
					MarkdownDescription: "Required when selecting type GITHUB. The GitHub repo that the workflow belongs to",
					Required:            true,
				},
				"workflow": schema.StringAttribute{
					MarkdownDescription: "The GitHub workflow that the action belongs to",
					Optional:            true,
				},
				"omit_payload": schema.BoolAttribute{
					MarkdownDescription: "Omit the payload when invoking the action",
					Optional:            true,
				},
				"omit_user_inputs": schema.BoolAttribute{
					MarkdownDescription: "Omit the user inputs when invoking the action",
					Optional:            true,
				},
				"report_workflow_status": schema.BoolAttribute{
					MarkdownDescription: "Report the workflow status when invoking the action",
					Optional:            true,
				},
			},
		},
		"azure_method": schema.SingleNestedAttribute{
			MarkdownDescription: "The invocation method of the action",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"org": schema.StringAttribute{
					MarkdownDescription: "Required when selecting type AZURE. The Azure org that the workflow belongs to",
					Required:            true,
				},
				"webhook": schema.StringAttribute{
					MarkdownDescription: "Required when selecting type AZURE. The Azure webhook that the workflow belongs to",
					Required:            true,
				},
			},
		},
		"user_properties": schema.SingleNestedAttribute{
			MarkdownDescription: "User properties",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"identifier": schema.StringAttribute{
					MarkdownDescription: "Identifier",
					Required:            true,
				},
				"title": schema.StringAttribute{
					MarkdownDescription: "Title",
					Required:            true,
				},
				"type": schema.StringAttribute{
					MarkdownDescription: "Type",
					Required:            true,
				},
				"description": schema.StringAttribute{
					MarkdownDescription: "Description",
					Optional:            true,
				},
				"default": schema.StringAttribute{
					MarkdownDescription: "Default",
					Optional:            true,
				},
				"format": schema.StringAttribute{
					MarkdownDescription: "Format",
					Optional:            true,
				},
				"blueprint": schema.StringAttribute{
					MarkdownDescription: "Blueprint",
					Optional:            true,
				},
				"required": schema.BoolAttribute{
					MarkdownDescription: "Required",
					Optional:            true,
				},
			},
		},
	}

}

func (r *ActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Entity resource",
		Attributes:          ActionSchema(),
	}
}
