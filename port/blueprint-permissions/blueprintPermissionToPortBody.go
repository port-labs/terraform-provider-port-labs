package blueprint_permissions

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func blueprintPermissionsTFBlockToBlueprintPermissionsBlock(block BlueprintPermissionsTFBlock) cli.BlueprintPermissionsBlock {
	return cli.BlueprintPermissionsBlock{
		Users:       utils.TFStringListToStringArray(block.Users),
		Roles:       utils.TFStringListToStringArray(block.Roles),
		Teams:       utils.TFStringListToStringArray(block.Teams),
		OwnedByTeam: block.OwnedByTeam.ValueBoolPointer(),
	}
}

func blueprintPermissionsToPortBody(state *BlueprintPermissionsModel) (*cli.BlueprintPermissions, error) {
	if state == nil {
		return nil, nil
	}

	var updateRelations cli.BlueprintRolesOrPropertiesPermissionsBlock = nil
	if state.Entities.UpdateRelations != nil {
		updateRelations = make(cli.BlueprintRolesOrPropertiesPermissionsBlock, len(*state.Entities.UpdateRelations))
		for updateRelationKey, updateRelationValue := range *state.Entities.UpdateRelations {
			updateRelations[updateRelationKey] = blueprintPermissionsTFBlockToBlueprintPermissionsBlock(updateRelationValue)
		}
	}

	var updateMetadataProperties cli.BlueprintRolesOrPropertiesPermissionsBlock = nil
	if state.Entities.UpdateMetadataProperties != nil {
		updateMetadataProperties = make(cli.BlueprintRolesOrPropertiesPermissionsBlock)
		if state.Entities.UpdateMetadataProperties.Team != nil {
			updateMetadataProperties["$team"] = blueprintPermissionsTFBlockToBlueprintPermissionsBlock(*state.Entities.UpdateMetadataProperties.Team)
		}
		if state.Entities.UpdateMetadataProperties.Icon != nil {
			updateMetadataProperties["$icon"] = blueprintPermissionsTFBlockToBlueprintPermissionsBlock(*state.Entities.UpdateMetadataProperties.Icon)
		}
		if state.Entities.UpdateMetadataProperties.Identifier != nil {
			updateMetadataProperties["$identifier"] = blueprintPermissionsTFBlockToBlueprintPermissionsBlock(*state.Entities.UpdateMetadataProperties.Identifier)
		}
		if state.Entities.UpdateMetadataProperties.Title != nil {
			updateMetadataProperties["$title"] = blueprintPermissionsTFBlockToBlueprintPermissionsBlock(*state.Entities.UpdateMetadataProperties.Title)
		}
	}
	var updateProperties cli.BlueprintRolesOrPropertiesPermissionsBlock = nil
	if state.Entities.UpdateProperties != nil {
		updateProperties = make(cli.BlueprintRolesOrPropertiesPermissionsBlock, len(*state.Entities.UpdateProperties))
		for updatePropertiesKey, updatePropertiesValue := range *state.Entities.UpdateProperties {
			updateProperties[updatePropertiesKey] = blueprintPermissionsTFBlockToBlueprintPermissionsBlock(updatePropertiesValue)
		}
	}

	var registerBlock = cli.BlueprintPermissionsBlock{}
	if state.Entities.Register != nil {
		registerBlock = blueprintPermissionsTFBlockToBlueprintPermissionsBlock(*state.Entities.Register)
	}

	var unregisterBlock = cli.BlueprintPermissionsBlock{}
	if state.Entities.Unregister != nil {
		unregisterBlock = blueprintPermissionsTFBlockToBlueprintPermissionsBlock(*state.Entities.Unregister)
	}
	var updateBlock = cli.BlueprintPermissionsBlock{}
	if state.Entities.Update != nil {
		updateBlock = blueprintPermissionsTFBlockToBlueprintPermissionsBlock(*state.Entities.Update)
	}
	var finalUpdateProperties cli.BlueprintRolesOrPropertiesPermissionsBlock = nil
	if updateMetadataProperties != nil {
		finalUpdateProperties = updateMetadataProperties
	}
	if updateProperties != nil {
		if finalUpdateProperties == nil {
			finalUpdateProperties = updateProperties
		} else {
			utils.CopyGenericMaps(finalUpdateProperties, updateProperties)
		}
	}
	blueprintPermissions := cli.BlueprintPermissions{
		Entities: cli.BlueprintPermissionsEntities{
			Register:         registerBlock,
			Unregister:       unregisterBlock,
			Update:           updateBlock,
			UpdateProperties: finalUpdateProperties,
			UpdateRelations:  updateRelations,
		},
	}

	return &blueprintPermissions, nil
}
