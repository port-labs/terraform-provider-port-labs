package blueprint_permissions

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func refreshEntityScopePermissionsState(oldBlock *BlueprintPermissionsTFBlockWithPolicy, apiBlock cli.BlueprintPermissionsBlock, jsonEscapeHTML bool) (*BlueprintPermissionsTFBlockWithPolicy, error) {
	var users, roles, teams []types.String
	if oldBlock == nil {
		users = utils.Map(apiBlock.Users, types.StringValue)
		roles = utils.Map(apiBlock.Roles, types.StringValue)
		teams = utils.Map(apiBlock.Teams, types.StringValue)
	} else {
		users = utils.Map(utils.SortStringSliceByOther(apiBlock.Users, utils.TFStringListToStringArray(oldBlock.Users)), types.StringValue)
		roles = utils.Map(utils.SortStringSliceByOther(apiBlock.Roles, utils.TFStringListToStringArray(oldBlock.Roles)), types.StringValue)
		teams = utils.Map(utils.SortStringSliceByOther(apiBlock.Teams, utils.TFStringListToStringArray(oldBlock.Teams)), types.StringValue)
	}

	ownedByTeam := false
	if apiBlock.OwnedByTeam != nil {
		ownedByTeam = *apiBlock.OwnedByTeam
	}

	result := &BlueprintPermissionsTFBlockWithPolicy{
		Users:       users,
		Roles:       roles,
		Teams:       teams,
		OwnedByTeam: types.BoolValue(ownedByTeam),
		Policy:      types.StringNull(),
	}

	if apiBlock.Policy != nil {
		var oldPolicy types.String
		if oldBlock != nil {
			oldPolicy = oldBlock.Policy
		}
		policy, err := utils.GoObjectToTerraformStringPreferExisting(oldPolicy, apiBlock.Policy, jsonEscapeHTML)
		if err != nil {
			return nil, err
		}
		result.Policy = policy
	}

	return result, nil
}

