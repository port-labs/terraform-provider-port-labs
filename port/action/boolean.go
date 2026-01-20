package action

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func booleanPropResourceToBody(ctx context.Context, d *SelfServiceTriggerModel, props map[string]cli.ActionProperty, required *[]string) error {
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

			if !prop.DefaultJqQuery.IsNull() {
				defaultJqQuery := prop.DefaultJqQuery.ValueString()
				jqQueryMap := map[string]string{
					"jqQuery": defaultJqQuery,
				}
				property.Default = jqQueryMap
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

			if !prop.Visible.IsNull() {
				property.Visible = prop.Visible.ValueBoolPointer()
			}

			if !prop.VisibleJqQuery.IsNull() {
				VisibleJqQueryMap := map[string]string{
					"jqQuery": prop.VisibleJqQuery.ValueString(),
				}
				property.Visible = VisibleJqQueryMap
			}

			if !prop.Disabled.IsNull() {
				val := prop.Disabled.ValueBool()
				property.Disabled = &val
			}

			if !prop.DisabledJqQuery.IsNull() {
				DisabledJqQuery := map[string]string{
					"jqQuery": prop.DisabledJqQuery.ValueString(),
				}
				property.Disabled = DisabledJqQuery
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}
