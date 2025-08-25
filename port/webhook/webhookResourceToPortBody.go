package webhook

import (
	"context"
	"encoding/json"
	"fmt"

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
				Entity:    &cli.EntityProperty{},
			}

			// Handle complex identifier structure (search query object)
			if !v.Entity.Identifier.IsNull() {
				identifierValue := v.Entity.Identifier.ValueString()
				var parsed interface{}
				if err := json.Unmarshal([]byte(identifierValue), &parsed); err == nil {
					if identifierMap, ok := parsed.(map[string]interface{}); ok {
						// Check if it's a search query object with combinator and rules
						if _, hasCombinator := identifierMap["combinator"]; hasCombinator {
							if _, hasRules := identifierMap["rules"]; hasRules {
								// Validate the search query structure
								if rulesInterface, exists := identifierMap["rules"]; exists {
									if rules, ok := rulesInterface.([]interface{}); ok {
										for i, ruleInterface := range rules {
											if rule, ok := ruleInterface.(map[string]interface{}); ok {
												for _, field := range []string{"property", "operator", "value"} {
													if _, exists := rule[field]; !exists {
														return nil, fmt.Errorf("identifier rule at index %d missing required field '%s'", i, field)
													}
												}
											}
										}
									}
								}
								// Use the parsed object as identifier
								mapping.Entity.Identifier = parsed
							} else {
								// Has combinator but no rules - treat as string
								mapping.Entity.Identifier = identifierValue
							}
						} else {
							// Not a search query object - treat as string
							mapping.Entity.Identifier = identifierValue
						}
					} else {
						// JSON but not an object (e.g., array, string, number) - treat as string
						mapping.Entity.Identifier = identifierValue
					}
				} else {
					// Not valid JSON - treat as string
					mapping.Entity.Identifier = identifierValue
				}
			}

			if !v.Filter.IsNull() {
				filter := v.Filter.ValueString()
				mapping.Filter = &filter
			}
			if v.Operation != nil && !v.Operation.Type.IsNull() {
				operationType := v.Operation.Type.ValueString()
				operation := map[string]interface{}{
					"type": operationType,
				}
				if !v.Operation.DeleteDependents.IsNull() {
					deleteDependents := v.Operation.DeleteDependents.ValueBool()
					operation["deleteDependents"] = deleteDependents
				}
				mapping.Operation = operation
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
				relations := make(map[string]any)
				for k, relationValue := range v.Entity.Relations {
					var parsed interface{}
					if err := json.Unmarshal([]byte(relationValue), &parsed); err == nil {
						if relationMap, ok := parsed.(map[string]interface{}); ok {
							if _, exists := relationMap["combinator"]; !exists {
								return nil, fmt.Errorf("relation '%s' missing required field 'combinator'", k)
							}
							if rulesInterface, exists := relationMap["rules"]; !exists {
								return nil, fmt.Errorf("relation '%s' missing required field 'rules'", k)
							} else if rules, ok := rulesInterface.([]interface{}); ok {
								for i, ruleInterface := range rules {
									if rule, ok := ruleInterface.(map[string]interface{}); ok {
										for _, field := range []string{"property", "operator", "value"} {
											if _, exists := rule[field]; !exists {
												return nil, fmt.Errorf("relation '%s' rule at index %d missing required field '%s'", k, i, field)
											}
										}
									}
								}
							}
							relations[k] = parsed
						} else {
							// JSON but not an object (e.g., array, string, number) - treating as string
							relations[k] = relationValue
						}
					} else {
						// Not valid JSON - treating as string relation
						relations[k] = relationValue
					}
				}
				mapping.Entity.Relations = relations
			}

			w.Mappings = append(w.Mappings, mapping)
		}
	}

	return w, nil
}
