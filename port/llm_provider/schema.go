package llm_provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func supportedFeaturesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"temperature": schema.BoolAttribute{
			Optional: true,
		},
		"caching": schema.BoolAttribute{
			Optional: true,
		},
		"guardrails": schema.BoolAttribute{
			Optional: true,
		},
		"reasoning_effort": schema.BoolAttribute{
			Optional: true,
		},
		"native_structured_output": schema.BoolAttribute{
			MarkdownDescription: "Whether this model accepts Anthropic's native `output_format` for structured output. Leave unset to let Port infer it from the model identifier; set to `false` when the identifier resembles a Claude model but the endpoint does not support `output_format`.",
			Optional:            true,
		},
		"extended_thinking": schema.BoolAttribute{
			Optional: true,
		},
		"adaptive_thinking": schema.BoolAttribute{
			Optional: true,
		},
		"supports_structured_outputs": schema.BoolAttribute{
			Optional: true,
		},
	}
}

func customModelsBlock() schema.Block {
	return schema.ListNestedBlock{
		MarkdownDescription: "Model catalog for this provider. Each `name` is the model id used in the `model` field on invocations.",
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
		},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					MarkdownDescription: "Model identifier as exposed by the provider.",
					Required:            true,
				},
				"display_name": schema.StringAttribute{
					Optional: true,
				},
				"context_window": schema.Int64Attribute{
					Optional: true,
				},
			},
			Blocks: map[string]schema.Block{
				"supported_features": schema.SingleNestedBlock{
					Attributes: supportedFeaturesSchema(),
				},
			},
		},
	}
}

func LLMProviderSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The LLM provider type (same as `provider_type`).",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"provider_type": schema.StringAttribute{
			MarkdownDescription: "The LLM provider type. Supported values include `openai`, `anthropic`, `azure-openai`, `azure-anthropic`, `openai-compatible`, `anthropic-compatible`, `vertex-gemini`, and `vertex-anthropic`.",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(
					"openai",
					"anthropic",
					"azure-openai",
					"azure-anthropic",
					"openai-compatible",
					"anthropic-compatible",
					"vertex-gemini",
					"vertex-anthropic",
				),
			},
		},
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "Whether this provider is enabled.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
		"validate_connection": schema.BoolAttribute{
			MarkdownDescription: "When true, Port validates the provider configuration on create and update.",
			Optional:            true,
		},
	}
}

