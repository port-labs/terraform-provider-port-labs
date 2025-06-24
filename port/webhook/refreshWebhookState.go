package webhook

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
)

func validateRelationStructure(relation map[string]interface{}) error {
	// Combinator validation
	_, exists := relation["combinator"]
	if !exists {
		return fmt.Errorf("missing required field 'combinator'")
	}

	// Rules validation
	rulesInterface, exists := relation["rules"]
	if !exists {
		return fmt.Errorf("missing required field 'rules'")
	}

	rules, ok := rulesInterface.([]interface{})
	if !ok {
		return fmt.Errorf("field 'rules' must be an array, got %T", rulesInterface)
	}

	for i, ruleInterface := range rules {
		rule, ok := ruleInterface.(map[string]interface{})
		if !ok {
			return fmt.Errorf("rule at index %d must be an object, got %T", i, ruleInterface)
		}

		requiredFields := []string{"property", "operator", "value"}
		for _, field := range requiredFields {
			_, exists := rule[field]
			if !exists {
				return fmt.Errorf("rule at index %d is missing required field '%s'", i, field)
			}
		}
	}

	return nil
}

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
				for k, relationValue := range v.Entity.Relations {
					switch val := relationValue.(type) {
					case string:
						mapping.Entity.Relations[k] = val
					case map[string]interface{}:
						if err := validateRelationStructure(val); err != nil {
							return fmt.Errorf("invalid relation structure for key '%s': %w", k, err)
						}
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
