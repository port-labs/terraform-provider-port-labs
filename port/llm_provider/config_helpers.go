package llm_provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func customModelsToPortBody(models []CustomModelModel) ([]map[string]any, error) {
	if len(models) == 0 {
		return nil, fmt.Errorf("at least one model must be specified")
	}

	result := make([]map[string]any, 0, len(models))
	for _, model := range models {
		name := model.Name.ValueString()
		if len(name) < 3 {
			return nil, fmt.Errorf("model name must be at least 3 characters")
		}

		entry := map[string]any{
			"name": name,
		}

		if !model.DisplayName.IsNull() && !model.DisplayName.IsUnknown() {
			entry["displayName"] = model.DisplayName.ValueString()
		}

		if !model.ContextWindow.IsNull() && !model.ContextWindow.IsUnknown() {
			if model.ContextWindow.ValueInt64() <= 0 {
				return nil, fmt.Errorf("context_window must be positive")
			}
			entry["contextWindow"] = model.ContextWindow.ValueInt64()
		}

		if features := supportedFeaturesToPortBody(model.SupportedFeatures); features != nil {
			entry["supportedFeatures"] = features
		}

		result = append(result, entry)
	}

	return result, nil
}

func supportedFeaturesToPortBody(features *SupportedFeaturesModel) map[string]bool {
	if features == nil {
		return nil
	}

	result := map[string]bool{}
	setBoolIfKnown(result, "temperature", features.Temperature)
	setBoolIfKnown(result, "caching", features.Caching)
	setBoolIfKnown(result, "guardrails", features.Guardrails)
	setBoolIfKnown(result, "reasoningEffort", features.ReasoningEffort)
	setBoolIfKnown(result, "nativeStructuredOutput", features.NativeStructuredOutput)
	setBoolIfKnown(result, "extendedThinking", features.ExtendedThinking)
	setBoolIfKnown(result, "adaptiveThinking", features.AdaptiveThinking)
	setBoolIfKnown(result, "supportsStructuredOutputs", features.SupportsStructuredOutputs)

	if len(result) == 0 {
		return nil
	}
	return result
}

func setBoolIfKnown(target map[string]bool, key string, value types.Bool) {
	if value.IsNull() || value.IsUnknown() {
		return
	}
	target[key] = value.ValueBool()
}

func setStringIfKnown(target map[string]any, key string, value types.String) {
	if value.IsNull() || value.IsUnknown() {
		return
	}
	target[key] = value.ValueString()
}

func customHeadersToPortBody(headers types.Map) (map[string]string, error) {
	if headers.IsNull() || headers.IsUnknown() {
		return nil, nil
	}

	elements := headers.Elements()
	result := make(map[string]string, len(elements))
	for key, value := range elements {
		strValue, ok := value.(types.String)
		if !ok {
			return nil, fmt.Errorf("custom_headers values must be strings")
		}
		result[key] = strValue.ValueString()
	}
	return result, nil
}

func customModelsFromPort(config map[string]any) []CustomModelModel {
	modelsRaw, ok := config["models"]
	if !ok {
		return nil
	}

	modelsSlice, ok := modelsRaw.([]any)
	if !ok {
		return nil
	}

	models := make([]CustomModelModel, 0, len(modelsSlice))
	for _, item := range modelsSlice {
		modelMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		model := CustomModelModel{}
		if name, ok := modelMap["name"].(string); ok {
			model.Name = types.StringValue(name)
		}
		if displayName, ok := modelMap["displayName"].(string); ok {
			model.DisplayName = types.StringValue(displayName)
		}
		if contextWindow, ok := asInt64(modelMap["contextWindow"]); ok {
			model.ContextWindow = types.Int64Value(contextWindow)
		}
		if featuresRaw, ok := modelMap["supportedFeatures"].(map[string]any); ok {
			model.SupportedFeatures = supportedFeaturesFromPort(featuresRaw)
		}
		models = append(models, model)
	}

	return models
}

func supportedFeaturesFromPort(features map[string]any) *SupportedFeaturesModel {
	if len(features) == 0 {
		return nil
	}

	result := &SupportedFeaturesModel{}
	setFrameworkBoolIfPresent(result, &result.Temperature, features, "temperature")
	setFrameworkBoolIfPresent(result, &result.Caching, features, "caching")
	setFrameworkBoolIfPresent(result, &result.Guardrails, features, "guardrails")
	setFrameworkBoolIfPresent(result, &result.ReasoningEffort, features, "reasoningEffort")
	setFrameworkBoolIfPresent(result, &result.NativeStructuredOutput, features, "nativeStructuredOutput")
	setFrameworkBoolIfPresent(result, &result.ExtendedThinking, features, "extendedThinking")
	setFrameworkBoolIfPresent(result, &result.AdaptiveThinking, features, "adaptiveThinking")
	setFrameworkBoolIfPresent(result, &result.SupportsStructuredOutputs, features, "supportsStructuredOutputs")

	if result.Temperature.IsNull() &&
		result.Caching.IsNull() &&
		result.Guardrails.IsNull() &&
		result.ReasoningEffort.IsNull() &&
		result.NativeStructuredOutput.IsNull() &&
		result.ExtendedThinking.IsNull() &&
		result.AdaptiveThinking.IsNull() &&
		result.SupportsStructuredOutputs.IsNull() {
		return nil
	}

	return result
}

func setFrameworkBoolIfPresent(_ *SupportedFeaturesModel, target *types.Bool, source map[string]any, key string) {
	value, ok := source[key]
	if !ok {
		return
	}
	if boolValue, ok := value.(bool); ok {
		*target = types.BoolValue(boolValue)
	}
}

func asInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int64:
		return v, true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}

func stringFromConfig(config map[string]any, key string) types.String {
	if value, ok := config[key].(string); ok {
		return types.StringValue(value)
	}
	return types.StringNull()
}

func boolFromConfig(config map[string]any, key string) types.Bool {
	if value, ok := config[key].(bool); ok {
		return types.BoolValue(value)
	}
	return types.BoolNull()
}

func customHeadersFromPort(config map[string]any) types.Map {
	headersRaw, ok := config["customHeaders"]
	if !ok {
		return types.MapNull(types.StringType)
	}

	headersMap, ok := headersRaw.(map[string]any)
	if !ok {
		return types.MapNull(types.StringType)
	}

	elements := make(map[string]attr.Value, len(headersMap))
	for key, value := range headersMap {
		if strValue, ok := value.(string); ok {
			elements[key] = types.StringValue(strValue)
		}
	}

	if len(elements) == 0 {
		return types.MapNull(types.StringType)
	}

	result, diags := types.MapValue(types.StringType, elements)
	if diags.HasError() {
		return types.MapNull(types.StringType)
	}
	return result
}