func LLMProviderBlocks() map[string]schema.Block {
	return map[string]schema.Block{
		"anthropic_compatible": schema.SingleNestedBlock{
			MarkdownDescription: "Configuration for an Anthropic Messages API-compatible endpoint.",
			Attributes: map[string]schema.Attribute{
				"base_url": schema.StringAttribute{
					MarkdownDescription: "Base URL of the Anthropic Messages API endpoint — the prefix before `/messages`.",
					Required:            true,
				},
				"api_key_secret_name": schema.StringAttribute{
					MarkdownDescription: "Name of a Port secret holding the key sent as the `x-api-key` header.",
					Optional:            true,
				},
				"auth_token_secret_name": schema.StringAttribute{
					MarkdownDescription: "Name of a Port secret holding the token sent as the `Authorization: Bearer` header.",
					Optional:            true,
				},
				"custom_headers": schema.MapAttribute{
					MarkdownDescription: "Optional HTTP headers to add to every request to this provider.",
					Optional:            true,
					ElementType:         types.StringType,
				},
			},
			Blocks: map[string]schema.Block{
				"models": customModelsBlock(),
			},
		},
		"openai_compatible": schema.SingleNestedBlock{
			MarkdownDescription: "Configuration for an OpenAI chat completions-compatible endpoint.",
			Attributes: map[string]schema.Attribute{
				"base_url": schema.StringAttribute{
					MarkdownDescription: "Base URL of the OpenAI-compatible API endpoint.",
					Required:            true,
				},
				"api_key_secret_name": schema.StringAttribute{
					Optional: true,
				},
				"custom_headers": schema.MapAttribute{
					Optional:    true,
					ElementType: types.StringType,
				},
				"supports_structured_outputs": schema.BoolAttribute{
					Optional: true,
				},
			},
			Blocks: map[string]schema.Block{
				"models": customModelsBlock(),
			},
		},
		"vertex_anthropic": schema.SingleNestedBlock{
			MarkdownDescription: "Configuration for Claude models hosted on Google Vertex AI.",
			Attributes: map[string]schema.Attribute{
				"client_email_secret_name": schema.StringAttribute{
					Optional: true,
				},
				"private_key_secret_name": schema.StringAttribute{
					Optional: true,
				},
				"project": schema.StringAttribute{
					Optional: true,
				},
				"location": schema.StringAttribute{
					Optional: true,
				},
			},
			Blocks: map[string]schema.Block{
				"models": customModelsBlock(),
			},
		},
		"vertex_gemini": schema.SingleNestedBlock{
			MarkdownDescription: "Configuration for Gemini models hosted on Google Vertex AI.",
			Attributes: map[string]schema.Attribute{
				"client_email_secret_name": schema.StringAttribute{
					Optional: true,
				},
				"private_key_secret_name": schema.StringAttribute{
					Optional: true,
				},
				"project": schema.StringAttribute{
					Optional: true,
				},
				"location": schema.StringAttribute{
					Optional: true,
				},
				"api_key_secret_name": schema.StringAttribute{
					MarkdownDescription: "Vertex AI Express Mode API key secret name.",
					Optional:            true,
				},
			},
			Blocks: map[string]schema.Block{
				"models": customModelsBlock(),
			},
		},
		"api_key": schema.SingleNestedBlock{
			MarkdownDescription: "Simple API key configuration for built-in providers such as `openai` and `anthropic`.",
			Attributes: map[string]schema.Attribute{
				"api_key_secret_name": schema.StringAttribute{
					Required: true,
				},
			},
		},
		"azure_anthropic": schema.SingleNestedBlock{
			MarkdownDescription: "Configuration for Claude models hosted in Azure.",
			Attributes: map[string]schema.Attribute{
				"api_key_secret_name": schema.StringAttribute{
					Optional: true,
				},
				"endpoint": schema.StringAttribute{
					Optional: true,
				},
			},
		},
		"azure_openai": schema.SingleNestedBlock{
			MarkdownDescription: "Configuration for OpenAI models hosted in Azure.",
			Attributes: map[string]schema.Attribute{
				"api_key_secret_name": schema.StringAttribute{
					Optional: true,
				},
				"endpoint": schema.StringAttribute{
					Optional: true,
				},
				"deployment": schema.StringAttribute{
					Optional: true,
				},
			},
		},
	}
}

var ResourceMarkdownDescription = `
# LLM Provider resource

Register or update an LLM provider configuration for Port AI.

See the [LLM provider management documentation](https://docs.port.io/port-ai/llm-providers-management/overview/) and the [Create or connect an LLM provider API](https://docs.port.io/api-reference/create-or-connect-an-llm-provider/).

## Example: Anthropic-compatible provider

` + "```hcl" + `
resource "port_llm_provider" "anthropic_compatible" {
  provider_type         = "anthropic-compatible"
  enabled             = true
  validate_connection = true

  anthropic_compatible {
    base_url            = "https://litellm.corp.example.com/anthropic"
    api_key_secret_name = "LITELLM_API_KEY"

    models {
      name           = "claude-sonnet"
      display_name   = "Claude Sonnet"
      context_window = 200000

      supported_features {
        caching                  = true
        native_structured_output = false
      }
    }
  }
}
` + "\n```"

func (r *LLMProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: ResourceMarkdownDescription,
		Attributes:          LLMProviderSchema(),
		Blocks:              LLMProviderBlocks(),
	}
}
