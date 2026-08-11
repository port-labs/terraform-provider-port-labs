package cli

type Workflow struct {
	ID                    string               `json:"id,omitempty"`
	Identifier            string               `json:"identifier"`
	Title                 *string              `json:"title,omitempty"`
	Icon                  *string              `json:"icon,omitempty"`
	Description           *string              `json:"description,omitempty"`
	Category              *string              `json:"category,omitempty"`
	AllowAnyoneToViewRuns *bool                `json:"allowAnyoneToViewRuns,omitempty"`
	Nodes                 []WorkflowNode       `json:"nodes"`
	Connections           []WorkflowConnection `json:"connections"`
}

type WorkflowNode struct {
	Identifier  string             `json:"identifier"`
	Title       *string            `json:"title,omitempty"`
	Icon        *string            `json:"icon,omitempty"`
	Description *string            `json:"description,omitempty"`
	Variables   map[string]string  `json:"variables,omitempty"`
	Links       []string           `json:"links,omitempty"`
	Verbose     *bool              `json:"verbose,omitempty"`
	Config      WorkflowNodeConfig `json:"config"`
}

// WorkflowNodeConfig is the flattened union of every node config. Type selects
// the active variant; every other field is omitted from the payload.
type WorkflowNodeConfig struct {
	Type string `json:"type"`

	// Shared by every trigger node type.
	Published *bool `json:"published,omitempty"`

	// EVENT_TRIGGER
	Event     *WorkflowTriggerEvent  `json:"event,omitempty"`
	Condition *WorkflowNodeCondition `json:"condition,omitempty"`

	// SCHEDULE_TRIGGER
	Cron *string `json:"cron,omitempty"`

	// SELF_SERVE_TRIGGER
	ActionCardButtonText    *string                  `json:"actionCardButtonText,omitempty"`
	ExecuteActionButtonText *string                  `json:"executeActionButtonText,omitempty"`
	Variant                 *string                  `json:"variant,omitempty"`
	Contexts                []WorkflowTriggerContext `json:"contexts,omitempty"`
	Permissions             *WorkflowNodePermissions `json:"permissions,omitempty"`

	// SELF_SERVE_TRIGGER and INPUT
	UserInputs *WorkflowUserInputs `json:"userInputs,omitempty"`

	// KAFKA
	Payload any `json:"payload,omitempty"`

	// WEBHOOK
	Url          *string           `json:"url,omitempty"`
	Agent        *bool             `json:"agent,omitempty"`
	Synchronized *bool             `json:"synchronized,omitempty"`
	Method       *string           `json:"method,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	Body         any               `json:"body,omitempty"`
	OnTimeout    *string           `json:"onTimeout,omitempty"`

	// INTEGRATION_ACTION
	InstallationId                       *string `json:"installationId,omitempty"`
	IntegrationProvider                  *string `json:"integrationProvider,omitempty"`
	IntegrationInvocationType            *string `json:"integrationInvocationType,omitempty"`
	IntegrationActionExecutionProperties any     `json:"integrationActionExecutionProperties,omitempty"`

	// UPSERT_ENTITY
	BlueprintIdentifier *string                `json:"blueprintIdentifier,omitempty"`
	Mapping             *WorkflowUpsertMapping `json:"mapping,omitempty"`

	// AI and AI_AGENT
	UserPrompt      *string             `json:"userPrompt,omitempty"`
	AgentIdentifier *string             `json:"agentIdentifier,omitempty"`
	Provider        *string             `json:"provider,omitempty"`
	Model           *string             `json:"model,omitempty"`
	SystemPrompt    *string             `json:"systemPrompt,omitempty"`
	Tools           []string            `json:"tools,omitempty"`
	McpServers      []WorkflowMcpServer `json:"mcpServers,omitempty"`
	OutputSchema    any                 `json:"outputSchema,omitempty"`

	// CONDITION and INPUT
	Outlets []WorkflowOutlet `json:"outlets,omitempty"`

	// INPUT
	Description   *string                     `json:"description,omitempty"`
	Notifications []WorkflowInputNotification `json:"notifications,omitempty"`
	Responders    *WorkflowNodePermissions    `json:"responders,omitempty"`

	// Shared by most invocation node types.
	OnFailure *string `json:"onFailure,omitempty"`
}

type WorkflowTriggerEvent struct {
	Type                string  `json:"type"`
	BlueprintIdentifier string  `json:"blueprintIdentifier"`
	PropertyIdentifier  *string `json:"propertyIdentifier,omitempty"`
}

type WorkflowNodeCondition struct {
	Type        string   `json:"type"`
	Expressions []string `json:"expressions"`
	Combinator  *string  `json:"combinator,omitempty"`
}

type WorkflowTriggerContext struct {
	On                  string  `json:"on"`
	BlueprintIdentifier *string `json:"blueprintIdentifier,omitempty"`
	UserInput           *string `json:"userInput,omitempty"`
}

// Self serve trigger permissions and input node responders share users/roles/teams
// but differ in their dynamic query: permissions take an RBAC `policy`, responders
// take an entity search `usersQuery` over _user.
type WorkflowNodePermissions struct {
	Users      []string `json:"users,omitempty"`
	Roles      []string `json:"roles,omitempty"`
	Teams      []string `json:"teams,omitempty"`
	Policy     any      `json:"policy,omitempty"`
	UsersQuery any      `json:"usersQuery,omitempty"`
}

type WorkflowUserInputs struct {
	Properties map[string]ActionProperty `json:"properties"`
	Required   any                       `json:"required,omitempty"`
	Order      []string                  `json:"order,omitempty"`
	Steps      []Step                    `json:"steps,omitempty"`
	Titles     map[string]ActionTitle    `json:"titles,omitempty"`
	Buttons    []WorkflowInputButton     `json:"buttons,omitempty"`
}

type WorkflowInputButton struct {
	Identifier string  `json:"identifier"`
	Label      string  `json:"label"`
	Variant    string  `json:"variant"`
	Icon       *string `json:"icon,omitempty"`
}

type WorkflowUpsertMapping struct {
	Identifier *string        `json:"identifier,omitempty"`
	Title      *string        `json:"title,omitempty"`
	Team       any            `json:"team,omitempty"`
	Icon       *string        `json:"icon,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
	Relations  map[string]any `json:"relations,omitempty"`
}

