package cli

import (
	"context"
	"encoding/json"
	"fmt"
)

const llmProvidersURL = "/v1/llm-providers"

func (c *PortClient) ListLLMProviders(ctx context.Context) ([]LLMProvider, int, error) {
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		Get(llmProvidersURL)
	if err != nil {
		return nil, resp.StatusCode(), err
	}

	var pb LLMProvidersListBody
	if err := json.Unmarshal(resp.Body(), &pb); err != nil {
		return nil, resp.StatusCode(), err
	}
	if !pb.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to list llm providers, got: %s", resp.Body())
	}
	return pb.Providers, resp.StatusCode(), nil
}

func (c *PortClient) ReadLLMProvider(ctx context.Context, providerName string) (*LLMProvider, int, error) {
	providers, statusCode, err := c.ListLLMProviders(ctx)
	if err != nil {
		return nil, statusCode, err
	}

	for i := range providers {
		if providers[i].Provider == providerName {
			return &providers[i], statusCode, nil
		}
	}

	return nil, 404, fmt.Errorf("llm provider %q not found", providerName)
}

func (c *PortClient) UpsertLLMProvider(ctx context.Context, provider *LLMProviderUpsert, validateConnection bool) (*LLMProvider, error) {
	req := c.Client.R().
		SetBody(provider).
		SetContext(ctx)

	if validateConnection {
		req.SetQueryParam("validate_connection", "true")
	}

	resp, err := req.Post(llmProvidersURL)
	if err != nil {
		return nil, err
	}

	var pb LLMProviderBody
	if err := json.Unmarshal(resp.Body(), &pb); err != nil {
		return nil, err
	}
	if !pb.OK {
		return nil, fmt.Errorf("failed to upsert llm provider, got: %s", resp.Body())
	}

	if pb.Provider.Provider != "" {
		return &pb.Provider, nil
	}

	readProvider, _, err := c.ReadLLMProvider(ctx, provider.Provider)
	if err != nil {
		return nil, err
	}
	return readProvider, nil
}

func (c *PortClient) DeleteLLMProvider(ctx context.Context, providerName string) error {
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("provider", providerName).
		Delete(llmProvidersURL + "/{provider}")
	if err != nil {
		return err
	}

	if resp.StatusCode() == 204 {
		return nil
	}

	var pb PortBodyDelete
	if err := json.Unmarshal(resp.Body(), &pb); err != nil {
		return err
	}
	if !pb.Ok {
		return fmt.Errorf("failed to delete llm provider, got: %s", string(resp.Body()))
	}
	return nil
}
