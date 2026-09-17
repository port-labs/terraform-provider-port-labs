package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
)

func TestCreateAwaitParams(t *testing.T) {
	saas := &IntegrationModel{InstallationType: types.StringValue(consts.InstallationTypeSaas)}
	waitReady, waitProvisioned := createAwaitParams(saas)
	assert.True(t, waitReady, "SaaS create should wait for operation ready")
	assert.True(t, waitProvisioned, "SaaS create should wait for default resource provisioning")

	onPrem := &IntegrationModel{InstallationType: types.StringValue(consts.InstallationTypeOnPrem)}
	waitReady, waitProvisioned = createAwaitParams(onPrem)
	assert.False(t, waitReady, "self-hosted create should not wait for SaaS operation status")
	assert.True(t, waitProvisioned, "self-hosted create should wait for default resource provisioning")

	defaultType := &IntegrationModel{}
	waitReady, waitProvisioned = createAwaitParams(defaultType)
	assert.False(t, waitReady)
	assert.True(t, waitProvisioned)
}

