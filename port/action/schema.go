package action

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
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
				boolvalidator.ConflictsWith(path.MatchRoot("self_service_trigger").AtName("required_jq_query")),
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
	}
}

func StringBooleanOrJQTemplateValidator() []validator.String {
	return []validator.String{
		stringvalidator.Any(
			stringvalidator.OneOf("true", "false"),
			stringvalidator.RegexMatches(
				regexp.MustCompile(`^[\n\r\s]*{{.*}}[\n\r\s]*$`),
				"must be a valid jq template: {{JQ_EXPRESSION}}",
			)),
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
			Optional:            true,
			DeprecationMessage:  "Action is not attached to blueprint anymore. This value is ignored",
			Validators:          []validator.String{stringvalidator.OneOf("")},
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "Title",
			Optional:            true,
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "Icon",
			Optional:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "Description",
			Optional:            true,
		},
		"self_service_trigger": schema.SingleNestedAttribute{
			MarkdownDescription: "Self service trigger for the action. Note: you can define only one of `order_properties` and `steps`",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"blueprint_identifier": schema.StringAttribute{
					Description: "The ID of the blueprint",
					Optional:    true,
				},
				"operation": schema.StringAttribute{
					MarkdownDescription: "The operation type of the action",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("CREATE", "DAY-2", "DELETE"),
					},
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
				"order_properties": schema.ListAttribute{
					MarkdownDescription: "Order properties",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"steps": schema.ListNestedAttribute{
					MarkdownDescription: "The steps of the action",
					Optional:            true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"title": schema.StringAttribute{
								MarkdownDescription: "The step's title",
								Required:            true,
							},
							"order": schema.ListAttribute{
								MarkdownDescription: "The order of the properties in this step",
								Required:            true,
								ElementType:         types.StringType,
							},
						},
					},
					Validators: []validator.List{
						listvalidator.ConflictsWith(
							path.MatchRelative().AtParent().AtName("order_properties"),
						),
					},
				},
				"condition": schema.StringAttribute{
					MarkdownDescription: "The `condition` field allows you to define rules using Port's [search & query syntax](https://docs.getport.io/search-and-query/#rules) to determine which entities the action will be available for.",
					Optional:            true,
				},
			},
			Validators: []validator.Object{
				objectvalidator.ExactlyOneOf(
					path.MatchRoot("self_service_trigger"),
					path.MatchRoot("automation_trigger"),
				),
			},
		},
		"automation_trigger": schema.SingleNestedAttribute{
			MarkdownDescription: "Automation trigger for the action",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"entity_created_event": schema.SingleNestedAttribute{
					MarkdownDescription: "Entity created event trigger",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"blueprint_identifier": schema.StringAttribute{
							MarkdownDescription: "The blueprint identifier of the created entity",
							Required:            true,
						},
					},
					Validators: []validator.Object{
						objectvalidator.ExactlyOneOf(
							path.MatchRelative().AtParent().AtName("entity_created_event"),
							path.MatchRelative().AtParent().AtName("entity_updated_event"),
							path.MatchRelative().AtParent().AtName("entity_deleted_event"),
							path.MatchRelative().AtParent().AtName("any_entity_change_event"),
							path.MatchRelative().AtParent().AtName("timer_property_expired_event"),
							path.MatchRelative().AtParent().AtName("run_created_event"),
							path.MatchRelative().AtParent().AtName("run_updated_event"),
							path.MatchRelative().AtParent().AtName("any_run_change_event"),
						),
					},
				},
				"entity_updated_event": schema.SingleNestedAttribute{
					MarkdownDescription: "Entity updated event trigger",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"blueprint_identifier": schema.StringAttribute{
							MarkdownDescription: "The blueprint identifier of the updated entity",
							Required:            true,
						},
					},
				},
				"entity_deleted_event": schema.SingleNestedAttribute{
					MarkdownDescription: "Entity deleted event trigger",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"blueprint_identifier": schema.StringAttribute{
							MarkdownDescription: "The blueprint identifier of the deleted entity",
							Required:            true,
						},
					},
				},
				"any_entity_change_event": schema.SingleNestedAttribute{
					MarkdownDescription: "Any entity change event trigger",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"blueprint_identifier": schema.StringAttribute{
							MarkdownDescription: "The blueprint identifier of the changed entity",
							Required:            true,
						},
					},
				},
				"timer_property_expired_event": schema.SingleNestedAttribute{
					MarkdownDescription: "Timer property expired event trigger",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"blueprint_identifier": schema.StringAttribute{
							MarkdownDescription: "The blueprint identifier of the expired timer property",
							Required:            true,
						},
						"property_identifier": schema.StringAttribute{
							MarkdownDescription: "The property identifier of the expired timer property",
							Required:            true,
						},
					},
				},
				"run_created_event": schema.SingleNestedAttribute{
					MarkdownDescription: "Run created event trigger",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"action_identifier": schema.StringAttribute{
							MarkdownDescription: "The action identifier of the created run",
							Required:            true,
						},
					},
				},
				"run_updated_event": schema.SingleNestedAttribute{
					MarkdownDescription: "Run updated event trigger",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"action_identifier": schema.StringAttribute{
							MarkdownDescription: "The action identifier of the updated run",
							Required:            true,
						},
					},
				},
				"any_run_change_event": schema.SingleNestedAttribute{
					MarkdownDescription: "Any run change event trigger",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"action_identifier": schema.StringAttribute{
							MarkdownDescription: "The action identifier of the changed run",
							Required:            true,
						},
					},
				},
				"jq_condition": schema.SingleNestedAttribute{
					MarkdownDescription: "JQ condition for automation trigger",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"expressions": schema.ListAttribute{
							MarkdownDescription: "The jq expressions of the condition",
							ElementType:         types.StringType,
							Required:            true,
						},
						"combinator": schema.StringAttribute{
							MarkdownDescription: "The combinator of the condition",
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString("and"),
							Validators: []validator.String{
								stringvalidator.OneOf("and", "or"),
							},
						},
					},
				},
			},
		},
		"kafka_method": schema.SingleNestedAttribute{
			MarkdownDescription: "Kafka invocation method",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"payload": schema.StringAttribute{
					MarkdownDescription: "The Kafka message [payload](https://docs.getport.io/create-self-service-experiences/setup-backend/#define-the-actions-payload) should be in `JSON` format, encoded as a string. Use [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode) to encode arrays or objects. Learn about how to [define the action payload](https://docs.getport.io/create-self-service-experiences/setup-backend/#define-the-actions-payload).",
					Optional:            true,
				},
			},
			Validators: []validator.Object{
				objectvalidator.ExactlyOneOf(
					path.MatchRoot("kafka_method"),
					path.MatchRoot("webhook_method"),
					path.MatchRoot("github_method"),
					path.MatchRoot("gitlab_method"),
					path.MatchRoot("azure_method"),
					path.MatchRoot("upsert_entity_method"),
				),
			},
		},
		"webhook_method": schema.SingleNestedAttribute{
			MarkdownDescription: "Webhook invocation method",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"url": schema.StringAttribute{
					MarkdownDescription: "Required when selecting type WEBHOOK. The URL to invoke the action",
					Required:            true,
				},
				"agent": schema.StringAttribute{
					MarkdownDescription: "Specifies whether to use an agent to invoke the action. This can be a boolean value (`'true''` or `'false'`) or a JQ if dynamic evaluation is needed.",
					Optional:            true,
					Validators:          StringBooleanOrJQTemplateValidator(),
				},
				"synchronized": schema.StringAttribute{
					MarkdownDescription: "Synchronize the action",
					Optional:            true,
					Validators:          StringBooleanOrJQTemplateValidator(),
				},
				"method": schema.StringAttribute{
					MarkdownDescription: "The HTTP method to invoke the action",
					Optional:            true,
				},
				"headers": schema.MapAttribute{
					MarkdownDescription: "The HTTP headers for invoking the action. They should be encoded as a key-value object to a string using [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode). Learn about how to [define the action payload](https://docs.getport.io/create-self-service-experiences/setup-backend/#define-the-actions-payload).",
					ElementType:         types.StringType,
					Optional:            true,
				},
				"body": schema.StringAttribute{
					MarkdownDescription: "The Webhook body should be in `JSON` format, encoded as a string. Use [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode) to encode arrays or objects. Learn about how to [define the action payload](https://docs.getport.io/create-self-service-experiences/setup-backend/#define-the-actions-payload).",
					Optional:            true,
				},
			},
		},
		"github_method": schema.SingleNestedAttribute{
			MarkdownDescription: "GitHub invocation method",
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
				"workflow_inputs": schema.StringAttribute{
					MarkdownDescription: "The GitHub [workflow inputs](https://docs.getport.io/create-self-service-experiences/setup-backend/#define-the-actions-payload) should be in `JSON` format, encoded as a string. Use [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode) to encode arrays or objects. Learn about how to [define the action payload](https://docs.getport.io/create-self-service-experiences/setup-backend/#define-the-actions-payload).",
					Optional:            true,
				},
				"report_workflow_status": schema.StringAttribute{
					MarkdownDescription: "Report the workflow status when invoking the action",
					Optional:            true,
					Validators:          StringBooleanOrJQTemplateValidator(),
				},
			},
		},
		"gitlab_method": schema.SingleNestedAttribute{
			MarkdownDescription: "Gitlab invocation method",
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
				"default_ref": schema.StringAttribute{
					MarkdownDescription: "The default ref of the action",
					Optional:            true,
				},
				"pipeline_variables": schema.StringAttribute{
					MarkdownDescription: "The Gitlab pipeline variables should be in `JSON` format, encoded as a string. Use [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode) to encode arrays or objects. Learn about how to [define the action payload](https://docs.getport.io/create-self-service-experiences/setup-backend/#define-the-actions-payload).",
					Optional:            true,
				},
			},
		},
		"azure_method": schema.SingleNestedAttribute{
			MarkdownDescription: "Azure DevOps invocation method",
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
				"payload": schema.StringAttribute{
					MarkdownDescription: "The Azure Devops workflow [payload](https://docs.getport.io/create-self-service-experiences/setup-backend/#define-the-actions-payload) should be in `JSON` format, encoded as a string. Use [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode) to encode arrays or objects. Learn about how to [define the action payload](https://docs.getport.io/create-self-service-experiences/setup-backend/#define-the-actions-payload).",
					Optional:            true,
				},
			},
		},
		"upsert_entity_method": schema.SingleNestedAttribute{
			MarkdownDescription: "Upsert Entity invocation method",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"title": schema.StringAttribute{
					MarkdownDescription: "The title of the entity",
					Optional:            true,
				},
				"blueprint_identifier": schema.StringAttribute{
					MarkdownDescription: "Required when selecting type Upsert Entity. The blueprint identifier of the entity for the upsert",
					Required:            true,
				},
				"mapping": schema.SingleNestedAttribute{
					MarkdownDescription: "Upsert Entity invocation method",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"identifier": schema.StringAttribute{
							MarkdownDescription: "Required when selecting type Upsert Entity. The entity identifier for the upsert",
							Optional:            true,
						},
						"teams": schema.ListAttribute{
							MarkdownDescription: "The teams the entity belongs to",
							ElementType:         types.StringType,
							Optional:            true,
						},
						"icon": schema.StringAttribute{
							MarkdownDescription: "The icon of the entity",
							Optional:            true,
						},
						"properties": schema.StringAttribute{
							MarkdownDescription: "The properties of the entity (key-value object encoded to a string)",
							Optional:            true,
						},
						"relations": schema.StringAttribute{
							MarkdownDescription: "The relations of the entity (key-value object encoded to a string)",
							Optional:            true,
						},
					},
				},
			},
		},
		"required_approval": schema.StringAttribute{
			MarkdownDescription: "Require approval before invoking the action. Can be one of \"true\", \"false\", \"ANY\" or \"ALL\"",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("true", "false", "ANY", "ALL"),
			},
		},
		"approval_webhook_notification": schema.SingleNestedAttribute{
			MarkdownDescription: "The webhook notification of the approval",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"url": schema.StringAttribute{
					MarkdownDescription: "The URL to invoke the webhook",
					Required:            true,
				},
				"format": schema.StringAttribute{
					MarkdownDescription: "The format to invoke the webhook",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("json", "slack"),
					},
				},
			},
		},
		"approval_email_notification": schema.ObjectAttribute{
			MarkdownDescription: "The email notification of the approval",
			Optional:            true,
			AttributeTypes:      map[string]attr.Type{},
			Validators: []validator.Object{
				objectvalidator.ConflictsWith(path.MatchRoot("approval_webhook_notification")),
			},
		},
		"publish": schema.BoolAttribute{
			MarkdownDescription: "Publish action",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
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
			MarkdownDescription: "The format of the string property, Accepted values include `date-time`, `url`, `email`, `ipv4`, `ipv6`, `yaml`, `entity`, `user`, `team`, `proto`, `markdown`",
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
			MarkdownDescription: "The algorithm to encrypt the property with. Accepted value: `aes256-gcm`",
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
		"dataset": schema.SingleNestedAttribute{
			MarkdownDescription: "The dataset of an the entity-format property",
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
		"sort": schema.SingleNestedAttribute{
			MarkdownDescription: "How to sort entities when in the self service action form in the UI",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"property": schema.StringAttribute{
					MarkdownDescription: "The property to sort the entities by",
					Required:            true,
				},
				"order": schema.StringAttribute{
					MarkdownDescription: "The order to sort the entities in",
					Computed:            true,
					Optional:            true,
					Default:             stringdefault.StaticString("ASC"),
					Validators: []validator.String{
						stringvalidator.OneOf("ASC", "DESC"),
					},
				},
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
			MarkdownDescription: "The algorithm to encrypt the property with. Accepted value: `aes256-gcm`",
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
			MarkdownDescription: "An array of string items within the property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"format": schema.StringAttribute{
					MarkdownDescription: "The format of the string property, Accepted values include `date-time`, `url`, `email`, `ipv4`, `ipv6`, `yaml`, `entity`, `user`, `team`, `proto`, `markdown`",
					Optional:            true,
				},
				"blueprint": schema.StringAttribute{
					MarkdownDescription: "The blueprint identifier related to each string item",
					Optional:            true,
				},
				"default": schema.ListAttribute{
					MarkdownDescription: "The default value of the items",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"enum": schema.ListAttribute{
					MarkdownDescription: "The enum of possible values for the string items",
					Optional:            true,
					ElementType:         types.StringType,
					Validators: []validator.List{
						listvalidator.UniqueValues(),
						listvalidator.SizeAtLeast(1),
					},
				},
				"enum_jq_query": schema.StringAttribute{
					MarkdownDescription: "The jq query for the enum of string items",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("enum")),
					},
				},
				"dataset": schema.StringAttribute{
					MarkdownDescription: "The dataset of the entity-format items",
					Optional:            true,
				},
			},
		},
		"number_items": schema.SingleNestedAttribute{
			MarkdownDescription: "An array of number items within the property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"default": schema.ListAttribute{
					MarkdownDescription: "The default values for the number items",
					Optional:            true,
					ElementType:         types.Float64Type,
				},
				"enum": schema.ListAttribute{
					MarkdownDescription: "The enum of possible values for the number items",
					Optional:            true,
					ElementType:         types.Float64Type,
					Validators: []validator.List{
						listvalidator.UniqueValues(),
						listvalidator.SizeAtLeast(1),
					},
				},
				"enum_jq_query": schema.StringAttribute{
					MarkdownDescription: "The jq query for the enum number items",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("enum")),
					},
				},
			},
		},
		"boolean_items": schema.SingleNestedAttribute{
			MarkdownDescription: "An array of boolean items within the property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"default": schema.ListAttribute{
					MarkdownDescription: "The default values for the boolean items",
					Optional:            true,
					ElementType:         types.BoolType,
				},
			},
		},
		"object_items": schema.SingleNestedAttribute{
			MarkdownDescription: "An array of object items within the property",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"default": schema.ListAttribute{
					MarkdownDescription: "The default values for the object items",
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
		"sort": schema.SingleNestedAttribute{
			MarkdownDescription: "How to sort entities when in the self service action form in the UI",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"property": schema.StringAttribute{
					MarkdownDescription: "The property to sort the entities by",
					Required:            true,
				},
				"order": schema.StringAttribute{
					MarkdownDescription: "The order to sort the entities in",
					Computed:            true,
					Optional:            true,
					Default:             stringdefault.StaticString("ASC"),
					Validators: []validator.String{
						stringvalidator.OneOf("ASC", "DESC"),
					},
				},
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
		MarkdownDescription: ResourceMarkdownDescription,
		Attributes:          ActionSchema(),
	}
}

func (r *ActionResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var state *ActionValidationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	validateUserInputRequiredNotSetToFalse(ctx, state, resp)
}

func validateUserInputRequiredNotSetToFalse(ctx context.Context, state *ActionValidationModel, resp *resource.ValidateConfigResponse) {
	// go over all the properties and check if required is set to false, it is false, raise an error that false is not
	// supported anymore
	const errorString = "required is set to false, this is not supported anymore, if you don't want to make the stringProp required, remove the required stringProp"

	if state.SelfServiceTrigger.IsNull() {
		return
	}

	var sst = state.SelfServiceTrigger.Attributes()
	if sst == nil {
		return
	}

	var up, _ = sst["user_properties"]
	if up == nil {
		return
	}

	var val, err = up.ToTerraformValue(ctx)
	if err != nil {
		return
	}

	userProperties := map[string]tftypes.Value{}

	err = val.As(&userProperties)
	if err != nil {
		return
	}

	var stringProperties, _ = userProperties["string_props"]

	if !stringProperties.IsNull() {
		v := map[string]tftypes.Value{}

		err = val.As(&v)
		if err != nil {
			return
		}

		stringPropValidationsObjects := make(map[string]StringPropValidationModel, len(v))
		for key := range v {
			var val StringPropValidationModel
			err = v[key].As(&val)

			if err != nil {
				return
			}

			stringPropValidationsObjects[key] = val
		}

		for _, stringProp := range stringPropValidationsObjects {
			if stringProp.Required != nil && !*stringProp.Required {
				resp.Diagnostics.AddError(errorString, fmt.Sprint(`Error in User Property: `, stringProp.Title, ` in action: `, state.Identifier))
			}
		}
	}

	var numberProperties, _ = userProperties["number_props"]

	if !numberProperties.IsNull() {
		v := map[string]tftypes.Value{}

		err = val.As(&v)
		if err != nil {
			return
		}

		numberPropValidationsObjects := make(map[string]NumberPropValidationModel, len(v))
		for key := range v {
			var val NumberPropValidationModel
			err = v[key].As(&val)

			if err != nil {
				return
			}

			numberPropValidationsObjects[key] = val
		}

		for _, numberProp := range numberPropValidationsObjects {
			if numberProp.Required != nil && !*numberProp.Required {
				resp.Diagnostics.AddError(errorString, fmt.Sprint(`Error in User Property: `, numberProp.Title, ` in action: `, state.Identifier))
			}
		}
	}

	var booleanProperties, _ = userProperties["boolean_props"]

	if !booleanProperties.IsNull() {
		v := map[string]tftypes.Value{}

		err = val.As(&v)
		if err != nil {
			return
		}

		booleanPropValidationsObjects := make(map[string]BooleanPropValidationModel, len(v))
		for key := range v {
			var val BooleanPropValidationModel
			err = v[key].As(&val)

			if err != nil {
				return
			}

			booleanPropValidationsObjects[key] = val
		}

		for _, booleanProp := range booleanPropValidationsObjects {
			if booleanProp.Required != nil && !*booleanProp.Required {
				resp.Diagnostics.AddError(errorString, fmt.Sprint(`Error in User Property: `, booleanProp.Title, ` in action: `, state.Identifier))
			}
		}
	}

	var objectProperties, _ = userProperties["object_props"]

	if !objectProperties.IsNull() {
		v := map[string]tftypes.Value{}

		err = val.As(&v)
		if err != nil {
			return
		}

		objectPropValidationsObjects := make(map[string]ObjectPropValidationModel, len(v))
		for key := range v {
			var val ObjectPropValidationModel
			err = v[key].As(&val)

			if err != nil {
				return
			}

			objectPropValidationsObjects[key] = val
		}

		for _, objectProp := range objectPropValidationsObjects {
			if objectProp.Required != nil && !*objectProp.Required {
				resp.Diagnostics.AddError(errorString, fmt.Sprint(`Error in User Property: `, objectProp.Title, ` in action: `, state.Identifier))
			}
		}
	}

	var arrayProperties, _ = userProperties["array_props"]

	if !arrayProperties.IsNull() {
		v := map[string]tftypes.Value{}

		err = val.As(&v)
		if err != nil {
			return
		}

		arrayPropValidationsObjects := make(map[string]ArrayPropValidationModel, len(v))
		for key := range v {
			var val ArrayPropValidationModel
			err = v[key].As(&val)

			if err != nil {
				return
			}

			arrayPropValidationsObjects[key] = val
		}

		for _, arrayProp := range arrayPropValidationsObjects {
			if arrayProp.Required != nil && !*arrayProp.Required {
				resp.Diagnostics.AddError(errorString, fmt.Sprint(`Error in User Property: `, arrayProp.Title, ` in action: `, state.Identifier))
			}
		}
	}
}

var ResourceMarkdownDescription = `


# Action resource

Docs for the Action resource can be found [here](https://docs.getport.io/create-self-service-experiences/).

## Example Usage

` + "```hcl" + `
resource "port_action" "create_microservice" {
	title = "Create Microservice"
	identifier = "create-microservice"
	icon = "Terraform"
	self_service_trigger = {
		operation = "CREATE"
		blueprint_identifier = port_blueprint.microservice.identifier
		user_properties = {
			string_props = {
				myStringIdentifier = {
					title = "My String Identifier"
					required = true
                    format = "entity"
                    blueprint = port_blueprint.parent.identifier
                    dataset = {
                        combinator = "and"
                        rules = [{
                            property = "$title"
                            operator = "contains"
                            value = {
                                jq_query = "\"specificValue\""
                            }
                        }]
                    }
                    sort = {
                        property = "$updatedAt"
                        order = "DESC"
                    }
				}
			}
			number_props = {
				myNumberIdentifier = {
					title = "My Number Identifier"
					required = true
					maximum = 100
					minimum = 0
				}
			}
			boolean_props = {
				myBooleanIdentifier = {
					title = "My Boolean Identifier"
					required = true
				}
			}
			object_props = {
				myObjectIdentifier = {
					title = "My Object Identifier"
					required = true
				}
			}
			array_props = {
				myArrayIdentifier = {
					title = "My Array Identifier"
					required = true
					string_items = {
						format = "entity"
                        blueprint = port_blueprint.parent.identifier
                        dataset = jsonencode({
                            combinator = "and"
                            rules = [{
                                property = "$title"
                                operator = "contains"
                                value    = "specificValue"
                            }]
                        })
					}
                    sort = {
                        property = "$updatedAt"
                        order = "DESC"
                    }
				}
			}
		}
	}
	kafka_method = {
		payload = jsonencode({
		  runId: "{{"{{.run.id}}"}}"
		})
	}
}` + "\n```" + `

## Example Usage with Automation trigger

Port allows setting an automation trigger to an action, for executing an action based on event occurred to an entity in Port.

` + "```hcl" + `
resource "port_action" "delete_temporary_microservice" {
	title = "Delete Temporary Microservice"
	identifier = "delete-temp-microservice"
	icon = "Terraform"
	automation_trigger = {
		timer_property_expired_event = {
			blueprint_identifier = port_blueprint.microservice.identifier
			property_identifier = "ttl"
		}
	}
	kafka_method = {
		payload = jsonencode({
		  runId: "{{"{{.run.id}}"}}"
		})
	}
}
` + "\n```" + `

## Example Usage With Condition

` + "```hcl" + `
resource "port_action" "create_microservice" {
	title = "Create Microservice"
	identifier = "create-microservice"
	icon = "Terraform"
	self_service_trigger = {
		operation = "CREATE"
		blueprint_identifier = port_blueprint.microservice.identifier
		condition = jsonencode({
			type = "SEARCH"
			combinator = "and"
			rules = [
				{
					property = "$title"
					operator = "!="
					value = "Test"
				}
			]
		})
		user_properties = {
			string_props = {
				myStringIdentifier = {
					title = "My String Identifier"
					required = true
				}
			}
		}
	}
	kafka_method = {
		payload = jsonencode({
		  runId: "{{"{{.run.id}}"}}"
		})
	}
	
` + "```" + `

`
