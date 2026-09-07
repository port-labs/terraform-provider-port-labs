package integration

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

func TestChangelogDestinationSchema(t *testing.T) {
	attrs := IntegrationSchema()

	webhook, ok := attrs["webhook_changelog_destination"].(schema.SingleNestedAttribute)
	if !ok || !webhook.Computed {
		t.Fatal("webhook_changelog_destination must be computed")
	}

	kafka, ok := attrs["kafka_changelog_destination"].(schema.ObjectAttribute)
	if !ok || !kafka.Computed {
		t.Fatal("kafka_changelog_destination must be computed")
	}

	wantWebhookType := types.ObjectType{AttrTypes: webhookChangelogDestinationType}
	if got := webhook.GetType(); !got.Equal(wantWebhookType) {
		t.Errorf("webhook_changelog_destination type = %s, want %s", got, wantWebhookType)
	}

	wantKafkaType := types.ObjectType{AttrTypes: kafkaChangelogDestinationType}
	if got := kafka.GetType(); !got.Equal(wantKafkaType) {
		t.Errorf("kafka_changelog_destination type = %s, want %s", got, wantKafkaType)
	}
}

func TestChangelogDestinationRoundTrip(t *testing.T) {
	const portManagedURL = "https://internal.port.io/webhooks/changelog"
	agent := true

	tests := []struct {
		name        string
		dest        *cli.ChangelogDestination
		wantKafka   types.Object
		wantWebhook types.Object
		wantBody    *cli.ChangelogDestination
	}{
		{
			name:        "no destination",
			dest:        nil,
			wantKafka:   types.ObjectNull(kafkaChangelogDestinationType),
			wantWebhook: types.ObjectNull(webhookChangelogDestinationType),
		},
		{
			name:        "empty destination",
			dest:        &cli.ChangelogDestination{},
			wantKafka:   types.ObjectNull(kafkaChangelogDestinationType),
			wantWebhook: types.ObjectNull(webhookChangelogDestinationType),
		},
		{
			name:        "kafka",
			dest:        &cli.ChangelogDestination{Type: consts.Kafka},
			wantKafka:   types.ObjectValueMust(kafkaChangelogDestinationType, map[string]attr.Value{}),
			wantWebhook: types.ObjectNull(webhookChangelogDestinationType),
			wantBody:    &cli.ChangelogDestination{Type: consts.Kafka},
		},
		{
			name:      "port-managed webhook, no agent",
			dest:      &cli.ChangelogDestination{Type: consts.Webhook, Url: portManagedURL},
			wantKafka: types.ObjectNull(kafkaChangelogDestinationType),
			wantWebhook: types.ObjectValueMust(webhookChangelogDestinationType, map[string]attr.Value{
				"url":   types.StringValue(portManagedURL),
				"agent": types.BoolNull(),
			}),
			wantBody: &cli.ChangelogDestination{Type: consts.Webhook, Url: portManagedURL},
		},
		{
			name:      "webhook through the agent",
			dest:      &cli.ChangelogDestination{Type: consts.Webhook, Url: "https://google.com", Agent: &agent},
			wantKafka: types.ObjectNull(kafkaChangelogDestinationType),
			wantWebhook: types.ObjectValueMust(webhookChangelogDestinationType, map[string]attr.Value{
				"url":   types.StringValue("https://google.com"),
				"agent": types.BoolValue(true),
			}),
			wantBody: &cli.ChangelogDestination{Type: consts.Webhook, Url: "https://google.com", Agent: &agent},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &IntegrationModel{InstallationId: types.StringValue("my-gitlab")}

			err := (&IntegrationResource{}).refreshIntegrationState(state, &cli.Integration{ChangelogDestination: tt.dest}, "my-gitlab")
			if err != nil {
				t.Fatalf("refreshIntegrationState: %v", err)
			}

			if !state.KafkaChangelogDestination.Equal(tt.wantKafka) {
				t.Errorf("kafka_changelog_destination = %s, want %s", state.KafkaChangelogDestination, tt.wantKafka)
			}
			if !state.WebhookChangelogDestination.Equal(tt.wantWebhook) {
				t.Errorf("webhook_changelog_destination = %s, want %s", state.WebhookChangelogDestination, tt.wantWebhook)
			}

			body, err := integrationToPortBody(state)
			if err != nil {
				t.Fatalf("integrationToPortBody: %v", err)
			}
			if !reflect.DeepEqual(body.ChangelogDestination, tt.wantBody) {
				t.Errorf("changelog destination sent to Port = %+v, want %+v", body.ChangelogDestination, tt.wantBody)
			}
		})
	}
}
