package consts

const (
	// SystemBlueprintAIAgentIdentifier is the Port system blueprint for AI agents.
	SystemBlueprintAIAgentIdentifier = "_ai_agent"

	// LegacyAIAgentPromptMaxLength was the previous maxLength shipped for long prompt
	// fields on the AI agent blueprint before the API raised the limit.
	LegacyAIAgentPromptMaxLength = 2500

	// AIAgentPromptMaxLength matches the current Port API limit for AI agent prompt text.
	AIAgentPromptMaxLength = 5000
)
