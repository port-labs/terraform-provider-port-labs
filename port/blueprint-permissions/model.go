package blueprint_permissions

import "github.com/hashicorp/terraform-plugin-framework/types"

type BlueprintPermissionsTFBlock struct {
	Users       []types.String `tfsdk:"users"`
	Roles       []types.String `tfsdk:"roles"`
	Teams       []types.String `tfsdk:"teams"`
	OwnedByTeam types.Bool     `tfsdk:"owned_by_team"`
}

type BlueprintPermissionsTFBlockWithPolicy struct {
	Users       []types.String `tfsdk:"users"`
	Roles       []types.String `tfsdk:"roles"`
	Teams       []types.String `tfsdk:"teams"`
	OwnedByTeam types.Bool     `tfsdk:"owned_by_team"`
	Policy      types.String   `tfsdk:"policy"`
}

type BlueprintMetadataPermissionsTFBlock struct {
	Team       *BlueprintPermissionsTFBlock `tfsdk:"team"`
	Icon       *BlueprintPermissionsTFBlock `tfsdk:"icon"`
	Identifier *BlueprintPermissionsTFBlock `tfsdk:"identifier"`
	Title      *BlueprintPermissionsTFBlock `tfsdk:"title"`
}

type BlueprintReadMetadataPermissionsTFBlock struct {
	Team  *BlueprintPermissionsTFBlock `tfsdk:"team"`
	Icon  *BlueprintPermissionsTFBlock `tfsdk:"icon"`
	Title *BlueprintPermissionsTFBlock `tfsdk:"title"`
}

type BlueprintRelationsPermissionsTFBlock map[string]BlueprintPermissionsTFBlock

type EntitiesBlueprintPermissionsModel struct {
	Read                     *BlueprintPermissionsTFBlockWithPolicy `tfsdk:"read"`
	Register                 *BlueprintPermissionsTFBlockWithPolicy `tfsdk:"register"`
	Unregister               *BlueprintPermissionsTFBlockWithPolicy `tfsdk:"unregister"`
	Update                   *BlueprintPermissionsTFBlockWithPolicy `tfsdk:"update"`
	ReadProperties           *BlueprintRelationsPermissionsTFBlock  `tfsdk:"read_properties"`
	ReadMetadataProperties   *BlueprintReadMetadataPermissionsTFBlock `tfsdk:"read_metadata_properties"`
	UpdateProperties         *BlueprintRelationsPermissionsTFBlock  `tfsdk:"update_properties"`
	UpdateMetadataProperties *BlueprintMetadataPermissionsTFBlock   `tfsdk:"update_metadata_properties"`
	UpdateRelations          *BlueprintRelationsPermissionsTFBlock  `tfsdk:"update_relations"`
}

type BlueprintPermissionsModel struct {
	ID                  types.String                       `tfsdk:"id"`
	BlueprintIdentifier types.String                       `tfsdk:"blueprint_identifier"`
	Entities            *EntitiesBlueprintPermissionsModel `tfsdk:"entities"`
}
