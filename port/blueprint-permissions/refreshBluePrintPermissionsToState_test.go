package blueprint_permissions

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TestRefreshBlueprintPermissionsStateWithNilUpdateMetadataProperties(t *testing.T) {
	state := &BlueprintPermissionsModel{
		ID:                  types.StringValue("testBlueprint"),
		BlueprintIdentifier: types.StringValue("testBlueprint"),
		Entities:            nil,
	}

	ownedByTeam := false
	apiResponse := &cli.BlueprintPermissions{
		Entities: cli.BlueprintPermissionsEntities{
			Register: cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Admin"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
			},
			Unregister: cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Admin"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
			},
			Update: cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Admin"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
			},
			UpdateProperties: cli.BlueprintRolesOrPropertiesPermissionsBlock{
				"$title": cli.BlueprintPermissionsBlock{
					Users:       []string{},
					Roles:       []string{"Member"},
					Teams:       []string{},
					OwnedByTeam: &ownedByTeam,
				},
				"$identifier": cli.BlueprintPermissionsBlock{
					Users:       []string{},
					Roles:       []string{"Member"},
					Teams:       []string{},
					OwnedByTeam: &ownedByTeam,
				},
			},
		},
	}

	err := refreshBlueprintPermissionsState(state, apiResponse, "testBlueprint", false)
	if err != nil {
		t.Fatalf("refreshBlueprintPermissionsState failed: %v", err)
	}
	if state.Entities == nil {
		t.Fatal("Entities should not be nil after refresh")
	}

	if state.Entities.UpdateMetadataProperties == nil {
		t.Fatal("UpdateMetadataProperties should not be nil after refresh")
	}

	if state.Entities.UpdateMetadataProperties.Title == nil {
		t.Fatal("UpdateMetadataProperties.Title should not be nil")
	}

	if state.Entities.UpdateMetadataProperties.Identifier == nil {
		t.Fatal("UpdateMetadataProperties.Identifier should not be nil")
	}

	if len(state.Entities.UpdateMetadataProperties.Title.Roles) != 1 {
		t.Errorf("Expected 1 role for Title, got %d", len(state.Entities.UpdateMetadataProperties.Title.Roles))
	}

	if state.Entities.UpdateMetadataProperties.Title.Roles[0].ValueString() != "Member" {
		t.Errorf("Expected role 'Member' for Title, got %s", state.Entities.UpdateMetadataProperties.Title.Roles[0].ValueString())
	}
}

func TestRefreshBlueprintPermissionsStateWithExistingUpdateMetadataProperties(t *testing.T) {
	ownedByTeam := false
	state := &BlueprintPermissionsModel{
		ID:                  types.StringValue("testBlueprint"),
		BlueprintIdentifier: types.StringValue("testBlueprint"),
		Entities: &EntitiesBlueprintPermissionsModel{
			UpdateMetadataProperties: &BlueprintMetadataPermissionsTFBlock{
				Title: &BlueprintPermissionsTFBlock{
					Users:       []types.String{types.StringValue("user1@example.com")},
					Roles:       []types.String{types.StringValue("Admin")},
					Teams:       []types.String{},
					OwnedByTeam: types.BoolValue(ownedByTeam),
				},
			},
		},
	}

	apiResponse := &cli.BlueprintPermissions{
		Entities: cli.BlueprintPermissionsEntities{
			Register: cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Admin"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
			},
			Unregister: cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Admin"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
			},
			Update: cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Admin"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
			},
			UpdateProperties: cli.BlueprintRolesOrPropertiesPermissionsBlock{
				"$title": cli.BlueprintPermissionsBlock{
					Users:       []string{"user1@example.com", "user2@example.com"},
					Roles:       []string{"Member"},
					Teams:       []string{},
					OwnedByTeam: &ownedByTeam,
				},
			},
		},
	}

	err := refreshBlueprintPermissionsState(state, apiResponse, "testBlueprint", false)
	if err != nil {
		t.Fatalf("refreshBlueprintPermissionsState failed: %v", err)
	}

	if state.Entities.UpdateMetadataProperties.Title == nil {
		t.Fatal("UpdateMetadataProperties.Title should not be nil")
	}

	if len(state.Entities.UpdateMetadataProperties.Title.Users) != 2 {
		t.Errorf("Expected 2 users for Title, got %d", len(state.Entities.UpdateMetadataProperties.Title.Users))
	}
}

