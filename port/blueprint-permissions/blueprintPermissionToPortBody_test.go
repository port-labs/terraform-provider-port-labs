package blueprint_permissions

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBlueprintPermissionsToPortBodyWithPolicy(t *testing.T) {
	ownedByTeam := false
	state := &BlueprintPermissionsModel{
		BlueprintIdentifier: types.StringValue("svc"),
		Entities: &EntitiesBlueprintPermissionsModel{
			Read: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin"), types.StringValue("Member")},
				Users:       []types.String{},
				Teams:       []types.String{},
				OwnedByTeam: types.BoolValue(ownedByTeam),
				Policy:      types.StringValue(`{"combinator":"and","rules":[{"operator":"=","property":"$blueprint","value":"svc"}]}`),
			},
			Register: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				Users:       []types.String{},
				Teams:       []types.String{},
				OwnedByTeam: types.BoolValue(ownedByTeam),
				Policy:      types.StringValue(`{"combinator":"and","rules":[]}`),
			},
			Unregister: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				Users:       []types.String{},
				Teams:       []types.String{},
				OwnedByTeam: types.BoolValue(ownedByTeam),
			},
			Update: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				Users:       []types.String{},
				Teams:       []types.String{},
				OwnedByTeam: types.BoolValue(ownedByTeam),
			},
			UpdateMetadataProperties: &BlueprintMetadataPermissionsTFBlock{
				Icon: &BlueprintPermissionsTFBlock{
					Roles:       []types.String{types.StringValue("Admin")},
					OwnedByTeam: types.BoolValue(ownedByTeam),
				},
				Title: &BlueprintPermissionsTFBlock{
					Roles:       []types.String{types.StringValue("Admin")},
					OwnedByTeam: types.BoolValue(ownedByTeam),
				},
				Team: &BlueprintPermissionsTFBlock{
					Roles:       []types.String{types.StringValue("Admin")},
					OwnedByTeam: types.BoolValue(ownedByTeam),
				},
				Identifier: &BlueprintPermissionsTFBlock{
					Roles:       []types.String{types.StringValue("Admin")},
					OwnedByTeam: types.BoolValue(ownedByTeam),
				},
			},
		},
	}

	body, err := blueprintPermissionsToPortBody(state)
	if err != nil {
		t.Fatalf("blueprintPermissionsToPortBody failed: %v", err)
	}
	if body == nil {
		t.Fatal("expected non-nil body")
	}
	if body.Entities.Read == nil {
		t.Fatal("expected read block")
	}
	if body.Entities.Read.Policy == nil {
		t.Fatal("expected read policy")
	}
	if body.Entities.Register.Policy == nil {
		t.Fatal("expected register policy")
	}
	if body.Entities.Unregister.Policy != nil {
		t.Fatal("expected unregister policy to be omitted when unset")
	}
	if _, ok := body.Entities.UpdateProperties["$icon"]; !ok {
		t.Fatal("expected $icon in updateProperties")
	}
}

func TestBlueprintPermissionsToPortBodyOmitsReadWhenUnconfigured(t *testing.T) {
	ownedByTeam := false
	state := &BlueprintPermissionsModel{
		Entities: &EntitiesBlueprintPermissionsModel{
			Register: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				OwnedByTeam: types.BoolValue(ownedByTeam),
			},
			Unregister: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				OwnedByTeam: types.BoolValue(ownedByTeam),
			},
			Update: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				OwnedByTeam: types.BoolValue(ownedByTeam),
			},
		},
	}

	body, err := blueprintPermissionsToPortBody(state)
	if err != nil {
		t.Fatalf("blueprintPermissionsToPortBody failed: %v", err)
	}
	if body.Entities.Read != nil {
		t.Fatal("expected read to be omitted when unconfigured")
	}
}

func TestBlueprintPermissionsToPortBodyInvalidPolicyJSON(t *testing.T) {
	ownedByTeam := false
	state := &BlueprintPermissionsModel{
		Entities: &EntitiesBlueprintPermissionsModel{
			Register: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				OwnedByTeam: types.BoolValue(ownedByTeam),
				Policy:      types.StringValue(`{not-json`),
			},
			Unregister: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				OwnedByTeam: types.BoolValue(ownedByTeam),
			},
			Update: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				OwnedByTeam: types.BoolValue(ownedByTeam),
			},
		},
	}

	_, err := blueprintPermissionsToPortBody(state)
	if err == nil {
		t.Fatal("expected error for invalid policy JSON")
	}
}
