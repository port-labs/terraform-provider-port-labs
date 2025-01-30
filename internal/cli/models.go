package cli

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type (
	Meta struct {
		CreatedAt *time.Time `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		CreatedBy string     `json:"createdBy,omitempty"`
		UpdatedBy string     `json:"updatedBy,omitempty"`
	}
	AccessTokenResponse struct {
		Ok          bool   `json:"ok"`
		AccessToken string `json:"accessToken"`
		ExpiresIn   int64  `json:"expiresIn"`
		TokenType   string `json:"tokenType"`
	}

	ScorecardRulesModel struct {
		Identifier string `json:"identifier"`
		Status     string `json:"status"`
		Level      string `json:"level"`
	}

	ScorecardModel struct {
		Rules []ScorecardRulesModel `json:"rules"`
		Level string                `json:"level"`
	}

	Entity struct {
		Meta
		Identifier string                    `json:"identifier,omitempty"`
		Title      string                    `json:"title"`
		Blueprint  string                    `json:"blueprint"`
		Team       []string                  `json:"team,omitempty"`
		Properties map[string]any            `json:"properties"`
		Relations  map[string]any            `json:"relations"`
		Scorecards map[string]ScorecardModel `json:"scorecards,omitempty"`
		// TODO: add the rest of the fields.
	}

	BlueprintProperty struct {
		Type               string              `json:"type,omitempty"`
		Title              *string             `json:"title,omitempty"`
		Identifier         string              `json:"identifier,omitempty"`
		Items              map[string]any      `json:"items,omitempty"`
		Default            any                 `json:"default,omitempty"`
		Icon               *string             `json:"icon,omitempty"`
		Format             *string             `json:"format,omitempty"`
		MaxLength          *int                `json:"maxLength,omitempty"`
		MinLength          *int                `json:"minLength,omitempty"`
		MaxItems           *int                `json:"maxItems,omitempty"`
		MinItems           *int                `json:"minItems,omitempty"`
		Maximum            *float64            `json:"maximum,omitempty"`
		Minimum            *float64            `json:"minimum,omitempty"`
		Description        *string             `json:"description,omitempty"`
		Blueprint          *string             `json:"blueprint,omitempty"`
		Pattern            *string             `json:"pattern,omitempty"`
		Enum               []any               `json:"enum,omitempty"`
		Spec               *string             `json:"spec,omitempty"`
		SpecAuthentication *SpecAuthentication `json:"specAuthentication,omitempty"`
		EnumColors         map[string]string   `json:"enumColors,omitempty"`
	}

	EntitiesSortModel struct {
		Property string `json:"property"`
		Order    string `json:"order"`
	}

	ActionProperty struct {
		Type               string              `json:"type,omitempty"`
		Title              *string             `json:"title,omitempty"`
		Identifier         string              `json:"identifier,omitempty"`
		Items              map[string]any      `json:"items,omitempty"`
		Default            any                 `json:"default,omitempty"`
		Icon               *string             `json:"icon,omitempty"`
		Format             *string             `json:"format,omitempty"`
		MaxLength          *int                `json:"maxLength,omitempty"`
		MinLength          *int                `json:"minLength,omitempty"`
		MaxItems           *int                `json:"maxItems,omitempty"`
		MinItems           *int                `json:"minItems,omitempty"`
		Maximum            *float64            `json:"maximum,omitempty"`
		Minimum            *float64            `json:"minimum,omitempty"`
		Description        *string             `json:"description,omitempty"`
		Blueprint          *string             `json:"blueprint,omitempty"`
		Pattern            *string             `json:"pattern,omitempty"`
		Enum               any                 `json:"enum,omitempty"`
		Spec               *string             `json:"spec,omitempty"`
		SpecAuthentication *SpecAuthentication `json:"specAuthentication,omitempty"`
		EnumColors         map[string]string   `json:"enumColors,omitempty"`
		DependsOn          []string            `json:"dependsOn,omitempty"`
		Dataset            *Dataset            `json:"dataset,omitempty"`
		Encryption         *string             `json:"encryption,omitempty"`
		Visible            any                 `json:"visible,omitempty"`
		Sort               *EntitiesSortModel  `json:"sort,omitempty"`
	}

	SpecAuthentication struct {
		ClientId         string `json:"clientId,omitempty"`
		AuthorizationUrl string `json:"authorizationUrl,omitempty"`
		TokenUrl         string `json:"tokenUrl,omitempty"`
	}

	DatasetValue struct {
		JqQuery string `json:"jqQuery,omitempty"`
	}
	DatasetRule struct {
		Blueprint *string       `json:"blueprint,omitempty"`
		Property  *string       `json:"property,omitempty"`
		Operator  string        `json:"operator,omitempty"`
		Value     *DatasetValue `json:"value,omitempty"`
	}
	Dataset struct {
		Combinator string        `json:"combinator,omitempty"`
		Rules      []DatasetRule `json:"rules,omitempty"`
	}

	BlueprintCalculationProperty struct {
		Type        string            `json:"type,omitempty"`
		Title       *string           `json:"title,omitempty"`
		Identifier  string            `json:"identifier,omitempty"`
		Calculation string            `json:"calculation,omitempty"`
		Default     any               `json:"default,omitempty"`
		Icon        *string           `json:"icon,omitempty"`
		Format      *string           `json:"format,omitempty"`
		Description *string           `json:"description,omitempty"`
		Colorized   *bool             `json:"colorized,omitempty"`
		Colors      map[string]string `json:"colors,omitempty"`
	}

	BlueprintAggregationProperty struct {
		Title           *string           `json:"title,omitempty"`
		Description     *string           `json:"description,omitempty"`
		Icon            *string           `json:"icon,omitempty"`
		Target          string            `json:"target,omitempty"`
		CalculationSpec map[string]string `json:"calculationSpec,omitempty"`
		Query           any               `json:"query,omitempty"`
	}

	BlueprintMirrorProperty struct {
		Identifier string  `json:"identifier,omitempty"`
		Title      *string `json:"title,omitempty"`
		Path       string  `json:"path,omitempty"`
	}

	BlueprintSchema struct {
		Properties map[string]BlueprintProperty `json:"properties"`
		Required   []string                     `json:"required,omitempty"`
	}

	MappingSchema struct {
		Identifier *string        `json:"identifier,omitempty"`
		Title      *string        `json:"title,omitempty"`
		Team       any            `json:"team,omitempty"`
		Icon       *string        `json:"icon,omitempty"`
		Properties map[string]any `json:"properties,omitempty"`
		Relations  map[string]any `json:"relations,omitempty"`
	}

	InvocationMethod struct {
		Type                 string            `json:"type"`
		Payload              any               `json:"payload,omitempty"`
		Url                  *string           `json:"url,omitempty"`
		Agent                any               `json:"agent,omitempty"`
		Synchronized         any               `json:"synchronized,omitempty"`
		Method               *string           `json:"method,omitempty"`
		Headers              map[string]string `json:"headers,omitempty"`
		Body                 any               `json:"body,omitempty"`
		Org                  *string           `json:"org,omitempty"`
		Repo                 *string           `json:"repo,omitempty"`
		Workflow             *string           `json:"workflow,omitempty"`
		WorkflowInputs       map[string]any    `json:"workflowInputs,omitempty"`
		ReportWorkflowStatus any               `json:"reportWorkflowStatus,omitempty"`
		Branch               *string           `json:"branch,omitempty"`
		ProjectName          *string           `json:"projectName,omitempty"`
		GroupName            *string           `json:"groupName,omitempty"`
		DefaultRef           *string           `json:"defaultRef,omitempty"`
		PipelineVariables    map[string]any    `json:"pipelineVariables,omitempty"`
		Webhook              *string           `json:"webhook,omitempty"`
		BlueprintIdentifier  *string           `json:"blueprintIdentifier,omitempty"`
		Mapping              *MappingSchema    `json:"mapping,omitempty"`
	}

	ApprovalNotification struct {
		Type   string  `json:"type,omitempty"`
		Url    string  `json:"url,omitempty"`
		Format *string `json:"format,omitempty"`
	}

	ChangelogDestination struct {
		Type  string `json:"type,omitempty"`
		Url   string `json:"url,omitempty"`
		Agent *bool  `json:"agent,omitempty"`
	}

	TeamInheritance struct {
		Path string `json:"path,omitempty"`
	}

	Ownership struct {
		Type string  `json:"type"`
		Path *string `json:"path,omitempty"`
	}

	Step = struct {
		Title string   `json:"title"`
		Order []string `json:"order"`
	}

	ActionUserInputs = struct {
		Properties map[string]ActionProperty `json:"properties"`
		Required   any                       `json:"required,omitempty"`
		Order      []string                  `json:"order,omitempty"`
		Steps      []Step                    `json:"steps,omitempty"`
	}

	TriggerEvent struct {
		Type                string  `json:"type"`
		BlueprintIdentifier *string `json:"blueprintIdentifier,omitempty"`
		PropertyIdentifier  *string `json:"propertyIdentifier,omitempty"`
		ActionIdentifier    *string `json:"actionIdentifier,omitempty"`
	}

	TriggerCondition struct {
		Expressions []string `json:"expressions,omitempty"`
		Combinator  *string  `json:"combinator,omitempty"`
		Rules       []any    `json:"rules,omitempty"`
		Type        string   `json:"type"`
	}

	Trigger struct {
		Type                string            `json:"type"`
		BlueprintIdentifier *string           `json:"blueprintIdentifier,omitempty"`
		Operation           *string           `json:"operation,omitempty"`
		UserInputs          *ActionUserInputs `json:"userInputs,omitempty"`
		Event               *TriggerEvent     `json:"event,omitempty"`
		Condition           *TriggerCondition `json:"condition,omitempty"`
	}

	Blueprint struct {
		Meta
		Identifier            string                                  `json:"identifier,omitempty"`
		Title                 string                                  `json:"title,omitempty"`
		Icon                  *string                                 `json:"icon,omitempty"`
		Description           *string                                 `json:"description,omitempty"`
		Schema                BlueprintSchema                         `json:"schema"`
		MirrorProperties      map[string]BlueprintMirrorProperty      `json:"mirrorProperties"`
		CalculationProperties map[string]BlueprintCalculationProperty `json:"calculationProperties"`
		AggregationProperties map[string]BlueprintAggregationProperty `json:"aggregationProperties,omitempty"`
		ChangelogDestination  *ChangelogDestination                   `json:"changelogDestination,omitempty"`
		TeamInheritance       *TeamInheritance                        `json:"teamInheritance,omitempty"`
		Relations             map[string]Relation                     `json:"relations"`
		Ownership             *Ownership                              `json:"ownership,omitempty"`
	}

	Action struct {
		ID                   string                `json:"id,omitempty"`
		Identifier           string                `json:"identifier"`
		Title                *string               `json:"title,omitempty"`
		Icon                 *string               `json:"icon,omitempty"`
		Description          *string               `json:"description,omitempty"`
		Trigger              *Trigger              `json:"trigger"`
		InvocationMethod     *InvocationMethod     `json:"invocationMethod,omitempty"`
		RequiredApproval     any                   `json:"requiredApproval,omitempty"`
		ApprovalNotification *ApprovalNotification `json:"approvalNotification,omitempty"`
		Publish              *bool                 `json:"publish,omitempty"`
	}

	ActionExecutePermissions struct {
		Users       []string        `json:"users"`
		Roles       []string        `json:"roles"`
		Teams       []string        `json:"teams"`
		OwnedByTeam *bool           `json:"ownedByTeam"`
		Policy      *map[string]any `json:"policy"`
	}

	ActionApprovePermissions struct {
		Users  []string        `json:"users"`
		Roles  []string        `json:"roles"`
		Teams  []string        `json:"teams"`
		Policy *map[string]any `json:"policy"`
	}

	ActionPermissions struct {
		Execute ActionExecutePermissions `json:"execute"`
		Approve ActionApprovePermissions `json:"approve"`
	}

	Page struct {
		Meta
		Identifier  string            `json:"identifier,omitempty"`
		Type        string            `json:"type,omitempty"`
		Icon        *string           `json:"icon,omitempty"`
		Parent      *string           `json:"parent,omitempty"`
		After       *string           `json:"after,omitempty"`
		Title       *string           `json:"title,omitempty"`
		Locked      *bool             `json:"locked,omitempty"`
		Blueprint   *string           `json:"blueprint,omitempty"`
		Widgets     *[]map[string]any `json:"widgets,omitempty"`
		Description *string           `json:"description,omitempty"`
	}

	PageReadPermissions struct {
		Users []string `json:"users"`
		Roles []string `json:"roles"`
		Teams []string `json:"teams"`
	}

	PagePermissions struct {
		Read PageReadPermissions `json:"read"`
	}

	BlueprintPermissionsBlock struct {
		Users       []string `json:"users"`
		Roles       []string `json:"roles"`
		Teams       []string `json:"teams"`
		OwnedByTeam *bool    `json:"ownedByTeam"`
	}

	BlueprintRolesOrPropertiesPermissionsBlock map[string]BlueprintPermissionsBlock

	BlueprintPermissionsEntities struct {
		Register         BlueprintPermissionsBlock                  `json:"register"`
		Unregister       BlueprintPermissionsBlock                  `json:"unregister"`
		Update           BlueprintPermissionsBlock                  `json:"update"`
		UpdateProperties BlueprintRolesOrPropertiesPermissionsBlock `json:"updateProperties,omitempty"`
		UpdateRelations  BlueprintRolesOrPropertiesPermissionsBlock `json:"updateRelations,omitempty"`
	}

	BlueprintPermissions struct {
		Entities BlueprintPermissionsEntities `json:"entities"`
	}

	Relation struct {
		Identifier  *string `json:"identifier,omitempty"`
		Title       *string `json:"title,omitempty"`
		Description *string `json:"description,omitempty"`
		Target      *string `json:"target,omitempty"`
		Required    *bool   `json:"required,omitempty"`
		Many        *bool   `json:"many,omitempty"`
	}

	Scorecard struct {
		Meta
		Identifier string  `json:"identifier,omitempty"`
		Title      string  `json:"title,omitempty"`
		Blueprint  string  `json:"blueprint,omitempty"`
		Filter     *Query  `json:"filter,omitempty"`
		Levels     []Level `json:"levels,omitempty"`
		Rules      []Rule  `json:"rules,omitempty"`
	}

	Rule struct {
		Identifier  string `json:"identifier,omitempty"`
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
		Level       string `json:"level,omitempty"`
		Query       Query  `json:"query,omitempty"`
	}

	Level struct {
		Title string `json:"title,omitempty"`
		Color string `json:"color,omitempty"`
	}
	Query struct {
		Combinator string `json:"combinator,omitempty"`
		Conditions []any  `json:"conditions,omitempty"`
	}

	Webhook struct {
		Meta
		Identifier  string     `json:"identifier,omitempty"`
		Title       *string    `json:"title,omitempty"`
		Icon        *string    `json:"icon,omitempty"`
		Description *string    `json:"description,omitempty"`
		Enabled     *bool      `json:"enabled,omitempty"`
		Security    *Security  `json:"security,omitempty"`
		Mappings    []Mappings `json:"mappings,omitempty"`
		WebhookKey  string     `json:"webhookKey,omitempty"`
		Url         string     `json:"url,omitempty"`
	}

	Security struct {
		Secret                *string `json:"secret,omitempty"`
		SignatureHeaderName   *string `json:"signatureHeaderName,omitempty"`
		SignatureAlgorithm    *string `json:"signatureAlgorithm,omitempty"`
		SignaturePrefix       *string `json:"signaturePrefix,omitempty"`
		RequestIdentifierPath *string `json:"requestIdentifierPath,omitempty"`
	}

	EntityProperty struct {
		Identifier string            `json:"identifier,omitempty"`
		Title      *string           `json:"title,omitempty"`
		Icon       *string           `json:"icon,omitempty"`
		Team       *string           `json:"team,omitempty"`
		Properties map[string]string `json:"properties,omitempty"`
		Relations  map[string]string `json:"relations,omitempty"`
	}

	Mappings struct {
		Blueprint    string          `json:"blueprint,omitempty"`
		Filter       *string         `json:"filter,omitempty"`
		Operation    *string         `json:"operation,omitempty"`
		ItemsToParse *string         `json:"itemsToParse,omitempty"`
		Entity       *EntityProperty `json:"entity,omitempty"`
	}

	Team struct {
		CreatedAt   *time.Time `json:"createdAt,omitempty"`
		UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
		Name        string     `json:"name,omitempty"`
		Description *string    `json:"description,omitempty"`
		Users       []string   `json:"users,omitempty"`
		Provider    string     `json:"provider,omitempty"`
	}

	Migration struct {
		Meta
		Id              string `json:"id,omitempty"`
		Actor           string `json:"actor,omitempty"`
		SourceBlueprint string `json:"sourceBlueprint,omitempty"`
		Mapping         any    `json:"mapping,omitempty"`
		Status          string `json:"status,omitempty"`
		DeleteBlueprint bool   `json:"deleteBlueprint,omitempty"`
		DeleteEntities  bool   `json:"deleteEntities,omitempty"`
		FailureCount    int    `json:"failureCount,omitempty"`
		SuccessCount    int    `json:"successCount,omitempty"`
	}

	SearchRequestQuery struct {
		Query                       *map[string]any `json:"query"`
		ExcludeCalculatedProperties *bool           `json:"exclude_calculated_properties,omitempty"`
		Include                     []string        `json:"include,omitempty"`
		Exclude                     []string        `json:"exclude,omitempty"`
		AttachTitleToRelation       *bool           `json:"attach_title_to_relation,omitempty"`
	}
)

type PortBody struct {
	OK                   bool              `json:"ok"`
	Entity               Entity            `json:"entity"`
	Blueprint            Blueprint         `json:"blueprint"`
	BlueprintPermissions Blueprint         `json:"blueprint_permissions"`
	Action               Action            `json:"action"`
	ActionPermissions    ActionPermissions `json:"permissions"`
	Webhook              Webhook           `json:"integration"`
	Scorecard            Scorecard         `json:"Scorecard"`
	Team                 Team              `json:"team"`
	Page                 Page              `json:"page"`
	MigrationId          string            `json:"migrationId"`
	Migration            Migration         `json:"migration"`
}

type SearchEntityResult struct {
	Meta
	Identifier string         `json:"identifier,omitempty"`
	Title      string         `json:"title,omitempty"`
	Icon       *string        `json:"icon,omitempty"`
	Team       []string       `json:"team,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
	Relations  map[string]any `json:"relations,omitempty"`
}

type SearchResult struct {
	OK                 bool     `json:"ok"`
	MatchingBlueprints []string `json:"matchingBlueprints"`
	Entities           []Entity `json:"entities"`
}

type PortPagePermissionsBody struct {
	OK              bool            `json:"ok"`
	PagePermissions PagePermissions `json:"permissions"`
}

type PortBlueprintPermissionsBody struct {
	OK                   bool                 `json:"ok"`
	BlueprintPermissions BlueprintPermissions `json:"permissions"`
}

type TeamUserBody struct {
	Email string `json:"email"`
}

type TeamPortBody struct {
	CreatedAt   *time.Time     `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time     `json:"updatedAt,omitempty"`
	Name        string         `json:"name,omitempty"`
	Description *string        `json:"description,omitempty"`
	Users       []TeamUserBody `json:"users,omitempty"`
	Provider    string         `json:"provider,omitempty"`
}

type PortTeamBody struct {
	OK   bool         `json:"ok"`
	Team TeamPortBody `json:"team"`
}

type PortProviderModel struct {
	ClientId types.String `tfsdk:"client_id"`
	Secret   types.String `tfsdk:"secret"`
	Token    types.String `tfsdk:"token"`
	BaseUrl  types.String `tfsdk:"base_url"`
}

type PortBodyDelete struct {
	Ok bool `json:"ok"`
}

type Integration struct {
	InstallationId       string                `json:"installationId"`
	Title                *string               `json:"title"`
	InstallationAppType  *string               `json:"installationAppType"`
	Version              *string               `json:"version"`
	Config               *map[string]any       `json:"config"`
	ChangelogDestination *ChangelogDestination `json:"changelogDestination,omitempty"`
}
