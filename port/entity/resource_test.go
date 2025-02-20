package entity_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func TestAccPortEntity(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
			"number_props" = {
				"myNumberIdentifier" =  {
					"title" = "My Number Identifier"
				}
			}
			"boolean_props" = {
				"myBooleanIdentifier" =  {
					"title" = "My Boolean Identifier"
				}
			}
			"object_props" = {
				"myObjectIdentifier" =  {
					"title" = "My Object Identifier"
				}
			}
			"array_props" = {
				"myStringArrayIdentifier" =  {
					"title" = "My String Array Identifier"
					"string_items" = {}
				}
				"myNumberArrayIdentifier" =  {
					"title" = "My Number Array Identifier"
					"number_items" = {}
				}
				"myBooleanArrayIdentifier" =  {
					"title" = "My Boolean Array Identifier"
					"boolean_items" = {}
				}
				"myObjectArrayIdentifier" =  {
					"title" = "My Object Array Identifier"
					"object_items" = {}
				}
			}
		}
	}
	resource "port_entity" "microservice" {
		title = "TF Provider Test Entity0"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.identifier
		properties = {
			"string_props" = {
				"myStringIdentifier" =  "My String Value"
			}
			"number_props" = {
				"myNumberIdentifier" =  123.456
			}
			"boolean_props" = {
				"myBooleanIdentifier" =  true
			}
			"object_props" = {
				"myObjectIdentifier" =  jsonencode({"foo": "bar"})
			}
			"array_props" = {
				string_items = {
					"myStringArrayIdentifier" =  ["My Array Value", "My Array Value2"]
				}
				number_items = {
					"myNumberArrayIdentifier" =  [123, 456]
				}
				boolean_items = {
					"myBooleanArrayIdentifier" =  [true, false]
				}
				object_items = {
					"myObjectArrayIdentifier" =  [jsonencode({"foo": "bar"}), jsonencode({"foo": "bar2"})]
				}
			}
		}
	}
	`, identifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier", "My String Value"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.number_props.myNumberIdentifier", "123.456"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.boolean_props.myBooleanIdentifier", "true"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.object_props.myObjectIdentifier", "{\"foo\":\"bar\"}"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.string_items.myStringArrayIdentifier.0", "My Array Value"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.string_items.myStringArrayIdentifier.1", "My Array Value2"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.number_items.myNumberArrayIdentifier.0", "123"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.number_items.myNumberArrayIdentifier.1", "456"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.boolean_items.myBooleanArrayIdentifier.0", "true"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.boolean_items.myBooleanArrayIdentifier.1", "false"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.object_items.myObjectArrayIdentifier.0", "{\"foo\":\"bar\"}"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.object_items.myObjectArrayIdentifier.1", "{\"foo\":\"bar2\"}"),
				),
			},
		},
	})
}

func TestAccPortEntityWithNulls(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
			"number_props" = {
				"myNumberIdentifier" =  {
					"title" = "My Number Identifier"
				}
			}
			"boolean_props" = {
				"myBooleanIdentifier" =  {
					"title" = "My Boolean Identifier"
				}
			}
			"object_props" = {
				"myObjectIdentifier" =  {
					"title" = "My Object Identifier"
				}
			}
			"array_props" = {
				"myStringArrayIdentifier" =  {
					"title" = "My String Array Identifier"
					"string_items" = {}
				}
				"myNumberArrayIdentifier" =  {
					"title" = "My Number Array Identifier"
					"number_items" = {}
				}
				"myBooleanArrayIdentifier" =  {
					"title" = "My Boolean Array Identifier"
					"boolean_items" = {}
				}
				"myObjectArrayIdentifier" =  {
					"title" = "My Object Array Identifier"
					"object_items" = {}
				}
			}
		}
	}
	resource "port_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port_blueprint.microservice.identifier
		properties = {
			"string_props" = {
				"myStringIdentifier" = null
			}
			"number_props" = {
				"myNumberIdentifier" =  null
			}
			"boolean_props" = {
				"myBooleanIdentifier" =  null
			}
			"object_props" = {
				"myObjectIdentifier" = null
			}
			"array_props" = {
				string_items = {
					"myStringArrayIdentifier" = null
				}
				number_items = {
					"myNumberArrayIdentifier" = null
				}
				boolean_items = {
					"myBooleanArrayIdentifier" = null
				}
				object_items = {
					"myObjectArrayIdentifier" = null
				}
			}
		}
	}
	`, identifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", identifier),
					resource.TestCheckNoResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier"),
					resource.TestCheckNoResourceAttr("port_entity.microservice", "properties.number_props.myNumberIdentifier"),
					resource.TestCheckNoResourceAttr("port_entity.microservice", "properties.boolean_props.myBooleanIdentifier"),
					resource.TestCheckNoResourceAttr("port_entity.microservice", "properties.object_props.myObjectIdentifier"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.string_items.myStringArrayIdentifier.#", "0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.number_items.myNumberArrayIdentifier.#", "0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.boolean_items.myBooleanArrayIdentifier.#", "0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.object_items.myObjectArrayIdentifier.#", "0"),
				),
			},
		},
	})
}

