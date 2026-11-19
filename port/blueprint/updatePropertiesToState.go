package blueprint

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/flex"
	"github.com/samber/lo"
)

func SetCommonProperties(v cli.BlueprintProperty, prop interface{}, jsonEscapeHTML bool) {
	properties := []string{"Description", "Icon", "Default", "Title"}
	for _, property := range properties {
		switch property {
		case "Description":
			switch p := prop.(type) {
			case *StringPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			case *NumberPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			case *BooleanPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			case *ArrayPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			case *ObjectPropModel:
				p.Description = flex.GoStringToFramework(v.Description)
			}
		case "Icon":
			switch p := prop.(type) {
			case *StringPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			case *NumberPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			case *BooleanPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			case *ArrayPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			case *ObjectPropModel:
				p.Icon = flex.GoStringToFramework(v.Icon)
			}
		case "Title":
			switch p := prop.(type) {
			case *StringPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			case *NumberPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			case *BooleanPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			case *ArrayPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			case *ObjectPropModel:
				p.Title = flex.GoStringToFramework(v.Title)
			}
		case "Default":
			if v.Default != nil {
				switch p := prop.(type) {
				case *StringPropModel:
					p.Default = types.StringValue(v.Default.(string))
				case *NumberPropModel:
					p.Default = types.Float64Value(v.Default.(float64))
				case *BooleanPropModel:
					p.Default = types.BoolValue(v.Default.(bool))
				}
			}
		}
	}
}

func SetObjectDefault(objectProp *ObjectPropModel, defaultValue interface{}, jsonEscapeHTML bool, oldObjectDefault types.String) {
	if defaultValue == nil {
		return
	}
	objectProp.Default, _ = utils.GoObjectToTerraformStringPreferExisting(oldObjectDefault, defaultValue, jsonEscapeHTML)
}

func (r *BlueprintResource) updatePropertiesToState(ctx context.Context, b *cli.Blueprint, bm *BlueprintModel) error {
	properties := &PropertiesModel{}
	oldProperties := bm.Properties

	for k, v := range b.Schema.Properties {
		switch v.Type {
		case "string":
			if properties.StringProps == nil {
				properties.StringProps = make(map[string]StringPropModel)
			}
			stringProp := AddStringPropertiesToState(ctx, &v)

			if lo.Contains(b.Schema.Required, k) {
				stringProp.Required = types.BoolValue(true)
			} else {
				stringProp.Required = types.BoolValue(false)
			}

			SetCommonProperties(v, stringProp, r.portClient.JSONEscapeHTML)

			properties.StringProps[k] = *stringProp

		case "number":
			if properties.NumberProps == nil {
				properties.NumberProps = make(map[string]NumberPropModel)
			}

			numberProp := AddNumberPropertiesToState(ctx, &v)

			if lo.Contains(b.Schema.Required, k) {
				numberProp.Required = types.BoolValue(true)
			} else {
				numberProp.Required = types.BoolValue(false)
			}

			SetCommonProperties(v, numberProp, r.portClient.JSONEscapeHTML)

			properties.NumberProps[k] = *numberProp

		case "array":
			if properties.ArrayProps == nil {
				properties.ArrayProps = make(map[string]ArrayPropModel)
			}

			var oldArrayProp *ArrayPropModel
			if oldProperties != nil && oldProperties.ArrayProps != nil {
				if old, ok := oldProperties.ArrayProps[k]; ok {
					oldArrayProp = &old
				}
			}
			arrayProp := AddArrayPropertiesToState(ctx, &v, r.portClient.JSONEscapeHTML, oldArrayProp)

			if lo.Contains(b.Schema.Required, k) {
				arrayProp.Required = types.BoolValue(true)
			} else {
				arrayProp.Required = types.BoolValue(false)
			}

			SetCommonProperties(v, arrayProp, r.portClient.JSONEscapeHTML)

			properties.ArrayProps[k] = *arrayProp

		case "boolean":
			if properties.BooleanProps == nil {
				properties.BooleanProps = make(map[string]BooleanPropModel)
			}

			booleanProp := &BooleanPropModel{}

			SetCommonProperties(v, booleanProp, r.portClient.JSONEscapeHTML)

			if lo.Contains(b.Schema.Required, k) {
				booleanProp.Required = types.BoolValue(true)
			} else {
				booleanProp.Required = types.BoolValue(false)
			}

			properties.BooleanProps[k] = *booleanProp

		case "object":
			if properties.ObjectProps == nil {
				properties.ObjectProps = make(map[string]ObjectPropModel)
			}

			objectProp := AddObjectPropertiesToState(&v)

			if lo.Contains(b.Schema.Required, k) {
				objectProp.Required = types.BoolValue(true)
			} else {
				objectProp.Required = types.BoolValue(false)
			}

			var oldObjectDefault types.String
			if oldProperties != nil && oldProperties.ObjectProps != nil {
				if old, ok := oldProperties.ObjectProps[k]; ok {
					oldObjectDefault = old.Default
				}
			}

			SetCommonProperties(v, objectProp, r.portClient.JSONEscapeHTML)
			SetObjectDefault(objectProp, v.Default, r.portClient.JSONEscapeHTML, oldObjectDefault)

			properties.ObjectProps[k] = *objectProp
		}
	}

	bm.Properties = properties

	return nil
}

