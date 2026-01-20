package blueprint_permissions

import "github.com/hashicorp/terraform-plugin-framework/types"

type BlueprintPermissionsTFBlock struct {
	Users       []types.String `tfsdk:"users"`
	Roles       []types.String `tfsdk:"roles"`
	Teams       []types.String `tfsdk:"teams"`
	OwnedByTeam types.Bool     `tfsdk:"owned_by_team"`
}

type BlueprintMetadataPermissionsTFBlock struct {
	Team       *BlueprintPermissionsTFBlock `tfsdk:"team"`
	Icon       *BlueprintPermissionsTFBlock `tfsdk:"icon"`
	Identifier *BlueprintPermissionsTFBlock `tfsdk:"identifier"`
	Title      *BlueprintPermissionsTFBlock `tfsdk:"title"`
}

type BlueprintRelationsPermissionsTFBlock map[string]BlueprintPermissionsTFBlock

type EntitiesBlueprintPermissionsModel struct {
	Register                 *BlueprintPermissionsTFBlock          `tfsdk:"register"`
	Unregister               *BlueprintPermissionsTFBlock          `tfsdk:"unregister"`
	Update                   *BlueprintPermissionsTFBlock          `tfsdk:"update"`
	UpdateProperties         *BlueprintRelationsPermissionsTFBlock `tfsdk:"update_properties"`
	UpdateMetadataProperties *BlueprintMetadataPermissionsTFBlock  `tfsdk:"update_metadata_properties"`
	UpdateRelations          *BlueprintRelationsPermissionsTFBlock `tfsdk:"update_relations"`
}

type BlueprintPermissionsModel struct {
	ID                  types.String                       `tfsdk:"id"`
	BlueprintIdentifier types.String                       `tfsdk:"blueprint_identifier"`
	Entities            *EntitiesBlueprintPermissionsModel `tfsdk:"entities"`
}
