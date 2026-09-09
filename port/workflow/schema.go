package workflow

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

// CURSOR_AGENT is intentionally absent: it is still a live member of the service
// invocation union, but is deprecated and not exposed through Terraform.
//
// The integration action `deferIntegrationInstallation` field is omitted for a
// different reason. It lets a workflow be saved before its integration is installed,
// but Terraform users are expected to install the integration first, so an
// integration_action node always needs an installation_id that resolves against the
// integrations installed in the organization.
var nodeTypeBlockNames = []string{
	"self_serve_trigger",
	"event_trigger",
	"schedule_trigger",
	"kafka",
	"webhook",
	"integration_action",
	"upsert_entity",
	"ai",
	"ai_agent",
	"condition",
	"input",
}

const (
	maxIdentifierLength  = 60
	maxTitleLength       = 100
	maxDescriptionLength = 1100
	maxCategoryLength    = 40
	maxNodesCount        = 100
)

var identifierPattern = regexp.MustCompile(`^[\p{L}0-9@_:-]+$`)

func identifierValidators() []validator.String {
	return []validator.String{
		stringvalidator.LengthAtMost(maxIdentifierLength),
		stringvalidator.RegexMatches(identifierPattern, "must contain only letters, digits, and the characters @ _ : -"),
	}
}

func titleValidators() []validator.String {
	return []validator.String{stringvalidator.LengthAtMost(maxTitleLength)}
}

func descriptionValidators() []validator.String {
	return []validator.String{stringvalidator.LengthAtMost(maxDescriptionLength)}
}

func nodeTypeValidators() []validator.Object {
	expressions := make([]path.Expression, 0, len(nodeTypeBlockNames))
	for _, name := range nodeTypeBlockNames {
		expressions = append(expressions, path.MatchRelative().AtParent().AtName(name))
	}
	return []validator.Object{objectvalidator.ExactlyOneOf(expressions...)}
}

func onFailureAttribute() schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: "The action to take if the node fails. One of `continue`, `terminate`.",
		Optional:            true,
		Computed:            true,
		Default:             stringdefault.StaticString("terminate"),
		Validators: []validator.String{
			stringvalidator.OneOf("continue", "terminate"),
		},
	}
}

func publishedAttribute() schema.Attribute {
	return schema.BoolAttribute{
		MarkdownDescription: "Whether the trigger is published.",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(true),
	}
}

func statusLabelBlock(description string) schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: description,
		Attributes: map[string]schema.Attribute{
			"text": schema.StringAttribute{
				MarkdownDescription: "The label text. Supports JQ expressions for dynamic content.",
				Optional:            true,
			},
			"variant": schema.StringAttribute{
				MarkdownDescription: "Semantic variant controlling the label color/style. One of `success`, `alert`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("success", "alert"),
				},
			},
		},
	}
}