// addRelationsToState fully rebuilds bm.Relations from the live Port blueprint b.
//
// Previously this only ever ADDED keys found in b.Relations into bm.Relations, and never
// reset/cleared bm.Relations first. That is harmless when the Terraform resource's prior state
// already matches the real relation set, but it silently masks drift whenever bm already carries
// a relation key that is no longer (or not yet) present in Port: that stale key simply survives
// the refresh untouched.
//
// This surfaces concretely the first time a resource is reconciled with no cached Terraform state
// (a freshly created or freshly imported resource): callers that reconstruct a synthetic "prior
// state" from the resource's desired configuration before calling Read (as upjet's Terraform
// Plugin Framework external client does, see terraformPluginFrameworkExternalClient.Connect in
// crossplane/upjet when GetObservation() returns an empty/missing status) will pass in a `bm`
// that already contains a relation the user just added but Port has not received yet. Because
// this function only adds keys from `b` and never removes ones that are absent from it, that
// not-yet-persisted relation survives the refresh, making the freshly-read "observed" state look
// identical to the desired config. Any caller that compares observed vs. desired state to decide
// whether an update is needed (e.g. Crossplane's provider-port, or `terraform plan` itself in
// unusual state-reconstruction scenarios) then concludes nothing changed and never calls
// UpdateBlueprint -- a pure relations addition, with nothing else in the blueprint changing, is
// dropped without any error.
//
// updatePropertiesToState (above) does not have this problem: it already does a full
// `bm.Properties = properties` rebuild rather than merging into whatever bm.Properties previously
// held.
func addRelationsToState(b *cli.Blueprint, bm *BlueprintModel) {
	bm.Relations = make(map[string]RelationModel, len(b.Relations))
	for k, v := range b.Relations {
		relationModel := &RelationModel{
			Target:      types.StringValue(*v.Target),
			Title:       flex.GoStringToFramework(v.Title),
			Description: flex.GoStringToFramework(v.Description),
			Many:        flex.GoBoolToFramework(v.Many),
			Required:    flex.GoBoolToFramework(v.Required),
		}

		bm.Relations[k] = *relationModel

	}
}

// addMirrorPropertiesToState fully rebuilds bm.MirrorProperties from b. See the comment on
// addRelationsToState above -- same additive-only-merge issue applied here too.
func addMirrorPropertiesToState(b *cli.Blueprint, bm *BlueprintModel) {
	bm.MirrorProperties = make(map[string]MirrorPropertyModel, len(b.MirrorProperties))
	for k, v := range b.MirrorProperties {

		mirrorPropertyModel := &MirrorPropertyModel{
			Path:  types.StringValue(v.Path),
			Title: flex.GoStringToFramework(v.Title),
		}

		bm.MirrorProperties[k] = *mirrorPropertyModel

	}
}

// addCalculationPropertiesToState fully rebuilds bm.CalculationProperties from b. See the comment
// on addRelationsToState above -- same additive-only-merge issue applied here too.
func addCalculationPropertiesToState(ctx context.Context, b *cli.Blueprint, bm *BlueprintModel) {
	bm.CalculationProperties = make(map[string]CalculationPropertyModel, len(b.CalculationProperties))
	for k, v := range b.CalculationProperties {
		calculationPropertyModel := &CalculationPropertyModel{
			Calculation: types.StringValue(v.Calculation),
			Type:        types.StringValue(v.Type),
			Title:       flex.GoStringToFramework(v.Title),
			Icon:        flex.GoStringToFramework(v.Icon),
			Description: flex.GoStringToFramework(v.Description),
			Format:      flex.GoStringToFramework(v.Format),
			DateFormat:  flex.GoStringToFramework(v.DateFormat),
			Colorized:   flex.GoBoolToFramework(v.Colorized),
			Spec:        flex.GoStringToFramework(v.Spec),
		}

		if v.Colors != nil {
			calculationPropertyModel.Colors, _ = types.MapValueFrom(ctx, types.StringType, v.Colors)
		} else {
			calculationPropertyModel.Colors = types.MapNull(types.StringType)
		}

		if v.SpecAuthentication != nil {
			calculationPropertyModel.SpecAuthentication = &SpecAuthenticationModel{
				AuthorizationUrl: types.StringValue(v.SpecAuthentication.AuthorizationUrl),
				TokenUrl:         types.StringValue(v.SpecAuthentication.TokenUrl),
				ClientId:         types.StringValue(v.SpecAuthentication.ClientId),
			}
		}

		if v.Items != nil {
			calculationPropertyModel.Items = &CalculationItems{
				Type:   flex.GoStringToFramework(v.Items.Type),
				Format: flex.GoStringToFramework(v.Items.Format),
			}
		}

		bm.CalculationProperties[k] = *calculationPropertyModel

	}
}
