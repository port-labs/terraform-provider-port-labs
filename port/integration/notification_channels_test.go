package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func notificationChannelsValue(t *testing.T, enabledByChannel map[string]bool) types.Map {
	t.Helper()
	channels := make(map[string]NotificationChannelModel, len(enabledByChannel))
	for name, enabled := range enabledByChannel {
		channels[name] = NotificationChannelModel{Enabled: types.BoolValue(enabled)}
	}
	value, diags := types.MapValueFrom(context.Background(), notificationChannelElementType(), channels)
	require.False(t, diags.HasError(), diags)
	return value
}

func TestNotificationChannelsChanged(t *testing.T) {
	elementType := notificationChannelElementType()
	emailOn := notificationChannelsValue(t, map[string]bool{consts.NotificationChannelEmail: true})
	emailOff := notificationChannelsValue(t, map[string]bool{consts.NotificationChannelEmail: false})
	both := notificationChannelsValue(t, map[string]bool{
		consts.NotificationChannelEmail: true,
		consts.NotificationChannelSlack: false,
	})

	tests := []struct {
		name  string
		state types.Map
		plan  types.Map
		want  bool
	}{
		{name: "both null", state: types.MapNull(elementType), plan: types.MapNull(elementType), want: false},
		{name: "same channels", state: emailOn, plan: emailOn, want: false},
		{name: "enabled flipped", state: emailOn, plan: emailOff, want: true},
		{name: "channel added", state: emailOn, plan: both, want: true},
		{name: "channel removed", state: both, plan: emailOn, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, notificationChannelsChanged(tt.state, tt.plan))
		})
	}
}

func TestApplyNotificationChannelsPostsEnabledAndKeepsAudience(t *testing.T) {
	var posts int
	var posted cli.NotificationOptInSetting
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
	err := resource.applyNotificationChannels(context.Background(), "github-prod", notificationChannelsValue(t, map[string]bool{
		consts.NotificationChannelEmail: true,
	}))
	require.NoError(t, err)
	assert.Equal(t, 1, posts)
	assert.Equal(t, cli.NotificationOptInSetting{
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
	err := resource.applyNotificationChannels(context.Background(), "github-prod", notificationChannelsValue(t, map[string]bool{
		consts.NotificationChannelSlack: true,
	}))
	require.NoError(t, err)
	assert.Equal(t, 0, posts)
}

func TestApplyNotificationChannelsPostsOnlyChangedChannel(t *testing.T) {
	var posted []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			channel := r.URL.Query().Get("channel")
			enabled := channel == consts.NotificationChannelSlack
			_, _ = fmt.Fprintf(w, `{"ok":true,"data":[{"channel":%q,"enabled":%t,"emails":[],"teams":[]}]}`, channel, enabled)
		case http.MethodPost:
			var body cli.NotificationOptInSetting
			raw, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(raw, &body))
			posted = append(posted, body.Channel)
			_, _ = w.Write([]byte(`{"ok":true,"data":{"channel":"email","enabled":true}}`))
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	}))
	t.Cleanup(server.Close)

	resource := &IntegrationResource{portClient: &cli.PortClient{Client: resty.New().SetBaseURL(server.URL)}}
	err := resource.applyNotificationChannels(context.Background(), "github-prod", notificationChannelsValue(t, map[string]bool{
		consts.NotificationChannelEmail: true,
		consts.NotificationChannelSlack: true,
	}))
	require.NoError(t, err)
	assert.Equal(t, []string{consts.NotificationChannelEmail}, posted)
}

