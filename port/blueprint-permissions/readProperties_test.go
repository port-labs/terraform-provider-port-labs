package blueprint_permissions

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestValidateReadPropertiesAgainstBlueprintRejectsMirrorProperty(t *testing.T) {
	resp := &resource.ValidateConfigResponse{}
	readProperties := BlueprintRelationsPermissionsTFBlock{
		"mirrorProp": {
			Roles: []types.String{types.StringValue("Member")},
		},
	}

	validateReadPropertiesAgainstBlueprint(
		resp,
		path.Root("entities"),
		&EntitiesBlueprintPermissionsModel{
			ReadProperties: &readProperties,
		},
		map[string]struct{}{"mirrorProp": {}},
	)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected validation error for mirror property in read_properties")
	}
}

func TestValidateReadPropertiesAgainstBlueprintAllowsRegularProperty(t *testing.T) {
	resp := &resource.ValidateConfigResponse{}
	readProperties := BlueprintRelationsPermissionsTFBlock{
		"myProp": {
			Roles: []types.String{types.StringValue("Member")},
		},
	}

	validateReadPropertiesAgainstBlueprint(
		resp,
		path.Root("entities"),
		&EntitiesBlueprintPermissionsModel{
			ReadProperties: &readProperties,
		},
		map[string]struct{}{"mirrorProp": {}},
	)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected validation error: %v", resp.Diagnostics.Errors())
	}
}

func TestMirrorPropertyKeysFromBlueprint(t *testing.T) {
	keys := mirrorPropertyKeysFromBlueprint(&cli.Blueprint{
		MirrorProperties: map[string]cli.BlueprintMirrorProperty{
			"mirrorProp": {},
		},
	})
	if len(keys) != 1 {
		t.Fatalf("expected one mirror property key, got %d", len(keys))
	}
	if _, ok := keys["mirrorProp"]; !ok {
		t.Fatal("expected mirrorProp in mirror property keys")
	}
}
