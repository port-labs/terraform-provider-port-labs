package action

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func objectPropResourceToBody(ctx context.Context, d *SelfServiceTriggerModel, props map[string]cli.ActionProperty, required *[]string) error {
	for propIdentifier, prop := range d.UserProperties.ObjectProps {
		props[propIdentifier] = cli.ActionProperty{
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

			if !prop.DefaultJqQuery.IsNull() {
				defaultJqQuery := prop.DefaultJqQuery.ValueString()
				jqQueryMap := map[string]string{
					"jqQuery": defaultJqQuery,
				}
				property.Default = jqQueryMap
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

			if !prop.DependsOn.IsNull() {
				dependsOn, err := utils.TerraformListToGoArray(ctx, prop.DependsOn, "string")
				if err != nil {
					return err
				}
				property.DependsOn = utils.InterfaceToStringArray(dependsOn)
			}

			if !prop.Encryption.IsNull() {
				encryption := prop.Encryption.ValueString()
				property.Encryption = encryption
			}

			if prop.ClientSideEncryption != nil {
				property.Encryption = map[string]string{
					"algorithm": prop.ClientSideEncryption.Algorithm.ValueString(),
					"key":       prop.ClientSideEncryption.Key.ValueString(),
				}
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

func addObjectPropertiesToResource(v *cli.ActionProperty) *ObjectPropModel {
	objectProp := &ObjectPropModel{}

	// Handle encryption - can be either a string or an object
	objectProp.Encryption = types.StringNull()
	objectProp.ClientSideEncryption = nil
	if v.Encryption != nil {
		switch enc := v.Encryption.(type) {
		case string:
			objectProp.Encryption = types.StringValue(enc)
		case map[string]interface{}:
			algorithm, _ := enc["algorithm"].(string)
			key, _ := enc["key"].(string)
			if algorithm != "" && key != "" {
				objectProp.ClientSideEncryption = &ClientSideEncryptionModel{
					Algorithm: types.StringValue(algorithm),
					Key:       types.StringValue(key),
				}
			}
		}
	}

	return objectProp
}
