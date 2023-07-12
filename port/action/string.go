package action

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func stringPropResourceToBody(ctx context.Context, d *ActionModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range d.UserProperties.StringProps {
		property := cli.BlueprintProperty{
			Type: "string",
		}

		if !prop.Title.IsNull() {
			title := prop.Title.ValueString()
			property.Title = &title
		}

		if !prop.Default.IsNull() {
			property.Default = prop.Default.ValueString()
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

		if !prop.MinLength.IsNull() {
			minLength := int(prop.MinLength.ValueInt64())
			property.MinLength = &minLength
		}

		if !prop.MaxLength.IsNull() {
			maxLength := int(prop.MaxLength.ValueInt64())
			property.MaxLength = &maxLength
		}

		if !prop.Pattern.IsNull() {
			pattern := prop.Pattern.ValueString()
			property.Pattern = &pattern
		}

		if !prop.Description.IsNull() {
			description := prop.Description.ValueString()
			property.Description = &description
		}

		if !prop.Enum.IsNull() {
			enumList, err := utils.TerraformListToGoArray(ctx, prop.Enum, "string")
			if err != nil {
				return err
			}

			property.Enum = enumList
		}

		if !prop.DependsOn.IsNull() {
			dependsOn, err := utils.TerraformListToGoArray(ctx, prop.DependsOn, "string")
			if err != nil {
				return err
			}
			property.DependsOn = utils.InterfaceToStringArray(dependsOn)

		}

		props[propIdentifier] = property

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}

func addStringPropertiesToResource(ctx context.Context, v *cli.BlueprintProperty) *StringPropModel {
	stringProp := &StringPropModel{
		MinLength: flex.GoInt64ToFramework(v.MinLength),
		MaxLength: flex.GoInt64ToFramework(v.MaxLength),
		Pattern:   flex.GoStringToFramework(v.Pattern),
		Format:    flex.GoStringToFramework(v.Format),
		Blueprint: flex.GoStringToFramework(v.Blueprint),
	}

	if v.Enum != nil {
		attrs := make([]attr.Value, 0, len(v.Enum))
		for _, value := range v.Enum {
			attrs = append(attrs, basetypes.NewStringValue(value.(string)))
		}

		stringProp.Enum, _ = types.ListValue(types.StringType, attrs)
	} else {
		stringProp.Enum = types.ListNull(types.StringType)
	}

	return stringProp
}