func TestRefreshBlueprintPermissionsStateWithPolicy(t *testing.T) {
	ownedByTeam := false
	policy := map[string]any{
		"combinator": "and",
		"rules": []any{
			map[string]any{
				"property": map[string]any{
					"property": "$identifier",
					"context":  "user",
				},
				"operator": "=",
				"value":    "true",
			},
		},
	}

	state := &BlueprintPermissionsModel{
		ID:                  types.StringValue("testBlueprint"),
		BlueprintIdentifier: types.StringValue("testBlueprint"),
		Entities: &EntitiesBlueprintPermissionsModel{
			Read: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Member")},
				OwnedByTeam: types.BoolValue(false),
				Policy:      types.StringValue(`{"combinator":"and"}`),
			},
			Register: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				OwnedByTeam: types.BoolValue(false),
			},
			Unregister: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				OwnedByTeam: types.BoolValue(false),
			},
			Update: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				OwnedByTeam: types.BoolValue(false),
			},
		},
	}

	apiResponse := &cli.BlueprintPermissions{
		Entities: cli.BlueprintPermissionsEntities{
			Read: &cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Member"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
				Policy:      &policy,
			},
			Register: cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Admin"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
				Policy:      &policy,
			},
			Unregister: cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Admin"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
			},
			Update: cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Admin"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
			},
		},
	}

	err := refreshBlueprintPermissionsState(state, apiResponse, "testBlueprint", false)
	if err != nil {
		t.Fatalf("refreshBlueprintPermissionsState failed: %v", err)
	}

	if state.Entities.Read == nil {
		t.Fatal("Read should be populated when previously configured")
	}
	if state.Entities.Read.Policy.IsNull() {
		t.Fatal("Read.Policy should not be null")
	}
	if state.Entities.Register.Policy.IsNull() {
		t.Fatal("Register.Policy should not be null")
	}
	if !state.Entities.Unregister.Policy.IsNull() {
		t.Fatal("Unregister.Policy should be null when API has no policy")
	}
}

func TestRefreshBlueprintPermissionsStateDoesNotAdoptReadWhenUnconfigured(t *testing.T) {
	ownedByTeam := false
	policy := map[string]any{
		"combinator": "and",
		"rules":      []any{},
	}

	state := &BlueprintPermissionsModel{
		ID:                  types.StringValue("testBlueprint"),
		BlueprintIdentifier: types.StringValue("testBlueprint"),
		Entities: &EntitiesBlueprintPermissionsModel{
			Register: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				OwnedByTeam: types.BoolValue(false),
			},
			Unregister: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				OwnedByTeam: types.BoolValue(false),
			},
			Update: &BlueprintPermissionsTFBlockWithPolicy{
				Roles:       []types.String{types.StringValue("Admin")},
				OwnedByTeam: types.BoolValue(false),
			},
		},
	}

	apiResponse := &cli.BlueprintPermissions{
		Entities: cli.BlueprintPermissionsEntities{
			Read: &cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Admin", "Member"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
				Policy:      &policy,
			},
			Register: cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Admin"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
			},
			Unregister: cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Admin"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
			},
			Update: cli.BlueprintPermissionsBlock{
				Users:       []string{},
				Roles:       []string{"Admin"},
				Teams:       []string{},
				OwnedByTeam: &ownedByTeam,
			},
		},
	}

	err := refreshBlueprintPermissionsState(state, apiResponse, "testBlueprint", false)
	if err != nil {
		t.Fatalf("refreshBlueprintPermissionsState failed: %v", err)
	}

	if state.Entities.Read != nil {
		t.Fatal("Read should remain nil when it was not previously configured")
	}
}
