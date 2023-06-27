package utils

import (
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func SpreadMaps(target map[string]schema.Attribute, source map[string]schema.Attribute) {
	for key, value := range source {
		target[key] = value
	}
}

func GenID() string {
	id, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("t-%s", id[:18])
}
