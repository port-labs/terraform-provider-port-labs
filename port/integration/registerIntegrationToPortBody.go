package integration

import (
	"fmt"
	"strings"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

func integrationToRegisterRequest(state *IntegrationModel, shouldUpdate bool) (*cli.RegisterIntegrationRequest, error) {
	integration, err := integrationToPortBody(state)
	if err != nil {
		return nil, err
	}

	if integration.InstallationAppType == nil || *integration.InstallationAppType == "" {
		return nil, fmt.Errorf("installation_app_type is required when registering an integration")
	}
	if integration.Version == nil || *integration.Version == "" {
		return nil, fmt.Errorf("version is required when registering an integration")
	}
	if integration.Config == nil {
		return nil, fmt.Errorf("config is required when registering an integration")
	}
	if err := validateIntegrationConfigResources(*integration.Config); err != nil {
		return nil, err
	}

	return &cli.RegisterIntegrationRequest{
		InstallationId:      integration.InstallationId,
		InstallationAppType: *integration.InstallationAppType,
		Version:             *integration.Version,
		Config:              *integration.Config,
		Title:               integration.Title,
		ShouldUpdate:        shouldUpdate,
	}, nil
}

func validateIntegrationConfigResources(config map[string]any) error {
	resources, ok := config["resources"]
	if !ok || resources == nil {
		return fmt.Errorf("config.resources must contain at least one mapping")
	}

	switch typedResources := resources.(type) {
	case []any:
		if len(typedResources) == 0 {
			return fmt.Errorf("config.resources must contain at least one mapping")
		}
	case []map[string]any:
		if len(typedResources) == 0 {
			return fmt.Errorf("config.resources must contain at least one mapping")
		}
	default:
		return fmt.Errorf("config.resources must contain at least one mapping")
	}

	return nil
}

func integrationHasChangelogDestination(state *IntegrationModel) bool {
	return state.WebhookChangelogDestination != nil || !state.KafkaChangelogDestination.IsNull()
}

func isRegisterUpdateFallbackError(err error) bool {
	if err == nil {
		return false
	}

	message := err.Error()
	return strings.Contains(message, "invalid_installation_type") ||
		strings.Contains(message, "Register can only update OnPrem integrations")
}