type WorkflowMcpServer struct {
	Identifier string `json:"identifier"`
}

type WorkflowStatusLabel struct {
	Text    string  `json:"text"`
	Variant *string `json:"variant,omitempty"`
}

type WorkflowOutlet struct {
	Identifier          string               `json:"identifier"`
	Title               *string              `json:"title,omitempty"`
	Expression          *string              `json:"expression,omitempty"`
	EvaluationMethod    *string              `json:"evaluationMethod,omitempty"`
	NumOfResponders     *int64               `json:"numOfResponders,omitempty"`
	StatusLabel         *WorkflowStatusLabel `json:"statusLabel,omitempty"`
	WorkflowStatusLabel *WorkflowStatusLabel `json:"workflowStatusLabel,omitempty"`
}

type WorkflowInputNotification struct {
	Target  string                           `json:"target"`
	Fields  []WorkflowInputNotificationField `json:"fields,omitempty"`
	Url     *string                          `json:"url,omitempty"`
	Method  *string                          `json:"method,omitempty"`
	Headers map[string]string                `json:"headers,omitempty"`
	Body    map[string]any                   `json:"body,omitempty"`
	Agent   *bool                            `json:"agent,omitempty"`
}

type WorkflowInputNotificationField struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type WorkflowConnection struct {
	SourceIdentifier       string  `json:"sourceIdentifier"`
	TargetIdentifier       string  `json:"targetIdentifier"`
	Description            *string `json:"description,omitempty"`
	SourceOutletIdentifier *string `json:"sourceOutletIdentifier,omitempty"`
	Fallback               *bool   `json:"fallback,omitempty"`
}
