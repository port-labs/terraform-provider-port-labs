package organization

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OrganizationSecretModel struct {
	ID          types.String `tfsdk:"id"`
	SecretName  types.String `tfsdk:"secret_name"`
	SecretValue types.String `tfsdk:"secret_value"`
	Description types.String `tfsdk:"description"`
}
