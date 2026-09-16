package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

const (
	integrationPollInterval = 5 * time.Second
	integrationPollJitter   = 2 * time.Second
	integrationMaxAttempts  = 30
	integrationMaxDelay     = 15 * time.Second
)

var (
	errIntegrationOperationPending = utils.StringErr("integration operation is not finished yet")
	errIntegrationNotProvisioned   = utils.StringErr("integration resources are not provisioned yet")
	errIntegrationNotDeleted       = utils.StringErr("integration is not deleted yet")
)

type PortBodyForIntegration struct {
	OK          bool        `json:"ok"`
	Integration Integration `json:"integration"`
}

func (c *PortClient) GetIntegration(ctx context.Context, id string) (*Integration, int, error) {
	pb := &PortBodyForIntegration{}
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		SetPathParam("identifier", id).
		SetQueryParam("byField", "installationId").
		Get("v1/integration/{identifier}")
	if err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read integration, got: %s", resp.Body())
	}
	return &pb.Integration, resp.StatusCode(), nil
}

func (c *PortClient) CreateIntegration(ctx context.Context, integration *Integration) (*Integration, error) {
	resp, err := c.Client.R().
		SetBody(integration).
		SetContext(ctx).
		Post("v1/integration")
	if err != nil {
		return nil, err
	}

	var pb PortBodyForIntegration
	if err := json.Unmarshal(resp.Body(), &pb); err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to create integration, got: %s", resp.Body())
	}
	return &pb.Integration, nil
}

func (c *PortClient) UpdateIntegration(ctx context.Context, id string, integration *Integration) (*Integration, error) {
	resp, err := c.Client.R().
		SetBody(integration).
		SetContext(ctx).
		SetPathParam("identifier", id).
		Patch("v1/integration/{identifier}")
	if err != nil {
		return nil, err
	}

	var pb PortBodyForIntegration
	if err := json.Unmarshal(resp.Body(), &pb); err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to update integration, got: %s", resp.Body())
	}
	return &pb.Integration, nil
}

