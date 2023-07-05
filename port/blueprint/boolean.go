package blueprint

import "github.com/port-labs/terraform-provider-port-labs/internal/cli"

func booleanPropResourceToBody(state *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) {
	for propIdentifier, prop := range state.Properties.BooleanProp {
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
