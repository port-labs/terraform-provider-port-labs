package cli

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetChannelNotificationEnabledSendsAudience(t *testing.T) {
	var got NotificationOptInSetting
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/integration/github-prod/notifications/opt-ins", r.URL.Path)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(body, &got))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"data":{"channel":"slack","enabled":true,"emails":["admin@example.com"],"teams":["team-1"]}}`))
	}))
	t.Cleanup(server.Close)

	client := &PortClient{Client: resty.New().SetBaseURL(server.URL)}
	err := client.SetChannelNotificationEnabled(context.Background(), "github-prod", NotificationOptInSetting{
		Channel: "slack",
		Enabled: true,
		Emails:  []string{"admin@example.com"},
		Teams:   []string{"team-1"},
	})
	require.NoError(t, err)
	assert.Equal(t, NotificationOptInSetting{
		Channel: "slack",
		Enabled: true,
		Emails:  []string{"admin@example.com"},
		Teams:   []string{"team-1"},
	}, got)
}

func TestGetChannelNotificationOptIn(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/integration/github-prod/notifications/opt-ins", r.URL.Path)
		assert.Equal(t, "email", r.URL.Query().Get("channel"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"channel":"email","enabled":false,"emails":["admin@example.com"],"teams":[]}]}`))
	}))
	t.Cleanup(server.Close)

	client := &PortClient{Client: resty.New().SetBaseURL(server.URL)}
	got, err := client.GetChannelNotificationOptIn(context.Background(), "github-prod", "email")
	require.NoError(t, err)
	assert.Equal(t, &NotificationOptInSetting{
		Channel: "email",
		Enabled: false,
		Emails:  []string{"admin@example.com"},
		Teams:   []string{},
	}, got)
}

func TestSetChannelNotificationEnabledSendsEmptyAudienceWhenUnset(t *testing.T) {
	var raw string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		raw = string(body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"data":{"channel":"email","enabled":true}}`))
	}))
	t.Cleanup(server.Close)

	client := &PortClient{Client: resty.New().SetBaseURL(server.URL)}
	err := client.SetChannelNotificationEnabled(context.Background(), "github-prod", NotificationOptInSetting{
		Channel: "email",
		Enabled: true,
	})
	require.NoError(t, err)
	assert.JSONEq(t, `{"channel":"email","enabled":true,"emails":[],"teams":[]}`, raw)
}

func TestSetChannelNotificationEnabledReportsAPIMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"ok":false,"error":"bad_request","message":"invalid channel \"sms\""}`))
	}))
	t.Cleanup(server.Close)

	client := &PortClient{Client: resty.New().SetBaseURL(server.URL)}
	err := client.SetChannelNotificationEnabled(context.Background(), "github-prod", NotificationOptInSetting{
		Channel: "sms",
		Enabled: true,
	})
	require.EqualError(t, err, `invalid channel "sms"`)
}
