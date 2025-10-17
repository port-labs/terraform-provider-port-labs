package organization

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func organizationSecretResourceToPortBody(ctx context.Context, state *OrganizationSecretModel) (*cli.OrganizationSecret, error) {
	secret := &cli.OrganizationSecret{
		SecretName: state.SecretName.ValueString(),
	}

	if !state.SecretValue.IsNull() {
		secretValue := state.SecretValue.ValueString()
		secret.SecretValue = &secretValue
	}

	if !state.Description.IsNull() {
		description := state.Description.ValueString()
		secret.Description = &description
	}

	return secret, nil
}
