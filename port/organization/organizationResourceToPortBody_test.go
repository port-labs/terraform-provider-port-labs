package organization

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganizationResourceToPortBody_IdleTimeMS(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("includes idle time value when set", func(t *testing.T) {
		t.Parallel()

		state := &OrganizationModel{
			IdleTimeMS: types.Int64Value(1_800_000),
		}

		update, err := organizationResourceToPortBody(ctx, state)
		require.NoError(t, err)
		require.True(t, update.IncludeIdleTimeMS)
		require.NotNil(t, update.IdleTimeMS)
		assert.Equal(t, 1_800_000, *update.IdleTimeMS)
	})

	t.Run("includes null when cleared", func(t *testing.T) {
		t.Parallel()

		state := &OrganizationModel{
			IdleTimeMS: types.Int64Null(),
		}

		update, err := organizationResourceToPortBody(ctx, state)
		require.NoError(t, err)
		require.True(t, update.IncludeIdleTimeMS)
		assert.Nil(t, update.IdleTimeMS)
	})

	t.Run("omits idle time when unknown", func(t *testing.T) {
		t.Parallel()

		state := &OrganizationModel{
			IdleTimeMS: types.Int64Unknown(),
		}

		update, err := organizationResourceToPortBody(ctx, state)
		require.NoError(t, err)
		assert.False(t, update.IncludeIdleTimeMS)
	})
}
