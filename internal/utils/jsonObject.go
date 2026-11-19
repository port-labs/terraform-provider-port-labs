package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// DecodeJSONObject decodes a JSON encoded string that must hold an object.
func DecodeJSONObject(value string) (map[string]any, error) {
	var decoded any
	if err := json.Unmarshal([]byte(value), &decoded); err != nil {
		return nil, fmt.Errorf("must be valid JSON: %s", err)
	}
	object, ok := decoded.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("must be a JSON object")
	}
	return object, nil
}

type jsonObjectValidator struct{}

// JSONObject validates at plan time that a string attribute holds a JSON object.
// Unknown and null values are skipped.
func JSONObject() validator.String { return jsonObjectValidator{} }

func (jsonObjectValidator) Description(context.Context) string {
	return "must be a JSON encoded object"
}

func (v jsonObjectValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (jsonObjectValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if _, err := DecodeJSONObject(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid JSON object", err.Error())
	}
}

type notJSONObjectValidator struct{ useInstead string }

// NotJSONObject rejects a string that looks like a JSON object, which is
// never a valid identifier, and points the user at the attribute to use.
func NotJSONObject(useInstead string) validator.String {
	return notJSONObjectValidator{useInstead: useInstead}
}

func (v notJSONObjectValidator) Description(context.Context) string {
	return fmt.Sprintf("value must be an identifier, not a JSON object; use `%s` for object values", v.useInstead)
}

func (v notJSONObjectValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v notJSONObjectValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if strings.HasPrefix(strings.TrimSpace(req.ConfigValue.ValueString()), "{") {
		resp.Diagnostics.AddAttributeError(req.Path, "Property must be an identifier",
			fmt.Sprintf("This attribute takes a property identifier, but the value looks like a JSON object. Use `%s` for contextual or path rules.", v.useInstead))
	}
}
