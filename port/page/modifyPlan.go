package page

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

var _ resource.ResourceWithModifyPlan = &PageResource{}

// ModifyPlan checks, before apply, whether a page about to be created already exists in Port, so that
// taking it over is never silent: it fails unless adopt_existing is set, and warns when it is.
func (r *PageResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if r.portClient == nil || !req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var identifier, pageType types.String
	var adoptExisting types.Bool
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("identifier"), &identifier)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("type"), &pageType)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("adopt_existing"), &adoptExisting)...)
	if resp.Diagnostics.HasError() || identifier.IsUnknown() || pageType.ValueString() == entityPageType {
		return
	}

	existing, statusCode, err := r.portClient.GetPage(ctx, identifier.ValueString())
	if err != nil {
		if statusCode != http.StatusNotFound {
			resp.Diagnostics.AddWarning("failed to check whether the page already exists", err.Error())
		}
		return
	}

	adoption.TakeOver(adoption.AtPlan, r.portClient, existingPage(identifier.ValueString()), adoptExisting,
		pageTypeConflict(identifier.ValueString(), existing.Type, pageType.ValueString()), &resp.Diagnostics)
}

func existingPage(identifier string) *cli.ResourceExistsError {
	return &cli.ResourceExistsError{ResourceKind: consts.ResourceKindPage, Identifier: identifier}
}

// pageTypeConflict reports a page that cannot be taken over because its type differs from the configuration, as
// a page's type cannot be changed
func pageTypeConflict(identifier, existingType, configuredType string) *cli.AdoptionConflictError {
	if existingType == configuredType {
		return nil
	}
	return &cli.AdoptionConflictError{
		ResourceKind: consts.ResourceKindPage,
		Identifier:   identifier,
		Field:        "type",
		Existing:     existingType,
		Configured:   configuredType,
	}
}
