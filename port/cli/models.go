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
	}

	BlueprintSchema struct {
		Properties map[string]BlueprintProperty `json:"properties"`
	}

	Blueprint struct {
		Meta
		Identifier string          `json:"identifier,omitempty"`
		Title      string          `json:"title"`
		Icon       string          `json:"icon"`
		Schema     BlueprintSchema `json:"schema"`
		// TODO: relations
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
}
