package action

import (
	"encoding/json"

	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

func objectPropResourceToBody(d *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range d.UserProperties.ObjectProps {
		props[propIdentifier] = cli.BlueprintProperty{
			Type: "object",
		}

		if property, ok := props[propIdentifier]; ok {
			if !prop.Default.IsNull() {
				defaultAsString := prop.Default.ValueString()
				defaultObj := make(map[string]interface{})
				err := json.Unmarshal([]byte(defaultAsString), &defaultObj)
				if err != nil {
					return err
				} else {
					property.Default = defaultObj
				}
			}

			if !prop.Title.IsNull() {
				title := prop.Title.ValueString()
				property.Title = &title
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
	return nil
}

func addObjectPropertiesToResource(v *cli.BlueprintProperty) *ObjectPropModel {
	objectProp := &ObjectPropModel{}

	return objectProp
}
