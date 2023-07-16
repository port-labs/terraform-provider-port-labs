package action

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func booleanPropResourceToBody(ctx context.Context, d *ActionModel, props map[string]cli.ActionProperty, required *[]string) error {
	for propIdentifier, prop := range d.UserProperties.BooleanProps {
		props[propIdentifier] = cli.ActionProperty{
			Type: "boolean",
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
			}

			if !prop.Default.IsNull() {
				property.Default = prop.Default.ValueBool()
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}

			if !prop.DependsOn.IsNull() {
				dependsOn, err := utils.TerraformListToGoArray(ctx, prop.DependsOn, "string")
				if err != nil {
					return err
				}
				property.DependsOn = utils.InterfaceToStringArray(dependsOn)

			}
			if prop.Dataset != nil {
				property.Dataset = actionDataSetToPortBody(prop.Dataset)
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}
