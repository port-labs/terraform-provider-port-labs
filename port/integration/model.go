package integration

import "github.com/hashicorp/terraform-plugin-framework/types"

type WebhookChangelogDestinationModel struct {
	Url   types.String `tfsdk:"url"`
	Agent types.Bool   `tfsdk:"agent"`
}

type IntegrationModel struct {
	ID                          types.String                      `tfsdk:"id"`
	InstallationId              types.String                      `tfsdk:"installation_id"`
	InstallationAppType         types.String                      `tfsdk:"installation_app_type"`
	Title                       types.String                      `tfsdk:"title"`
	Version                     types.String                      `tfsdk:"version"`
	Config                      types.String                      `tfsdk:"config"`
	KafkaChangelogDestination   types.Object                      `tfsdk:"kafka_changelog_destination"`
	WebhookChangelogDestination *WebhookChangelogDestinationModel `tfsdk:"webhook_changelog_destination"`
}
