package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// JSONObjectToGo decodes a JSON encoded string that must hold an object.
func JSONObjectToGo(s string) (map[string]any, error) {
	var object map[string]any
	if err := json.Unmarshal([]byte(s), &object); err != nil {
		return nil, fmt.Errorf("must be a JSON encoded object: %w", err)
	}
	if object == nil {
		return nil, fmt.Errorf("must be a JSON encoded object, got null")
	}
	return object, nil
}

type jsonObjectValidator struct{}

// JSONObject validates at plan time that a string attribute holds a JSON object.
// Unknown and null values are skipped.
func JSONObject() validator.String { return jsonObjectValidator{} }

func (jsonObjectValidator) Description(context.Context) string {
	return "value must be a JSON encoded object"
}

func (v jsonObjectValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (jsonObjectValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if _, err := JSONObjectToGo(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid JSON object", err.Error())
	}
}

// PropertyToTerraformString encodes an object-typed dataset rule property for
// state, keeping the prior value when it is semantically equal so formatting
// and HTML escaping never cause drift.
func PropertyToTerraformString(prior types.String, property any, jsonEscapeHTML bool) (types.String, error) {
	return GoObjectToTerraformStringPreferExisting(prior, property, jsonEscapeHTML)
}
