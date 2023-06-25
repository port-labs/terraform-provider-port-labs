package entity_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/acctest"
)

func genID() string {
	id, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("t-%s", id[:18])
}

func TestAccPortEntity(t *testing.T) {
	identifier := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_prop" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
			"number_prop" = {
				"myNumberIdentifier" =  {
					"title" = "My Number Identifier"
				}
			}
			"boolean_prop" = {
				"myBooleanIdentifier" =  {
					"title" = "My Boolean Identifier"
				}
			}
			"object_prop" = {
				"myObjectIdentifier" =  {
					"title" = "My Object Identifier"
				}
			}
			"array_prop" = {
				"myArrayIdentifier" =  {
					"title" = "My Array Identifier"
					"string_items" = {}
				}
			}
		}
	}
	resource "port-labs_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port-labs_blueprint.microservice.id
		properties = {
			"string_prop" = {
				"myStringIdentifier" =  "My String Value"
			}
			"number_prop" = {
				"myNumberIdentifier" =  123
			}
			"boolean_prop" = {
				"myBooleanIdentifier" =  true
			}
			"object_prop" = {
				"myObjectIdentifier" =  jsonencode({"foo": "bar"})
			}
			"array_prop" = {
				string_items = {
					"myArrayIdentifier" =  ["My Array Value"]
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
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "properties.string_prop.myStringIdentifier", "My String Value"),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "properties.number_prop.myNumberIdentifier", "123"),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "properties.boolean_prop.myBooleanIdentifier", "true"),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "properties.object_prop.myObjectIdentifier", "{\"foo\":\"bar\"}"),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "properties.array_prop.string_items.myArrayIdentifier.0", "My Array Value"),
				),
			},
		},
	})
}
func TestAccPortEntityWithRelation(t *testing.T) {
	identifier := genID()
	identifier2 := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_prop" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
		relations = {
			"tfRelation" = {
				"title" = "Test Relation"
				"target" = port-labs_blueprint.microservice2.identifier
			}
		}	
	}
	resource "port-labs_blueprint" "microservice2" {
		title = "TF Provider Test BP1"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_prop" = {
				"myStringIdentifier2" =  {
					"title" = "My String Identifier2"
				}
			}
		}
	}

	resource "port-labs_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port-labs_blueprint.microservice.identifier
		properties = {
			"string_prop" = {
				"myStringIdentifier" =  "My String Value"
			}
		}
		relations = {
			"tfRelation" = [port-labs_entity.microservice2.id]
		}
	}
	
	resource "port-labs_entity" "microservice2" {
		title = "TF Provider Test Entity1"
		identifier = "tf-entity-2"
		blueprint = port-labs_blueprint.microservice2.identifier
		properties = {
			"string_prop" = {
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
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "properties.string_prop.myStringIdentifier", "My String Value"),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "relations.tfRelation.0", "tf-entity-2"),
				),
			},
		},
	})
}

func TestAccPortEntityImport(t *testing.T) {
	blueprintIdentifier := genID()
	entityIdentifier := genID()

	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_prop" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
	}
	resource "port-labs_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port-labs_blueprint.microservice.id
		identifier = "%s"
		properties = {
			"string_prop" = {
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
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "properties.string_prop.myStringIdentifier", "My String Value"),
				),
			},
			{
				ResourceName:            "port-labs_entity.microservice",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           fmt.Sprintf("%s:%s", blueprintIdentifier, entityIdentifier),
				ImportStateVerifyIgnore: []string{"identifier"},
			},
		},
	})
}

func TestAccPortEntityUpdateProp(t *testing.T) {

	identifier := genID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_prop" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
	}
	resource "port-labs_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port-labs_blueprint.microservice.id
		properties = {
			"string_prop" = {
				"myStringIdentifier" =  "My String Value"
			}
		}
	}`, identifier)

	var testAccActionConfigUpdate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			"string_prop" = {
				"myStringIdentifier" =  {
					"title" = "My String Identifier"
				}
			}
		}
	}
	resource "port-labs_entity" "microservice" {
		title = "TF Provider Test Entity0"
		blueprint = port-labs_blueprint.microservice.id
		properties = {
			"string_prop" = {
				"myStringIdentifier" =  "My String Value2"
			}
		}
	}`, identifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "properties.string_prop.myStringIdentifier", "My String Value"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "title", "TF Provider Test Entity0"),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port-labs_entity.microservice", "properties.string_prop.myStringIdentifier", "My String Value2"),
				),
			},
		},
	})
}
