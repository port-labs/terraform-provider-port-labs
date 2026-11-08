package entity

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestWriteUnionArrayResourceToBody(t *testing.T) {
	state := &EntityModel{
		Properties: &EntityPropertiesModel{
			ArrayProps: &ArrayPropsModel{
				UnionStringSlices: map[string]UnionStringSliceModel{
					"vulnerabilities": {
						SourceKey: types.StringValue("scanner-a"),
						Items:     types.ListValueMust(types.StringType, []attr.Value{types.StringValue("CVE-2024-1"), types.StringValue("CVE-2024-2")}),
					},
				},
				UnionNumberSlices: map[string]UnionNumberSliceModel{
					"scores": {
						SourceKey: types.StringValue("tool-b"),
						Items:     types.ListValueMust(types.Float64Type, []attr.Value{types.Float64Value(1), types.Float64Value(2)}),
					},
				},
			},
		},
	}

	properties := make(map[string]interface{})
	err := writeUnionArrayResourceToBody(context.Background(), state, properties)
	if err != nil {
		t.Fatalf("writeUnionArrayResourceToBody returned error: %v", err)
	}

	vulnerabilities, ok := properties["vulnerabilities"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected vulnerabilities map, got %#v", properties["vulnerabilities"])
	}

	slice, ok := vulnerabilities["scanner-a"].([]interface{})
	if !ok || len(slice) != 2 {
		t.Fatalf("expected scanner-a slice with 2 items, got %#v", vulnerabilities["scanner-a"])
	}

	scores, ok := properties["scores"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected scores map, got %#v", properties["scores"])
	}

	numberSlice, ok := scores["tool-b"].([]interface{})
	if !ok || len(numberSlice) != 2 {
		t.Fatalf("expected tool-b slice with 2 items, got %#v", scores["tool-b"])
	}
}

func TestEntityResourceToBodyIncludesUnionArraySlices(t *testing.T) {
	state := &EntityModel{
		Title:     types.StringValue("Service"),
		Blueprint: types.StringValue("service"),
		Properties: &EntityPropertiesModel{
			ArrayProps: &ArrayPropsModel{
				UnionStringSlices: map[string]UnionStringSliceModel{
					"tags": {
						SourceKey: types.StringValue("terraform"),
						Items:     types.ListValueMust(types.StringType, []attr.Value{types.StringValue("alpha")}),
					},
				},
			},
		},
	}

	entity, err := entityResourceToBody(context.Background(), state, &cli.Blueprint{Identifier: "service"})
	if err != nil {
		t.Fatalf("entityResourceToBody returned error: %v", err)
	}

	tags, ok := entity.Properties["tags"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected tags map, got %#v", entity.Properties["tags"])
	}

	if _, ok := tags["terraform"]; !ok {
		t.Fatalf("expected terraform source key in tags, got %#v", tags)
	}
}