func principalAttributes(noun string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"users": schema.ListAttribute{
			MarkdownDescription: "The emails of the users " + noun + ". They must exist in the organization.",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"roles": schema.ListAttribute{
			MarkdownDescription: "The roles " + noun + ".",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"teams": schema.ListAttribute{
			MarkdownDescription: "The identifiers of the teams " + noun + ". They must exist in the organization.",
			Optional:            true,
			ElementType:         types.StringType,
		},
	}
}

func permissionsBlock(description string) schema.Block {
	attributes := principalAttributes("the permission applies to")
	attributes["policy"] = schema.StringAttribute{
		MarkdownDescription: "A JSON encoded RBAC query that dynamically resolves who is permitted, of the form " +
			"`{\"combinator\":\"and\",\"rules\":[{\"property\":{\"context\":\"user\",\"property\":\"department\"},\"operator\":\"=\",\"value\":\"engineering\"}]}`. " +
			"`context` is one of `user`, `userTeams`, `form`, `workflowRun`.",
		Optional:   true,
		Validators: []validator.String{queryValidator("Invalid permissions policy")},
	}

	return schema.SingleNestedBlock{
		MarkdownDescription: description,
		Attributes:          attributes,
	}
}

func respondersBlock(description string) schema.Block {
	attributes := principalAttributes("allowed to respond")
	attributes["users_query"] = schema.StringAttribute{
		MarkdownDescription: "A JSON encoded entity search query, run against the `_user` blueprint, " +
			"resolving additional responders.",
		Optional:   true,
		Validators: []validator.String{queryValidator("Invalid responders query")},
	}

	return schema.SingleNestedBlock{
		MarkdownDescription: description,
		Attributes:          attributes,
	}
}

func userInputsAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"user_properties": userPropertiesSchema(),
		"titles": schema.MapNestedAttribute{
			MarkdownDescription: "Static titles rendered between the inputs of the form.",
			Optional:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"title": schema.StringAttribute{
						MarkdownDescription: "The title text.",
						Required:            true,
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "The title description.",
						Optional:            true,
					},
					"visible": schema.BoolAttribute{
						MarkdownDescription: "The visibility of the title.",
						Optional:            true,
					},
					"visible_jq_query": schema.StringAttribute{
						MarkdownDescription: "The visibility condition jq query of the title.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("visible")),
						},
					},
				},
			},
		},
		"required_jq_query": schema.StringAttribute{
			MarkdownDescription: "A jq query resolving which inputs are required.",
			Optional:            true,
		},
		"order_properties": schema.ListAttribute{
			MarkdownDescription: "The order the inputs are rendered in. Cannot be combined with `steps`.",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"steps": schema.ListNestedAttribute{
			MarkdownDescription: "Splits the form into steps. Cannot be combined with `order_properties`.",
			Optional:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"title": schema.StringAttribute{
						MarkdownDescription: "The step's title (max 25 characters).",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtMost(25),
						},
					},
					"order": schema.ListAttribute{
						MarkdownDescription: "The order of the inputs in this step.",
						Required:            true,
						ElementType:         types.StringType,
					},
					"visible": schema.BoolAttribute{
						MarkdownDescription: "The visibility of the step.",
						Optional:            true,
					},
					"visible_jq_query": schema.StringAttribute{
						MarkdownDescription: "The visibility condition jq query of the step.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("visible")),
						},
					},
					"validations": validationsSchema("Validation rules evaluated when the step is submitted."),
				},
			},
			Validators: []validator.List{
				listvalidator.ConflictsWith(
					path.MatchRelative().AtParent().AtName("order_properties"),
					path.MatchRelative().AtParent().AtName("validations"),
				),
			},
		},
		"validations": validationsSchema(
			"Validation rules evaluated against the whole form when it is submitted. " +
				"Cannot be combined with `steps`, add the rules to the individual steps instead.",
		),
	}
}

func validationsSchema(description string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		MarkdownDescription: description + " Up to 10 rules are allowed.",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"constraint": schema.StringAttribute{
					MarkdownDescription: "A jq expression that has to evaluate to `true` for the form to be valid.",
					Required:            true,
				},
				"message": schema.StringAttribute{
					MarkdownDescription: "The error message shown when the constraint evaluates to `false` (max 100 characters).",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.LengthAtMost(100),
					},
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.SizeAtMost(10),
		},
	}
}

func selfServeTriggerBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "A self service trigger node that starts the workflow from a user submitted form.",
		Attributes: map[string]schema.Attribute{
			"action_card_button_text": schema.StringAttribute{
				MarkdownDescription: "The text of the button displayed on the self service card (max 15 characters).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 15),
				},
			},
			"execute_action_button_text": schema.StringAttribute{
				MarkdownDescription: "The text of the button that executes the workflow (max 15 characters).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 15),
				},
			},
			"variant": schema.StringAttribute{
				MarkdownDescription: "The trigger variant. One of `DEFAULT`, `ALERT`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("DEFAULT", "ALERT"),
				},
			},
			"published": publishedAttribute(),
		},
		Blocks: map[string]schema.Block{
			"user_inputs": schema.SingleNestedBlock{
				MarkdownDescription: "The form presented to the user when triggering the workflow.",
				Attributes:          userInputsAttributes(),
			},
			"permissions": permissionsBlock("Who is allowed to execute this trigger."),
			"contexts": schema.ListNestedBlock{
				MarkdownDescription: "Where the trigger is surfaced in the UI.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"on": schema.StringAttribute{
							MarkdownDescription: "The context type. One of `CREATE_ENTITY`, `ENTITY`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("CREATE_ENTITY", "ENTITY"),
							},
						},
						"blueprint_identifier": schema.StringAttribute{
							MarkdownDescription: "The blueprint the trigger creates an entity for. Required when `on` is `CREATE_ENTITY`.",
							Optional:            true,
						},
						"user_input": schema.StringAttribute{
							MarkdownDescription: "The user input the trigger is bound to. Required when `on` is `ENTITY`.",
							Optional:            true,
						},
					},
				},
			},
		},
		Validators: nodeTypeValidators(),
	}
}

func eventTriggerBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "An event trigger node that starts the workflow when an entity event occurs.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				MarkdownDescription: "The event type that triggers the workflow. One of `ENTITY_CREATED`, `ENTITY_UPDATED`, `ENTITY_DELETED`, `TIMER_EXPIRED`, `ANY_ENTITY_CHANGE`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						consts.EntityCreated,
						consts.EntityUpdated,
						consts.EntityDeleted,
						consts.WorkflowTimerExpired,
						consts.AnyEntityChange,
					),
				},
			},
			"blueprint_identifier": schema.StringAttribute{
				MarkdownDescription: "The blueprint identifier the event relates to.",
				Optional:            true,
			},
			"property_identifier": schema.StringAttribute{
				MarkdownDescription: "The property identifier the timer event relates to. Required for the `TIMER_EXPIRED` event type.",
				Optional:            true,
			},
			"next_expire_at": schema.StringAttribute{
				MarkdownDescription: "A JQ expression returning the next expiration time as an ISO date. " +
					"The timer property is set to it after the trigger fires, so the timer keeps running. " +
					"Only valid for the `TIMER_EXPIRED` event type.",
				Optional: true,
			},
			"published": publishedAttribute(),
		},
		Blocks: map[string]schema.Block{
			"condition": schema.SingleNestedBlock{
				MarkdownDescription: "A JQ condition gating whether the event starts the workflow.",
				Attributes: map[string]schema.Attribute{
					"expressions": schema.ListAttribute{
						MarkdownDescription: "The JQ expressions evaluated against the event.",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"combinator": schema.StringAttribute{
						MarkdownDescription: "How the expressions are combined. One of `and`, `or`.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("and", "or"),
						},
					},
				},
			},
		},
		Validators: nodeTypeValidators(),
	}
}

func scheduleTriggerBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "A schedule trigger node that starts the workflow on a cron schedule.",
		Attributes: map[string]schema.Attribute{
			"cron": schema.StringAttribute{
				MarkdownDescription: "The cron expression defining when the workflow triggers (e.g. `0 9 * * 1-5`), evaluated in UTC.",
				Optional:            true,
				Validators:          []validator.String{cronValidator()},
			},
			"published": publishedAttribute(),
		},
		Validators: nodeTypeValidators(),
	}
}

func kafkaBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "A Kafka node that publishes a message to the organization's Kafka topic.",
		Attributes: map[string]schema.Attribute{
			"payload": schema.StringAttribute{
				MarkdownDescription: "The Kafka message payload as a JSON encoded string.",
				Optional:            true,
			},
			"on_failure": onFailureAttribute(),
		},
		Validators: nodeTypeValidators(),
	}
}

func webhookBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "A webhook node that sends an HTTP request.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the webhook.",
				Optional:            true,
			},
			"agent": schema.BoolAttribute{
				MarkdownDescription: "Whether the request is routed through the Port agent.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"synchronized": schema.BoolAttribute{
				MarkdownDescription: "Whether the request is sent synchronously.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"method": schema.StringAttribute{
				MarkdownDescription: "The HTTP method. One of `GET`, `POST`, `PUT`, `PATCH`, `DELETE`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("POST"),
				Validators: []validator.String{
					stringvalidator.OneOf("GET", "POST", "PUT", "PATCH", "DELETE"),
				},
			},
			"headers": schema.MapAttribute{
				MarkdownDescription: "The HTTP headers of the request.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"body": schema.StringAttribute{
				MarkdownDescription: "The request body as a JSON encoded string.",
				Optional:            true,
			},
			"on_timeout": schema.StringAttribute{
				MarkdownDescription: "The action to take if the webhook times out. One of `fail`, `continue`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("fail"),
				Validators: []validator.String{
					stringvalidator.OneOf("fail", "continue"),
				},
			},
			"on_failure": onFailureAttribute(),
		},
		Validators: nodeTypeValidators(),
	}
}

func integrationActionBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "An integration action node that invokes an installed integration.",
		Attributes: map[string]schema.Attribute{
			"installation_id": schema.StringAttribute{
				MarkdownDescription: "The installation id of the integration.",
				Optional:            true,
			},
			"integration_provider": schema.StringAttribute{
				MarkdownDescription: "The provider of the integration action.",
				Optional:            true,
			},
			"integration_invocation_type": schema.StringAttribute{
				MarkdownDescription: "The invocation type of the integration action.",
				Optional:            true,
			},
			"execution_properties": schema.StringAttribute{
				MarkdownDescription: "The integration action execution properties as a JSON encoded string.",
				Optional:            true,
			},
			"on_failure": onFailureAttribute(),
		},
		Validators: nodeTypeValidators(),
	}
}

func upsertEntityBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "An upsert entity node that creates or updates an entity in the catalog.",
		Attributes: map[string]schema.Attribute{
			"blueprint_identifier": schema.StringAttribute{
				MarkdownDescription: "The identifier of the blueprint to upsert into.",
				Optional:            true,
			},
			"on_failure": onFailureAttribute(),
		},
		Blocks: map[string]schema.Block{
			"mapping": schema.SingleNestedBlock{
				MarkdownDescription: "The entity fields to upsert.",
				Attributes: map[string]schema.Attribute{
					"identifier": schema.StringAttribute{
						MarkdownDescription: "The identifier of the entity to upsert.",
						Optional:            true,
					},
					"title": schema.StringAttribute{
						MarkdownDescription: "The title of the entity to upsert.",
						Optional:            true,
					},
					"icon": schema.StringAttribute{
						MarkdownDescription: "The icon of the entity to upsert.",
						Optional:            true,
					},
					"teams": schema.ListAttribute{
						MarkdownDescription: "The teams of the entity to upsert. Values may contain " +
							"`{{ }}` template expressions that are resolved when the node runs.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"properties": schema.StringAttribute{
						MarkdownDescription: "The properties of the entity as a JSON encoded string.",
						Optional:            true,
					},
					"relations": schema.StringAttribute{
						MarkdownDescription: "The relations of the entity as a JSON encoded string.",
						Optional:            true,
					},
				},
			},
		},
		Validators: nodeTypeValidators(),
	}
}

func aiBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "An AI node that runs a prompt through Port AI.",
		Attributes: map[string]schema.Attribute{
			"user_prompt": schema.StringAttribute{
				MarkdownDescription: "The message or query processed by Port AI.",
				Optional:            true,
			},
			"system_prompt": schema.StringAttribute{
				MarkdownDescription: "Instructions describing the AI's role and operational rules.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"provider": schema.StringAttribute{
				MarkdownDescription: "The AI provider to use. Must be set together with `model`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("model")),
				},
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "The AI model to use. Must be set together with `provider`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("provider")),
				},
			},
			"tools": schema.ListAttribute{
				MarkdownDescription: "Regex patterns matched against the available tool names.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"mcp_servers": schema.ListNestedAttribute{
				MarkdownDescription: "The MCP servers available to the AI node (max 5).",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"identifier": schema.StringAttribute{
							MarkdownDescription: "The identifier of the MCP server.",
							Required:            true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(5),
				},
			},
			"output_schema": schema.StringAttribute{
				MarkdownDescription: "A JSON schema, encoded as a JSON string, the AI response is validated against.",
				Optional:            true,
				Validators:          []validator.String{outputSchemaValidator()},
			},
		},
		Validators: nodeTypeValidators(),
	}
}

func aiAgentBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "An AI agent node that invokes a configured Port AI agent.",
		Attributes: map[string]schema.Attribute{
			"user_prompt": schema.StringAttribute{
				MarkdownDescription: "The message or query processed by the agent.",
				Optional:            true,
			},
			"agent_identifier": schema.StringAttribute{
				MarkdownDescription: "The identifier of the agent to invoke.",
				Optional:            true,
			},
			"provider": schema.StringAttribute{
				MarkdownDescription: "The AI provider to use. Must be set together with `model`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("model")),
				},
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "The AI model to use. Must be set together with `provider`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("provider")),
				},
			},
			"output_schema": schema.StringAttribute{
				MarkdownDescription: "A JSON schema, encoded as a JSON string, the agent response is validated against.",
				Optional:            true,
				Validators:          []validator.String{outputSchemaValidator()},
			},
		},
		Validators: nodeTypeValidators(),
	}
}

func conditionBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "A condition node that branches the workflow based on JQ expressions. Connections leaving this node must set `source_outlet_identifier` or `fallback`.",
		Blocks: map[string]schema.Block{
			"outlets": schema.ListNestedBlock{
				MarkdownDescription: "The branches of the condition, evaluated in order.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"identifier": schema.StringAttribute{
							MarkdownDescription: "The identifier of the outlet, referenced by a connection's `source_outlet_identifier`.",
							Required:            true,
							Validators:          identifierValidators(),
						},
						"title": schema.StringAttribute{
							MarkdownDescription: "The title of the outlet.",
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString(""),
							Validators:          titleValidators(),
						},
						"expression": schema.StringAttribute{
							MarkdownDescription: "The JQ expression that selects this outlet.",
							Required:            true,
						},
					},
					Blocks: map[string]schema.Block{
						"status_label":          statusLabelBlock("A custom status label displayed on the node run."),
						"workflow_status_label": statusLabelBlock("A custom status label displayed on the workflow run."),
					},
				},
			},
		},
		Validators: nodeTypeValidators(),
	}
}

func inputBlock() schema.Block {
	return schema.SingleNestedBlock{
		MarkdownDescription: "An input node that pauses the workflow and waits for a human response. Connections leaving this node must set `source_outlet_identifier`.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				MarkdownDescription: "The description shown on the response form.",
				Optional:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"user_inputs": schema.SingleNestedBlock{
				MarkdownDescription: "The form presented to the responders.",
				Attributes: func() map[string]schema.Attribute {
					attributes := userInputsAttributes()
					attributes["buttons"] = schema.ListNestedAttribute{
						MarkdownDescription: "The buttons rendered on the response form. Each outlet must reference one of these identifiers.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"identifier": schema.StringAttribute{
									MarkdownDescription: "The identifier of the button.",
									Required:            true,
									Validators:          identifierValidators(),
								},
								"label": schema.StringAttribute{
									MarkdownDescription: "The label of the button.",
									Required:            true,
								},
								"variant": schema.StringAttribute{
									MarkdownDescription: "The button variant. One of `PRIMARY`, `SECONDARY`, `DANGER`.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("PRIMARY", "SECONDARY", "DANGER"),
									},
								},
								"icon": schema.StringAttribute{
									MarkdownDescription: "The icon of the button.",
									Optional:            true,
								},
							},
						},
					}
					return attributes
				}(),
			},
			"outlets": schema.ListNestedBlock{
				MarkdownDescription: "The branches of the input node, each bound to a button.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"identifier": schema.StringAttribute{
							MarkdownDescription: "The identifier of the outlet. Must match a button identifier.",
							Required:            true,
							Validators:          identifierValidators(),
						},
						"title": schema.StringAttribute{
							MarkdownDescription: "The title of the outlet.",
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString(""),
							Validators:          titleValidators(),
						},
						"num_of_responders": schema.Int64Attribute{
							MarkdownDescription: "How many responders must press the button before the workflow continues.",
							Required:            true,
							Validators: []validator.Int64{
								int64validator.AtLeast(1),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"status_label":          statusLabelBlock("A custom status label displayed on the node run."),
						"workflow_status_label": statusLabelBlock("A custom status label displayed on the workflow run."),
					},
				},
			},
			"notifications": schema.ListNestedBlock{
				MarkdownDescription: "Notifications sent when the input node starts waiting.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"target": schema.StringAttribute{
							MarkdownDescription: "The notification target. One of `email`, `webhook`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("email", "webhook"),
							},
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "The webhook URL. Required when `target` is `webhook`.",
							Optional:            true,
						},
						"method": schema.StringAttribute{
							MarkdownDescription: "The webhook HTTP method. One of `GET`, `POST`, `PUT`, `PATCH`, `DELETE`.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("GET", "POST", "PUT", "PATCH", "DELETE"),
							},
						},
						"headers": schema.MapAttribute{
							MarkdownDescription: "The webhook headers.",
							Optional:            true,
							ElementType:         types.StringType,
						},
						"body": schema.StringAttribute{
							MarkdownDescription: "The webhook body as a JSON encoded string.",
							Optional:            true,
						},
						"agent": schema.BoolAttribute{
							MarkdownDescription: "Whether the webhook is routed through the Port agent.",
							Optional:            true,
						},
					},
					Blocks: map[string]schema.Block{
						"fields": schema.ListNestedBlock{
							MarkdownDescription: "The fields rendered in the email notification. Only valid when `target` is `email`.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"label": schema.StringAttribute{
										MarkdownDescription: "The label of the field.",
										Required:            true,
									},
									"value": schema.StringAttribute{
										MarkdownDescription: "The value of the field.",
										Required:            true,
									},
								},
							},
						},
					},
				},
			},
			"responders": respondersBlock("Who is allowed to respond to this input node."),
		},
		Validators: nodeTypeValidators(),
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
			Validators:          identifierValidators(),
		},
		"title": schema.StringAttribute{
			MarkdownDescription: "The title of the workflow",
			Optional:            true,
			Validators:          titleValidators(),
		},
		"icon": schema.StringAttribute{
			MarkdownDescription: "The icon of the workflow",
			Optional:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the workflow",
			Optional:            true,
			Validators:          descriptionValidators(),
		},
		"category": schema.StringAttribute{
			MarkdownDescription: "A free-form category used to group the workflow (max 40 characters)",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(maxCategoryLength),
			},
		},
		"allow_anyone_to_view_runs": schema.BoolAttribute{
			MarkdownDescription: "Whether any user can view this workflow's runs",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
	}
}

