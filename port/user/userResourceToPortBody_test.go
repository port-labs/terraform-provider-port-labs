package user

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserResourceToPortBody_InactivityTimeout(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("includes timeout value when set", func(t *testing.T) {
		t.Parallel()

		state := &UserModel{
			InactivityTimeout: types.Int64Value(30),
		}

		update, err := userResourceToPortBody(ctx, state)
		require.NoError(t, err)
		require.True(t, update.IncludeInactivityTimeout)
		require.NotNil(t, update.InactivityTimeout)
		assert.Equal(t, 30, *update.InactivityTimeout)
	})

	t.Run("includes null when cleared", func(t *testing.T) {
		t.Parallel()

		state := &UserModel{
			InactivityTimeout: types.Int64Null(),
		}

		update, err := userResourceToPortBody(ctx, state)
		require.NoError(t, err)
		require.True(t, update.IncludeInactivityTimeout)
		assert.Nil(t, update.InactivityTimeout)
	})

	t.Run("omits timeout when unknown", func(t *testing.T) {
		t.Parallel()

		state := &UserModel{
			InactivityTimeout: types.Int64Unknown(),
		}

		update, err := userResourceToPortBody(ctx, state)
		require.NoError(t, err)
		assert.False(t, update.IncludeInactivityTimeout)
	})
}
