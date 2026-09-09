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
	if want := (types.ObjectType{AttrTypes: webhookChangelogDestinationType}); !webhook.GetType().Equal(want) {
		t.Errorf("webhook_changelog_destination type = %s, want %s", webhook.GetType(), want)
	}

	kafka, ok := attrs["kafka_changelog_destination"].(schema.ObjectAttribute)
	if !ok || !kafka.Computed {
		t.Fatal("kafka_changelog_destination must be computed")
	}
	if want := (types.ObjectType{AttrTypes: kafkaChangelogDestinationType}); !kafka.GetType().Equal(want) {
		t.Errorf("kafka_changelog_destination type = %s, want %s", kafka.GetType(), want)
	}
}

func TestIntegrationToPortBodySkipsUnknownChangelogDestination(t *testing.T) {
	plan := &IntegrationModel{
		InstallationId:              types.StringValue("my-kafka"),
		InstallationAppType:         types.StringValue("kafka"),
		Version:                     types.StringValue("1.0.0"),
		WebhookChangelogDestination: types.ObjectUnknown(webhookChangelogDestinationType),
		KafkaChangelogDestination:   types.ObjectUnknown(kafkaChangelogDestinationType),
	}

	body, err := integrationToPortBody(plan)
	if err != nil {
		t.Fatalf("integrationToPortBody: %v", err)
	}
	if body.ChangelogDestination != nil {
		t.Fatalf("changelog destination = %+v, want nil", body.ChangelogDestination)
	}
}

func TestChangelogDestinationRoundTrip(t *testing.T) {
	const portManagedURL = "https://internal.port.io/webhooks/changelog"
	agent := true

	noKafka := types.ObjectNull(kafkaChangelogDestinationType)
	noWebhook := types.ObjectNull(webhookChangelogDestinationType)
	webhookObject := func(url string, agent *bool) types.Object {
		return types.ObjectValueMust(webhookChangelogDestinationType, map[string]attr.Value{
			"url":   types.StringValue(url),
			"agent": types.BoolPointerValue(agent),
		})
	}

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
			wantKafka:   noKafka,
			wantWebhook: noWebhook,
		},
		{
			name:        "empty destination",
			dest:        &cli.ChangelogDestination{},
			wantKafka:   noKafka,
			wantWebhook: noWebhook,
		},
		{
			name:        "kafka",
			dest:        &cli.ChangelogDestination{Type: consts.Kafka},
			wantKafka:   types.ObjectValueMust(kafkaChangelogDestinationType, map[string]attr.Value{}),
			wantWebhook: noWebhook,
			wantBody:    &cli.ChangelogDestination{Type: consts.Kafka},
		},
		{
			name:        "port-managed webhook, no agent",
			dest:        &cli.ChangelogDestination{Type: consts.Webhook, Url: portManagedURL},
			wantKafka:   noKafka,
			wantWebhook: webhookObject(portManagedURL, nil),
			wantBody:    &cli.ChangelogDestination{Type: consts.Webhook, Url: portManagedURL},
		},
		{
			name:        "webhook through the agent",
			dest:        &cli.ChangelogDestination{Type: consts.Webhook, Url: "https://google.com", Agent: &agent},
			wantKafka:   noKafka,
			wantWebhook: webhookObject("https://google.com", &agent),
			wantBody:    &cli.ChangelogDestination{Type: consts.Webhook, Url: "https://google.com", Agent: &agent},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &IntegrationModel{InstallationId: types.StringValue("my-gitlab")}

			(&IntegrationResource{}).refreshIntegrationState(state, &cli.Integration{ChangelogDestination: tt.dest}, "my-gitlab")

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
