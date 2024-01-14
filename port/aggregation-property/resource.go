package aggregation_property

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"strings"
)

var _ resource.Resource = &AggregationPropertyResource{}
var _ resource.ResourceWithImportState = &AggregationPropertyResource{}

func NewAggregationPropertyResource() resource.Resource {
	return &AggregationPropertyResource{}
}

type AggregationPropertyResource struct {
	portClient *cli.PortClient
}

func (r *AggregationPropertyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aggregation_property"
}

func (r *AggregationPropertyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.portClient = req.ProviderData.(*cli.PortClient)
}

func (r *AggregationPropertyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError("invalid import ID", "import ID must be in the format <blueprint_identifier>:<aggregation_property_identifier>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("blueprint_identifier"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("aggregation_identifier"), idParts[1])...)

}

func (r *AggregationPropertyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *AggregationPropertyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	blueprintIdentifier := state.BlueprintIdentifier.ValueString()
	aggregationPropertyIdentifier := state.AggregationIdentifier.ValueString()

	b, statusCode, err := r.portClient.ReadBlueprint(ctx, blueprintIdentifier)

	if err != nil {
		if statusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
		return
	}

	aggregationProperty, ok := b.AggregationProperties[aggregationPropertyIdentifier]
	// another way to check if aggregationProperty exists
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}

	err = refreshAggregationPropertyState(state, aggregationProperty, blueprintIdentifier, aggregationPropertyIdentifier)
	if err != nil {
		resp.Diagnostics.AddError("failed writing aggregation property fields to resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AggregationPropertyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *AggregationPropertyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	aggr, err := aggregationPropertyToBody(state)

	if err != nil {
		resp.Diagnostics.AddError("failed to convert aggregation property to port valid request", err.Error())
		return
	}

	b, statusCode, err := r.portClient.ReadBlueprint(ctx, state.BlueprintIdentifier.ValueString())

	if err != nil {
		if statusCode == 404 {
			resp.Diagnostics.AddError("Blueprint doesn't exists, it is required to create aggregation property", err.Error())
			return
		}
		resp.Diagnostics.AddError("failed reading blueprint", err.Error())
	}

	_, ok := b.AggregationProperties[state.AggregationIdentifier.ValueString()]
	if ok {
		resp.Diagnostics.AddError("Aggregation property already exists", `Aggregation property with identifier "`+state.AggregationIdentifier.ValueString()+`" already exists`)
		return
	}

	b.AggregationProperties[state.AggregationIdentifier.ValueString()] = *aggr

	bp, err := r.portClient.UpdateBlueprint(ctx, b, state.BlueprintIdentifier.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to create aggregation property", err.Error())
		return
	}

	aggregationProperty, ok := bp.AggregationProperties[state.AggregationIdentifier.ValueString()]
	if !ok {
		resp.Diagnostics.AddError("failed to create aggregation property", `Aggregation property with identifier "`+state.AggregationIdentifier.ValueString()+`" doesn't exists after creation`)
		return
	}

	err = refreshAggregationPropertyState(state, aggregationProperty, state.BlueprintIdentifier.ValueString(), state.AggregationIdentifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed writing aggregation property fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AggregationPropertyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state *AggregationPropertyModel
	var previousState *AggregationPropertyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	aggr, err := aggregationPropertyToBody(state)

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

	b.AggregationProperties[state.AggregationIdentifier.ValueString()] = *aggr

	bp, err := r.portClient.UpdateBlueprint(ctx, b, state.BlueprintIdentifier.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to update aggregation property", err.Error())
		return
	}

	aggregationProperty, ok := bp.AggregationProperties[state.AggregationIdentifier.ValueString()]
	if !ok {
		resp.Diagnostics.AddError("failed to update aggregation property", `Aggregation property with identifier "`+state.AggregationIdentifier.ValueString()+`" doesn't exists after update`)
		return
	}

	err = refreshAggregationPropertyState(state, aggregationProperty, state.BlueprintIdentifier.ValueString(), state.AggregationIdentifier.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed writing aggregation property fields to resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AggregationPropertyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *AggregationPropertyModel

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

	delete(b.AggregationProperties, state.AggregationIdentifier.ValueString())

	bp, err := r.portClient.UpdateBlueprint(ctx, b, state.BlueprintIdentifier.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to delete aggregation property", err.Error())
		return
	}

	_, ok := bp.AggregationProperties[state.AggregationIdentifier.ValueString()]
	if ok {
		resp.Diagnostics.AddError("failed to delete aggregation property", `Aggregation property with identifier "`+state.AggregationIdentifier.ValueString()+`" still exists after deletion`)
		return
	}

	resp.State.RemoveResource(ctx)
}
