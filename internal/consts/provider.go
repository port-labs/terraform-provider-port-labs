package consts

const (
	ProviderName         = "port"
	DefaultBaseUrl       = "https://api.getport.io"
	Kafka                = "KAFKA"
	Webhook              = "WEBHOOK"
	Github               = "GITHUB"
	Gitlab               = "GITLAB"
	AzureDevops          = "AZURE_DEVOPS"
	UpsertEntity         = "UPSERT_ENTITY"
	IntegrationAction    = "INTEGRATION_ACTION"
	SelfService          = "self-service"
	Automation           = "automation"
	EntityCreated        = "ENTITY_CREATED"
	EntityUpdated        = "ENTITY_UPDATED"
	EntityDeleted        = "ENTITY_DELETED"
	TimerPropertyExpired = "TIMER_PROPERTY_EXPIRED"
	AnyEntityChange      = "ANY_ENTITY_CHANGE"
	RunCreated           = "RUN_CREATED"
	RunUpdated           = "RUN_UPDATED"
	AnyRunChange         = "ANY_RUN_CHANGE"
	JqCondition          = "JQ"

	// Workflow trigger node config types
	SelfServeTrigger = "SELF_SERVE_TRIGGER"
	EventTrigger     = "EVENT_TRIGGER"
	ScheduleTrigger  = "SCHEDULE_TRIGGER"

	// Workflow invocation node config types. Kafka, Webhook, UpsertEntity and
	// IntegrationAction reuse the constants declared above.
	AI            = "AI"
	AIAgent       = "AI_AGENT"
	ConditionNode = "CONDITION"
	InputNode     = "INPUT"

	// Workflow trigger event types
	WorkflowTimerExpired = "TIMER_EXPIRED"

	// Workflow self serve trigger context types
	CreateEntityContext = "CREATE_ENTITY"
	EntityContext       = "ENTITY"

	InstallationTypeOnPrem = "OnPrem"
	InstallationTypeSaas   = "Saas"

	IntegrationStatusCreating = "Creating"
	IntegrationStatusUpdating = "Updating"
	IntegrationStatusDeleting = "Deleting"
	IntegrationStatusDeleted  = "Deleted"
	IntegrationStatusError    = "Error"
)
