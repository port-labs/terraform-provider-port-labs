package integration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithValidateConfig = &IntegrationResource{}

func (r *IntegrationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model IntegrationModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateIntegrationModel(&model); err != nil {
		resp.Diagnostics.AddError("invalid integration configuration", err.Error())
	}
}
