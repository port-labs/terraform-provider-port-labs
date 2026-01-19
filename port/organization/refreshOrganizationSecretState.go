package organization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshOrganizationSecretState(ctx context.Context, state *OrganizationSecretModel, secret *cli.OrganizationSecret) error {
	state.ID = types.StringValue(secret.SecretName)
	state.SecretName = types.StringValue(secret.SecretName)

	if secret.Description != nil && *secret.Description != "" {
		state.Description = types.StringValue(*secret.Description)
	} else {
		state.Description = types.StringNull()
	}

	// Note: SecretValue is not returned by the API for security reasons
	// We keep the existing state value
	return nil
}
