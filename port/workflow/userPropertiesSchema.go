package workflow

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

var (
	stringPropertyFormats = []string{
		"multi-line", "date-time", "email", "entity", "team", "user", "url", "markdown", "yaml", "proto",
	}
	stringItemFormats = []string{"entity", "team", "user", "date-time", "email", "url", "yaml"}
	objectFormats     = []string{"labeled-url"}
	objectItemFormats = []string{"labeled-url"}
)

func userPropertiesSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "The user inputs the form collects.",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"string_props":  stringPropertySchema(),
			"number_props":  numberPropertySchema(),
			"boolean_props": booleanPropertySchema(),
			"object_props":  objectPropertySchema(),
			"array_props":   arrayPropertySchema(),
		},
	}
}

func propertyMetadataSchema(propertyType string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the " + propertyType + " input.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(140),
			},
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the " + propertyType + " input.",
			Optional:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the " + propertyType + " input.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(1000),
			},
		},
		"required": schema.BoolAttribute{
			MarkdownDescription: "Whether the input has to be filled in. Only `true` is accepted, and it cannot be combined with `required_jq_query`.",
			Optional:            true,
			Validators: []validator.Bool{
				isTrueValidator{},
				boolvalidator.ConflictsWith(
					path.MatchRelative().AtParent().AtParent().AtParent().AtParent().AtName("required_jq_query"),
				),
			},
		},
		"depends_on": schema.ListAttribute{
			MarkdownDescription: "The inputs this input depends on.",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"visible": schema.BoolAttribute{
			MarkdownDescription: "The visibility of the " + propertyType + " input.",
			Optional:            true,
		},
		"visible_jq_query": schema.StringAttribute{
			MarkdownDescription: "The visibility condition jq query of the " + propertyType + " input.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("visible")),
			},
		},
		"read_only": schema.BoolAttribute{
			MarkdownDescription: "Shows the value of the " + propertyType + " input without letting the user change it. " +
				"The value is still submitted with the form, unlike `disabled`.",
			Optional: true,
		},
		"read_only_jq_query": schema.StringAttribute{
			MarkdownDescription: "The read only condition jq query of the " + propertyType + " input.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("read_only")),
			},
		},
		"disabled": schema.BoolAttribute{
			MarkdownDescription: "Greys out the " + propertyType + " input. A disabled input is excluded from the submitted " +
				"data and from `required` validation, which makes it a way to drop an input from the form based on the other answers. " +
				"Use `read_only` to keep the value in the submission.",
			Optional: true,
		},
		"disabled_jq_query": schema.StringAttribute{
			MarkdownDescription: "The disabled condition jq query of the " + propertyType + " input.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("disabled")),
			},
		},
	}
}

func stringPropertySchema() schema.Attribute {
	properties := map[string]schema.Attribute{
		"default": schema.StringAttribute{
			MarkdownDescription: "The default of the string input.",
			Optional:            true,
		},
		"default_jq_query": schema.StringAttribute{
			MarkdownDescription: "The default jq query of the string input.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("default")),
			},
		},
		"format": schema.StringAttribute{
			MarkdownDescription: "The format of the string input.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf(stringPropertyFormats...),
			},
		},
		"blueprint": schema.StringAttribute{
			MarkdownDescription: "The blueprint the entities are taken from. Required when `format` is `entity`.",
			Optional:            true,
		},
		"dataset": datasetSchema(),
		"sort":    entitiesSortSchema(),
		"min_length": schema.Int64Attribute{
			MarkdownDescription: "The min length of the string input.",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"max_length": schema.Int64Attribute{
			MarkdownDescription: "The max length of the string input.",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"pattern": schema.StringAttribute{
			MarkdownDescription: "The regex pattern the value has to match.",
			Optional:            true,
		},
		"pattern_jq_query": schema.StringAttribute{
			MarkdownDescription: "A jq query resolving the pattern of the string input, either a regex string or a list of allowed values.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("pattern")),
				stringvalidator.LengthAtLeast(1),
			},
		},
		"enum": schema.ListAttribute{
			MarkdownDescription: "The values the user can pick from.",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.UniqueValues(),
				listvalidator.SizeAtLeast(1),
			},
		},
		"enum_jq_query": schema.StringAttribute{
			MarkdownDescription: "A jq query resolving the values the user can pick from.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("enum")),
			},
		},
		"enum_colors": schema.MapAttribute{
			MarkdownDescription: "The colors of the enum values.",
			Optional:            true,
			ElementType:         types.StringType,
		},
	}

	utils.CopyMaps(properties, propertyMetadataSchema("string"))
	return schema.MapNestedAttribute{
		MarkdownDescription: "The string inputs of the form.",
		Optional:            true,
		NestedObject:        schema.NestedAttributeObject{Attributes: properties},
	}
}