func WorkflowBlocks() map[string]schema.Block {
	return map[string]schema.Block{
		"node": schema.ListNestedBlock{
			MarkdownDescription: "A node of the workflow graph. Exactly one node config block must be set.",
			Validators: []validator.List{
				listvalidator.SizeAtMost(maxNodesCount),
			},
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"identifier": schema.StringAttribute{
						MarkdownDescription: "The identifier of the node",
						Required:            true,
						Validators:          identifierValidators(),
					},
					"title": schema.StringAttribute{
						MarkdownDescription: "The title of the node",
						Optional:            true,
						Validators:          titleValidators(),
					},
					"icon": schema.StringAttribute{
						MarkdownDescription: "The icon of the node",
						Optional:            true,
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "The description of the node",
						Optional:            true,
						Validators:          descriptionValidators(),
					},
					"verbose": schema.BoolAttribute{
						MarkdownDescription: "When true, the workflow service writes extended per-node run logs",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"links": schema.ListAttribute{
						MarkdownDescription: "Link templates (supporting `{{ .result.field }}` interpolation) evaluated when the node run completes (max 3)",
						Optional:            true,
						ElementType:         types.StringType,
						Validators: []validator.List{
							listvalidator.SizeAtMost(3),
						},
					},
					"variables": schema.MapAttribute{
						MarkdownDescription: "Named expressions made available to the node at run time",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
				Blocks: map[string]schema.Block{
					"self_serve_trigger": selfServeTriggerBlock(),
					"event_trigger":      eventTriggerBlock(),
					"schedule_trigger":   scheduleTriggerBlock(),
					"kafka":              kafkaBlock(),
					"webhook":            webhookBlock(),
					"integration_action": integrationActionBlock(),
					"upsert_entity":      upsertEntityBlock(),
					"ai":                 aiBlock(),
					"ai_agent":           aiAgentBlock(),
					"condition":          conditionBlock(),
					"input":              inputBlock(),
				},
			},
		},
		"connections": schema.ListNestedBlock{
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
					"description": schema.StringAttribute{
						MarkdownDescription: "The description of the connection",
						Optional:            true,
					},
					"source_outlet_identifier": schema.StringAttribute{
						MarkdownDescription: "The outlet of the source node this connection leaves from. Required for `condition` and `input` nodes, and not allowed for any other node type.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("fallback")),
						},
					},
					"fallback": schema.BoolAttribute{
						MarkdownDescription: "Marks this connection as the fallback branch of a `condition` node, taken when no outlet matches. Cannot be combined with `source_outlet_identifier`.",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (r *WorkflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Workflow resource for managing Port workflows and their node/connection graph",
		Attributes:          WorkflowSchema(),
		Blocks:              WorkflowBlocks(),
	}
}
