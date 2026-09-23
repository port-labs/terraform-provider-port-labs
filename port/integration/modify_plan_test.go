package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestInstallationAppTypeImmutableWhenOmittedFromPlan(t *testing.T) {
	state := &IntegrationModel{InstallationAppType: types.StringValue("kafka")}
	plan := &IntegrationModel{InstallationAppType: types.StringNull()}

	field := immutableFields[1]
	assert.False(t, field.changed(state, plan), "null planned app type should not count as a change")

	plan.InstallationAppType = types.StringUnknown()
	assert.False(t, field.changed(state, plan), "unknown planned app type should not count as a change")
}

func TestInstallationAppTypeImmutableWhenChanged(t *testing.T) {
	state := &IntegrationModel{InstallationAppType: types.StringValue("kafka")}
	plan := &IntegrationModel{InstallationAppType: types.StringValue("gitlab")}

	field := immutableFields[1]
	assert.True(t, field.changed(state, plan))
}