func numberPropertySchema() schema.Attribute {
	properties := map[string]schema.Attribute{
		"default": schema.Float64Attribute{
			MarkdownDescription: "The default of the number input.",
			Optional:            true,
		},
		"default_jq_query": schema.StringAttribute{
			MarkdownDescription: "The default jq query of the number input.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("default")),
			},
		},
		"minimum": schema.Float64Attribute{
			MarkdownDescription: "The smallest value the input accepts.",
			Optional:            true,
		},
		"maximum": schema.Float64Attribute{
			MarkdownDescription: "The largest value the input accepts.",
			Optional:            true,
		},
		"exclusive_minimum": schema.Float64Attribute{
			MarkdownDescription: "The value the input has to be strictly greater than.",
			Optional:            true,
		},
		"exclusive_maximum": schema.Float64Attribute{
			MarkdownDescription: "The value the input has to be strictly smaller than.",
			Optional:            true,
		},
		"enum": schema.ListAttribute{
			MarkdownDescription: "The values the user can pick from.",
			Optional:            true,
			ElementType:         types.Float64Type,
			Validators: []validator.List{
				listvalidator.UniqueValues(),
				listvalidator.SizeAtLeast(1),
			},
		},
		"enum_jq_query": schema.StringAttribute{
			MarkdownDescription: "A jq query resolving the values the user can pick from.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("enum")),
			},
		},
	}

	utils.CopyMaps(properties, propertyMetadataSchema("number"))
	return schema.MapNestedAttribute{
		MarkdownDescription: "The number inputs of the form.",
		Optional:            true,
		NestedObject:        schema.NestedAttributeObject{Attributes: properties},
	}
}

func booleanPropertySchema() schema.Attribute {
	properties := map[string]schema.Attribute{
		"default": schema.BoolAttribute{
			MarkdownDescription: "The default of the boolean input.",
			Optional:            true,
		},
		"default_jq_query": schema.StringAttribute{
			MarkdownDescription: "The default jq query of the boolean input.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("default")),
			},
		},
	}

	utils.CopyMaps(properties, propertyMetadataSchema("boolean"))
	return schema.MapNestedAttribute{
		MarkdownDescription: "The boolean inputs of the form.",
		Optional:            true,
		NestedObject:        schema.NestedAttributeObject{Attributes: properties},
	}
}

func objectPropertySchema() schema.Attribute {
	properties := map[string]schema.Attribute{
		"default": schema.StringAttribute{
			MarkdownDescription: "The default of the object input, as a JSON encoded string.",
			Optional:            true,
		},
		"default_jq_query": schema.StringAttribute{
			MarkdownDescription: "The default jq query of the object input.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("default")),
			},
		},
		"format": schema.StringAttribute{
			MarkdownDescription: "The format of the object input. `labeled-url` renders a url with a display text. " +
				"Leave it out for a free form object.",
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf(objectFormats...),
			},
		},
	}

	utils.CopyMaps(properties, propertyMetadataSchema("object"))
	return schema.MapNestedAttribute{
		MarkdownDescription: "The object inputs of the form.",
		Optional:            true,
		NestedObject:        schema.NestedAttributeObject{Attributes: properties},
	}
}

