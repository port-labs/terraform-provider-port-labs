package cli

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
)

func TestAdoptExisting(t *testing.T) {
	testCases := []struct {
		name             string
		providerAdoptees AdoptionPolicy
		resourceSetting  types.Bool
		expected         bool
	}{
		{"unset and type not listed", nil, types.BoolNull(), false},
		{"unset and type listed", AdoptionPolicy{consts.ResourceKindPage: true}, types.BoolNull(), true},
		{"unset and only another type listed", AdoptionPolicy{consts.ResourceKindBlueprint: true}, types.BoolNull(), false},
		{"unknown falls back to provider", AdoptionPolicy{consts.ResourceKindPage: true}, types.BoolUnknown(), true},
		{"resource true overrides type not listed", nil, types.BoolValue(true), true},
		{"resource false overrides type listed", AdoptionPolicy{consts.ResourceKindPage: true}, types.BoolValue(false), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := &PortClient{AdoptExistingResources: tc.providerAdoptees}
			assert.Equal(t, tc.expected, c.AdoptExisting(consts.ResourceKindPage, tc.resourceSetting))
		})
	}
}

func TestResourceExistsError(t *testing.T) {
	pageErr := &ResourceExistsError{ResourceKind: consts.ResourceKindPage, Identifier: "services"}
	assert.Equal(t, "page already exists", pageErr.Summary())
	assert.Contains(t, pageErr.Error(), `page "services" already exists in Port.`)
	assert.Contains(t, pageErr.Error(), "`terraform import <address> \"services\"`")
	assert.Contains(t, pageErr.Error(), `add "port_page" to the provider's`)

	relationErr := &ResourceExistsError{ResourceKind: consts.ResourceKindBlueprintRelation, Identifier: "group", Blueprint: "_team"}
	assert.Equal(t, "relation already exists", relationErr.Summary())
	assert.Contains(t, relationErr.Error(), `relation "group" already exists on blueprint "_team".`)
	assert.Contains(t, relationErr.Error(), "Remove it from the blueprint's relations block")
	assert.Contains(t, relationErr.Error(), "`terraform import <address> \"_team:group\"`")

	var existsErr *ResourceExistsError
	assert.True(t, errors.As(fmt.Errorf("wrapped: %w", pageErr), &existsErr))
}

func TestAdoptionConflictError(t *testing.T) {
	conflictErr := &AdoptionConflictError{
		ResourceKind: consts.ResourceKindIntegration,
		Identifier:   "my-gitlab",
		Field:        "installation_type",
		Existing:     "OnPrem",
		Configured:   "Saas",
	}
	assert.Equal(t, "cannot take over existing integration", conflictErr.Summary())
	assert.Contains(t, conflictErr.Error(), "integration \"my-gitlab\" already exists in Port with `installation_type` = \"OnPrem\", but the configuration sets \"Saas\".")
	assert.Contains(t, conflictErr.Error(), "Set `installation_type` to \"OnPrem\"")
}
