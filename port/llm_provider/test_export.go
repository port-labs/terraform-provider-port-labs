package llm_provider

import (
	"context"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
)

// LlmProviderToPortBodyForTest exposes llmProviderToPortBody for unit tests.
func LlmProviderToPortBodyForTest(state *LLMProviderModel) (*cli.LLMProviderUpsert, error) {
	return llmProviderToPortBody(context.Background(), state)
}
