package cli

import (
	"testing"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestOrganizationUpdate_ToPatchBody(t *testing.T) {
	t.Parallel()

	t.Run("includes idle time value", func(t *testing.T) {
		t.Parallel()

		update := &OrganizationUpdate{
			IdleTimeMS:        utils.PtrTo(1_800_000),
			IncludeIdleTimeMS: true,
		}

		body := update.ToPatchBody()
		assert.Equal(t, 1_800_000, body["idleTimeMS"])
	})

	t.Run("includes null idle time", func(t *testing.T) {
		t.Parallel()

		update := &OrganizationUpdate{
			IncludeIdleTimeMS: true,
		}

		body := update.ToPatchBody()
		assert.Nil(t, body["idleTimeMS"])
	})

	t.Run("omits idle time when not included", func(t *testing.T) {
		t.Parallel()

		update := &OrganizationUpdate{
			Name: utils.PtrTo("example"),
		}

		body := update.ToPatchBody()
		assert.Equal(t, "example", body["name"])
		_, ok := body["idleTimeMS"]
		assert.False(t, ok)
	})
}
