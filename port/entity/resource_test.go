package entity_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
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

func TestAccPortEntityImport(t *testing.T) {
	blueprintIdentifier := utils.GenID()
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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port_entity.microservice", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_entity.microservice", "properties.string_props.myStringIdentifier", "My String Value"),
				),
			},
			{
				ResourceName:            "port_entity.microservice",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           fmt.Sprintf("%s:%s", blueprintIdentifier, entityIdentifier),
				ImportStateVerifyIgnore: []string{"identifier"},
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
