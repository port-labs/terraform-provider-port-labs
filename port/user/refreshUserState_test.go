package user

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshUserState_InactivityTimeout(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	timeout := 45

	state := &UserModel{
		Email: types.StringValue("user@example.com"),
	}

	err := refreshUserState(ctx, state, &cli.User{
		Email:             "user@example.com",
		Status:            "Active",
		InactivityTimeout: &timeout,
		ManagedByScim:     &[]bool{false}[0],
	})
	require.NoError(t, err)
	assert.Equal(t, int64(45), state.InactivityTimeout.ValueInt64())
	assert.Equal(t, "Active", state.Status.ValueString())
	assert.False(t, state.ManagedByScim.ValueBool())
}

func TestRefreshUserState_NullInactivityTimeout(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	state := &UserModel{
		Email: types.StringValue("user@example.com"),
	}

	err := refreshUserState(ctx, state, &cli.User{
		Email:  "user@example.com",
		Status: "Active",
	})
	require.NoError(t, err)
	assert.True(t, state.InactivityTimeout.IsNull())
}
