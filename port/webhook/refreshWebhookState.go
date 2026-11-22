package webhook

import (
	"context"
	"encoding/json"
	"fmt"

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

	if security := securityFromAPI(state.Security, w.Security); security != nil {
		state.Security = security
	}

	if len(w.Mappings) > 0 {
		state.Mappings = []MappingsModel{}
		for _, v := range w.Mappings {
			mapping := &MappingsModel{
				Blueprint: types.StringValue(v.Blueprint),
				Entity:    &EntityModel{},
			}

			switch identifier := v.Entity.Identifier.(type) {
			case string:
				mapping.Entity.Identifier = types.StringValue(identifier)
			case map[string]interface{}:
				if jsonBytes, err := json.Marshal(identifier); err == nil {
					mapping.Entity.Identifier = types.StringValue(string(jsonBytes))
				} else {
					return fmt.Errorf("failed to marshal identifier to JSON: %w", err)
				}
			default:
				return fmt.Errorf("invalid identifier type: expected string or object, got %T", identifier)
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
				for k, relationValue := range v.Entity.Relations {
					switch val := relationValue.(type) {
					case string:
						mapping.Entity.Relations[k] = val
					case map[string]interface{}:
						if jsonBytes, err := json.Marshal(val); err == nil {
							mapping.Entity.Relations[k] = string(jsonBytes)
						} else {
							return fmt.Errorf("failed to marshal relation '%s' to JSON: %w", k, err)
						}
					default:
						return fmt.Errorf("invalid relation type for key '%s': expected string or object, got %T", k, val)
					}
				}
			}

			var operationModel OperationModel
			if v.Operation != nil {
				switch operation := v.Operation.(type) {
				// If the operation is a simple string.
				case string:
					operationModel.Type = types.StringValue(operation)
				// If the operation is an object.
				case map[string]interface{}:
					// Extract the "type" field.
					if t, ok := operation["type"].(string); ok {
						operationModel.Type = types.StringValue(t)
					} else {
						return fmt.Errorf("operation object missing 'type' field")
					}
					// Extract the "delete_dependant" field if present.
					if dd, exists := operation["deleteDependents"]; exists {
						if boolVal, ok := dd.(bool); ok {
							operationModel.DeleteDependents = types.BoolValue(boolVal)
						} else {
							return fmt.Errorf("invalid type for 'delete_dependants' field: %T", dd)
						}

						fmt.Println("delete_dependants: ", operationModel.DeleteDependents)
					}
				default:
					return fmt.Errorf("unexpected type for operation: %T", operation)
				}

				mapping.Operation = &operationModel
			}

			state.Mappings = append(state.Mappings, *mapping)
		}
	}

	return nil
}

// GET /v1/webhooks/{id} never returns security.secret. Keep the value already
// in state so a configured secret is not planned as an in-place change on every refresh.
func securityFromAPI(prior *SecurityModel, api *cli.Security) *SecurityModel {
	if api == nil || !securityReturned(api) {
		return nil
	}

	secret := flex.GoStringToFramework(api.Secret)
	if api.Secret == nil && prior != nil {
		secret = prior.Secret
	}

	return &SecurityModel{
		Secret:                secret,
		SignatureHeaderName:   flex.GoStringToFramework(api.SignatureHeaderName),
		SignatureAlgorithm:    flex.GoStringToFramework(api.SignatureAlgorithm),
		SignaturePrefix:       flex.GoStringToFramework(api.SignaturePrefix),
		RequestIdentifierPath: flex.GoStringToFramework(api.RequestIdentifierPath),
	}
}

func securityReturned(api *cli.Security) bool {
	return api.RequestIdentifierPath != nil || api.Secret != nil || api.SignatureHeaderName != nil || api.SignatureAlgorithm != nil || api.SignaturePrefix != nil
}