func TestAccPortEntityWithRelation(t *testing.T) {
	identifier := utils.GenID()
	identifier2 := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
		relations = {
			"tfRelation" = {
				"title" = "Test Relation"
				"target" = port_blueprint.microservice2.identifier
			}
		}
	}
	resource "port_blueprint" "microservice2" {
		title = "TF Provider Test BP1"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier2" =  {
					"title" = "My String Identifier2"
				}
			}
		}
	}

	resource "port_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port_blueprint.microservice.identifier
		properties = {
			"string_props" = {
				"myStringIdentifier" =  "My String Value"
			}
		}
		relations = {
			single_relations = {
				"tfRelation" = port_entity.microservice2.identifier
			}
		}
	}

	resource "port_entity" "microservice2" {
		title = "TF Provider Test Entity1"
		identifier = "tf-entity-2"
		blueprint = port_blueprint.microservice2.identifier
		properties = {
			"string_props" = {
				"myStringIdentifier2" =  "My String Value2"
			}
		}
	}
	`, identifier, identifier2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier", "My String Value"),
					resource.TestCheckResourceAttr("port_entity.microservice", "relations.single_relations.tfRelation", "tf-entity-2"),
				),
			},
		},
	})
}

func TestAccPortEntityWithManyRelation(t *testing.T) {
	identifier1 := utils.GenID()
	identifier2 := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
		relations = {
			"tfRelation" = {
				"title" = "Test Relation"
				"target" = port_blueprint.microservice2.identifier
				"many" = true
			}
		}
	}
	resource "port_blueprint" "microservice2" {
		title = "TF Provider Test BP1"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier2" =  {
					"title" = "My String Identifier2"
				}
			}
		}
	}

	resource "port_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port_blueprint.microservice.identifier
		properties = {
			"string_props" = {
				"myStringIdentifier" =  "My String Value"
			}
		}
		relations = {
			"many_relations" = {
				"tfRelation" = [port_entity.microservice2.identifier, port_entity.microservice3.identifier]
			}
		}
	}

	resource "port_entity" "microservice2" {
		title = "TF Provider Test Entity1"
		identifier = "tf-entity-2"
		blueprint = port_blueprint.microservice2.identifier
		properties = {
			"string_props" = {
				"myStringIdentifier2" =  "My String Value2"
			}
		}
	}

	resource "port_entity" "microservice3" {
		title = "TF Provider Test Entity2"
		identifier = "tf-entity-3"
		blueprint = port_blueprint.microservice2.identifier
		properties = {
			"string_props" = {
				"myStringIdentifier2" =  "My String Value3"
			}
		}
	}
	`, identifier1, identifier2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", identifier1),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier", "My String Value"),
					resource.TestCheckResourceAttr("port_entity.microservice", "relations.many_relations.tfRelation.0", "tf-entity-2"),
					resource.TestCheckResourceAttr("port_entity.microservice", "relations.many_relations.tfRelation.1", "tf-entity-3"),
				),
			},
		},
	})
}

func TestAccPortEntityWithEmptyRelation(t *testing.T) {
	identifier := utils.GenID()
	identifier2 := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
		relations = {
			"tfRelation" = {
				"title" = "Test Relation"
				"target" = port_blueprint.microservice2.identifier
			}
		}
	}
	resource "port_blueprint" "microservice2" {
		title = "TF Provider Test BP1"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier2" =  {
					"title" = "My String Identifier2"
				}
			}
		}
	}

	resource "port_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port_blueprint.microservice.identifier
		properties = {
			"string_props" = {
				"myStringIdentifier" =  "My String Value"
			}
		}
		relations = {
			single_relations = {
				"tfRelation" = null
			}
		}
	}
	`, identifier, identifier2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier", "My String Value"),
					resource.TestCheckNoResourceAttr("port_entity.microservice", "relations.single_relations.tfRelation"),
				),
			},
		},
	})
}

func TestAccPortEntityImport(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	entityIdentifier := utils.GenID()
	entityId := fmt.Sprintf("%s:%s", blueprintIdentifier, entityIdentifier)
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
	}
	resource "port_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port_blueprint.microservice.identifier
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  "My String Value"
			}
		}
	}`, blueprintIdentifier, entityIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "id", entityId),
					resource.TestCheckResourceAttr("port_entity.microservice", "identifier", entityIdentifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier", "My String Value"),
				),
			},
			{
				ResourceName:            "port_entity.microservice",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           entityId,
				ImportStateVerifyIgnore: []string{"identifier", "create_missing_related_entities"},
			},
		},
	})
}

