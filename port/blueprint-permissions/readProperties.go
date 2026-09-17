package blueprint_permissions

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func mirrorPropertyKeysFromBlueprint(blueprint *cli.Blueprint) map[string]struct{} {
	if blueprint == nil || len(blueprint.MirrorProperties) == 0 {
		return nil
	}

	keys := make(map[string]struct{}, len(blueprint.MirrorProperties))
	for key := range blueprint.MirrorProperties {
		keys[key] = struct{}{}
	}
	return keys
}

func validateReadPropertiesAgainstBlueprint(
	resp *resource.ValidateConfigResponse,
	entitiesPath path.Path,
	entities *EntitiesBlueprintPermissionsModel,
	mirrorPropertyKeys map[string]struct{},
) {
	if entities == nil || entities.ReadProperties == nil || len(mirrorPropertyKeys) == 0 {
		return
	}

	for propertyKey := range *entities.ReadProperties {
		if _, isMirrorProperty := mirrorPropertyKeys[propertyKey]; isMirrorProperty {
			resp.Diagnostics.AddAttributeError(
				entitiesPath.AtName("read_properties").AtMapKey(propertyKey),
				"Mirror property cannot have read restrictions",
				fmt.Sprintf("The property %q is a mirror property and cannot have read restrictions.", propertyKey),
			)
		}
	}
}

func (r *BlueprintPermissionsResource) validateReadPropertiesConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var state BlueprintPermissionsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() || state.Entities == nil || state.Entities.ReadProperties == nil {
		return
	}

	if state.BlueprintIdentifier.IsUnknown() || state.BlueprintIdentifier.IsNull() {
		return
	}

	blueprint, statusCode, err := r.portClient.ReadBlueprint(ctx, state.BlueprintIdentifier.ValueString())
	if err != nil {
		if statusCode == 404 {
			return
		}
		resp.Diagnostics.AddError("failed to read blueprint for read_properties validation", err.Error())
		return
	}

	validateReadPropertiesAgainstBlueprint(
		resp,
		path.Root("entities"),
		state.Entities,
		mirrorPropertyKeysFromBlueprint(blueprint),
	)
}
