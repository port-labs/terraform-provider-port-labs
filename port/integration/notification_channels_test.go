package integration

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyNotificationChannelsPostsEnabledAndKeepsAudience(t *testing.T) {
	var posts int
	var posted cli.OptInSetting
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`{"ok":true,"data":[{"channel":"email","enabled":false,"emails":["admin@example.com"],"teams":["ops"]}]}`))
		case http.MethodPost:
			posts++
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(body, &posted))
			_, _ = w.Write([]byte(`{"ok":true,"data":{"channel":"email","enabled":true}}`))
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	}))
	t.Cleanup(server.Close)

	resource := &IntegrationResource{portClient: &cli.PortClient{Client: resty.New().SetBaseURL(server.URL)}}
	err := resource.applyNotificationChannels(context.Background(), "github-prod", []NotificationChannelModel{
		{Channel: types.StringValue(consts.NotificationChannelEmail), Enabled: types.BoolValue(true)},
	})
	require.NoError(t, err)
	assert.Equal(t, 1, posts)
	assert.Equal(t, cli.OptInSetting{
		Channel: consts.NotificationChannelEmail,
		Enabled: true,
		Emails:  []string{"admin@example.com"},
		Teams:   []string{"ops"},
	}, posted)
}

func TestApplyNotificationChannelsSkipsUnchangedChannel(t *testing.T) {
	posts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			posts++
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"channel":"slack","enabled":true,"emails":["admin@example.com"],"teams":[]}]}`))
	}))
	t.Cleanup(server.Close)

	resource := &IntegrationResource{portClient: &cli.PortClient{Client: resty.New().SetBaseURL(server.URL)}}
	err := resource.applyNotificationChannels(context.Background(), "github-prod", []NotificationChannelModel{
		{Channel: types.StringValue(consts.NotificationChannelSlack), Enabled: types.BoolValue(true)},
	})
	require.NoError(t, err)
	assert.Equal(t, 0, posts)
}

func TestRefreshNotificationChannelsReadsEnabled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "slack", r.URL.Query().Get("channel"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"channel":"slack","enabled":false}]}`))
	}))
	t.Cleanup(server.Close)

	resource := &IntegrationResource{portClient: &cli.PortClient{Client: resty.New().SetBaseURL(server.URL)}}
	channels := []NotificationChannelModel{
		{Channel: types.StringValue(consts.NotificationChannelSlack), Enabled: types.BoolValue(true)},
	}
	err := resource.refreshNotificationChannels(context.Background(), "github-prod", channels)
	require.NoError(t, err)
	assert.Equal(t, "slack", channels[0].Channel.ValueString())
	assert.False(t, channels[0].Enabled.ValueBool())
}

func TestValidateNotificationChannelsRejectsDuplicates(t *testing.T) {
	err := validateIntegrationConfig(&IntegrationModel{
		InstallationType: types.StringValue(consts.InstallationTypeSaas),
		NotificationChannels: []NotificationChannelModel{
			{Channel: types.StringValue(consts.NotificationChannelEmail), Enabled: types.BoolValue(true)},
			{Channel: types.StringValue(consts.NotificationChannelEmail), Enabled: types.BoolValue(false)},
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), `lists "email" more than once`)
}

func TestValidateNotificationChannelsRejectsMissingChannel(t *testing.T) {
	for name, channel := range map[string]types.String{
		"null":    types.StringNull(),
		"unknown": types.StringUnknown(),
	} {
		t.Run(name, func(t *testing.T) {
			err := validateIntegrationConfig(&IntegrationModel{
				InstallationType: types.StringValue(consts.InstallationTypeSaas),
				NotificationChannels: []NotificationChannelModel{
					{Channel: channel, Enabled: types.BoolValue(true)},
				},
			})
			assert.ErrorContains(t, err, "notification_channels entries must have a known channel")
		})
	}
}

func TestValidateNotificationChannelsRequiresHostedInstallation(t *testing.T) {
	channels := []NotificationChannelModel{
		{Channel: types.StringValue(consts.NotificationChannelEmail), Enabled: types.BoolValue(true)},
	}

	for _, installationType := range []types.String{types.StringNull(), types.StringValue(consts.InstallationTypeOnPrem)} {
		err := validateIntegrationConfig(&IntegrationModel{InstallationType: installationType, NotificationChannels: channels})
		assert.ErrorContains(t, err, "notification_channels is only supported when installation_type is")
	}
	for _, installationType := range []types.String{
		types.StringValue(consts.InstallationTypeSaas),
		types.StringValue(consts.InstallationTypeSaasOAuth2),
		types.StringUnknown(),
	} {
		err := validateIntegrationConfig(&IntegrationModel{InstallationType: installationType, NotificationChannels: channels})
		assert.NoError(t, err, installationType.String())
	}
}

func TestValidateNotificationChannelsAllowsOnPremWithoutChannels(t *testing.T) {
	err := validateIntegrationConfig(&IntegrationModel{
		InstallationType: types.StringValue(consts.InstallationTypeOnPrem),
	})
	assert.NoError(t, err)
}

func TestNotificationOptInKeepsSavedAudience(t *testing.T) {
	current := &cli.OptInSetting{
		Channel: consts.NotificationChannelSlack,
		Enabled: false,
		Emails:  []string{"admin@example.com"},
		Teams:   []string{"team-1"},
	}

	got := notificationOptIn(consts.NotificationChannelSlack, true, current)

	assert.Equal(t, cli.OptInSetting{
		Channel: consts.NotificationChannelSlack,
		Enabled: true,
		Emails:  []string{"admin@example.com"},
		Teams:   []string{"team-1"},
	}, got)
	assert.Equal(t, []string{"admin@example.com"}, current.Emails)
}

func TestNotificationOptInWithoutSavedAudienceNotifiesAdmins(t *testing.T) {
	got := notificationOptIn(consts.NotificationChannelEmail, true, nil)

	assert.Equal(t, cli.OptInSetting{
		Channel: consts.NotificationChannelEmail,
		Enabled: true,
		Emails:  []string{},
		Teams:   []string{},
	}, got)
}
