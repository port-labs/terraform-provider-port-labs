package llm_provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LLMProviderModel struct {
	ID                 types.String `tfsdk:"id"`
	ProviderType       types.String `tfsdk:"provider_type"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	ValidateConnection types.Bool   `tfsdk:"validate_connection"`

	AnthropicCompatible *AnthropicCompatibleModel `tfsdk:"anthropic_compatible"`
	OpenAICompatible    *OpenAICompatibleModel    `tfsdk:"openai_compatible"`
	VertexAnthropic     *VertexAnthropicModel     `tfsdk:"vertex_anthropic"`
	VertexGemini        *VertexGeminiModel        `tfsdk:"vertex_gemini"`
	ApiKey              *ApiKeyConfigModel        `tfsdk:"api_key"`
	AzureAnthropic      *AzureAnthropicModel      `tfsdk:"azure_anthropic"`
	AzureOpenAI         *AzureOpenAIConfigModel   `tfsdk:"azure_openai"`
}

type AnthropicCompatibleModel struct {
	BaseUrl             types.String         `tfsdk:"base_url"`
	ApiKeySecretName    types.String         `tfsdk:"api_key_secret_name"`
	AuthTokenSecretName types.String         `tfsdk:"auth_token_secret_name"`
	CustomHeaders       types.Map            `tfsdk:"custom_headers"`
	Models              []CustomModelModel   `tfsdk:"models"`
}

type OpenAICompatibleModel struct {
	BaseUrl                  types.String       `tfsdk:"base_url"`
	ApiKeySecretName         types.String       `tfsdk:"api_key_secret_name"`
	CustomHeaders            types.Map          `tfsdk:"custom_headers"`
	Models                   []CustomModelModel `tfsdk:"models"`
	SupportsStructuredOutputs types.Bool         `tfsdk:"supports_structured_outputs"`
}

type VertexAnthropicModel struct {
	ClientEmailSecretName types.String       `tfsdk:"client_email_secret_name"`
	PrivateKeySecretName  types.String       `tfsdk:"private_key_secret_name"`
	Project               types.String       `tfsdk:"project"`
	Location              types.String       `tfsdk:"location"`
	Models                []CustomModelModel `tfsdk:"models"`
}

type VertexGeminiModel struct {
	ClientEmailSecretName types.String       `tfsdk:"client_email_secret_name"`
	PrivateKeySecretName  types.String       `tfsdk:"private_key_secret_name"`
	Project               types.String       `tfsdk:"project"`
	Location              types.String       `tfsdk:"location"`
	ApiKeySecretName      types.String       `tfsdk:"api_key_secret_name"`
	Models                []CustomModelModel `tfsdk:"models"`
}

type ApiKeyConfigModel struct {
	ApiKeySecretName types.String `tfsdk:"api_key_secret_name"`
}

type AzureAnthropicModel struct {
	ApiKeySecretName types.String `tfsdk:"api_key_secret_name"`
	Endpoint         types.String `tfsdk:"endpoint"`
}

type AzureOpenAIConfigModel struct {
	ApiKeySecretName types.String `tfsdk:"api_key_secret_name"`
	Endpoint         types.String `tfsdk:"endpoint"`
	Deployment       types.String `tfsdk:"deployment"`
}

type CustomModelModel struct {
	Name              types.String            `tfsdk:"name"`
	DisplayName       types.String            `tfsdk:"display_name"`
	ContextWindow     types.Int64             `tfsdk:"context_window"`
	SupportedFeatures *SupportedFeaturesModel `tfsdk:"supported_features"`
}

type SupportedFeaturesModel struct {
	Temperature               types.Bool `tfsdk:"temperature"`
	Caching                   types.Bool `tfsdk:"caching"`
	Guardrails                types.Bool `tfsdk:"guardrails"`
	ReasoningEffort           types.Bool `tfsdk:"reasoning_effort"`
	NativeStructuredOutput    types.Bool `tfsdk:"native_structured_output"`
	ExtendedThinking          types.Bool `tfsdk:"extended_thinking"`
	AdaptiveThinking          types.Bool `tfsdk:"adaptive_thinking"`
	SupportsStructuredOutputs types.Bool `tfsdk:"supports_structured_outputs"`
}
