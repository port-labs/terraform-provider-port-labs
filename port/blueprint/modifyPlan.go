package blueprint

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/adoption"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

var _ resource.ResourceWithModifyPlan = &BlueprintResource{}

// ModifyPlan checks, before apply, whether a blueprint about to be created already exists in Port, so that
// taking it over is never silent: it fails unless adopt_existing is set, and warns when it is.
func (r *BlueprintResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if r.portClient == nil || !req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var identifier types.String
	var adoptExisting types.Bool
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("identifier"), &identifier)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("adopt_existing"), &adoptExisting)...)
	if resp.Diagnostics.HasError() || identifier.IsUnknown() {
		return
	}

	_, statusCode, err := r.portClient.ReadBlueprint(ctx, identifier.ValueString())
	if err != nil {
		if statusCode != http.StatusNotFound {
			resp.Diagnostics.AddWarning("failed to check whether the blueprint already exists", err.Error())
		}
		return
	}

	existing := &cli.ResourceExistsError{ResourceKind: consts.ResourceKindBlueprint, Identifier: identifier.ValueString()}
	adoption.TakeOver(adoption.AtPlan, r.portClient, existing, adoptExisting, nil, &resp.Diagnostics)
}
