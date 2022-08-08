package port

import "time"

type Entity struct {
	Identifier string                 `json:"identifier"`
	Title      string                 `json:"title"`
	CreatedAt  time.Time              `json:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt"`
	CreatedBy  string                 `json:"createdBy"`
	UpdatedBy  string                 `json:"updatedBy"`
	Properties map[string]interface{} `json:"properties"`

	// TODO: add the rest of the fields
}

type PortBody struct {
	OK     bool   `json:"ok"`
	Entity Entity `json:"entity"`
}
