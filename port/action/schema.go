package action

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func MetadataProperties() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the property",
			Optional:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the property",
			Optional:            true,
		},
		"required": schema.BoolAttribute{
			MarkdownDescription: "Whether the property is required, by default not required, this property can't be set at the same time if `required_jq_query` is set, and only supports true as value",
			Optional:            true,
			Validators: []validator.Bool{
				boolvalidator.ConflictsWith(path.MatchRoot("required_jq_query")),
			},
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the property",
			Optional:            true,
		},
		"depends_on": schema.ListAttribute{
			MarkdownDescription: "The properties that this property depends on",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"dataset": schema.SingleNestedAttribute{
			MarkdownDescription: "The dataset of the property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"combinator": schema.StringAttribute{
					MarkdownDescription: "The combinator of the dataset",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("and", "or"),
					},
				},
				"rules": schema.ListNestedAttribute{
					MarkdownDescription: "The rules of the dataset",
					Required:            true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"blueprint": schema.StringAttribute{
								MarkdownDescription: "The blueprint identifier of the rule",
								Optional:            true,
							},
							"property": schema.StringAttribute{
								MarkdownDescription: "The property identifier of the rule",
								Optional:            true,
							},
							"operator": schema.StringAttribute{
								MarkdownDescription: "The operator of the rule",
								Required:            true,
							},
							"value": schema.ObjectAttribute{
								MarkdownDescription: "The value of the rule",
								Required:            true,
								AttributeTypes: map[string]attr.Type{
									"jq_query": types.StringType,
								},
							},
						},
					},
				},
			},
		},
	}
}

func ActionSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
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
		"approval_webhook_notification": schema.SingleNestedAttribute{
			MarkdownDescription: "The webhook notification of the approval",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"url": schema.StringAttribute{
					MarkdownDescription: "The URL to invoke the webhook",
					Required:            true,
				},
			},
		},
		"approval_email_notification": schema.ObjectAttribute{
			MarkdownDescription: "The email notification of the approval",
			Optional:            true,
			AttributeTypes:      map[string]attr.Type{},
		},
		"trigger": schema.StringAttribute{
			MarkdownDescription: "The trigger type of the action",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("CREATE", "DAY-2", "DELETE"),
			},
		},
		"kafka_method": schema.ObjectAttribute{
			MarkdownDescription: "The invocation method of the action",
			Optional:            true,
			AttributeTypes:      map[string]attr.Type{},
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
				"synchronized": schema.BoolAttribute{
					MarkdownDescription: "Synchronize the action",
					Optional:            true,
				},
				"method": schema.StringAttribute{
					MarkdownDescription: "The HTTP method to invoke the action",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("POST", "PUT", "PATCH", "DELETE"),
					},
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
					Required:            true,
				},
				"omit_payload": schema.BoolAttribute{
					MarkdownDescription: "Omit the payload when invoking the action",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(false),
				},
				"omit_user_inputs": schema.BoolAttribute{
					MarkdownDescription: "Omit the user inputs when invoking the action",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(false),
				},
				"report_workflow_status": schema.BoolAttribute{
					MarkdownDescription: "Report the workflow status when invoking the action",
					Optional:            true,
				},
			},
		},
		"gitlab_method": schema.SingleNestedAttribute{
			MarkdownDescription: "The invocation method of the action",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"project_name": schema.StringAttribute{
					MarkdownDescription: "Required when selecting type GITLAB. The GitLab project name that the workflow belongs to",
					Required:            true,
				},
				"group_name": schema.StringAttribute{
					MarkdownDescription: "Required when selecting type GITLAB. The GitLab group name that the workflow belongs to",
					Required:            true,
				},
				"omit_payload": schema.BoolAttribute{
					MarkdownDescription: "Omit the payload when invoking the action",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(false),
				},
				"omit_user_inputs": schema.BoolAttribute{
					MarkdownDescription: "Omit the user inputs when invoking the action",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(false),
				},
				"default_ref": schema.StringAttribute{
					MarkdownDescription: "The default ref of the action",
					Optional:            true,
				},
				"agent": schema.BoolAttribute{
					MarkdownDescription: "Use the agent to invoke the action",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(true),
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
		"order_properties": schema.ListAttribute{
			MarkdownDescription: "Order properties",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"user_properties": schema.SingleNestedAttribute{
			MarkdownDescription: "User properties",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"string_props":  StringPropertySchema(),
				"number_props":  NumberPropertySchema(),
				"boolean_props": BooleanPropertySchema(),
				"object_props":  ObjectPropertySchema(),
				"array_props":   ArrayPropertySchema(),
			},
		},
		"required_jq_query": schema.StringAttribute{
			MarkdownDescription: "The required jq query of the property",
			Optional:            true,
		},
	}

}