func (c *PortClient) ValidateIntegrationSpec(ctx context.Context, integrationType string, body ValidateIntegrationSpecBody) error {
	resp, err := c.Client.R().
		SetBody(body).
		SetContext(ctx).
		SetPathParam("integration_type", integrationType).
		Post("v1/integration/{integration_type}/spec/validate")
	if err != nil {
		return err
	}

	var result struct {
		OK      bool   `json:"ok"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("failed to validate integration spec, got: %s", resp.Body())
	}
	if resp.IsError() {
		if result.Message != "" {
			return fmt.Errorf("%s", result.Message)
		}
		return fmt.Errorf("failed to validate integration spec (HTTP %d): %s", resp.StatusCode(), resp.Body())
	}
	if !result.OK {
		if result.Message != "" {
			return fmt.Errorf("%s", result.Message)
		}
		return fmt.Errorf("failed to validate integration spec, got: %s", resp.Body())
	}
	return nil
}

func (c *PortClient) DeleteIntegration(ctx context.Context, id string) (int, error) {
	resp, err := c.Client.R().
		SetContext(ctx).
		SetPathParam("identifier", id).
		Delete("v1/integration/{identifier}")
	if err != nil {
		return resp.StatusCode(), err
	}

	var pb PortBodyForIntegration
	if err := json.Unmarshal(resp.Body(), &pb); err != nil {
		return resp.StatusCode(), err
	}
	if !pb.OK {
		return resp.StatusCode(), fmt.Errorf("failed to delete integration, got: %s", resp.Body())
	}
	return resp.StatusCode(), nil
}

func integrationPollOptions(ctx context.Context, retryOn error) []retry.Option {
	return []retry.Option{
		retry.Context(ctx),
		retry.LastErrorOnly(true),
		retry.Attempts(integrationMaxAttempts),
		retry.Delay(integrationPollInterval),
		retry.DelayType(retry.BackOffDelay),
		retry.MaxDelay(integrationMaxDelay),
		retry.MaxJitter(integrationPollJitter),
		retry.RetryIf(func(err error) bool {
			return errors.Is(err, retryOn)
		}),
	}
}

// WaitForIntegrationReady polls until the integration operation finishes (e.g. Ocean
// deployment reaches Running). Returns (nil, nil) on timeout — the caller decides
// whether that's a warning or error.
func (c *PortClient) WaitForIntegrationReady(ctx context.Context, installationId string) (*Integration, error) {
	integration, err := retry.DoWithData(
		func() (*Integration, error) {
			integration, statusCode, err := c.GetIntegration(ctx, installationId)
			if err != nil {
				if statusCode == 404 {
					return nil, fmt.Errorf("integration %q disappeared while waiting for operation", installationId)
				}
				return nil, err
			}

			status := integrationStatus(integration)
			switch status {
			case "", consts.IntegrationStatusRunning:
				return integration, nil
			case consts.IntegrationStatusCreating, consts.IntegrationStatusUpdating:
				return nil, errIntegrationOperationPending
			case consts.IntegrationStatusDeleting:
				return nil, fmt.Errorf("integration %q is being deleted", installationId)
			case consts.IntegrationStatusError, consts.IntegrationStatusUnHealthy:
				return nil, fmt.Errorf("integration operation failed (status: %s%s)", status, statusMessage(integration))
			default:
				return nil, fmt.Errorf("integration %q unexpected status: %s%s", installationId, status, statusMessage(integration))
			}
		},
		integrationPollOptions(ctx, errIntegrationOperationPending)...,
	)
	if errors.Is(err, errIntegrationOperationPending) {
		return nil, nil
	}
	return integration, err
}

// WaitForIntegrationProvisioned polls until default blueprints and mappings are
// provisioned (config is no longer {}). Returns (nil, nil) on timeout.
func (c *PortClient) WaitForIntegrationProvisioned(ctx context.Context, installationId string) (*Integration, error) {
	integration, err := retry.DoWithData(
		func() (*Integration, error) {
			integration, statusCode, err := c.GetIntegration(ctx, installationId)
			if err != nil {
				if statusCode == 404 {
					return nil, fmt.Errorf("integration %q disappeared while waiting for provisioning", installationId)
				}
				return nil, err
			}
			if IsIntegrationConfigProvisioned(integration.Config) {
				return integration, nil
			}
			return nil, errIntegrationNotProvisioned
		},
		integrationPollOptions(ctx, errIntegrationNotProvisioned)...,
	)
	if errors.Is(err, errIntegrationNotProvisioned) {
		return nil, nil
	}
	return integration, err
}

func IsIntegrationConfigProvisioned(config *map[string]any) bool {
	return config != nil && len(*config) > 0
}

func (c *PortClient) WaitForIntegrationDeleted(ctx context.Context, installationId string) error {
	err := retry.Do(
		func() error {
			_, statusCode, err := c.GetIntegration(ctx, installationId)
			if statusCode == 404 {
				return nil
			}
			if err != nil {
				return err
			}
			return errIntegrationNotDeleted
		},
		integrationPollOptions(ctx, errIntegrationNotDeleted)...,
	)
	if errors.Is(err, errIntegrationNotDeleted) {
		return fmt.Errorf("timed out waiting for integration %q to be deleted (still exists after %d polling attempts)", installationId, integrationMaxAttempts)
	}
	return err
}

func integrationStatus(i *Integration) string {
	if i == nil || i.StatusInfo == nil {
		return ""
	}
	return i.StatusInfo.IntegrationStatus.Status
}

func statusMessage(i *Integration) string {
	if i == nil || i.StatusInfo == nil || i.StatusInfo.IntegrationStatus.Message == nil {
		return ""
	}
	return ": " + *i.StatusInfo.IntegrationStatus.Message
}
