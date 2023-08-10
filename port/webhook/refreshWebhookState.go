package webhook

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

func refreshWebhookState(ctx context.Context, state *WebhookModel, w *cli.Webhook) error {
	state.ID = types.StringValue(w.Identifier)
	state.Identifier = types.StringValue(w.Identifier)
	state.CreatedAt = types.StringValue(w.CreatedAt.String())
	state.CreatedBy = types.StringValue(w.CreatedBy)
	state.UpdatedAt = types.StringValue(w.UpdatedAt.String())
	state.UpdatedBy = types.StringValue(w.UpdatedBy)

	if w.Icon != nil {
		state.Icon = types.StringValue(*w.Icon)
	}

	if w.Title != nil {
		state.Title = types.StringValue(*w.Title)
	}

	if w.Description != nil {
		state.Description = types.StringValue(*w.Description)
	}

	if w.Enabled != nil {
		state.Enabled = types.BoolValue(*w.Enabled)
	}

	return nil
}
