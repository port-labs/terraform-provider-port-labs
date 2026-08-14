package entity

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func EntitySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier of the entity",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the entity",
			Required:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the entity",
			Optional:            true,
		},
		"run_id": schema.StringAttribute{
			MarkdownDescription: "The runID of the action run that created the entity",
			Optional:            true,
		},
		"create_missing_related_entities": schema.BoolAttribute{
			MarkdownDescription: "Whether to create missing related entities",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"teams": schema.SetAttribute{
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
				"string_props": schema.MapAttribute{
					MarkdownDescription: "The string properties of the entity",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"number_props": schema.MapAttribute{
					MarkdownDescription: "The number properties of the entity",
					Optional:            true,
					ElementType:         types.Float64Type,
				},
				"boolean_props": schema.MapAttribute{
					MarkdownDescription: "The bool properties of the entity",
					Optional:            true,
					ElementType:         types.BoolType,
				},
				"object_props": schema.MapAttribute{
					MarkdownDescription: "The object properties of the entity",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"array_props": schema.SingleNestedAttribute{
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
						"union_string_slices": schema.MapNestedAttribute{
							MarkdownDescription: "Write-only slices for union array string properties. Each entry names a property and the source key whose array value should be written. On read, assembled values are returned in `string_items` for non-union properties only.",
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"source_key": schema.StringAttribute{
										MarkdownDescription: "Stable source key for this writer's slice. Must match `^[A-Za-z0-9._:@/-]{1,128}$`.",
										Required:            true,
									},
									"items": schema.ListAttribute{
										MarkdownDescription: "The string array values for this source slice.",
										Optional:            true,
										ElementType:         types.StringType,
									},
								},
							},
						},
						"union_number_slices": schema.MapNestedAttribute{
							MarkdownDescription: "Write-only slices for union array number properties. Each entry names a property and the source key whose array value should be written. On read, assembled values are returned in `number_items` for non-union properties only.",
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"source_key": schema.StringAttribute{
										MarkdownDescription: "Stable source key for this writer's slice. Must match `^[A-Za-z0-9._:@/-]{1,128}$`.",
										Required:            true,
									},
									"items": schema.ListAttribute{
										MarkdownDescription: "The number array values for this source slice.",
										Optional:            true,
										ElementType:         types.Float64Type,
									},
								},
							},
						},
					},
				},
			},
		},
		"relations": schema.SingleNestedAttribute{
			MarkdownDescription: "The relations of the entity",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"single_relations": schema.MapAttribute{
					MarkdownDescription: "The single relation of the entity",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"many_relations": schema.MapAttribute{
					MarkdownDescription: "The many relation of the entity",
					Optional:            true,
					ElementType:         types.ListType{ElemType: types.StringType},
				},
			},
		},
		"created_at": schema.StringAttribute{
			MarkdownDescription: "The creation date of the entity",
			Computed:            true,
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
		},
	}
}

func (r *EntityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Entity resource",
		Attributes:          EntitySchema(),
	}
}
