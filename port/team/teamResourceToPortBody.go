package team

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
)

func TeamResourceToPortBody(ctx context.Context, state *TeamModel) (*cli.Team, error) {
	w := &cli.Team{
		Name:        state.Name.ValueString(),
		Description: state.Description.ValueString(),
	}
	if state.Users != nil {
		w.Users = make([]string, len(state.Users))
		for i, t := range state.Users {
			w.Users[i] = t.ValueString()
		}
	}

	return w, nil
}
