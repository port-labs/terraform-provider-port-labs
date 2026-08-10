package cli

type Workflow struct {
	ID          string               `json:"id,omitempty"`
	Identifier  string               `json:"identifier"`
	Title       *string              `json:"title,omitempty"`
	Icon        *string              `json:"icon,omitempty"`
	Description *string              `json:"description,omitempty"`
	Category    *string              `json:"category,omitempty"`
	Nodes       []WorkflowNode       `json:"nodes"`
	Connections []WorkflowConnection `json:"connections"`
}

type WorkflowNode struct {
	Identifier string             `json:"identifier"`
	Title      *string            `json:"title,omitempty"`
	Config     WorkflowNodeConfig `json:"config"`
}

type WorkflowNodeConfig struct {
	Type string `json:"type"`

	Event *WorkflowTriggerEvent `json:"event,omitempty"`

	ApiKey *string            `json:"apiKey,omitempty"`
	Prompt *CursorAgentPrompt `json:"prompt,omitempty"`
	Source *CursorAgentSource `json:"source,omitempty"`

	InstallationId                       *string `json:"installationId,omitempty"`
	IntegrationProvider                  *string `json:"integrationProvider,omitempty"`
	IntegrationInvocationType            *string `json:"integrationInvocationType,omitempty"`
	IntegrationActionExecutionProperties any     `json:"integrationActionExecutionProperties,omitempty"`
	DeferIntegrationInstallation         *bool   `json:"deferIntegrationInstallation,omitempty"`

	OnFailure *string `json:"onFailure,omitempty"`
}

type WorkflowTriggerEvent struct {
	Type                string  `json:"type"`
	BlueprintIdentifier string  `json:"blueprintIdentifier"`
	PropertyIdentifier  *string `json:"propertyIdentifier,omitempty"`
}

type CursorAgentPrompt struct {
	Text *string `json:"text,omitempty"`
}

type CursorAgentSource struct {
	Repository *string `json:"repository,omitempty"`
	Ref        *string `json:"ref,omitempty"`
	PrUrl      *string `json:"prUrl,omitempty"`
}

type WorkflowConnection struct {
	SourceIdentifier string `json:"sourceIdentifier"`
	TargetIdentifier string `json:"targetIdentifier"`
}
