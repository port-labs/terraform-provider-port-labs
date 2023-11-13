package blueprint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func addStringPropertiesToState(ctx context.Context, v *cli.BlueprintProperty) *StringPropModel {
	stringProp := &StringPropModel{
		MinLength: flex.GoInt64ToFramework(v.MinLength),
		MaxLength: flex.GoInt64ToFramework(v.MaxLength),
		Format:    types.StringPointerValue(v.Format),
		Spec:      types.StringPointerValue(v.Spec),
		Pattern:   types.StringPointerValue(v.Pattern),
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

	if v.EnumColors != nil {
		stringProp.EnumColors, _ = types.MapValueFrom(ctx, types.StringType, v.EnumColors)
	} else {
		stringProp.EnumColors = types.MapNull(types.StringType)
	}

	if v.SpecAuthentication != nil {
		stringProp.SpecAuthentication = &SpecAuthenticationModel{
			AuthorizationUrl: types.StringValue(v.SpecAuthentication.AuthorizationUrl),
			TokenUrl:         types.StringValue(v.SpecAuthentication.TokenUrl),
			ClientId:         types.StringValue(v.SpecAuthentication.ClientId),
		}
	}

	return stringProp
}

func stringPropResourceToBody(ctx context.Context, state *BlueprintModel, props map[string]cli.BlueprintProperty, required *[]string) error {
	for propIdentifier, prop := range state.Properties.StringProps {
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

		if !prop.Spec.IsNull() {
			spec := prop.Spec.ValueString()
			property.Spec = &spec
		}

		if prop.SpecAuthentication != nil {
			specAuth := &cli.SpecAuthentication{
				AuthorizationUrl: prop.SpecAuthentication.AuthorizationUrl.ValueString(),
				TokenUrl:         prop.SpecAuthentication.TokenUrl.ValueString(),
				ClientId:         prop.SpecAuthentication.ClientId.ValueString(),
			}
			property.SpecAuthentication = specAuth
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

		if !prop.EnumColors.IsNull() {
			enumColor := map[string]string{}
			for k, v := range prop.EnumColors.Elements() {
				value, _ := v.ToTerraformValue(ctx)
				var keyValue string
				err := value.As(&keyValue)
				if err != nil {
					return err
				}
				enumColor[k] = keyValue
			}

			property.EnumColors = enumColor
		}

		props[propIdentifier] = property

		if prop.Required.ValueBool() {
			*required = append(*required, propIdentifier)
		}
	}
	return nil
}
