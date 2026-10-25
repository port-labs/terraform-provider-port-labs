package integration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
)

func validateNotificationChannels(m *IntegrationModel) error {
	channels := m.NotificationChannels
	if len(channels) == 0 {
		return nil
	}
	if !m.InstallationType.IsUnknown() && !m.isHosted() {
		return fmt.Errorf(
			"notification_channels is only supported when installation_type is %q or %q",
			consts.InstallationTypeSaas,
			consts.InstallationTypeSaasOAuth2,
		)
	}

	seen := make(map[string]struct{}, len(channels))
	for _, channel := range channels {
		if channel.Channel.IsNull() || channel.Channel.IsUnknown() {
			return fmt.Errorf("notification_channels entries must have a known channel")
		}
		name := channel.Channel.ValueString()
		if _, dup := seen[name]; dup {
			return fmt.Errorf("notification_channels lists %q more than once", name)
		}
		seen[name] = struct{}{}
	}
	return nil
}

// notificationOptIn builds the opt-in to save. Enabled comes from the
// resource. Emails and teams stay as Port last saved them, because the
// opt-in API replaces the whole audience on every write.
func notificationOptIn(channel string, enabled bool, current *cli.OptInSetting) cli.OptInSetting {
	setting := cli.OptInSetting{
		Channel: channel,
		Enabled: enabled,
		Emails:  []string{},
		Teams:   []string{},
	}
	if current == nil {
		return setting
	}
	if current.Emails != nil {
		setting.Emails = current.Emails
	}
	if current.Teams != nil {
		setting.Teams = current.Teams
	}
	return setting
}

func (r *IntegrationResource) applyNotificationChannels(ctx context.Context, installationID string, channels []NotificationChannelModel) error {
	for _, channel := range channels {
		if channel.Channel.IsNull() || channel.Channel.IsUnknown() || channel.Enabled.IsNull() || channel.Enabled.IsUnknown() {
			continue
		}
		name := channel.Channel.ValueString()
		enabled := channel.Enabled.ValueBool()

		current, err := r.portClient.GetChannelNotificationOptIn(ctx, installationID, name)
		if err != nil {
			return fmt.Errorf("read notification channel %q: %w", name, err)
		}
		if current.Channel == name && current.Enabled == enabled {
			continue
		}
		if err := r.portClient.SetChannelNotificationEnabled(ctx, installationID, notificationOptIn(name, enabled, current)); err != nil {
			return fmt.Errorf("set notification channel %q: %w", name, err)
		}
	}
	return nil
}

func (r *IntegrationResource) refreshNotificationChannels(ctx context.Context, installationID string, channels []NotificationChannelModel) error {
	for i, channel := range channels {
		if channel.Channel.IsNull() || channel.Channel.IsUnknown() {
			continue
		}
		name := channel.Channel.ValueString()
		current, err := r.portClient.GetChannelNotificationOptIn(ctx, installationID, name)
		if err != nil {
			return fmt.Errorf("read notification channel %q: %w", name, err)
		}
		channels[i].Channel = types.StringValue(current.Channel)
		channels[i].Enabled = types.BoolValue(current.Enabled)
	}
	return nil
}
