package team

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/internal/flex"
)

func refreshTeamState(ctx context.Context, state *TeamModel, t *cli.Team) error {
	state.CreatedAt = types.StringValue(t.CreatedAt.String())
	state.UpdatedAt = types.StringValue(t.UpdatedAt.String())
	state.ID = types.StringValue(t.Name)
	state.Name = types.StringValue(t.Name)
	state.Description = flex.GoStringToFramework(&t.Description)
	state.ProviderName = flex.GoStringToFramework(&t.Provider)

	if len(t.Users) != 0 {
		state.Users = make([]types.String, len(t.Users))
		for i, u := range t.Users {
			state.Users[i] = types.StringValue(u)
		}
	}

	return nil
}