func refreshBlueprintPermissionsState(state *BlueprintPermissionsModel, a *cli.BlueprintPermissions, blueprintId string, jsonEscapeHTML bool) error {
	oldPermissions := state.Entities
	if oldPermissions == nil {
		oldPermissions = &EntitiesBlueprintPermissionsModel{}
	}

	state.ID = types.StringValue(blueprintId)
	state.BlueprintIdentifier = types.StringValue(blueprintId)
	state.Entities = &EntitiesBlueprintPermissionsModel{}

	updateBlock, err := refreshEntityScopePermissionsState(oldPermissions.Update, a.Entities.Update, jsonEscapeHTML)
	if err != nil {
		return err
	}
	state.Entities.Update = updateBlock

	unregisterBlock, err := refreshEntityScopePermissionsState(oldPermissions.Unregister, a.Entities.Unregister, jsonEscapeHTML)
	if err != nil {
		return err
	}
	state.Entities.Unregister = unregisterBlock

	registerBlock, err := refreshEntityScopePermissionsState(oldPermissions.Register, a.Entities.Register, jsonEscapeHTML)
	if err != nil {
		return err
	}
	state.Entities.Register = registerBlock

	// Only manage read in state when it was previously configured, to avoid
	// forcing existing configs without entities.read to adopt API defaults.
	if oldPermissions.Read != nil && a.Entities.Read != nil {
		readBlock, err := refreshEntityScopePermissionsState(oldPermissions.Read, *a.Entities.Read, jsonEscapeHTML)
		if err != nil {
			return err
		}
		state.Entities.Read = readBlock
	} else {
		state.Entities.Read = nil
	}

	if oldPermissions.UpdateProperties == nil {
		oldPermissions.UpdateProperties = &BlueprintRelationsPermissionsTFBlock{}
	}

	if oldPermissions.UpdateMetadataProperties == nil {
		oldPermissions.UpdateMetadataProperties = &BlueprintMetadataPermissionsTFBlock{}
	}

	state.Entities.UpdateProperties = nil
	var mappedUpdateProperties BlueprintRelationsPermissionsTFBlock = nil
	if len(a.Entities.UpdateProperties) > 0 {
		state.Entities.UpdateMetadataProperties = &BlueprintMetadataPermissionsTFBlock{}
		mappedUpdateProperties = make(BlueprintRelationsPermissionsTFBlock)
		for updatePropertyKey, updatePropertyValue := range a.Entities.UpdateProperties {
			var oldPropValue *BlueprintPermissionsTFBlock
			if strings.HasPrefix(updatePropertyKey, "$") {
				switch updatePropertyKey {
				case "$title":
					oldPropValue = oldPermissions.UpdateMetadataProperties.Title
				case "$identifier":
					oldPropValue = oldPermissions.UpdateMetadataProperties.Identifier
				case "$icon":
					oldPropValue = oldPermissions.UpdateMetadataProperties.Icon
				case "$team":
					oldPropValue = oldPermissions.UpdateMetadataProperties.Team
				}
			} else if val, ok := (*oldPermissions.UpdateProperties)[updatePropertyKey]; ok {
				oldPropValue = &val
			}
			var current *BlueprintPermissionsTFBlock
			if oldPropValue == nil {
				current = &BlueprintPermissionsTFBlock{
					Users:       utils.Map(updatePropertyValue.Users, types.StringValue),
					Roles:       utils.Map(updatePropertyValue.Roles, types.StringValue),
					Teams:       utils.Map(updatePropertyValue.Teams, types.StringValue),
					OwnedByTeam: types.BoolValue(*updatePropertyValue.OwnedByTeam),
				}
			} else {
				current = &BlueprintPermissionsTFBlock{
					Users:       utils.Map(utils.SortStringSliceByOther(updatePropertyValue.Users, utils.TFStringListToStringArray(oldPropValue.Users)), types.StringValue),
					Roles:       utils.Map(utils.SortStringSliceByOther(updatePropertyValue.Roles, utils.TFStringListToStringArray(oldPropValue.Roles)), types.StringValue),
					Teams:       utils.Map(utils.SortStringSliceByOther(updatePropertyValue.Teams, utils.TFStringListToStringArray(oldPropValue.Teams)), types.StringValue),
					OwnedByTeam: types.BoolValue(*updatePropertyValue.OwnedByTeam),
				}
			}

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

	if oldPermissions.UpdateRelations == nil {
		oldPermissions.UpdateRelations = &BlueprintRelationsPermissionsTFBlock{}
	}
	if len(a.Entities.UpdateRelations) > 0 {
		var mappedUpdateRelations = make(BlueprintRelationsPermissionsTFBlock, len(a.Entities.UpdateRelations))
		for updateRelationKey, updateRelationValue := range a.Entities.UpdateRelations {
			oldRelValue, hasOldRelValue := (*oldPermissions.UpdateRelations)[updateRelationKey]
			if !hasOldRelValue {
				mappedUpdateRelations[updateRelationKey] = BlueprintPermissionsTFBlock{
					Users:       utils.Map(updateRelationValue.Users, types.StringValue),
					Roles:       utils.Map(updateRelationValue.Roles, types.StringValue),
					Teams:       utils.Map(updateRelationValue.Teams, types.StringValue),
					OwnedByTeam: types.BoolValue(*updateRelationValue.OwnedByTeam),
				}
			} else {
				mappedUpdateRelations[updateRelationKey] = BlueprintPermissionsTFBlock{
					Users:       utils.Map(utils.SortStringSliceByOther(updateRelationValue.Users, utils.TFStringListToStringArray(oldRelValue.Users)), types.StringValue),
					Roles:       utils.Map(utils.SortStringSliceByOther(updateRelationValue.Roles, utils.TFStringListToStringArray(oldRelValue.Roles)), types.StringValue),
					Teams:       utils.Map(utils.SortStringSliceByOther(updateRelationValue.Teams, utils.TFStringListToStringArray(oldRelValue.Teams)), types.StringValue),
					OwnedByTeam: types.BoolValue(*updateRelationValue.OwnedByTeam),
				}
			}
		}
		state.Entities.UpdateRelations = &mappedUpdateRelations
	} else {
		state.Entities.UpdateRelations = nil
	}

	if oldPermissions.ReadProperties == nil {
		oldPermissions.ReadProperties = &BlueprintRelationsPermissionsTFBlock{}
	}

	if oldPermissions.ReadMetadataProperties == nil {
		oldPermissions.ReadMetadataProperties = &BlueprintReadMetadataPermissionsTFBlock{}
	}

	state.Entities.ReadProperties = nil
	var mappedReadProperties BlueprintRelationsPermissionsTFBlock = nil
	if len(a.Entities.ReadProperties) > 0 {
		state.Entities.ReadMetadataProperties = &BlueprintReadMetadataPermissionsTFBlock{}
		mappedReadProperties = make(BlueprintRelationsPermissionsTFBlock)
		for readPropertyKey, readPropertyValue := range a.Entities.ReadProperties {
			if readPropertyKey == "$identifier" {
				continue
			}

			var oldPropValue *BlueprintPermissionsTFBlock
			if strings.HasPrefix(readPropertyKey, "$") {
				switch readPropertyKey {
				case "$title":
					oldPropValue = oldPermissions.ReadMetadataProperties.Title
				case "$icon":
					oldPropValue = oldPermissions.ReadMetadataProperties.Icon
				case "$team":
					oldPropValue = oldPermissions.ReadMetadataProperties.Team
				}
			} else if val, ok := (*oldPermissions.ReadProperties)[readPropertyKey]; ok {
				oldPropValue = &val
			}
			var current *BlueprintPermissionsTFBlock
			if oldPropValue == nil {
				current = &BlueprintPermissionsTFBlock{
					Users:       utils.Map(readPropertyValue.Users, types.StringValue),
					Roles:       utils.Map(readPropertyValue.Roles, types.StringValue),
					Teams:       utils.Map(readPropertyValue.Teams, types.StringValue),
					OwnedByTeam: types.BoolValue(*readPropertyValue.OwnedByTeam),
				}
			} else {
				current = &BlueprintPermissionsTFBlock{
					Users:       utils.Map(utils.SortStringSliceByOther(readPropertyValue.Users, utils.TFStringListToStringArray(oldPropValue.Users)), types.StringValue),
					Roles:       utils.Map(utils.SortStringSliceByOther(readPropertyValue.Roles, utils.TFStringListToStringArray(oldPropValue.Roles)), types.StringValue),
					Teams:       utils.Map(utils.SortStringSliceByOther(readPropertyValue.Teams, utils.TFStringListToStringArray(oldPropValue.Teams)), types.StringValue),
					OwnedByTeam: types.BoolValue(*readPropertyValue.OwnedByTeam),
				}
			}

			if strings.HasPrefix(readPropertyKey, "$") {
				switch readPropertyKey {
				case "$title":
					state.Entities.ReadMetadataProperties.Title = current
				case "$icon":
					state.Entities.ReadMetadataProperties.Icon = current
				case "$team":
					state.Entities.ReadMetadataProperties.Team = current
				}
			} else {
				mappedReadProperties[readPropertyKey] = *current
			}
		}
		if len(mappedReadProperties) > 0 {
			state.Entities.ReadProperties = &mappedReadProperties
		}
		if state.Entities.ReadMetadataProperties.Title == nil &&
			state.Entities.ReadMetadataProperties.Icon == nil &&
			state.Entities.ReadMetadataProperties.Team == nil {
			state.Entities.ReadMetadataProperties = nil
		}
	} else {
		state.Entities.ReadMetadataProperties = nil
	}

	return nil
}
