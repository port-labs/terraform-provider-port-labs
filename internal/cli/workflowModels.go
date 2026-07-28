package cli

// Workflow represents a Port workflow and its polymorphic node/connection graph.
type Workflow struct {
	ID          string  `json:"id,omitempty"`
	Identifier  string  `json:"identifier"`
	Title       *string `json:"title,omitempty"`
	Icon        *string `json:"icon,omitempty"`
	Description *string `json:"description,omitempty"`
	// Tags is a pointer to a slice so that an explicit empty list (`tags = []`)
	// clears the tags on update instead of being omitted by `omitempty`.
	Tags        *[]string            `json:"tags,omitempty"`
	Nodes       []WorkflowNode       `json:"nodes"`
	Connections []WorkflowConnection `json:"connections"`
}

// WorkflowNode is a discriminated union keyed by Type. Only the fields relevant
// to the node's Type are populated. Type holds either one of the entity event
// types (for event trigger nodes) or one of the workflow node type constants
// (CURSOR_AGENT, DELAY, INTEGRATION_ACTION).
type WorkflowNode struct {
	Identifier string  `json:"identifier"`
	Title      *string `json:"title,omitempty"`
	Type       string  `json:"type"`

	// Event trigger fields (Type is one of the entity event types).
	BlueprintIdentifier *string `json:"blueprintIdentifier,omitempty"`

	// CURSOR_AGENT fields.
	ApiKey *string            `json:"apiKey,omitempty"`
	Prompt *CursorAgentPrompt `json:"prompt,omitempty"`
	Source *CursorAgentSource `json:"source,omitempty"`

	// DELAY fields. Seconds is an int for a fixed delay or a string for a
	// dynamic Port expression ("{{ ... }}").
	Seconds any `json:"seconds,omitempty"`

	// INTEGRATION_ACTION fields.
	InstallationId                       *string                                       `json:"installationId,omitempty"`
	IntegrationActionType                *string                                       `json:"integrationActionType,omitempty"`
	IntegrationActionExecutionProperties *WorkflowIntegrationActionExecutionProperties `json:"integrationActionExecutionProperties,omitempty"`
}

type CursorAgentPrompt struct {
	Text *string `json:"text,omitempty"`
}

type CursorAgentSource struct {
	PrUrl *string `json:"prUrl,omitempty"`
}

type WorkflowIntegrationActionExecutionProperties struct {
	Org                  *string        `json:"org,omitempty"`
	Repo                 *string        `json:"repo,omitempty"`
	Workflow             *string        `json:"workflow,omitempty"`
	WorkflowInputs       map[string]any `json:"workflowInputs,omitempty"`
	ReportWorkflowStatus any            `json:"reportWorkflowStatus,omitempty"`
}

type WorkflowConnection struct {
	SourceIdentifier string `json:"sourceIdentifier"`
	TargetIdentifier string `json:"targetIdentifier"`
}
