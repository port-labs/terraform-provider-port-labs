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
		Identifier string                 `json:"identifier,omitempty"`
		Title      string                 `json:"title"`
		Blueprint  string                 `json:"blueprint"`
		Team       []string               `json:"team,omitempty"`
		Properties map[string]interface{} `json:"properties"`
		Relations  map[string]interface{} `json:"relations"`
		// TODO: add the rest of the fields.
	}

	BlueprintProperty struct {
		Type               string              `json:"type,omitempty"`
		Title              string              `json:"title,omitempty"`
		Identifier         string              `json:"identifier,omitempty"`
		Items              map[string]any      `json:"items,omitempty"`
		Default            interface{}         `json:"default,omitempty"`
		Icon               string              `json:"icon,omitempty"`
		Format             string              `json:"format,omitempty"`
		MaxLength          int                 `json:"maxLength,omitempty"`
		MinLength          int                 `json:"minLength,omitempty"`
		MaxItems           int                 `json:"maxItems,omitempty"`
		MinItems           int                 `json:"minItems,omitempty"`
		Maximum            float64             `json:"maximum,omitempty"`
		Minimum            float64             `json:"minimum,omitempty"`
		Description        string              `json:"description,omitempty"`
		Blueprint          string              `json:"blueprint,omitempty"`
		Pattern            string              `json:"pattern,omitempty"`
		Enum               []interface{}       `json:"enum,omitempty"`
		Spec               string              `json:"spec,omitempty"`
		SpecAuthentication *SpecAuthentication `json:"specAuthentication,omitempty"`
		EnumColors         map[string]string   `json:"enumColors,omitempty"`
	}

	SpecAuthentication struct {
		ClientId         string `json:"clientId,omitempty"`
		AuthorizationUrl string `json:"authorizationUrl,omitempty"`
		TokenUrl         string `json:"tokenUrl,omitempty"`
	}

	BlueprintCalculationProperty struct {
		Type        string            `json:"type,omitempty"`
		Title       string            `json:"title,omitempty"`
		Identifier  string            `json:"identifier,omitempty"`
		Calculation string            `json:"calculation,omitempty"`
		Default     interface{}       `json:"default,omitempty"`
		Icon        string            `json:"icon,omitempty"`
		Format      string            `json:"format,omitempty"`
		Description string            `json:"description,omitempty"`
		Colorized   bool              `json:"colorized,omitempty"`
		Colors      map[string]string `json:"colors,omitempty"`
	}

	BlueprintMirrorProperty struct {
		Identifier string `json:"identifier,omitempty"`
		Title      string `json:"title,omitempty"`
		Path       string `json:"path,omitempty"`
	}

	BlueprintSchema struct {
		Properties map[string]BlueprintProperty `json:"properties"`
		Required   []string                     `json:"required,omitempty"`
	}

	InvocationMethod struct {
		Type                 string `json:"type,omitempty"`
		Url                  string `json:"url,omitempty"`
		Agent                bool   `json:"agent,omitempty"`
		Org                  string `json:"org,omitempty"`
		Repo                 string `json:"repo,omitempty"`
		Webhook              string `json:"webhook,omitempty"`
		Workflow             string `json:"workflow,omitempty"`
		OmitPayload          bool   `json:"omitPayload,omitempty"`
		OmitUserInputs       bool   `json:"omitUserInputs,omitempty"`
		ReportWorkflowStatus *bool  `json:"reportWorkflowStatus,omitempty"`
	}

	ChangelogDestination struct {
		Type  string `json:"type,omitempty"`
		Url   string `json:"url,omitempty"`
		Agent bool   `json:"agent,omitempty"`
	}

	ActionUserInputs = BlueprintSchema

	Blueprint struct {
		Meta
		Identifier            string                                  `json:"identifier,omitempty"`
		Title                 string                                  `json:"title,omitempty"`
		Icon                  string                                  `json:"icon,omitempty"`
		Description           string                                  `json:"description,omitempty"`
		Schema                BlueprintSchema                         `json:"schema"`
		MirrorProperties      map[string]BlueprintMirrorProperty      `json:"mirrorProperties"`
		CalculationProperties map[string]BlueprintCalculationProperty `json:"calculationProperties"`
		ChangelogDestination  *ChangelogDestination                   `json:"changelogDestination,omitempty"`
		Relations             map[string]Relation                     `json:"relations"`
	}

	Action struct {
		ID               string            `json:"id,omitempty"`
		Identifier       string            `json:"identifier,omitempty"`
		Description      string            `json:"description,omitempty"`
		Title            string            `json:"title,omitempty"`
		Icon             string            `json:"icon,omitempty"`
		UserInputs       ActionUserInputs  `json:"userInputs"`
		Trigger          string            `json:"trigger"`
		RequiredApproval bool              `json:"requiredApproval,omitempty"`
		InvocationMethod *InvocationMethod `json:"invocationMethod,omitempty"`
	}

	Relation struct {
		Identifier string `json:"identifier,omitempty"`
		Title      string `json:"title,omitempty"`
		Target     string `json:"target,omitempty"`
		Required   bool   `json:"required,omitempty"`
		Many       bool   `json:"many,omitempty"`
	}
)

