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
	Entity struct {
		Meta
		Identifier string         `json:"identifier,omitempty"`
		Title      string         `json:"title"`
		Blueprint  string         `json:"blueprint"`
		Team       []string       `json:"team,omitempty"`
		Properties map[string]any `json:"properties"`
		Relations  map[string]any `json:"relations"`
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

	InvocationMethod struct {
		Type                 string  `json:"type,omitempty"`
		Url                  *string `json:"url,omitempty"`
		Agent                *bool   `json:"agent,omitempty"`
		Synchronized         *bool   `json:"synchronized,omitempty"`
		Method               *string `json:"method,omitempty"`
		Org                  *string `json:"org,omitempty"`
		Repo                 *string `json:"repo,omitempty"`
		Webhook              *string `json:"webhook,omitempty"`
		Workflow             *string `json:"workflow,omitempty"`
		OmitPayload          *bool   `json:"omitPayload,omitempty"`
		OmitUserInputs       *bool   `json:"omitUserInputs,omitempty"`
		ReportWorkflowStatus *bool   `json:"reportWorkflowStatus,omitempty"`
		Branch               *string `json:"branch,omitempty"`
		ProjectName          *string `json:"projectName,omitempty"`
		GroupName            *string `json:"groupName,omitempty"`
		DefaultRef           *string `json:"defaultRef,omitempty"`
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

	ActionUserInputs = struct {
		Properties map[string]ActionProperty `json:"properties"`
		Required   any                       `json:"required,omitempty"`
		Order      []string                  `json:"order,omitempty"`
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
	}

	Action struct {
		ID                   string                `json:"id,omitempty"`
		Identifier           string                `json:"identifier,omitempty"`
		Description          *string               `json:"description,omitempty"`
		Title                string                `json:"title,omitempty"`
		Icon                 *string               `json:"icon,omitempty"`
		UserInputs           ActionUserInputs      `json:"userInputs"`
		Trigger              string                `json:"trigger"`
		RequiredApproval     *bool                 `json:"requiredApproval,omitempty"`
		InvocationMethod     *InvocationMethod     `json:"invocationMethod,omitempty"`
		ApprovalNotification *ApprovalNotification `json:"approvalNotification,omitempty"`
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
		Identifier string            `json:"identifier,omitempty"`
		Type       string            `json:"type,omitempty"`
		Icon       *string           `json:"icon,omitempty"`
		Parent     *string           `json:"parent,omitempty"`
		After      *string           `json:"after,omitempty"`
		Title      *string           `json:"title,omitempty"`
		Locked     *bool             `json:"locked,omitempty"`
		Blueprint  *string           `json:"blueprint,omitempty"`
		Widgets    *[]map[string]any `json:"widgets,omitempty"`
	}

	PageReadPermissions struct {
		Users []string `json:"users"`
		Roles []string `json:"roles"`
		Teams []string `json:"teams"`
	}

	PagePermissions struct {
		Read PageReadPermissions `json:"read"`
	}

	Relation struct {
		Identifier *string `json:"identifier,omitempty"`
		Title      *string `json:"title,omitempty"`
		Target     *string `json:"target,omitempty"`
		Required   *bool   `json:"required,omitempty"`
		Many       *bool   `json:"many,omitempty"`
	}

	Scorecard struct {
		Meta
		Identifier string `json:"identifier,omitempty"`
		Title      string `json:"title,omitempty"`
		Blueprint  string `json:"blueprint,omitempty"`
		Rules      []Rule `json:"rules,omitempty"`
	}

	Rule struct {
		Identifier string `json:"identifier,omitempty"`
		Title      string `json:"title,omitempty"`
		Level      string `json:"level,omitempty"`
		Query      Query  `json:"query,omitempty"`
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
)

type PortBody struct {
	OK                bool              `json:"ok"`
	Entity            Entity            `json:"entity"`
	Blueprint         Blueprint         `json:"blueprint"`
	Action            Action            `json:"action"`
	ActionPermissions ActionPermissions `json:"permissions"`
	Integration       Webhook           `json:"integration"`
	Scorecard         Scorecard         `json:"Scorecard"`
	Team              Team              `json:"team"`
	Page              Page              `json:"page"`
	MigrationId       string            `json:"migrationId"`
	Migration         Migration         `json:"migration"`
}

type PortPagePermissionsBody struct {
	OK              bool            `json:"ok"`
	PagePermissions PagePermissions `json:"permissions"`
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