func TestRefreshNotificationChannelsReadsEnabled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "slack", r.URL.Query().Get("channel"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"channel":"Slack","enabled":false}]}`))
	}))
	t.Cleanup(server.Close)

	resource := &IntegrationResource{portClient: &cli.PortClient{Client: resty.New().SetBaseURL(server.URL)}}
	channels := notificationChannelsValue(t, map[string]bool{
		consts.NotificationChannelSlack: true,
	})
	updated, err := resource.refreshNotificationChannels(context.Background(), "github-prod", channels)
	require.NoError(t, err)
	var got map[string]NotificationChannelModel
	diags := updated.ElementsAs(context.Background(), &got, false)
	require.False(t, diags.HasError(), diags)
	require.Contains(t, got, consts.NotificationChannelSlack)
	assert.NotContains(t, got, "Slack")
	assert.False(t, got[consts.NotificationChannelSlack].Enabled.ValueBool())
}

func TestRefreshNotificationChannelsLeavesNullMap(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	resource := &IntegrationResource{portClient: &cli.PortClient{Client: resty.New().SetBaseURL(server.URL)}}
	channels := types.MapNull(notificationChannelElementType())
	updated, err := resource.refreshNotificationChannels(context.Background(), "github-prod", channels)
	require.NoError(t, err)
	assert.Equal(t, 0, calls)
	assert.True(t, updated.IsNull())
}

func TestUpdateWritesNotificationChannelsOnlyWhenChanged(t *testing.T) {
	elementType := notificationChannelElementType()
	enabled := notificationChannelsValue(t, map[string]bool{consts.NotificationChannelEmail: true})
	disabled := notificationChannelsValue(t, map[string]bool{consts.NotificationChannelEmail: false})

	tests := []struct {
		name     string
		state    types.Map
		plan     types.Map
		wantPost int
		wantCall int
	}{
		{name: "unchanged channels", state: enabled, plan: enabled, wantPost: 0, wantCall: 0},
		{name: "null channels", state: types.MapNull(elementType), plan: types.MapNull(elementType), wantPost: 0, wantCall: 0},
		{name: "enabled flipped", state: enabled, plan: disabled, wantPost: 1, wantCall: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			posts, calls, server := newRunningIntegrationServer(t)
			integration := &IntegrationResource{portClient: &cli.PortClient{Client: resty.New().SetBaseURL(server.URL)}}
			diags := updateIntegration(t, integration, "before", tt.state, "after", tt.plan)
			require.False(t, diags.HasError(), diags)
			assert.Equal(t, tt.wantPost, *posts, "opt-in posts")
			assert.Equal(t, tt.wantCall, *calls, "opt-in calls")
		})
	}
}

func TestValidateNotificationChannelsRequiresHostedInstallation(t *testing.T) {
	channels := notificationChannelsValue(t, map[string]bool{
		consts.NotificationChannelEmail: true,
	})

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
	current := &cli.NotificationOptInSetting{
		Channel: consts.NotificationChannelSlack,
		Enabled: false,
		Emails:  []string{"admin@example.com"},
		Teams:   []string{"team-1"},
	}

	got := notificationOptIn(consts.NotificationChannelSlack, true, current)

	assert.Equal(t, cli.NotificationOptInSetting{
		Channel: consts.NotificationChannelSlack,
		Enabled: true,
		Emails:  []string{"admin@example.com"},
		Teams:   []string{"team-1"},
	}, got)
	assert.Equal(t, []string{"admin@example.com"}, current.Emails)
}

func notificationChannelElementType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"enabled": types.BoolType,
	}}
}

const runningIntegrationJSON = `{"ok":true,"integration":{"installationId":"github-prod","installationType":"Saas","statusInfo":{"integrationStatus":{"status":"Running"}},"config":{"resources":[]}}}`

func newRunningIntegrationServer(t *testing.T) (posts, optInCalls *int, server *httptest.Server) {
	t.Helper()
	var postCount, callCount int
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(r.URL.Path, "/notifications/opt-ins") {
			callCount++
			if r.Method == http.MethodPost {
				postCount++
				_, _ = w.Write([]byte(`{"ok":true,"data":{"channel":"email","enabled":false}}`))
				return
			}
			_, _ = fmt.Fprintf(w, `{"ok":true,"data":[{"channel":%q,"enabled":true,"emails":[],"teams":[]}]}`, r.URL.Query().Get("channel"))
			return
		}
		_, _ = w.Write([]byte(runningIntegrationJSON))
	}))
	t.Cleanup(server.Close)
	return &postCount, &callCount, server
}

func integrationResourceSchema(t *testing.T) schema.Schema {
	t.Helper()
	var resp resource.SchemaResponse
	(&IntegrationResource{}).Schema(context.Background(), resource.SchemaRequest{}, &resp)
	require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
	return resp.Schema
}

func hostedIntegrationModel(title string, channels types.Map) *IntegrationModel {
	return &IntegrationModel{
		ID:                          types.StringValue("github-prod"),
		InstallationId:              types.StringValue("github-prod"),
		InstallationAppType:         types.StringValue("github-ocean"),
		InstallationType:            types.StringValue(consts.InstallationTypeSaas),
		Title:                       types.StringValue(title),
		Spec:                        types.StringValue(`{"integrationSpec":{"token":"secret"}}`),
		Config:                      types.StringNull(),
		Version:                     types.StringNull(),
		Status:                      types.StringNull(),
		CreatePortResourcesOrigin:   types.StringNull(),
		KafkaChangelogDestination:   types.ObjectNull(kafkaChangelogDestinationType),
		WebhookChangelogDestination: types.ObjectNull(webhookChangelogDestinationType),
		NotificationChannels:        channels,
	}
}

func updateIntegration(t *testing.T, integration *IntegrationResource, stateTitle string, stateChannels types.Map, planTitle string, planChannels types.Map) diag.Diagnostics {
	t.Helper()
	resourceSchema := integrationResourceSchema(t)
	req := resource.UpdateRequest{
		State: tfsdk.State{Schema: resourceSchema},
		Plan:  tfsdk.Plan{Schema: resourceSchema},
	}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: resourceSchema}}
	diags := req.State.Set(context.Background(), hostedIntegrationModel(stateTitle, stateChannels))
	diags.Append(req.Plan.Set(context.Background(), hostedIntegrationModel(planTitle, planChannels))...)
	require.False(t, diags.HasError(), diags)
	integration.Update(context.Background(), req, resp)
	return resp.Diagnostics
}

func TestNotificationOptInWithoutSavedAudienceNotifiesAdmins(t *testing.T) {
	got := notificationOptIn(consts.NotificationChannelEmail, true, nil)

	assert.Equal(t, cli.NotificationOptInSetting{
		Channel: consts.NotificationChannelEmail,
		Enabled: true,
		Emails:  []string{},
		Teams:   []string{},
	}, got)
}
