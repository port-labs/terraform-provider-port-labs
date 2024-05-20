package blueprint_permissions

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func goStringListToTFList(list []string) []types.String {
	var result = make([]types.String, len(list))
	for i, u := range list {
		result[i] = types.StringValue(u)
	}

	return result
}

func blueprintPermissionsBlockToBlueprintPermissionsTFBlock(block cli.BlueprintPermissionsBlock) *BlueprintPermissionsTFBlock {
	return &BlueprintPermissionsTFBlock{
		Users:       goStringListToTFList(block.Users),
		Roles:       goStringListToTFList(block.Roles),
		Teams:       goStringListToTFList(block.Teams),
		OwnedByTeam: types.BoolValue(*block.OwnedByTeam),
	}
}

func refreshBlueprintPermissionsState(state *BlueprintPermissionsModel, a *cli.BlueprintPermissions, blueprintId string) error {
	state.ID = types.StringValue(blueprintId)
	state.BlueprintIdentifier = types.StringValue(blueprintId)
	state.Entities = &EntitiesBlueprintPermissionsModel{}
	state.Entities.Update = blueprintPermissionsBlockToBlueprintPermissionsTFBlock(a.Entities.Update)
	state.Entities.Unregister = blueprintPermissionsBlockToBlueprintPermissionsTFBlock(a.Entities.Unregister)
	state.Entities.Register = blueprintPermissionsBlockToBlueprintPermissionsTFBlock(a.Entities.Register)

	state.Entities.UpdateProperties = nil
	var mappedUpdateProperties BlueprintRelationsPermissionsTFBlock = nil
	if len(a.Entities.UpdateProperties) > 0 {
		state.Entities.UpdateMetadataProperties = &BlueprintMetadataPermissionsTFBlock{}
		mappedUpdateProperties = make(BlueprintRelationsPermissionsTFBlock)
		for updatePropertyKey, updatePropertyValue := range a.Entities.UpdateProperties {
			var current = blueprintPermissionsBlockToBlueprintPermissionsTFBlock(updatePropertyValue)

			if strings.HasPrefix(updatePropertyKey, "$") {
				switch updatePropertyKey {
				case "$title":
					state.Entities.UpdateMetadataProperties.Title = current
				case "$identifier":
					state.Entities.UpdateMetadataProperties.Identifier = current
				case "$icon":
					state.Entities.UpdateMetadataProperties.Icon = current
				case "$team":
					state.Entities.UpdateMetadataProperties.Team = current
				}
			} else {
				mappedUpdateProperties[updatePropertyKey] = *current
			}
		}
		if len(mappedUpdateProperties) > 0 {
			state.Entities.UpdateProperties = &mappedUpdateProperties
		}
	}
	if len(a.Entities.UpdateRelations) > 0 {
		var mappedUpdateRelations = make(BlueprintRelationsPermissionsTFBlock, len(a.Entities.UpdateRelations))
		for updateRelationKey, updateRelationValue := range a.Entities.UpdateRelations {
			mappedUpdateRelations[updateRelationKey] = *blueprintPermissionsBlockToBlueprintPermissionsTFBlock(updateRelationValue)
		}
		state.Entities.UpdateRelations = &mappedUpdateRelations
	} else {
		state.Entities.UpdateRelations = nil
	}

	return nil
}