func StringPropertySchema() schema.Attribute {
	stringPropertySchema := map[string]schema.Attribute{
		"default": schema.StringAttribute{
			MarkdownDescription: "The default of the string property",
			Optional:            true,
		},
		"default_jq_query": schema.StringAttribute{
			MarkdownDescription: "The default jq query of the string property",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("default")),
			},
		},
		"blueprint": schema.StringAttribute{
			MarkdownDescription: "The blueprint identifier the string property relates to",
			Optional:            true,
		},
		"format": schema.StringAttribute{
			MarkdownDescription: "The format of the string property",
			Optional:            true,
		},
		"min_length": schema.Int64Attribute{
			MarkdownDescription: "The min length of the string property",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"max_length": schema.Int64Attribute{
			MarkdownDescription: "The max length of the string property",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"pattern": schema.StringAttribute{
			MarkdownDescription: "The pattern of the string property",
			Optional:            true,
		},
		"enum": schema.ListAttribute{
			MarkdownDescription: "The enum of the string property",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.UniqueValues(),
				listvalidator.SizeAtLeast(1),
			},
		},
		"enum_jq_query": schema.StringAttribute{
			MarkdownDescription: "The enum jq query of the string property",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("enum")),
			},
		},
		"encryption": schema.StringAttribute{
			MarkdownDescription: "The algorithm to encrypt the property with",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("aes256-gcm"),
			},
		},
		"visible": schema.BoolAttribute{
			MarkdownDescription: "The visibility of the string property",
			Optional:            true,
		},
		"visible_jq_query": schema.StringAttribute{
			MarkdownDescription: "The visibility condition jq query of the string property",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("visible")),
			},
		},
	}

	utils.CopyMaps(stringPropertySchema, MetadataProperties())
	return schema.MapNestedAttribute{
		MarkdownDescription: "The string property of the action",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: stringPropertySchema,
		},
	}
}

func NumberPropertySchema() schema.Attribute {
	numberPropertySchema := map[string]schema.Attribute{
		"default": schema.Float64Attribute{
			MarkdownDescription: "The default of the number property",
			Optional:            true,
		},
		"default_jq_query": schema.StringAttribute{
			MarkdownDescription: "The default jq query of the number property",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("default")),
			},
		},
		"maximum": schema.Float64Attribute{
			MarkdownDescription: "The min of the number property",
			Optional:            true,
		},
		"minimum": schema.Float64Attribute{
			MarkdownDescription: "The max of the number property",
			Optional:            true,
		},
		"enum": schema.ListAttribute{
			MarkdownDescription: "The enum of the number property",
			Optional:            true,
			ElementType:         types.Float64Type,
			Validators: []validator.List{
				listvalidator.UniqueValues(),
				listvalidator.SizeAtLeast(1),
			},
		},
		"enum_jq_query": schema.StringAttribute{
			MarkdownDescription: "The enum jq query of the string property",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("enum")),
			},
		},
		"visible": schema.BoolAttribute{
			MarkdownDescription: "The visibility of the number property",
			Optional:            true,
		},
		"visible_jq_query": schema.StringAttribute{
			MarkdownDescription: "The visibility condition jq query of the number property",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("visible")),
			},
		},
	}

	utils.CopyMaps(numberPropertySchema, MetadataProperties())
	return schema.MapNestedAttribute{
		MarkdownDescription: "The number property of the action",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: numberPropertySchema,
		},
	}
}

func BooleanPropertySchema() schema.Attribute {
	booleanPropertySchema := map[string]schema.Attribute{
		"default": schema.BoolAttribute{
			MarkdownDescription: "The default of the boolean property",
			Optional:            true,
		},
		"default_jq_query": schema.StringAttribute{
			MarkdownDescription: "The default jq query of the boolean property",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("default")),
			},
		},
		"visible": schema.BoolAttribute{
			MarkdownDescription: "The visibility of the boolean property",
			Optional:            true,
		},
		"visible_jq_query": schema.StringAttribute{
			MarkdownDescription: "The visibility condition jq query of the boolean property",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("visible")),
			},
		},
	}

	utils.CopyMaps(booleanPropertySchema, MetadataProperties())
	return schema.MapNestedAttribute{
		MarkdownDescription: "The boolean property of the action",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: booleanPropertySchema,
		},
	}
}

func ObjectPropertySchema() schema.Attribute {
	objectPropertySchema := map[string]schema.Attribute{
		"default": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The default of the object property",
		},
		"default_jq_query": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The default jq query of the object property",
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("default")),
			},
		},
		"encryption": schema.StringAttribute{
			MarkdownDescription: "The algorithm to encrypt the property with",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("aes256-gcm"),
			},
		},
		"visible": schema.BoolAttribute{
			MarkdownDescription: "The visibility of the object property",
			Optional:            true,
		},
		"visible_jq_query": schema.StringAttribute{
			MarkdownDescription: "The visibility condition jq query of the object property",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("visible")),
			},
		},
	}
	utils.CopyMaps(objectPropertySchema, MetadataProperties())
	return schema.MapNestedAttribute{
		MarkdownDescription: "The object property of the action",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: objectPropertySchema,
		},
	}
}

