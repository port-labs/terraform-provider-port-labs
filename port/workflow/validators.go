package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type stringValidator struct {
	summary     string
	description string
	validate    func(string) error
}

func (v stringValidator) Description(context.Context) string { return v.description }

func (v stringValidator) MarkdownDescription(context.Context) string { return v.description }

func (v stringValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if err := v.validate(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, v.summary, err.Error())
	}
}

// isTrueValidator rejects `required = false`. Leaving an input out of the
// required list is how it is marked optional, so an explicit false would be
// silently dropped.
type isTrueValidator struct{}

var _ validator.Bool = isTrueValidator{}

func (v isTrueValidator) Description(context.Context) string { return "value must be true" }

func (v isTrueValidator) MarkdownDescription(ctx context.Context) string { return v.Description(ctx) }

func (v isTrueValidator) ValidateBool(_ context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.ValueBool() {
		return
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid required value",
		"Only `true` is accepted. Leave `required` out to make the input optional.",
	)
}

func decodeJSONObject(value string) (map[string]any, error) {
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

// The service parses the output schema as a JSON Schema object whose `type` is
// `object`, so any other shape is rejected on save.
func outputSchemaValidator() validator.String {
	return stringValidator{
		summary:     "Invalid output schema",
		description: "must be a JSON encoded JSON Schema object whose `type` is `object`",
		validate: func(value string) error {
			object, err := decodeJSONObject(value)
			if err != nil {
				return err
			}
			if object["type"] != "object" {
				return fmt.Errorf("must set `type` to `object`, the only output schema shape the API accepts")
			}
			return nil
		},
	}
}

// Both the RBAC policy and the responders entity query are `{combinator, rules}`
// objects. The rules themselves are a large union that is left to the API, so
// only the envelope is checked here.
func queryValidator(summary string) validator.String {
	return stringValidator{
		summary:     summary,
		description: "must be a JSON encoded object with a `combinator` of `and` or `or` and a `rules` array",
		validate: func(value string) error {
			object, err := decodeJSONObject(value)
			if err != nil {
				return err
			}

			combinator, ok := object["combinator"].(string)
			if !ok || (combinator != "and" && combinator != "or") {
				return fmt.Errorf("must set `combinator` to `and` or `or`")
			}

			if _, ok := object["rules"].([]any); !ok {
				return fmt.Errorf("must set `rules` to an array")
			}
			return nil
		},
	}
}

var (
	cronFieldPattern  = regexp.MustCompile(`^[0-9A-Za-z*?,/#\-]+$`)
	cronDescriptors   = []string{"@yearly", "@annually", "@monthly", "@weekly", "@daily", "@midnight", "@hourly"}
	cronFieldCountMin = 5
	cronFieldCountMax = 6
)

// The service parses the expression with a full cron parser. Rather than
// reimplement it, this checks the shape only, so it catches the common mistakes
// without rejecting anything the parser would accept.
func cronValidator() validator.String {
	return stringValidator{
		summary:     "Invalid cron expression",
		description: fmt.Sprintf("must be a cron expression of %d or %d fields, or one of %s", cronFieldCountMin, cronFieldCountMax, strings.Join(cronDescriptors, ", ")),
		validate: func(value string) error {
			expression := strings.TrimSpace(value)
			if expression == "" {
				return fmt.Errorf("must not be empty")
			}

			if strings.HasPrefix(expression, "@") {
				for _, descriptor := range cronDescriptors {
					if expression == descriptor {
						return nil
					}
				}
				return fmt.Errorf("%q is not a recognized descriptor, expected one of %s", expression, strings.Join(cronDescriptors, ", "))
			}

			fields := strings.Fields(expression)
			if len(fields) < cronFieldCountMin || len(fields) > cronFieldCountMax {
				return fmt.Errorf("expected %d or %d space separated fields, got %d in %q", cronFieldCountMin, cronFieldCountMax, len(fields), expression)
			}

			for _, field := range fields {
				if !cronFieldPattern.MatchString(field) {
					return fmt.Errorf("field %q contains characters that are not valid in a cron expression", field)
				}
			}
			return nil
		},
	}
}
