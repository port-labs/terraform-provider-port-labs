package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

func TestPortManagedChangelogDestination(t *testing.T) {
	attrs := IntegrationSchema()

	webhook, ok := attrs["webhook_changelog_destination"].(schema.SingleNestedAttribute)
	if !ok || !webhook.Computed {
		t.Fatal("webhook_changelog_destination must be computed")
	}

	kafka, ok := attrs["kafka_changelog_destination"].(schema.ObjectAttribute)
	if !ok || !kafka.Computed {
		t.Fatal("kafka_changelog_destination must be computed")
	}

	state := &IntegrationModel{InstallationId: types.StringValue("my-gitlab")}
	err := (&IntegrationResource{}).refreshIntegrationState(state, &cli.Integration{
		ChangelogDestination: &cli.ChangelogDestination{
			Type: consts.Webhook,
			Url:  "https://internal.port.io/webhooks/changelog",
		},
	}, "my-gitlab")
	if err != nil {
		t.Fatalf("refreshIntegrationState failed: %v", err)
	}

	if state.WebhookChangelogDestination == nil {
		t.Fatal("expected port-managed webhook changelog destination in state")
	}
	if state.WebhookChangelogDestination.Url.ValueString() != "https://internal.port.io/webhooks/changelog" {
		t.Fatalf("unexpected webhook url: %s", state.WebhookChangelogDestination.Url.ValueString())
	}
}
