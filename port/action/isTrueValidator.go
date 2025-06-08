package action

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type isTrueValidator struct{}

var _ validator.Bool = isTrueValidator{}

// Description describes the validation in plain text formatting.
func (v isTrueValidator) Description(ctx context.Context) string {
	return "Value must be true"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v isTrueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateBool performs the validation.
func (v isTrueValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue
	if value.Equal(types.BoolValue(true)) {
		return
	}
	resp.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
		req.Path,
		v.Description(ctx),
		value.String(),
	))

}
