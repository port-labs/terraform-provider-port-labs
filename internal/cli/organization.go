package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

const orgUrl = "/v1/organization"
const orgSecretsUrl = "/v1/organization/secrets"

func (c *PortClient) ReadOrganization(ctx context.Context) (*Organization, int, error) {
	pb := &PortBody{}
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetResult(pb).
		Get(orgUrl)
	if err != nil {
		return nil, 0, err
	} else if resp.IsError() {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read organization, got: %s", resp.Body())
	}
	if pb.Organization == nil {
		return nil, 0, fmt.Errorf("port-api returned an invalid response: Organization is nil")
	}
	return pb.Organization, resp.StatusCode(), nil
}

func (c *PortClient) CreateOrganizationSecret(ctx context.Context, secret *OrganizationSecret) (*OrganizationSecret, error) {
	resp, err := c.Client.R().
		SetBody(secret).
		SetContext(ctx).
		Post(orgSecretsUrl)

	if err != nil {
		return nil, err
	}
	var pb OrganizationSecretBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to create organization secret, got: %s", resp.Body())
	}

	return &pb.Secret, nil
}

func (c *PortClient) ReadOrganizationSecret(ctx context.Context, secretName string) (*OrganizationSecret, int, error) {
	var pb OrganizationSecretBody
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("secretName", secretName).
		SetResult(&pb).
		Get(orgSecretsUrl + "/{secretName}")
	if err != nil {
		return nil, 0, err
	} else if resp.IsError() {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read organization secret, got: %s", resp.Body())
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read organization secret, got: %s", resp.Body())
	}

	return &pb.Secret, resp.StatusCode(), nil
}

func (c *PortClient) UpdateOrganizationSecret(ctx context.Context, secretName string, secret *OrganizationSecret) (*OrganizationSecret, error) {
	resp, err := c.Client.R().
		SetBody(secret).
		SetContext(ctx).
		SetPathParam("secretName", secretName).
		Patch(orgSecretsUrl + "/{secretName}")

	if err != nil {
		return nil, err
	}
	var pb OrganizationSecretBody
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to update organization secret, got: %s", resp.Body())
	}

	return &pb.Secret, nil
}

func (c *PortClient) DeleteOrganizationSecret(ctx context.Context, secretName string) error {
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("secretName", secretName).
		Delete(orgSecretsUrl + "/{secretName}")

	if err != nil {
		return err
	}
	var pb PortBodyDelete
	err = json.Unmarshal(resp.Body(), &pb)
	if err != nil {
		return err
	}

	if !pb.Ok {
		return fmt.Errorf("failed to delete organization secret, got: %s", string(resp.Body()))
	}
	return nil
}
