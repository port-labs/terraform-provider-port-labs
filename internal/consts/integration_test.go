package consts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallationTypeHelpers(t *testing.T) {
	assert.True(t, IsSaas(InstallationTypeSaas))
	assert.False(t, IsSaas(InstallationTypeSaasOAuth2))
	assert.False(t, IsSaas(InstallationTypeOnPrem))

	assert.True(t, IsSaasOAuth2(InstallationTypeSaasOAuth2))
	assert.False(t, IsSaasOAuth2(InstallationTypeSaas))

	assert.True(t, IsHosted(InstallationTypeSaas))
	assert.True(t, IsHosted(InstallationTypeSaasOAuth2))
	assert.False(t, IsHosted(InstallationTypeOnPrem))
}
