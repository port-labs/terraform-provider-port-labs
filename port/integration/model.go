package integration

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	kafkaChangelogDestinationType   = map[string]attr.Type{}
	webhookChangelogDestinationType = map[string]attr.Type{
		"url":   types.StringType,
		"agent": types.BoolType,
	}
)

type IntegrationModel struct {
	ID                          types.String `tfsdk:"id"`
	InstallationId              types.String `tfsdk:"installation_id"`
	InstallationAppType         types.String `tfsdk:"installation_app_type"`
	Title                       types.String `tfsdk:"title"`
	Version                     types.String `tfsdk:"version"`
	Config                      types.String `tfsdk:"config"`
	KafkaChangelogDestination   types.Object `tfsdk:"kafka_changelog_destination"`
	WebhookChangelogDestination types.Object `tfsdk:"webhook_changelog_destination"`
}
