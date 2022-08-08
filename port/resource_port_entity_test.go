package port

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccPortEntityUpdateProp(t *testing.T) {
	var testAccActionConfigCreate = `
	provider "port" {}
	resource "port_entity" "microservice" {
		title = "monolith"
		blueprint = "tf-provider-test-bp"
		properties {
			name = "text"
			value = "hedwig"
			type = "string"
		}
	}
`
	var testAccActionConfigUpdate = `
	provider "port" {}
	resource "port_entity" "microservice" {
		title = "monolith"
		blueprint = "tf-provider-test-bp"
		properties {
			name = "text"
			value = "hedwig2"
			type = "string"
		}
	}
`
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"port": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccActionConfigCreate,
			},
			{
				Config: testAccActionConfigUpdate,
				Check:  resource.TestCheckResourceAttr("port_entity.microservice", "properties.0.value", "hedwig2"),
			},
		},
	})
}

func TestAccPortEntity(t *testing.T) {
	var testAccActionConfigCreate = `
	provider "port" {}
	resource "port_entity" "microservice" {
		title = "monolith"
		blueprint = "tf-provider-test-bp"
		properties {
			name = "text"
			value = "hedwig"
			type = "string"
		}
		properties {
			name = "bool"
			value = "true"
			type = "boolean"
		}
		properties {
			name = "num"
			value = 123
			type = "number"
		}
		properties {
			name = "arr"
			items = [1,2,3]
			type = "array"
		}
		properties {
			name = "obj"
			value = jsonencode({"a":"b"})
			type = "object"
		}
	}
`
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"port": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccActionConfigCreate,
			},
		},
	})
}

func TestAccPortEntitiesRelation(t *testing.T) {
	var testAccActionConfigCreate = `
	provider "port" {}
	resource "port_entity" "microservice" {
		title = "monolith"
		blueprint = "tf-provider-test-bp"
		relations {
			name = "tf-relation"
			identifier = port_entity.microservice2.id
		}
		properties {
			name = "text"
			value = "test-relation"
			type = "string"
		}
	}
	resource "port_entity" "microservice2" {
		title = "monolith2"
		blueprint = "tf-provider-test-bp2"
		properties {
			name = "str"
			value = "test-relation"
			type = "string"
		}
	}
`
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"port": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccActionConfigCreate,
			},
		},
	})
}
