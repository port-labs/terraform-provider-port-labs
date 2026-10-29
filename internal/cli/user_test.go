package cli

import (
	"testing"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestUserUpdate_ToPatchBody(t *testing.T) {
	t.Parallel()

	t.Run("includes inactivity timeout value", func(t *testing.T) {
		t.Parallel()

		update := &UserUpdate{
			InactivityTimeout:        utils.PtrTo(15),
			IncludeInactivityTimeout: true,
		}

		body := update.ToPatchBody()
		assert.Equal(t, 15, body["inactivityTimeout"])
	})

	t.Run("includes null inactivity timeout", func(t *testing.T) {
		t.Parallel()

		update := &UserUpdate{
			IncludeInactivityTimeout: true,
		}

		body := update.ToPatchBody()
		assert.Nil(t, body["inactivityTimeout"])
	})

	t.Run("includes roles and teams when set", func(t *testing.T) {
		t.Parallel()

		update := &UserUpdate{
			Roles:        []string{"Admin"},
			Teams:        []string{"engineering"},
			IncludeRoles: true,
			IncludeTeams: true,
		}

		body := update.ToPatchBody()
		assert.Equal(t, []string{"Admin"}, body["roles"])
		assert.Equal(t, []string{"engineering"}, body["teams"])
	})

	t.Run("omits unset fields", func(t *testing.T) {
		t.Parallel()

		update := &UserUpdate{}

		body := update.ToPatchBody()
		_, hasRoles := body["roles"]
		_, hasTeams := body["teams"]
		_, hasTimeout := body["inactivityTimeout"]
		assert.False(t, hasRoles)
		assert.False(t, hasTeams)
		assert.False(t, hasTimeout)
	})
}
