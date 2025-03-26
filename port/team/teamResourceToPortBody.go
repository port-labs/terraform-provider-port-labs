package team

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func TeamResourceToPortBody(ctx context.Context, state *TeamModel) (*cli.Team, error) {
	portTeam := cli.PortTeam{
		Name: state.Name.ValueString(),
	}

	if !state.Description.IsNull() {
		description := state.Description.ValueString()
		portTeam.Description = &description
	}

	if state.Users != nil {
		portTeam.Users = make([]string, len(state.Users))
		for i, t := range state.Users {
			portTeam.Users[i] = t.ValueString()
		}
	}

	return &cli.Team{
		PortTeam:   portTeam,
		Identifier: state.Identifier.ValueStringPointer(),
	}, nil
}
