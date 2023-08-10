package webhook

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

func webhookResourceToBody(ctx context.Context, state *WebhookModel) (*cli.Webhook, error) {
	w := &cli.Webhook{
		Identifier: state.Identifier.ValueString(),
	}

	if !state.Icon.IsNull() {
		icon := state.Icon.ValueString()
		w.Icon = &icon
	}

	if !state.Title.IsNull() {
		title := state.Title.ValueString()
		w.Title = &title
	}

	if !state.Description.IsNull() {
		description := state.Description.ValueString()
		w.Description = &description
	}

	if !state.Enabled.IsNull() {
		enabled := state.Enabled.ValueBool()
		w.Enabled = &enabled
	}

	if state.Security != nil {
		w.Security = &cli.Security{}
		if !state.Security.Secret.IsNull() {
			secret := state.Security.Secret.ValueString()
			w.Security.Secret = &secret
		}
		if !state.Security.SignatureHeaderName.IsNull() {
			signatureHeaderName := state.Security.SignatureHeaderName.ValueString()
			w.Security.SignatureHeaderName = &signatureHeaderName
		}
		if !state.Security.SignatureAlgorithm.IsNull() {
			signatureAlgorithm := state.Security.SignatureAlgorithm.ValueString()
			w.Security.SignatureAlgorithm = &signatureAlgorithm
		}
		if !state.Security.SignaturePrefix.IsNull() {
			signaturePrefix := state.Security.SignaturePrefix.ValueString()
			w.Security.SignaturePrefix = &signaturePrefix
		}

		if !state.Security.RequestIdentifierPath.IsNull() {
			requestIdentifierPath := state.Security.RequestIdentifierPath.ValueString()
			w.Security.RequestIdentifierPath = &requestIdentifierPath
		}

	}

	return w, nil
}
