package team

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
)

func TeamResourceToPortBody(ctx context.Context, state *TeamModel) (*cli.Team, error) {
	tp := &cli.Team{
		Name:        state.Name.ValueString(),
		Description: state.Description.ValueString(),
	}
	if state.Users != nil {
		tp.Users = flex.TerraformStringListToGoArray(state.Users)
	}

	return tp, nil
}
