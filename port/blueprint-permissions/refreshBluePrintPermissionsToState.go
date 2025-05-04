package blueprint_permissions

import (
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func refreshBlueprintPermissionsState(state *BlueprintPermissionsModel, a *cli.BlueprintPermissions, blueprintId string) error {
	oldPermissions := state.Entities
	if oldPermissions == nil {
		oldPermissions = &EntitiesBlueprintPermissionsModel{}
	}

	state.ID = types.StringValue(blueprintId)
	state.BlueprintIdentifier = types.StringValue(blueprintId)
	state.Entities = &EntitiesBlueprintPermissionsModel{}

	if oldPermissions.Update == nil {
		state.Entities.Update = &BlueprintPermissionsTFBlock{
			Users:       utils.Map(a.Entities.Update.Users, types.StringValue),
			Roles:       utils.Map(a.Entities.Update.Roles, types.StringValue),
			Teams:       utils.Map(a.Entities.Update.Teams, types.StringValue),
			OwnedByTeam: types.BoolValue(*a.Entities.Update.OwnedByTeam),
		}
	} else {
		state.Entities.Update = &BlueprintPermissionsTFBlock{
			Users:       utils.Map(utils.SortStringSliceByOther(a.Entities.Update.Users, utils.TFStringListToStringArray(oldPermissions.Update.Users)), types.StringValue),
			Roles:       utils.Map(utils.SortStringSliceByOther(a.Entities.Update.Roles, utils.TFStringListToStringArray(oldPermissions.Update.Roles)), types.StringValue),
			Teams:       utils.Map(utils.SortStringSliceByOther(a.Entities.Update.Teams, utils.TFStringListToStringArray(oldPermissions.Update.Teams)), types.StringValue),
			OwnedByTeam: types.BoolValue(*a.Entities.Update.OwnedByTeam),
		}
	}

	if oldPermissions.Unregister == nil {
		state.Entities.Unregister = &BlueprintPermissionsTFBlock{
			Users:       utils.Map(a.Entities.Unregister.Users, types.StringValue),
			Roles:       utils.Map(a.Entities.Unregister.Roles, types.StringValue),
			Teams:       utils.Map(a.Entities.Unregister.Teams, types.StringValue),
			OwnedByTeam: types.BoolValue(*a.Entities.Unregister.OwnedByTeam),
		}
	} else {
		state.Entities.Unregister = &BlueprintPermissionsTFBlock{
			Users:       utils.Map(utils.SortStringSliceByOther(a.Entities.Unregister.Users, utils.TFStringListToStringArray(oldPermissions.Unregister.Users)), types.StringValue),
			Roles:       utils.Map(utils.SortStringSliceByOther(a.Entities.Unregister.Roles, utils.TFStringListToStringArray(oldPermissions.Unregister.Roles)), types.StringValue),
			Teams:       utils.Map(utils.SortStringSliceByOther(a.Entities.Unregister.Teams, utils.TFStringListToStringArray(oldPermissions.Unregister.Teams)), types.StringValue),
			OwnedByTeam: types.BoolValue(*a.Entities.Unregister.OwnedByTeam),
		}
	}

	if oldPermissions.Register == nil {
		state.Entities.Register = &BlueprintPermissionsTFBlock{
			Users:       utils.Map(a.Entities.Register.Users, types.StringValue),
			Roles:       utils.Map(a.Entities.Register.Roles, types.StringValue),
			Teams:       utils.Map(a.Entities.Register.Teams, types.StringValue),
			OwnedByTeam: types.BoolValue(*a.Entities.Register.OwnedByTeam),
		}
	} else {
		state.Entities.Register = &BlueprintPermissionsTFBlock{
			Users:       utils.Map(utils.SortStringSliceByOther(a.Entities.Register.Users, utils.TFStringListToStringArray(oldPermissions.Register.Users)), types.StringValue),
			Roles:       utils.Map(utils.SortStringSliceByOther(a.Entities.Register.Roles, utils.TFStringListToStringArray(oldPermissions.Register.Roles)), types.StringValue),
			Teams:       utils.Map(utils.SortStringSliceByOther(a.Entities.Register.Teams, utils.TFStringListToStringArray(oldPermissions.Register.Teams)), types.StringValue),
			OwnedByTeam: types.BoolValue(*a.Entities.Register.OwnedByTeam),
		}
	}

	if oldPermissions.UpdateProperties == nil {
		oldPermissions.UpdateProperties = &BlueprintRelationsPermissionsTFBlock{}
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

	return nil
}
