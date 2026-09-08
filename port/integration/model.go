package integration

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

type WebhookChangelogDestinationModel struct {
	Url   types.String `tfsdk:"url"`
	Agent types.Bool   `tfsdk:"agent"`
}

type IntegrationModel struct {
	ID                          types.String                      `tfsdk:"id"`
	InstallationId              types.String                      `tfsdk:"installation_id"`
	InstallationAppType         types.String                      `tfsdk:"installation_app_type"`
	InstallationType            types.String                      `tfsdk:"installation_type"`
	Title                       types.String                      `tfsdk:"title"`
	Version                     types.String                      `tfsdk:"version"`
	Config                      types.String                      `tfsdk:"config"`
	Spec                        types.String                      `tfsdk:"spec"`
	Status                      types.String                      `tfsdk:"status"`
	KafkaChangelogDestination   types.Object                      `tfsdk:"kafka_changelog_destination"`
	WebhookChangelogDestination *WebhookChangelogDestinationModel `tfsdk:"webhook_changelog_destination"`
}

func (m *IntegrationModel) installationType() string {
	if m.InstallationType.IsNull() || m.InstallationType.ValueString() == "" {
		return consts.InstallationTypeOnPrem
	}
	return m.InstallationType.ValueString()
}

func (m *IntegrationModel) isSaas() bool {
	return consts.IsSaas(m.installationType())
}
