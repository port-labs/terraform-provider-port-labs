package cli

import "encoding/json"

type FactoryLineStage struct {
	Identifier       string          `json:"identifier"`
	Title            *string         `json:"title,omitempty"`
	Description      *string         `json:"description,omitempty"`
	Instructions     *string         `json:"instructions,omitempty"`
	Team             []string        `json:"team,omitempty"`
	Executor         string          `json:"executor"`
	ExecutorSettings json.RawMessage `json:"executor_settings"`
	Checkpoints      json.RawMessage `json:"checkpoints,omitempty"`
	MaxTries         int             `json:"max_tries"`
	AllowTakeOver    bool            `json:"allow_take_over"`
	Timeout          *string         `json:"timeout,omitempty"`
}

type FactoryLine struct {
	Identifier   string             `json:"identifier"`
	Title        string             `json:"title"`
	Description  *string            `json:"description,omitempty"`
	Team         []string           `json:"team,omitempty"`
	MaxRunning   *int               `json:"max_running,omitempty"`
	Status       *string            `json:"status,omitempty"`
	StatusReason *string            `json:"status_reason,omitempty"`
	Stages       []FactoryLineStage `json:"stages,omitempty"`
}

type FactoryLineResponse struct {
	Identifier   string             `json:"identifier"`
	Title        string             `json:"title"`
	Description  *string            `json:"description"`
	Team         []string           `json:"team"`
	Status       string             `json:"status"`
	StatusReason *string            `json:"statusReason"`
	MaxRunning   int                `json:"maxRunning"`
	StageOrder   []string           `json:"stageOrder"`
	Stages       []FactoryLineStage `json:"stages"`
	CreatedAt    *string            `json:"createdAt"`
	UpdatedAt    *string            `json:"updatedAt"`
}

type FactoryLineBody struct {
	OK       bool                `json:"ok"`
	Line     FactoryLineResponse `json:"line"`
	Warnings []struct {
		Path    string `json:"path"`
		Message string `json:"message"`
	} `json:"warnings,omitempty"`
}

type FactoryLineChange struct {
	Title        *string            `json:"title,omitempty"`
	Description  *string            `json:"description,omitempty"`
	Team         []string           `json:"team,omitempty"`
	MaxRunning   *int               `json:"max_running,omitempty"`
	Status       *string            `json:"status,omitempty"`
	StatusReason *string            `json:"status_reason,omitempty"`
	Stages       []FactoryLineStage `json:"stages,omitempty"`
}
