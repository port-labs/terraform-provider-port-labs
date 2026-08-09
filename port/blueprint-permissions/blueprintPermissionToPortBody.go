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

func blueprintEntityScopePermissionsTFBlockToBlueprintPermissionsBlock(block BlueprintPermissionsTFBlockWithPolicy) (cli.BlueprintPermissionsBlock, error) {
	policy, err := utils.TerraformJsonStringToGoObject(block.Policy.ValueStringPointer())
	if err != nil {
		return cli.BlueprintPermissionsBlock{}, err
	}

	return cli.BlueprintPermissionsBlock{
		Users:       utils.TFStringListToStringArray(block.Users),
		Roles:       utils.TFStringListToStringArray(block.Roles),
		Teams:       utils.TFStringListToStringArray(block.Teams),
		OwnedByTeam: block.OwnedByTeam.ValueBoolPointer(),
		Policy:      policy,
	}, nil
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

	var readBlock *cli.BlueprintPermissionsBlock
	if state.Entities.Read != nil {
		block, err := blueprintEntityScopePermissionsTFBlockToBlueprintPermissionsBlock(*state.Entities.Read)
		if err != nil {
			return nil, err
		}
		readBlock = &block
	}

	var registerBlock = cli.BlueprintPermissionsBlock{}
	if state.Entities.Register != nil {
		block, err := blueprintEntityScopePermissionsTFBlockToBlueprintPermissionsBlock(*state.Entities.Register)
		if err != nil {
			return nil, err
		}
		registerBlock = block
	}

	var unregisterBlock = cli.BlueprintPermissionsBlock{}
	if state.Entities.Unregister != nil {
		block, err := blueprintEntityScopePermissionsTFBlockToBlueprintPermissionsBlock(*state.Entities.Unregister)
		if err != nil {
			return nil, err
		}
		unregisterBlock = block
	}
	var updateBlock = cli.BlueprintPermissionsBlock{}
	if state.Entities.Update != nil {
		block, err := blueprintEntityScopePermissionsTFBlockToBlueprintPermissionsBlock(*state.Entities.Update)
		if err != nil {
			return nil, err
		}
		updateBlock = block
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
			Read:             readBlock,
			Register:         registerBlock,
			Unregister:       unregisterBlock,
			Update:           updateBlock,
			UpdateProperties: finalUpdateProperties,
			UpdateRelations:  updateRelations,
		},
	}

	return &blueprintPermissions, nil
}
