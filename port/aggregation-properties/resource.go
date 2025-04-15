package aggregation_properties

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

var _ resource.Resource = &AggregationPropertiesResource{}
var _ resource.ResourceWithImportState = &AggregationPropertiesResource{}

func NewAggregationPropertiesResource() resource.Resource {
	return &AggregationPropertiesResource{}
}

type AggregationPropertiesResource struct {
	portClient *cli.PortClient
}

func (r *AggregationPropertiesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aggregation_properties"
}

func (r *AggregationPropertiesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *AggregationPropertiesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("blueprint_identifier"), req.ID,
	)...)

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), req.ID,
	)...)
}

func (r *AggregationPropertiesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *AggregationPropertiesModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	b, statusCode, err := r.portClient.ReadBlueprint(ctx, state.BlueprintIdentifier.ValueString())

	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
		return
	}

	err = r.refreshAggregationPropertiesState(state, b.AggregationProperties)
	if err != nil {
		resp.Diagnostics.AddError("failed writing aggregation property fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AggregationPropertiesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *AggregationPropertiesModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	aggregationProperties, err := aggregationPropertiesToBody(state)

	if err != nil {
		resp.Diagnostics.AddError("failed to convert aggregation property to port valid request", err.Error())
		return
	}

	b, statusCode, err := r.portClient.ReadBlueprint(ctx, state.BlueprintIdentifier.ValueString())

	if err != nil {
		if statusCode == 404 {
			resp.Diagnostics.AddError("Blueprint doesn't exists, it is required to create aggregation properties", err.Error())
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
	}

	// check if the aggregation properties already exists
	if b.AggregationProperties != nil {
		for aggregationPropertyIdentifier := range *aggregationProperties {
			if _, ok := b.AggregationProperties[aggregationPropertyIdentifier]; ok {
				resp.Diagnostics.AddError("aggregation property already exists", aggregationPropertyIdentifier)
				return
			}
		}
	}

	b.AggregationProperties = *aggregationProperties

	_, err = r.portClient.UpdateBlueprint(ctx, b, state.BlueprintIdentifier.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to create aggregation properties", err.Error())
		return
	}

	// set the ID to the blueprint identifier
	state.ID = state.BlueprintIdentifier

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AggregationPropertiesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *AggregationPropertiesModel
	var previousState *AggregationPropertiesModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	aggregationProperties, err := aggregationPropertiesToBody(state)

	if err != nil {
		resp.Diagnostics.AddError("failed to convert aggregation property to port valid request", err.Error())
		return
	}

	b, statusCode, err := r.portClient.ReadBlueprint(ctx, state.BlueprintIdentifier.ValueString())

	if err != nil {
		if statusCode == 404 {
			resp.Diagnostics.AddError("Blueprint doesn't exists, it is required to update the aggregation property", err.Error())
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
	}

	b.AggregationProperties = *aggregationProperties

	_, err = r.portClient.UpdateBlueprint(ctx, b, state.BlueprintIdentifier.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to update aggregation property", err.Error())
		return
	}

	// set the ID to the blueprint identifier
	state.ID = state.BlueprintIdentifier

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AggregationPropertiesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *AggregationPropertiesModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	b, statusCode, err := r.portClient.ReadBlueprint(ctx, state.BlueprintIdentifier.ValueString())

	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
	}

	b.AggregationProperties = make(map[string]cli.BlueprintAggregationProperty)

	_, err = r.portClient.UpdateBlueprint(ctx, b, state.BlueprintIdentifier.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to delete aggregation property", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