type PortBody struct {
	OK        bool      `json:"ok"`
	Entity    Entity    `json:"entity"`
	Blueprint Blueprint `json:"blueprint"`
	Action    Action    `json:"action"`
}

type PortProviderModel struct {
	ClientId types.String `tfsdk:"client_id"`
	Secret   types.String `tfsdk:"secret"`
	Token    types.String `tfsdk:"token"`
	BaseUrl  types.String `tfsdk:"base_url"`
}

type ChangelogDestinationModel struct {
	Type  types.String `tfsdk:"type"`
	Url   types.String `tfsdk:"url"`
	Agent types.Bool   `tfsdk:"agent"`
}

type StringPropModel struct {
	Title       types.String `tfsdk:"title"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
	Default     types.String `tfsdk:"default"`
	Required    types.Bool   `tfsdk:"required"`
	Format      types.String `tfsdk:"format"`
	MaxLength   types.Int64  `tfsdk:"max_length"`
	MinLength   types.Int64  `tfsdk:"min_length"`
	Pattern     types.String `tfsdk:"pattern"`
}

type NumberPropModel struct {
	Title       types.String  `tfsdk:"title"`
	Icon        types.String  `tfsdk:"icon"`
	Description types.String  `tfsdk:"description"`
	Default     types.Float64 `tfsdk:"default"`
	Required    types.Bool    `tfsdk:"required"`
	Maximum     types.Float64 `tfsdk:"maximum"`
	Minimum     types.Float64 `tfsdk:"minimum"`
}

type BooleanPropModel struct {
	Title       types.String `tfsdk:"title"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
	Default     types.Bool   `tfsdk:"default"`
	Required    types.Bool   `tfsdk:"required"`
}

type ItemsModal struct {
	Type    types.String `tfsdk:"type"`
	Format  types.String `tfsdk:"format"`
	Default types.List   `tfsdk:"default"`
}
type ArrayPropModel struct {
	Title       types.String `tfsdk:"title"`
	Icon        types.String `tfsdk:"icon"`
	Description types.String `tfsdk:"description"`
	MaxItems    types.Int64  `tfsdk:"max_items"`
	MinItems    types.Int64  `tfsdk:"min_items"`
	Required    types.Bool   `tfsdk:"required"`
	Items       *ItemsModal  `tfsdk:"items"`
	// Default     types.ListType `tfsdk:"default"`
}

type PropertiesModel struct {
	StringProp  map[string]StringPropModel  `tfsdk:"string_prop"`
	NumberProp  map[string]NumberPropModel  `tfsdk:"number_prop"`
	BooleanProp map[string]BooleanPropModel `tfsdk:"boolean_prop"`
	ArrayProp   map[string]ArrayPropModel   `tfsdk:"array_prop"`
}

type BlueprintModel struct {
	Identifier           types.String               `tfsdk:"identifier"`
	Title                types.String               `tfsdk:"title"`
	Icon                 types.String               `tfsdk:"icon"`
	Description          types.String               `tfsdk:"description"`
	CreatedAt            types.String               `tfsdk:"created_at"`
	CreatedBy            types.String               `tfsdk:"created_by"`
	UpdatedAt            types.String               `tfsdk:"updated_at"`
	UpdatedBy            types.String               `tfsdk:"updated_by"`
	ChangelogDestination *ChangelogDestinationModel `tfsdk:"changelog_destination"`
	Properties           *PropertiesModel           `tfsdk:"properties"`
}
