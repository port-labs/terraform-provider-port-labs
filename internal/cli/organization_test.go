package cli

import (
	"testing"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestOrganizationUpdate_ToPatchBody(t *testing.T) {
	t.Parallel()

	t.Run("includes inactivity timeout value", func(t *testing.T) {
		t.Parallel()

		update := &OrganizationUpdate{
			InactivityTimeout:        utils.PtrTo(15),
			IncludeInactivityTimeout: true,
		}

		body := update.ToPatchBody()
		assert.Equal(t, 15, body["inactivityTimeout"])
	})

	t.Run("includes null inactivity timeout", func(t *testing.T) {
		t.Parallel()

		update := &OrganizationUpdate{
			IncludeInactivityTimeout: true,
		}

		body := update.ToPatchBody()
		assert.Nil(t, body["inactivityTimeout"])
	})

	t.Run("omits inactivity timeout when not included", func(t *testing.T) {
		t.Parallel()

		update := &OrganizationUpdate{
			Name: utils.PtrTo("example"),
		}

		body := update.ToPatchBody()
		assert.Equal(t, "example", body["name"])
		_, ok := body["inactivityTimeout"]
		assert.False(t, ok)
	})
}