func TestAccPortEntityUpdateProp(t *testing.T) {

	identifier := utils.GenID()
	entityIdentifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
	}
	resource "port_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port_blueprint.microservice.identifier
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  "My String Value"
			}
		}
	}`, identifier, entityIdentifier)

	var testAccActionConfigUpdate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
	}
	resource "port_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port_blueprint.microservice.identifier
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  "My String Value2"
			}
		}
	}`, identifier, entityIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier", "My String Value"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier", "My String Value2"),
				),
			},
		},
	})
}

func TestAccPortEntityUpdateIdentifier(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	entityIdentifier := utils.GenID()
	entityUpdatedIdentifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
	}
	resource "port_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port_blueprint.microservice.identifier
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  "My String Value"
			}
		}
	}`, blueprintIdentifier, entityIdentifier)

	var testAccActionConfigUpdate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
	}
	resource "port_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port_blueprint.microservice.identifier
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  "My String Value2"
			}
		}
	}`, blueprintIdentifier, entityUpdatedIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "identifier", entityIdentifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier", "My String Value"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "identifier", entityUpdatedIdentifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier", "My String Value2"),
				),
			},
		},
	})

}

func TestAccPortEntityUpdateBlueprintIdentifier(t *testing.T) {

	blueprintIdentifier := utils.GenID()
	blueprintIdentifier2 := utils.GenID()
	entityIdentifier := utils.GenID()

	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
	}
	resource "port_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port_blueprint.microservice.identifier
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  "My String Value"
			}
		}
	}`, blueprintIdentifier, entityIdentifier)

	var testAccActionConfigUpdate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
	}

	resource "port_blueprint" "microservice2" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
	}

	resource "port_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port_blueprint.microservice2.identifier
		identifier = "%s"
		properties = {
			"string_props" = {
				"myStringIdentifier" =  "My String Value2"
			}
		}
	}`, blueprintIdentifier, blueprintIdentifier2, entityIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "identifier", entityIdentifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier", "My String Value"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "identifier", entityIdentifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", blueprintIdentifier2),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier", "My String Value2"),
				),
			},
		},
	})
}

func TestAccPortEntityWithDefaultArrayProp(t *testing.T) {
	identifier := utils.GenID()
	entityIdentifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"array_props" = {
				"myArrayIdentifier" =  {
					"title" = "My Array Identifier"
				}
			}
		}	
	}
	resource "port_entity" "microservice" {
		title = "TF Provider Test Entity0"	
		blueprint = port_blueprint.microservice.identifier	
		identifier = "%s"	
		properties = {	
			array_props = {
				string_items = {
					"myArrayIdentifier" =  ["My Array Value", "My Array Value2"]
				}
			}
		}
	}`, identifier, entityIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.string_items.myArrayIdentifier.0", "My Array Value"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.string_items.myArrayIdentifier.1", "My Array Value2"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.string_items.myArrayIdentifier.0", "My Array Value"),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.array_props.string_items.myArrayIdentifier.1", "My Array Value2"),
				),
			},
		},
	})
}

func TestAccPortEntityWithMissingRelation(t *testing.T) {
	identifier := utils.GenID()
	identifier2 := utils.GenID()
	testAccActionConfig := fmt.Sprintf(`
    resource "port_blueprint" "microservice" {
        title = "TF Provider Test BP0"
        icon = "Terraform"
        identifier = "%s"
        properties = {
            "string_props" = {
                "myStringIdentifier" =  {
                    "title" = "My String Identifier"
                }
            }
        }
        relations = {
            "tfRelation" = {
                "title" = "Test Relation"
                "target" = port_blueprint.microservice2.identifier
            }
        }
        force_delete_entities = true
    }
    resource "port_blueprint" "microservice2" {
        title = "TF Provider Test BP1"
        icon = "Terraform"
        identifier = "%s"
        properties = {
            "string_props" = {
                "myStringIdentifier2" =  {
                    "title" = "My String Identifier2"
                }
            }
        }
        force_delete_entities = true
    }

    // A: microservice
    resource "port_entity" "microservice" {
        title = "TF Provider Test Entity0"
        blueprint = port_blueprint.microservice.identifier
        properties = {
            "string_props" = {
                "myStringIdentifier" =  "My String Value"
            }
        }
        relations = {
            single_relations = {
                "tfRelation" = "new-identifier" // Non-existing identifier to trigger creation
            }
        }
        create_missing_related_entities = true
    }

    // B: microservice2
    resource "port_entity" "microservice2" {
        title = "TF Provider Test Entity1"
        identifier = "tf-entity-2"
        blueprint = port_blueprint.microservice2.identifier
        properties = {
            "string_props" = {
                "myStringIdentifier2" =  "My String Value2"
            }
        }
    }
    `, identifier, identifier2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier", "My String Value"),
					resource.TestCheckResourceAttr("port_entity.microservice", "relations.single_relations.tfRelation", "new-identifier"),
				),
			},
			{
				// Delete the blueprints with force_delete_entities set to true
				Config:  acctest.ProviderConfig + testAccActionConfig,
				Destroy: true,
			},
		},
	})
}
