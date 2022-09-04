package cli

import (
	"time"
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
		Properties map[string]interface{} `json:"properties"`
		Relations  map[string]string      `json:"relations"`
		// TODO: add the rest of the fields.
	}

	BlueprintProperty struct {
		Type        string `json:"type,omitempty"`
		Title       string `json:"title,omitempty"`
		Identifier  string `json:"identifier,omitempty"`
		Default     string `json:"default,omitempty"`
		Format      string `json:"format,omitempty"`
		Description string `json:"description,omitempty"`
		Pattern     string `json:"pattern,omitempty"`
	}

	BlueprintSchema struct {
		Properties map[string]BlueprintProperty `json:"properties"`
		Required   []string                     `json:"required,omitempty"`
	}

	ActionUserInputs = BlueprintSchema

	Blueprint struct {
		Meta
		Identifier string          `json:"identifier,omitempty"`
		Title      string          `json:"title"`
		Icon       string          `json:"icon"`
		Schema     BlueprintSchema `json:"schema"`
	}

	Action struct {
		ID               string           `json:"id,omitempty"`
		Identifier       string           `json:"identifier,omitempty"`
		Description      string           `json:"description,omitempty"`
		Title            string           `json:"title,omitempty"`
		Icon             string           `json:"icon,omitempty"`
		UserInputs       ActionUserInputs `json:"userInputs"`
		Trigger          string           `json:"trigger"`
		InvocationMethod string           `json:"invocationMethod"`
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
