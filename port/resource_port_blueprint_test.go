package port

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func genID() string {
	id, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("t-%s", id[:18])
}

func TestAccPortBlueprint(t *testing.T) {
	identifier := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	provider "port-labs" {}
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
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
`, identifier)
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
	identifier1 := genID()
	identifier2 := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	provider "port-labs" {}
	resource "port-labs_blueprint" "microservice1" {
		title = "TF Provider Test BP2"
		icon = "Terraform"
		identifier = "%s"
		properties {
			identifier = "text"
			type = "string"
			title = "text"
		}
	}
	resource "port-labs_blueprint" "microservice2" {
		title = "TF Provider Test BP3"
		icon = "Terraform"
		identifier = "%s"
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
`, identifier1, identifier2)
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

func TestAccPortBlueprintUpdate(t *testing.T) {
	identifier := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	provider "port-labs" {}
	resource "port-labs_blueprint" "microservice1" {
		title = "TF Provider Test BP2"
		icon = "Terraform"
		identifier = "%s"
		properties {
			identifier = "text"
			type = "string"
			title = "text"
		}
	}
`, identifier)
	var testAccActionConfigUpdate = fmt.Sprintf(`
	provider "port-labs" {}
	resource "port-labs_blueprint" "microservice1" {
		title = "TF Provider Test BP2"
		icon = "Terraform"
		identifier = "%s"
		properties {
			identifier = "text"
			type = "string"
			title = "text"
		}
		properties {
			identifier = "number"
			type = "number"
			title = "num"
		}
	}
`, identifier)
	var testAccActionConfigUpdateAgain = fmt.Sprintf(`
	provider "port-labs" {}
	resource "port-labs_blueprint" "microservice1" {
		title = "TF Provider Test BP2"
		icon = "Terraform"
		identifier = "%s"
		properties {
			identifier = "number"
			type = "number"
			title = "num"
		}
	}
`, identifier)
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"port-labs": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.0.title", "text"),
				),
			},
			{
				Config: testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.0.title", "num"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.1.title", "text"),
				),
			},
			{
				Config: testAccActionConfigUpdateAgain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.0.title", "num"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.#", "1"),
				),
			},
		},
	})
}
