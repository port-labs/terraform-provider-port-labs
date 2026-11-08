package integration

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func IntegrationSchemaForTest() schema.Schema {
	return schema.Schema{
		MarkdownDescription: IntegrationResourceMarkdownDescription,
		Attributes:          IntegrationSchema(),
	}
}
