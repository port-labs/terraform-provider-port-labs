package scorecard_group

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func propertiesToCLI(ctx context.Context, properties types.Map) (map[string]any, error) {
	if properties.IsNull() || properties.IsUnknown() {
		return nil, nil
	}

	stringMap := make(map[string]string, len(properties.Elements()))
	if diags := properties.ElementsAs(ctx, &stringMap, false); diags.HasError() {
		return nil, fmt.Errorf("failed to convert properties map: %s", diags.Errors()[0].Detail())
	}

	result := make(map[string]any, len(stringMap))
	for key, value := range stringMap {
		result[key] = terraformStringToAny(value)
	}

	return result, nil
}

func terraformStringToAny(value string) any {
	var parsed any
	if err := json.Unmarshal([]byte(value), &parsed); err == nil {
		return parsed
	}

	return value
}

func propertiesFromCLI(properties map[string]any, jsonEscapeHTML bool) types.Map {
	if len(properties) == 0 {
		return types.MapNull(types.StringType)
	}

	elements := make(map[string]attr.Value, len(properties))
	for key, value := range properties {
		switch typedValue := value.(type) {
		case string:
			elements[key] = types.StringValue(typedValue)
		default:
			encoded, _ := utils.GoObjectToTerraformString(value, jsonEscapeHTML)
			elements[key] = encoded
		}
	}

	result, diags := types.MapValue(types.StringType, elements)
	if diags.HasError() {
		return types.MapNull(types.StringType)
	}

	return result
}
