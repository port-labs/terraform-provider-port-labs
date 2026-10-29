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
	if channels.IsNull() || channels.IsUnknown() || len(channels.Elements()) == 0 {
		return nil
	}
	if !m.InstallationType.IsUnknown() && !m.isHosted() {
		return fmt.Errorf(
			"notification_channels is only supported when installation_type is %q or %q",
			consts.InstallationTypeSaas,
			consts.InstallationTypeSaasOAuth2,
		)
	}
	return nil
}

func notificationChannelsByName(ctx context.Context, channels types.Map) (map[string]NotificationChannelModel, error) {
	if channels.IsNull() || channels.IsUnknown() {
		return nil, nil
	}
	configured := make(map[string]NotificationChannelModel, len(channels.Elements()))
	diags := channels.ElementsAs(ctx, &configured, false)
	if diags.HasError() {
		return nil, fmt.Errorf("read notification_channels: %s", diags)
	}
	return configured, nil
}

// notificationOptIn builds the opt-in to save. Enabled comes from the
// resource. Emails and teams stay as Port last saved them, because the
// opt-in API replaces the whole audience on every write.
func notificationOptIn(channel string, enabled bool, current *cli.NotificationOptInSetting) cli.NotificationOptInSetting {
	setting := cli.NotificationOptInSetting{
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

// notificationChannelsChanged reports whether this apply is changing the
// channel map. Opt-in writes are a separate Port call, and repeating one
// after an integration update that left the channels alone makes Port
// schedule another integration update.
func notificationChannelsChanged(state, plan types.Map) bool {
	return !state.Equal(plan)
}

func (r *IntegrationResource) applyNotificationChannels(ctx context.Context, installationID string, channels types.Map) error {
	configured, err := notificationChannelsByName(ctx, channels)
	if err != nil {
		return err
	}
	for name, channel := range configured {
		if channel.Enabled.IsNull() || channel.Enabled.IsUnknown() {
			continue
		}
		enabled := channel.Enabled.ValueBool()

		current, err := r.portClient.GetChannelNotificationOptIn(ctx, installationID, name)
		if err != nil {
			return fmt.Errorf("read notification channel %q: %w", name, err)
		}
		if current.Enabled == enabled {
			continue
		}
		if err := r.portClient.SetChannelNotificationEnabled(ctx, installationID, notificationOptIn(name, enabled, current)); err != nil {
			return fmt.Errorf("set notification channel %q: %w", name, err)
		}
	}
	return nil
}

func (r *IntegrationResource) refreshNotificationChannels(ctx context.Context, installationID string, channels types.Map) (types.Map, error) {
	configured, err := notificationChannelsByName(ctx, channels)
	if err != nil || configured == nil {
		return channels, err
	}
	refreshed := make(map[string]NotificationChannelModel, len(configured))
	for name, channel := range configured {
		if channel.Enabled.IsNull() || channel.Enabled.IsUnknown() {
			refreshed[name] = channel
			continue
		}
		current, err := r.portClient.GetChannelNotificationOptIn(ctx, installationID, name)
		if err != nil {
			return channels, fmt.Errorf("read notification channel %q: %w", name, err)
		}
		// Port may return a different channel casing. The config key is the name Terraform stores.
		channel.Enabled = types.BoolValue(current.Enabled)
		refreshed[name] = channel
	}
	updated, diags := types.MapValueFrom(ctx, channels.ElementType(ctx), refreshed)
	if diags.HasError() {
		return channels, fmt.Errorf("read notification_channels: %s", diags)
	}
	return updated, nil
}
