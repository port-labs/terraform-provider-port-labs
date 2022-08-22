package port

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccPortEntityUpdateProp(t *testing.T) {
	var testAccActionConfigCreate = `
	provider "port-labs" {}
	resource "port-labs_entity" "microservice" {
		title = "monolith"
		blueprint = "tf-provider-test-bp"
		properties {
			name = "text"
			value = "hedwig"
		}
	}
`
	var testAccActionConfigUpdate = `
	provider "port-labs" {}
	resource "port-labs_entity" "microservice" {
		title = "monolith"
		blueprint = "tf-provider-test-bp"
		properties {
			name = "text"
			value = "hedwig2"
		}
	}
`
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"port-labs": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccActionConfigCreate,
			},
			{
				Config: testAccActionConfigUpdate,
				Check:  resource.TestCheckResourceAttr("port-labs_entity.microservice", "properties.0.value", "hedwig2"),
			},
		},
	})
}

func TestAccPortEntity(t *testing.T) {
	var testAccActionConfigCreate = `
	provider "port-labs" {}
	resource "port-labs_entity" "microservice" {
		title = "monolith"
		blueprint = "tf-provider-test-bp"
		properties {
			name = "text"
			value = "hedwig"
		}
		properties {
			name = "bool"
			value = "true"
		}
		properties {
			name = "num"
			value = 123
		}
		properties {
			name = "arr"
			items = [1,2,3]
		}
		properties {
			name = "obj"
			value = jsonencode({"a":"b"})
		}
	}
`
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"port-labs": Provider(),
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
	provider "port-labs" {}
	resource "port-labs_entity" "microservice" {
		title = "monolith"
		blueprint = "tf-provider-test-bp"
		relations {
			name = "tf-relation"
			identifier = port-labs_entity.microservice2.id
		}
		properties {
			name = "text"
			value = "test-relation"
		}
	}
	resource "port-labs_entity" "microservice2" {
		title = "monolith2"
		blueprint = "tf-provider-test-bp2"
		properties {
			name = "str"
			value = "test-relation"
		}
	}
`
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"port-labs": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccActionConfigCreate,
			},
		},
	})
}
