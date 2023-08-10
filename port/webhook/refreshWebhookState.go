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

	if w.Security != nil {
		state.Security = &SecurityModel{}
		if w.Security.Secret != nil {
			state.Security.Secret = types.StringValue(*w.Security.Secret)
		}
		if w.Security.SignatureHeaderName != nil {
			state.Security.SignatureHeaderName = types.StringValue(*w.Security.SignatureHeaderName)
		}
		if w.Security.SignatureAlgorithm != nil {
			state.Security.SignatureAlgorithm = types.StringValue(*w.Security.SignatureAlgorithm)
		}
		if w.Security.SignaturePrefix != nil {
			state.Security.SignaturePrefix = types.StringValue(*w.Security.SignaturePrefix)
		}
		if w.Security.RequestIdentifierPath != nil {
			state.Security.RequestIdentifierPath = types.StringValue(*w.Security.RequestIdentifierPath)
		}
	}

	if w.Mappings != nil {
		state.Mappings = []MappingsModel{}
		for _, v := range w.Mappings {
			mapping := MappingsModel{
				Blueprint: types.StringValue(v.Blueprint),
				Entity: &EntityModel{
					Identifier: types.StringValue(v.Entity.Identifier),
				},
			}

			if v.Filter != nil {
				mapping.Filter = types.StringValue(*v.Filter)
			}

			if v.ItemsToParse != nil {
				mapping.ItemsToParse = types.StringValue(*v.ItemsToParse)
			}

			if v.Entity.Icon != nil {
				mapping.Entity.Icon = types.StringValue(*v.Entity.Icon)
			}

			if v.Entity.Title != nil {
				mapping.Entity.Title = types.StringValue(*v.Entity.Title)
			}

			if v.Entity.Team != nil {
				mapping.Entity.Team = types.StringValue(*v.Entity.Team)
			}

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

			state.Mappings = append(state.Mappings, mapping)
		}
	}

	return nil
}