func arrayPropertySchema() schema.Attribute {
	properties := map[string]schema.Attribute{
		"default_jq_query": schema.StringAttribute{
			MarkdownDescription: "The default jq query of the array input.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("string_items").AtName("default")),
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("number_items").AtName("default")),
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("object_items").AtName("default")),
			},
		},
		"min_items": schema.Int64Attribute{
			MarkdownDescription: "The min items of the array input.",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"min_items_jq_query": schema.StringAttribute{
			MarkdownDescription: "The min items jq query of the array input.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("min_items")),
			},
		},
		"max_items": schema.Int64Attribute{
			MarkdownDescription: "The max items of the array input.",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"max_items_jq_query": schema.StringAttribute{
			MarkdownDescription: "The max items jq query of the array input.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("max_items")),
			},
		},
		"unique_items": schema.BoolAttribute{
			MarkdownDescription: "Whether the values of the array have to be unique.",
			Optional:            true,
		},
		"string_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The string items of the array input.",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"format": schema.StringAttribute{
					MarkdownDescription: "The format of each item.",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(stringItemFormats...),
					},
				},
				"blueprint": schema.StringAttribute{
					MarkdownDescription: "The blueprint the entities are taken from. Required when `format` is `entity`.",
					Optional:            true,
				},
				"default": schema.ListAttribute{
					MarkdownDescription: "The default value of the items.",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"enum": schema.ListAttribute{
					MarkdownDescription: "The values the user can pick from.",
					Optional:            true,
					ElementType:         types.StringType,
					Validators: []validator.List{
						listvalidator.UniqueValues(),
						listvalidator.SizeAtLeast(1),
					},
				},
				"enum_jq_query": schema.StringAttribute{
					MarkdownDescription: "A jq query resolving the values the user can pick from.",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("enum")),
					},
				},
				"enum_colors": schema.MapAttribute{
					MarkdownDescription: "The colors of the enum values.",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"dataset": schema.StringAttribute{
					MarkdownDescription: "The dataset filtering the entities of the items, as a JSON encoded string.",
					Optional:            true,
				},
			},
		},
		"number_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The number items of the array input.",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"default": schema.ListAttribute{
					MarkdownDescription: "The default value of the items.",
					Optional:            true,
					ElementType:         types.Float64Type,
				},
				"enum": schema.ListAttribute{
					MarkdownDescription: "The values the user can pick from.",
					Optional:            true,
					ElementType:         types.Float64Type,
					Validators: []validator.List{
						listvalidator.UniqueValues(),
						listvalidator.SizeAtLeast(1),
					},
				},
				"enum_jq_query": schema.StringAttribute{
					MarkdownDescription: "A jq query resolving the values the user can pick from.",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("enum")),
					},
				},
				"enum_colors": schema.MapAttribute{
					MarkdownDescription: "The colors of the enum values.",
					Optional:            true,
					ElementType:         types.StringType,
				},
			},
		},
		"object_items": schema.SingleNestedAttribute{
			MarkdownDescription: "The object items of the array input.",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"format": schema.StringAttribute{
					MarkdownDescription: "The format of each item.",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(objectItemFormats...),
					},
				},
				"default": schema.ListAttribute{
					MarkdownDescription: "The default value of the items.",
					Optional:            true,
					ElementType:         types.MapType{ElemType: types.StringType},
				},
			},
		},
		"sort": entitiesSortSchema(),
	}

	utils.CopyMaps(properties, propertyMetadataSchema("array"))
	return schema.MapNestedAttribute{
		MarkdownDescription: "The array inputs of the form.",
		Optional:            true,
		NestedObject:        schema.NestedAttributeObject{Attributes: properties},
	}
}

func entitiesSortSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "How the entities are sorted in the form.",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"property": schema.StringAttribute{
				MarkdownDescription: "The property to sort the entities by.",
				Required:            true,
			},
			"order": schema.StringAttribute{
				MarkdownDescription: "The order to sort the entities in.",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("ASC"),
				Validators: []validator.String{
					stringvalidator.OneOf("ASC", "DESC"),
				},
			},
		},
	}
}

func datasetSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "The dataset filtering the entities the user can pick from.",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"combinator": schema.StringAttribute{
				MarkdownDescription: "How the rules are combined.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("and", "or"),
				},
			},
			"rules": schema.ListNestedAttribute{
				MarkdownDescription: "The rules of the dataset. A rule either filters on a property or groups nested rules under a combinator.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: datasetRuleSchema(10),
				},
			},
		},
	}
}

func datasetRuleSchema(depth int) map[string]schema.Attribute {
	attributes := map[string]schema.Attribute{
		"blueprint": schema.StringAttribute{
			MarkdownDescription: "The blueprint identifier of the rule.",
			Optional:            true,
		},
		"property": schema.StringAttribute{
			MarkdownDescription: "The property identifier of the rule.",
			Optional:            true,
		},
		"operator": schema.StringAttribute{
			MarkdownDescription: "The operator of the rule. Set on filtering rules and left out on group rules.",
			Optional:            true,
		},
		"value": schema.ObjectAttribute{
			MarkdownDescription: "A value resolved from the form or the trigger when the form is rendered.",
			Optional:            true,
			AttributeTypes: map[string]attr.Type{
				"jq_query": types.StringType,
			},
		},
		"value_json": schema.StringAttribute{
			MarkdownDescription: "A fixed value, as a JSON encoded string. Use `value` for a value resolved by a jq query.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("value")),
			},
		},
		"combinator": schema.StringAttribute{
			MarkdownDescription: "How the nested rules of a group rule are combined.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("and", "or"),
			},
		},
	}

	if depth > 0 {
		attributes["rules"] = schema.ListNestedAttribute{
			MarkdownDescription: "The nested rules of a group rule.",
			Optional:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: datasetRuleSchema(depth - 1),
			},
		}
	}

	return attributes
}
