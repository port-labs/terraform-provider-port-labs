package action

import "github.com/port-labs/terraform-provider-port-labs/internal/cli"

func booleanPropResourceToBody(d *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) {
	for propIdentifier, prop := range d.UserProperties.BooleanProps {
		props[propIdentifier] = cli.BlueprintProperty{
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

			if !prop.Format.IsNull() {
				format := prop.Format.ValueString()
				property.Format = &format
			}

			if !prop.Blueprint.IsNull() {
				blueprint := prop.Blueprint.ValueString()
				property.Blueprint = &blueprint
			}

			if !prop.Icon.IsNull() {
				icon := prop.Icon.ValueString()
				property.Icon = &icon
			}

			if !prop.Description.IsNull() {
				description := prop.Description.ValueString()
				property.Description = &description
			}

			props[propIdentifier] = property
		}
		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
}
