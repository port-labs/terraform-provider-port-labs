package llm_provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func llmProviderToPortBody(ctx context.Context, state *LLMProviderModel) (*cli.LLMProviderUpsert, error) {
	providerName := state.ProviderType.ValueString()
	config, err := configToPortBody(providerName, state)
	if err != nil {
		return nil, err
	}

	enabled := true
	if !state.Enabled.IsNull() && !state.Enabled.IsUnknown() {
		enabled = state.Enabled.ValueBool()
	}

	return &cli.LLMProviderUpsert{
		Provider: providerName,
		Enabled:  enabled,
		Config:   config,
	}, nil
}

func configToPortBody(providerName string, state *LLMProviderModel) (map[string]any, error) {
	switch providerName {
	case "anthropic-compatible":
		return anthropicCompatibleToPortBody(state.AnthropicCompatible)
	case "openai-compatible":
		return openAICompatibleToPortBody(state.OpenAICompatible)
	case "vertex-anthropic":
		return vertexAnthropicToPortBody(state.VertexAnthropic)
	case "vertex-gemini":
		return vertexGeminiToPortBody(state.VertexGemini)
	case "openai", "anthropic":
		return apiKeyToPortBody(state.ApiKey)
	case "azure-anthropic":
		return azureAnthropicToPortBody(state.AzureAnthropic)
	case "azure-openai":
		return azureOpenAIToPortBody(state.AzureOpenAI)
	default:
		return nil, fmt.Errorf("unsupported provider %q", providerName)
	}
}

func anthropicCompatibleToPortBody(config *AnthropicCompatibleModel) (map[string]any, error) {
	if config == nil {
		return nil, fmt.Errorf("anthropic_compatible config is required for provider anthropic-compatible")
	}

	apiKeySet := !config.ApiKeySecretName.IsNull() && !config.ApiKeySecretName.IsUnknown()
	authTokenSet := !config.AuthTokenSecretName.IsNull() && !config.AuthTokenSecretName.IsUnknown()
	if apiKeySet == authTokenSet {
		return nil, fmt.Errorf("provide either api_key_secret_name or auth_token_secret_name for anthropic_compatible, but not both")
	}

	if config.BaseUrl.IsNull() || config.BaseUrl.IsUnknown() {
		return nil, fmt.Errorf("base_url is required for anthropic_compatible")
	}

	result := map[string]any{
		"baseUrl": config.BaseUrl.ValueString(),
	}

	if apiKeySet {
		result["apiKeySecretName"] = config.ApiKeySecretName.ValueString()
	} else {
		result["authTokenSecretName"] = config.AuthTokenSecretName.ValueString()
	}

	headers, err := customHeadersToPortBody(config.CustomHeaders)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		result["customHeaders"] = headers
	}

	models, err := customModelsToPortBody(config.Models)
	if err != nil {
		return nil, err
	}
	result["models"] = models

	return result, nil
}

func openAICompatibleToPortBody(config *OpenAICompatibleModel) (map[string]any, error) {
	if config == nil {
		return nil, fmt.Errorf("openai_compatible config is required for provider openai-compatible")
	}

	if config.BaseUrl.IsNull() || config.BaseUrl.IsUnknown() {
		return nil, fmt.Errorf("base_url is required for openai_compatible")
	}

	result := map[string]any{
		"baseUrl": config.BaseUrl.ValueString(),
	}

	setStringIfKnown(result, "apiKeySecretName", config.ApiKeySecretName)

	headers, err := customHeadersToPortBody(config.CustomHeaders)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		result["customHeaders"] = headers
	}

	models, err := customModelsToPortBody(config.Models)
	if err != nil {
		return nil, err
	}
	result["models"] = models

	if !config.SupportsStructuredOutputs.IsNull() && !config.SupportsStructuredOutputs.IsUnknown() {
		result["supportsStructuredOutputs"] = config.SupportsStructuredOutputs.ValueBool()
	}

	return result, nil
}

func vertexAnthropicToPortBody(config *VertexAnthropicModel) (map[string]any, error) {
	if config == nil {
		return nil, fmt.Errorf("vertex_anthropic config is required for provider vertex-anthropic")
	}

	result := map[string]any{}
	setStringIfKnown(result, "clientEmailSecretName", config.ClientEmailSecretName)
	setStringIfKnown(result, "privateKeySecretName", config.PrivateKeySecretName)
	setStringIfKnown(result, "project", config.Project)
	setStringIfKnown(result, "location", config.Location)

	models, err := customModelsToPortBody(config.Models)
	if err != nil {
		return nil, err
	}
	result["models"] = models

	return result, nil
}

func vertexGeminiToPortBody(config *VertexGeminiModel) (map[string]any, error) {
	if config == nil {
		return nil, fmt.Errorf("vertex_gemini config is required for provider vertex-gemini")
	}

	result := map[string]any{}
	setStringIfKnown(result, "clientEmailSecretName", config.ClientEmailSecretName)
	setStringIfKnown(result, "privateKeySecretName", config.PrivateKeySecretName)
	setStringIfKnown(result, "project", config.Project)
	setStringIfKnown(result, "location", config.Location)
	setStringIfKnown(result, "apiKeySecretName", config.ApiKeySecretName)

	models, err := customModelsToPortBody(config.Models)
	if err != nil {
		return nil, err
	}
	result["models"] = models

	return result, nil
}

