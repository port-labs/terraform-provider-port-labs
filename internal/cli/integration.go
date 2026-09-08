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
	provisioningPollInterval = 5 * time.Second
	provisioningPollJitter   = 2 * time.Second
	provisioningMaxAttempts  = 30
	provisioningMaxDelay     = 15 * time.Second
)

var (
	errIntegrationNotReady   = utils.StringErr("integration is not ready yet")
	errIntegrationNotDeleted = utils.StringErr("integration is not deleted yet")
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

func provisioningPollOptions(ctx context.Context, retryOn error) []retry.Option {
	return []retry.Option{
		retry.Context(ctx),
		retry.LastErrorOnly(true),
		retry.Attempts(provisioningMaxAttempts),
		retry.Delay(provisioningPollInterval),
		retry.DelayType(retry.BackOffDelay),
		retry.MaxDelay(provisioningMaxDelay),
		retry.MaxJitter(provisioningPollJitter),
		retry.RetryIf(func(err error) bool {
			return errors.Is(err, retryOn)
		}),
	}
}

// WaitForIntegrationReady polls until the integration reaches a terminal state.
// Returns (nil, nil) on timeout — the caller decides whether that's a warning or error.
func (c *PortClient) WaitForIntegrationReady(ctx context.Context, installationId string) (*Integration, error) {
	integration, err := retry.DoWithData(
		func() (*Integration, error) {
			integration, statusCode, err := c.GetIntegration(ctx, installationId)
			if err != nil {
				if statusCode == 404 {
					return nil, fmt.Errorf("integration %q disappeared during provisioning", installationId)
				}
				return nil, err
			}

			status := integrationStatus(integration)
			switch status {
			case "", consts.IntegrationStatusRunning:
				return integration, nil
			case consts.IntegrationStatusCreating, consts.IntegrationStatusUpdating:
				return nil, errIntegrationNotReady
			case consts.IntegrationStatusDeleting:
				return nil, fmt.Errorf("integration %q is being deleted", installationId)
			case consts.IntegrationStatusError, consts.IntegrationStatusUnHealthy:
				return nil, fmt.Errorf("integration provisioning failed (status: %s%s)", status, statusMessage(integration))
			default:
				return nil, fmt.Errorf("integration %q unexpected status: %s%s", installationId, status, statusMessage(integration))
			}
		},
		provisioningPollOptions(ctx, errIntegrationNotReady)...,
	)
	if errors.Is(err, errIntegrationNotReady) {
		return nil, nil
	}
	return integration, err
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
		provisioningPollOptions(ctx, errIntegrationNotDeleted)...,
	)
	if errors.Is(err, errIntegrationNotDeleted) {
		return fmt.Errorf("timed out waiting for integration %q to be deleted (still exists after %d polling attempts)", installationId, provisioningMaxAttempts)
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
