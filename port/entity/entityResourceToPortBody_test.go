package entity

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestEntityResourceToBodyUnionArrayProperty(t *testing.T) {
	union := true
	bp := &cli.Blueprint{
		Identifier: "service",
		Schema: cli.BlueprintSchema{
			Properties: map[string]cli.BlueprintProperty{
				"vulnerabilities": {
					Type:  "array",
					Union: &union,
					Items: map[string]any{"type": "string"},
				},
			},
		},
	}

	state := &EntityModel{
		Title:     types.StringValue("Service"),
		Blueprint: types.StringValue("service"),
		Properties: &EntityPropertiesModel{
			UnionArrayProps: &UnionArrayPropsModel{
				StringItems: map[string]UnionArrayItemModel{
					"vulnerabilities": {
						SourceKey: types.StringValue("terraform-sync"),
						Items: types.ListValueMust(types.StringType, []attr.Value{
							types.StringValue("CVE-2024-1"),
							types.StringValue("CVE-2024-2"),
						}),
					},
				},
			},
		},
	}

	entity, err := entityResourceToBody(context.Background(), state, bp)
	if err != nil {
		t.Fatalf("entityResourceToBody returned error: %v", err)
	}

	raw, ok := entity.Properties["vulnerabilities"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map value for union array property, got %#v", entity.Properties["vulnerabilities"])
	}

	items, ok := raw["terraform-sync"].([]interface{})
	if !ok {
		t.Fatalf("expected source slice array, got %#v", raw["terraform-sync"])
	}

	if len(items) != 2 || items[0].(string) != "CVE-2024-1" {
		t.Fatalf("unexpected union array items: %#v", items)
	}
}

func TestEntityResourceToBodyRejectsPlainArrayForUnionProperty(t *testing.T) {
	union := true
	bp := &cli.Blueprint{
		Identifier: "service",
		Schema: cli.BlueprintSchema{
			Properties: map[string]cli.BlueprintProperty{
				"vulnerabilities": {
					Type:  "array",
					Union: &union,
					Items: map[string]any{"type": "string"},
				},
			},
		},
	}

	state := &EntityModel{
		Title:     types.StringValue("Service"),
		Blueprint: types.StringValue("service"),
		Properties: &EntityPropertiesModel{
			ArrayProps: &ArrayPropsModel{
				StringItems: types.MapValueMust(types.ListType{ElemType: types.StringType}, map[string]attr.Value{
					"vulnerabilities": types.ListValueMust(types.StringType, []attr.Value{
						types.StringValue("CVE-2024-1"),
					}),
				}),
			},
		},
	}

	_, err := entityResourceToBody(context.Background(), state, bp)
	if err == nil {
		t.Fatal("expected error when writing plain array to union property")
	}
}