func ArrayPropertySchema() schema.Attribute {
	arrayPropertySchema := map[string]schema.Attribute{
		"min_items": schema.Int64Attribute{
			MarkdownDescription: "The min items of the array property",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"max_items": schema.Int64Attribute{
			MarkdownDescription: "The max items of the array property",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"default_jq_query": schema.StringAttribute{
			MarkdownDescription: "The default jq query of the array property",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("string_items").AtName("default")),
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("number_items").AtName("default")),
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("boolean_items").AtName("default")),
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("object_items").AtName("default")),
			},
		},
		"string_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The items of the array property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"format": schema.StringAttribute{
					MarkdownDescription: "The format of the items",
					Optional:            true,
				},
				"blueprint": schema.StringAttribute{
					MarkdownDescription: "The blueprint identifier the property relates to",
					Optional:            true,
				},
				"default": schema.ListAttribute{
					MarkdownDescription: "The default of the items",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"enum": schema.ListAttribute{
					MarkdownDescription: "The enum of the items",
					Optional:            true,
					ElementType:         types.StringType,
					Validators: []validator.List{
						listvalidator.UniqueValues(),
						listvalidator.SizeAtLeast(1),
					},
				},
				"enum_jq_query": schema.StringAttribute{
					MarkdownDescription: "The enum jq query of the string items",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("enum")),
					},
				},
			},
		},
		"number_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The items of the array property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"default": schema.ListAttribute{
					MarkdownDescription: "The default of the items",
					Optional:            true,
					ElementType:         types.Float64Type,
				},
				"enum": schema.ListAttribute{
					MarkdownDescription: "The enum of the items",
					Optional:            true,
					ElementType:         types.Float64Type,
					Validators: []validator.List{
						listvalidator.UniqueValues(),
						listvalidator.SizeAtLeast(1),
					},
				},
				"enum_jq_query": schema.StringAttribute{
					MarkdownDescription: "The enum jq query of the number items",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("enum")),
					},
				},
			},
		},
		"boolean_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The items of the array property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"default": schema.ListAttribute{
					MarkdownDescription: "The default of the items",
					Optional:            true,
					ElementType:         types.BoolType,
				},
			},
		},
		"object_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The items of the array property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"default": schema.ListAttribute{
					MarkdownDescription: "The default of the items",
					Optional:            true,
					ElementType:         types.MapType{ElemType: types.StringType},
				},
			},
		},
		"visible": schema.BoolAttribute{
			MarkdownDescription: "The visibility of the array property",
			Optional:            true,
		},
		"visible_jq_query": schema.StringAttribute{
			MarkdownDescription: "The visibility condition jq query of the array property",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("visible")),
			},
		},
	}

	utils.CopyMaps(arrayPropertySchema, MetadataProperties())
	return schema.MapNestedAttribute{
		MarkdownDescription: "The array property of the action",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: arrayPropertySchema,
		},
	}
}

func (r *ActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Action resource",
		Attributes:          ActionSchema(),
	}
}

func (r *ActionResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var state *ActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	validateUserInputRequiredNotSetToFalse(state, resp)
}

func validateUserInputRequiredNotSetToFalse(state *ActionModel, resp *resource.ValidateConfigResponse) {
	// go over all the properties and check if required is set to false, is its false, raise an error that false is not
	// supported anymore
	const errorString = "required is set to false, this is not supported anymore, if you don't want to make the stringProp required, remove the required stringProp"

	if state.UserProperties == nil {
		return
	}

	for _, stringProp := range state.UserProperties.StringProps {
		if !stringProp.Required.IsNull() && !stringProp.Required.ValueBool() && stringProp.Required == types.BoolValue(false) {
			resp.Diagnostics.AddError(errorString, fmt.Sprint(`Error in User Property: `, stringProp.Title, ` in action: `, state.Identifier))
		}
	}

	for _, numberProp := range state.UserProperties.NumberProps {
		if !numberProp.Required.IsNull() && !numberProp.Required.ValueBool() && numberProp.Required == types.BoolValue(false) {
			resp.Diagnostics.AddError(errorString, fmt.Sprint(`Error in User Property: `, numberProp.Title, ` in action: `, state.Identifier))
		}
	}

	for _, boolProp := range state.UserProperties.BooleanProps {
		if !boolProp.Required.IsNull() && !boolProp.Required.ValueBool() && boolProp.Required == types.BoolValue(false) {
			resp.Diagnostics.AddError(errorString, fmt.Sprint(`Error in User Property: `, boolProp.Title, ` in action: `, state.Identifier))
		}
	}

	for _, objectProp := range state.UserProperties.ObjectProps {
		if !objectProp.Required.IsNull() && !objectProp.Required.ValueBool() && objectProp.Required == types.BoolValue(false) {
			resp.Diagnostics.AddError(errorString, fmt.Sprint(`Error in User Property: `, objectProp.Title, ` in action: `, state.Identifier))
		}
	}

	for _, arrayProp := range state.UserProperties.ArrayProps {
		if !arrayProp.Required.IsNull() && !arrayProp.Required.ValueBool() && arrayProp.Required == types.BoolValue(false) {
			resp.Diagnostics.AddError(errorString, fmt.Sprint(`Error in User Property: `, arrayProp.Title, ` in action: `, state.Identifier))
		}
	}
}
