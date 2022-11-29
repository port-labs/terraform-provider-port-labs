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
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties {
			identifier = "bool"
			type = "boolean"
			title = "boolean"
			default = true
		}
		properties {
			identifier = "number"
			type = "number"
			title = "number"
			default = 1
		}
		properties {
			identifier = "obj"
			type = "object"
			title = "object"
			default = jsonencode({"a":"b"})
		}
		properties {
			identifier = "array"
			type = "array"
			title = "array"
			default_items = [1, 2, 3]
		}
		properties {
			identifier = "text"
			type = "string"
			title = "text"
			icon = "Terraform"
			enum = ["a", "b", "c"]
			enum_colors = {
				a = "red"
				b = "blue"
			}
			default = "a"
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
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.0.default_items.0", "1"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.0.default_items.#", "3"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.1.default", "1"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.2.identifier", "text"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.2.enum.0", "a"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.2.enum_colors.a", "red"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.2.default", "a"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.3.default", "true"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.4.default", "{\"a\":\"b\"}"),
				),
			},
		},
	})
}

func TestAccBlueprintWithChangelogDestination(t *testing.T) {
	identifier := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties {
			identifier = "text"
			type = "string"
			title = "text"
		}
		changelog_destination {
			type = "WEBHOOK"
			url = "https://google.com"
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
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "changelog_destination.0.type", "WEBHOOK"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "changelog_destination.0.url", "https://google.com"),
				),
			},
		},
	})
}

func TestAccPortBlueprintWithRelation(t *testing.T) {
	identifier1 := genID()
	identifier2 := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
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
	resource "port-labs_blueprint" "microservice1" {
		description = "Test Description"
		title = "TF Provider Test BP2"
		icon = "Terraform"
		identifier = "%s"
		properties {
			required = true
			identifier = "text"
			type = "string"
			icon = "Terraform"
			title = "text"
			enum = ["a", "b", "c"]
			enum_colors = {
				a = "red"
				b = "blue"
			}
		}
		formula_properties {
			identifier = "formula_id"
			formula = "{{$identifier}}formula"
		}
	}
`, identifier)
	var testAccActionConfigUpdate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice1" {
		title = "TF Provider Test BP2"
		icon = "Terraform"
		identifier = "%s"
		properties {
			required = false
			identifier = "text"
			type = "string"
			title = "text"
		}
		properties {
			identifier = "number"
			type = "number"
			title = "num"
		}
		formula_properties {
			identifier = "formula_id"
			formula = "{{$identifier}}formula-updated"
		}
	}
`, identifier)
	var testAccActionConfigUpdateAgain = fmt.Sprintf(`
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
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "description", "Test Description"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.0.title", "text"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.0.required", "true"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.0.icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "formula_properties.0.formula", "{{$identifier}}formula"),
				),
			},
			{
				Config: testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "description", ""),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.0.title", "num"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.1.title", "text"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.1.required", "false"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "formula_properties.0.formula", "{{$identifier}}formula-updated"),
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

func TestAccPortBlueprintUpdateRelation(t *testing.T) {
	envID := genID()
	vmID := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "Environment" {
		title = "Environment"
		icon = "Environment"
		identifier = "%s"
		properties {
			identifier = "env_name"
			type = "string"
			title = "Name"
		}
	}
	resource "port-labs_blueprint" "vm" {
		title = "Virtual Machine"
		icon = "Azure"
		identifier = "%s"
		properties {
			identifier = "image"
			type = "string"
			title = "Image"
		}
		relations {
			identifier = "vm-to-environment"
			title = "Related Environment"
			target = port-labs_blueprint.Environment.identifier
		}
	}
`, envID, vmID)
	var testAccActionConfigUpdate = fmt.Sprintf(`
	resource "port-labs_blueprint" "Environment" {
		title = "Environment"
		icon = "Environment"
		identifier = "%s"
		properties {
			identifier = "env_name"
			type = "string"
			title = "Name"
		}
	}
	resource "port-labs_blueprint" "vm" {
		title = "Virtual Machine"
		icon = "Azure"
		identifier = "%s"
		properties {
			identifier = "image"
			type = "string"
			title = "Image"
		}
		relations {
			identifier = "environment"
			title = "Related Environment"
			target = port-labs_blueprint.Environment.identifier
		}
	}
`, envID, vmID)
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"port-labs": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.#", "1"),
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.0.title", "Related Environment"),
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.0.target", envID),
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.0.identifier", "vm-to-environment"),
				),
			},
			{
				Config: testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.#", "1"),
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.0.title", "Related Environment"),
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.0.target", envID),
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.0.identifier", "environment"),
				),
			},
		},
	})
}

func TestAccPortBlueprintWithMirrorProperty(t *testing.T) {
	identifier1 := genID()
	identifier2 := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
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
		mirror_properties {
			identifier = "mirror-for-microservice1"
			title = "Mirror for microservice1"
			path = "test-rel.$identifier"
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

func TestAccPortBlueprintUpdateMirrorProperty(t *testing.T) {
	envID := genID()
	vmID := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "Environment" {
		title = "Environment"
		icon = "Environment"
		identifier = "%s"
		properties {
			identifier = "env_name"
			type = "string"
			title = "Name"
		}
	}
	resource "port-labs_blueprint" "vm" {
		title = "Virtual Machine"
		icon = "Azure"
		identifier = "%s"
		properties {
			identifier = "image"
			type = "string"
			title = "Image"
		}
		mirror_properties {
			identifier = "mirror-for-environment"
			title = "Mirror for environment"
			path = "vm-to-environment.$identifier"
		}
		relations {
			identifier = "vm-to-environment"
			title = "Related Environment"
			target = port-labs_blueprint.Environment.identifier
		}
	}
`, envID, vmID)
	var testAccActionConfigUpdate = fmt.Sprintf(`
	resource "port-labs_blueprint" "Environment" {
		title = "Environment"
		icon = "Environment"
		identifier = "%s"
		properties {
			identifier = "env_name"
			type = "string"
			title = "Name"
		}
	}
	resource "port-labs_blueprint" "vm" {
		title = "Virtual Machine"
		icon = "Azure"
		identifier = "%s"
		properties {
			identifier = "image"
			type = "string"
			title = "Image"
		}
		mirror_properties {
			identifier = "mirror-for-environment"
			title = "Mirror for environment2"
			path = "environment.$identifier"
		}
		relations {
			identifier = "environment"
			title = "Related Environment"
			target = port-labs_blueprint.Environment.identifier
		}
	}
`, envID, vmID)
	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"port-labs": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "mirror_properties.#", "1"),
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "mirror_properties.0.title", "Mirror for environment"),
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "mirror_properties.0.identifier", "mirror-for-environment"),
				),
			},
			{
				Config: testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "mirror_properties.#", "1"),
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "mirror_properties.0.title", "Mirror for environment2"),
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "mirror_properties.0.identifier", "mirror-for-environment"),
				),
			},
		},
	})
}