func apiKeyToPortBody(config *ApiKeyConfigModel) (map[string]any, error) {
	if config == nil {
		return nil, fmt.Errorf("api_key config is required for this provider")
	}

	if config.ApiKeySecretName.IsNull() || config.ApiKeySecretName.IsUnknown() {
		return nil, fmt.Errorf("api_key_secret_name is required")
	}

	return map[string]any{
		"apiKeySecretName": config.ApiKeySecretName.ValueString(),
	}, nil
}

func azureAnthropicToPortBody(config *AzureAnthropicModel) (map[string]any, error) {
	if config == nil {
		return nil, fmt.Errorf("azure_anthropic config is required for provider azure-anthropic")
	}

	result := map[string]any{}
	setStringIfKnown(result, "apiKeySecretName", config.ApiKeySecretName)
	setStringIfKnown(result, "endpoint", config.Endpoint)

	return result, nil
}

func azureOpenAIToPortBody(config *AzureOpenAIConfigModel) (map[string]any, error) {
	if config == nil {
		return nil, fmt.Errorf("azure_openai config is required for provider azure-openai")
	}

	result := map[string]any{}
	setStringIfKnown(result, "apiKeySecretName", config.ApiKeySecretName)
	setStringIfKnown(result, "endpoint", config.Endpoint)
	setStringIfKnown(result, "deployment", config.Deployment)

	return result, nil
}

func refreshLLMProviderState(ctx context.Context, state *LLMProviderModel, provider *cli.LLMProvider) error {
	state.ID = types.StringValue(provider.Provider)
	state.ProviderType = types.StringValue(provider.Provider)
	state.Enabled = types.BoolValue(provider.Enabled)

	config := provider.Config
	if config == nil {
		return nil
	}

	switch provider.Provider {
	case "anthropic-compatible":
		state.AnthropicCompatible = anthropicCompatibleFromPort(config)
	case "openai-compatible":
		state.OpenAICompatible = openAICompatibleFromPort(config)
	case "vertex-anthropic":
		state.VertexAnthropic = vertexAnthropicFromPort(config)
	case "vertex-gemini":
		state.VertexGemini = vertexGeminiFromPort(config)
	case "openai", "anthropic":
		state.ApiKey = apiKeyFromPort(config)
	case "azure-anthropic":
		state.AzureAnthropic = azureAnthropicFromPort(config)
	case "azure-openai":
		state.AzureOpenAI = azureOpenAIFromPort(config)
	}

	return nil
}

func anthropicCompatibleFromPort(config map[string]any) *AnthropicCompatibleModel {
	result := &AnthropicCompatibleModel{
		BaseUrl:       stringFromConfig(config, "baseUrl"),
		ApiKeySecretName:    stringFromConfig(config, "apiKeySecretName"),
		AuthTokenSecretName: stringFromConfig(config, "authTokenSecretName"),
		CustomHeaders:       customHeadersFromPort(config),
		Models:              customModelsFromPort(config),
	}
	return result
}

func openAICompatibleFromPort(config map[string]any) *OpenAICompatibleModel {
	return &OpenAICompatibleModel{
		BaseUrl:                   stringFromConfig(config, "baseUrl"),
		ApiKeySecretName:          stringFromConfig(config, "apiKeySecretName"),
		CustomHeaders:             customHeadersFromPort(config),
		Models:                    customModelsFromPort(config),
		SupportsStructuredOutputs: boolFromConfig(config, "supportsStructuredOutputs"),
	}
}

func vertexAnthropicFromPort(config map[string]any) *VertexAnthropicModel {
	return &VertexAnthropicModel{
		ClientEmailSecretName: stringFromConfig(config, "clientEmailSecretName"),
		PrivateKeySecretName:  stringFromConfig(config, "privateKeySecretName"),
		Project:               stringFromConfig(config, "project"),
		Location:              stringFromConfig(config, "location"),
		Models:                customModelsFromPort(config),
	}
}

func vertexGeminiFromPort(config map[string]any) *VertexGeminiModel {
	return &VertexGeminiModel{
		ClientEmailSecretName: stringFromConfig(config, "clientEmailSecretName"),
		PrivateKeySecretName:  stringFromConfig(config, "privateKeySecretName"),
		Project:               stringFromConfig(config, "project"),
		Location:              stringFromConfig(config, "location"),
		ApiKeySecretName:      stringFromConfig(config, "apiKeySecretName"),
		Models:                customModelsFromPort(config),
	}
}

func apiKeyFromPort(config map[string]any) *ApiKeyConfigModel {
	return &ApiKeyConfigModel{
		ApiKeySecretName: stringFromConfig(config, "apiKeySecretName"),
	}
}

func azureAnthropicFromPort(config map[string]any) *AzureAnthropicModel {
	return &AzureAnthropicModel{
		ApiKeySecretName: stringFromConfig(config, "apiKeySecretName"),
		Endpoint:         stringFromConfig(config, "endpoint"),
	}
}

func azureOpenAIFromPort(config map[string]any) *AzureOpenAIConfigModel {
	return &AzureOpenAIConfigModel{
		ApiKeySecretName: stringFromConfig(config, "apiKeySecretName"),
		Endpoint:         stringFromConfig(config, "endpoint"),
		Deployment:       stringFromConfig(config, "deployment"),
	}
}
