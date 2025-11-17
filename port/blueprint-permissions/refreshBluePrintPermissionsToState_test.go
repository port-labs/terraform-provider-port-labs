package blueprint_permissions

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

// TestRefreshBlueprintPermissionsStateWithNilUpdateMetadataProperties tests the bug
// that occurs during import when UpdateMetadataProperties is nil
func TestRefreshBlueprintPermissionsStateWithNilUpdateMetadataProperties(t *testing.T) {
	// Simulate the state during import - only ID and BlueprintIdentifier are set
	// This mimics what happens in ImportState method
	state := &BlueprintPermissionsModel{
		ID:                  types.StringValue("testBlueprint"),
		BlueprintIdentifier: types.StringValue("testBlueprint"),
		Entities:            nil, // This is nil during import
	}

	// Create API response with metadata properties permissions
	// This simulates what the Port API returns
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

	// This should panic with the bug, but pass after the fix
	err := refreshBlueprintPermissionsState(state, apiResponse, "testBlueprint")
	if err != nil {
		t.Fatalf("refreshBlueprintPermissionsState failed: %v", err)
	}

	// Verify the state was populated correctly
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

	// Verify the roles were set correctly
	if len(state.Entities.UpdateMetadataProperties.Title.Roles) != 1 {
		t.Errorf("Expected 1 role for Title, got %d", len(state.Entities.UpdateMetadataProperties.Title.Roles))
	}

	if state.Entities.UpdateMetadataProperties.Title.Roles[0].ValueString() != "Member" {
		t.Errorf("Expected role 'Member' for Title, got %s", state.Entities.UpdateMetadataProperties.Title.Roles[0].ValueString())
	}
}

// TestRefreshBlueprintPermissionsStateWithExistingUpdateMetadataProperties tests
// that the function works correctly when UpdateMetadataProperties already exists
func TestRefreshBlueprintPermissionsStateWithExistingUpdateMetadataProperties(t *testing.T) {
	// Simulate a state where UpdateMetadataProperties already exists (normal update scenario)
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

	// Create API response
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

	err := refreshBlueprintPermissionsState(state, apiResponse, "testBlueprint")
	if err != nil {
		t.Fatalf("refreshBlueprintPermissionsState failed: %v", err)
	}

	// Verify the state was updated correctly
	if state.Entities.UpdateMetadataProperties.Title == nil {
		t.Fatal("UpdateMetadataProperties.Title should not be nil")
	}

	// The users should be sorted according to the old order (user1 first)
	if len(state.Entities.UpdateMetadataProperties.Title.Users) != 2 {
		t.Errorf("Expected 2 users for Title, got %d", len(state.Entities.UpdateMetadataProperties.Title.Users))
	}
}
