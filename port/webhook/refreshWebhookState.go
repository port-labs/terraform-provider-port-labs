package webhook

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
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
	state.Icon = flex.GoStringToFramework(w.Icon)
	state.Title = flex.GoStringToFramework(w.Title)
	state.Description = flex.GoStringToFramework(w.Description)
	state.Enabled = flex.GoBoolToFramework(w.Enabled)

	if w.Security.RequestIdentifierPath != nil || w.Security.Secret != nil || w.Security.SignatureHeaderName != nil || w.Security.SignatureAlgorithm != nil || w.Security.SignaturePrefix != nil {
		state.Security = &SecurityModel{
			Secret:                flex.GoStringToFramework(w.Security.Secret),
			SignatureHeaderName:   flex.GoStringToFramework(w.Security.SignatureHeaderName),
			SignatureAlgorithm:    flex.GoStringToFramework(w.Security.SignatureAlgorithm),
			SignaturePrefix:       flex.GoStringToFramework(w.Security.SignaturePrefix),
			RequestIdentifierPath: flex.GoStringToFramework(w.Security.RequestIdentifierPath),
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

			mapping.Filter = flex.GoStringToFramework(v.Filter)
			mapping.ItemsToParse = flex.GoStringToFramework(v.ItemsToParse)
			mapping.Entity.Icon = flex.GoStringToFramework(v.Entity.Icon)
			mapping.Entity.Title = flex.GoStringToFramework(v.Entity.Title)
			mapping.Entity.Team = flex.GoStringToFramework(v.Entity.Team)

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
