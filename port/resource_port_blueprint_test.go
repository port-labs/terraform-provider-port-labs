package port

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccPortBlueprint(t *testing.T) {
	var testAccActionConfigCreate = `
	provider "port-labs" {}
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "tf-test-bp0"
		properties {
			identifier = "text"
			type = "string"
			title = "text"
		}
		properties {
			identifier = "bool"
			type = "boolean"
			title = "boolean"
		}
		properties {
			identifier = "number"
			type = "number"
			title = "number"
		}
		properties {
			identifier = "obj"
			type = "object"
			title = "object"
		}
		properties {
			identifier = "array"
			type = "array"
			title = "array"
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

func TestAccPortBlueprintWithRelation(t *testing.T) {
	var testAccActionConfigCreate = `
	provider "port-labs" {}
	resource "port-labs_blueprint" "microservice1" {
		title = "TF Provider Test BP2"
		icon = "Terraform"
		identifier = "tf-provider-bp2"
		properties {
			identifier = "text"
			type = "string"
			title = "text"
		}
	}
	resource "port-labs_blueprint" "microservice2" {
		title = "TF Provider Test BP3"
		icon = "Terraform"
		identifier = "tf-provider-bp3"
		properties {
			identifier = "text"
			type = "string"
			title = "text"
		}
		relations {
			identifier = "test-rel"
			title = "Test Relation"
			target = port-labs_blueprint.microservice1.identifier
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
