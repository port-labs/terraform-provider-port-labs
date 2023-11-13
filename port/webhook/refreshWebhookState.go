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
	state.Url = types.StringValue(w.Url)
	state.WebhookKey = types.StringValue(w.WebhookKey)
	state.Icon = types.StringPointerValue(w.Icon)
	state.Title = types.StringPointerValue(w.Title)
	state.Description = types.StringPointerValue(w.Description)
	state.Enabled = types.BoolPointerValue(w.Enabled)

	if w.Security.RequestIdentifierPath != nil || w.Security.Secret != nil || w.Security.SignatureHeaderName != nil || w.Security.SignatureAlgorithm != nil || w.Security.SignaturePrefix != nil {
		state.Security = &SecurityModel{
			Secret:                types.StringPointerValue(w.Security.Secret),
			SignatureHeaderName:   types.StringPointerValue(w.Security.SignatureHeaderName),
			SignatureAlgorithm:    types.StringPointerValue(w.Security.SignatureAlgorithm),
			SignaturePrefix:       types.StringPointerValue(w.Security.SignaturePrefix),
			RequestIdentifierPath: types.StringPointerValue(w.Security.RequestIdentifierPath),
		}
	}

	if len(w.Mappings) > 0 {
		state.Mappings = []MappingsModel{}
		for _, v := range w.Mappings {
			mapping := &MappingsModel{
				Blueprint: types.StringValue(v.Blueprint),
				Entity: &EntityModel{
					Identifier: types.StringValue(v.Entity.Identifier),
				},
			}

			mapping.Filter = types.StringPointerValue(v.Filter)
			mapping.ItemsToParse = types.StringPointerValue(v.ItemsToParse)
			mapping.Entity.Icon = types.StringPointerValue(v.Entity.Icon)
			mapping.Entity.Title = types.StringPointerValue(v.Entity.Title)
			mapping.Entity.Team = types.StringPointerValue(v.Entity.Team)

			if v.Entity.Properties != nil {
				mapping.Entity.Properties = map[string]string{}
				for k, v := range v.Entity.Properties {
					mapping.Entity.Properties[k] = v
				}
			}

			if v.Entity.Relations != nil {
				mapping.Entity.Relations = map[string]string{}
				for k, v := range v.Entity.Relations {
					mapping.Entity.Relations[k] = v
				}
			}
			state.Mappings = append(state.Mappings, *mapping)
		}
	}

	return nil
}
