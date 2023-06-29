package entity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func EntitySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "Identifier",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "Title",
			Required:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "Icon",
			Optional:            true,
		},
		"run_id": schema.StringAttribute{
			MarkdownDescription: "The runID of the action run that created the entity",
			Optional:            true,
		},
		"teams": schema.ListAttribute{
			MarkdownDescription: "The teams the entity belongs to",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"blueprint": schema.StringAttribute{
			MarkdownDescription: "The blueprint identifier the entity relates to",
			Required:            true,
		},
		"properties": schema.SingleNestedAttribute{
			MarkdownDescription: "The properties of the entity",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"string_prop": schema.MapAttribute{
					MarkdownDescription: "The string properties of the entity",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"number_prop": schema.MapAttribute{
					MarkdownDescription: "The number properties of the entity",
					Optional:            true,
					ElementType:         types.Float64Type,
				},
				"boolean_prop": schema.MapAttribute{
					MarkdownDescription: "The bool properties of the entity",
					Optional:            true,
					ElementType:         types.BoolType,
				},
				"object_prop": schema.MapAttribute{
					MarkdownDescription: "The object properties of the entity",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"array_prop": schema.SingleNestedAttribute{
					MarkdownDescription: "The array properties of the entity",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"string_items": schema.MapAttribute{
							ElementType: types.ListType{ElemType: types.StringType},
							Optional:    true,
						},
						"number_items": schema.MapAttribute{
							ElementType: types.ListType{ElemType: types.Float64Type},
							Optional:    true,
						},
						"boolean_items": schema.MapAttribute{
							ElementType: types.ListType{ElemType: types.BoolType},
							Optional:    true,
						},
						"object_items": schema.MapAttribute{
							ElementType: types.ListType{ElemType: types.StringType},
							Optional:    true,
						},
					},
				},
			},
		},
		"relations": schema.MapAttribute{
			MarkdownDescription: "The relations of the entity",
			Optional:            true,
			ElementType:         types.ListType{ElemType: types.StringType},
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the entity",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The creator of the entity",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"updated_at": schema.StringAttribute{
			MarkdownDescription: "The last update date of the entity",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The last updater of the entity",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}
}

func (r *EntityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Entity resource",
		Attributes:          EntitySchema(),
	}
}
