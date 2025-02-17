package webhook

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func webhookResourceToPortBody(ctx context.Context, state *WebhookModel) (*cli.Webhook, error) {
	w := &cli.Webhook{
		Identifier: state.Identifier.ValueString(),
		Security:   &cli.Security{},
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

	if len(state.Mappings) > 0 {
		w.Mappings = []cli.Mappings{}
		for _, v := range state.Mappings {
			mapping := cli.Mappings{
				Blueprint: v.Blueprint.ValueString(),
				Entity: &cli.EntityProperty{
					Identifier: v.Entity.Identifier.ValueString(),
				},
			}

			if !v.Filter.IsNull() {
				filter := v.Filter.ValueString()
				mapping.Filter = &filter
			}
			if !v.Operation.IsNull() {
				operation := v.Operation.ValueString()
				mapping.Operation = &operation
			}

			if !v.ItemsToParse.IsNull() {
				ItemsToParse := v.ItemsToParse.ValueString()
				mapping.ItemsToParse = &ItemsToParse
			}

			if !v.Entity.Icon.IsNull() {
				icon := v.Entity.Icon.ValueString()
				mapping.Entity.Icon = &icon
			}

			if !v.Entity.Title.IsNull() {
				title := v.Entity.Title.ValueString()
				mapping.Entity.Title = &title
			}

			if !v.Entity.Team.IsNull() {
				team := v.Entity.Team.ValueString()
				mapping.Entity.Team = &team
			}

			if v.Entity.Properties != nil {
				properties := make(map[string]string)
				for k, v := range v.Entity.Properties {
					properties[k] = v
				}
				mapping.Entity.Properties = properties
			}

			if v.Entity.Relations != nil {
				relations := make(map[string]string)
				for k, v := range v.Entity.Relations {
					relations[k] = v
				}
				mapping.Entity.Relations = relations
			}

			w.Mappings = append(w.Mappings, mapping)
		}
	}

	return w, nil
}
