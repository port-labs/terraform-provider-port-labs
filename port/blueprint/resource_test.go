package blueprint_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func TestAccPortBlueprintBasic(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test"
		icon = "Terraform"
		identifier = "%s"
		description = ""
	}
`, identifier)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckNoResourceAttr("port-labs_blueprint.microservice", "properties"),
				),
			},
		},
	})
}

func TestAccPortBlueprintStringProperty(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_prop = {
				myStringIdentifier = {
					description = "This is a string property"
					title = "text"
					icon = "Terraform"
					required = true
					min_length = 1
					max_length = 10
					default = "default"
					enum = ["default", "default2"]
					pattern = "^[a-zA-Z0-9]*$"
					format = "user"
					enum_colors = {
						default = "red"
						default2 = "green"
					}
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
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.myStringIdentifier.description", "This is a string property"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.myStringIdentifier.title", "text"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.myStringIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.myStringIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.myStringIdentifier.min_length", "1"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.myStringIdentifier.max_length", "10"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.myStringIdentifier.default", "default"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.myStringIdentifier.format", "user"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.myStringIdentifier.enum.0", "default"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.myStringIdentifier.enum.1", "default2"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.myStringIdentifier.pattern", "^[a-zA-Z0-9]*$"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.myStringIdentifier.enum_colors.default", "red"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.myStringIdentifier.enum_colors.default2", "green"),
				),
			},
		},
	})
}

func TestAccPortBlueprintNumberProperty(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			number_prop = {
				myNumberIdentifier = {
					description = "This is a number property"
					title = "number"
					icon = "Terraform"
					required = true
					minimum = 1
					maximum = 10
					default = 3
					enum = [1, 2, 3]
					enum_colors = {
						1 = "red"
						2 = "green"
						3 = "blue"
					}
				}
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
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.number_prop.myNumberIdentifier.description", "This is a number property"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.number_prop.myNumberIdentifier.title", "number"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.number_prop.myNumberIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.number_prop.myNumberIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.number_prop.myNumberIdentifier.minimum", "1"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.number_prop.myNumberIdentifier.maximum", "10"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.number_prop.myNumberIdentifier.default", "3"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.number_prop.myNumberIdentifier.enum.0", "1"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.number_prop.myNumberIdentifier.enum.1", "2"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.number_prop.myNumberIdentifier.enum.2", "3"),
					// resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.number_prop.myNumberIdentifier.enum_colors.1", "red"),
					// resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.number_prop.myNumberIdentifier.enum_colors.2", "green"),
					// resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.number_prop.myNumberIdentifier.enum_colors.3", "blue"),
				),
			},
		},
	})
}

func TestAccPortBlueprintBooleanProperty(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			boolean_prop = {
				myBooleanIdentifier = {
					description = "This is a boolean property"
					title = "boolean"
					icon = "Terraform"
					required = true
					default = true
				}
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
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.boolean_prop.myBooleanIdentifier.description", "This is a boolean property"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.boolean_prop.myBooleanIdentifier.title", "boolean"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.boolean_prop.myBooleanIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.boolean_prop.myBooleanIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.boolean_prop.myBooleanIdentifier.default", "true"),
				),
			},
		},
	})
}

func TestAccPortBlueprintArrayProperty(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			array_prop = {
				myArrayIdentifier = {
					description = "This is an array property"
					title = "array"
					icon = "Terraform"
					required = true
					min_items = 1
					max_items = 10
					string_items = {
						default = ["a", "b", "c"]
					}
				}
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
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.array_prop.myArrayIdentifier.description", "This is an array property"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.array_prop.myArrayIdentifier.title", "array"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.array_prop.myArrayIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.array_prop.myArrayIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.array_prop.myArrayIdentifier.min_items", "1"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.array_prop.myArrayIdentifier.max_items", "10"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.array_prop.myArrayIdentifier.string_items.default.0", "a"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.array_prop.myArrayIdentifier.string_items.default.1", "b"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.array_prop.myArrayIdentifier.string_items.default.2", "c"),
				),
			},
		},
	})
}

func TestAccPortBlueprintObjectProperty(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			object_prop = {
				myObjectIdentifier = {
					description = "This is an object property"
					title = "object"
					icon = "Terraform"
					required = true
					default = jsonencode({
						"key": "value"
					})
				}
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
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.object_prop.myObjectIdentifier.description", "This is an object property"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.object_prop.myObjectIdentifier.title", "object"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.object_prop.myObjectIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.object_prop.myObjectIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.object_prop.myObjectIdentifier.default", "{\"key\":\"value\"}"),
				),
			},
		},
	})
}

func TestAccPortBlueprintWithChangelogDestination(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test"
		icon = "Terraform"
		identifier = "%s"
		changelog_destination = {
						type = "WEBHOOK"
						url = "https://google.com"
						agent = true
					}
	}`, identifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "changelog_destination.type", "WEBHOOK"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "changelog_destination.url", "https://google.com"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "changelog_destination.agent", "true"),
				),
			},
		},
	})
}

func TestAccPortBlueprintWithRelation(t *testing.T) {
	identifier1 := utils.GenID()
	identifier2 := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice1" {
		title = "TF Provider Test BP2"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_prop = {
				"text" = {
					type = "string"
					title = "text"
				}
			}
		}
	}

	resource "port-labs_blueprint" "microservice2" {
		title = "TF Provider Test BP3"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_prop = {
			   "text" = {
					type = "string"
					title = "text"
				}
			}
		}
		relations = {
			"test-rel" = {
				title = "Test Relation"
				target = port-labs_blueprint.microservice1.identifier	
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
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "title", "TF Provider Test BP2"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "identifier", identifier1),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice2", "title", "TF Provider Test BP3"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice2", "identifier", identifier2),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice2", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice2", "relations.test-rel.title", "Test Relation"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice2", "relations.test-rel.target", identifier1),
				),
			},
		},
	})
}

// func TestAccPortBlueprintImport(t *testing.T) {
// 	var testAccActionConfigCreate = `
// 	resource "port-labs_blueprint" "microservice" {
// 		title      = "microservice"
// 		icon       = "Terraform"
// 		identifier = "import_microservice"
// 		properties {
// 			identifier = "bool"
// 			type       = "boolean"
// 			title      = "boolean"
// 		}
// 	}
// `
// 	resource.Test(t, resource.TestCase{
// 		Providers: map[string]*schema.Provider{
// 			"port-labs": Provider(),
// 		},
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccActionConfigCreate,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "title", "microservice"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "identifier", "import_microservice"),
// 				),
// 			},
// 			{
// 				ResourceName:            "port-labs_blueprint.microservice",
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"data_source"},
// 			},
// 		},
// 	})
// }

// func TestAccPortBlueprintWithDefaultValue(t *testing.T) {
// 	identifier := utils.GenID()
// 	var testAccActionConfigCreate = fmt.Sprintf(`
// 	resource "port-labs_blueprint" "microservice" {
// 		title = "TF Provider Test BP0"
// 		icon = "Terraform"
// 		identifier = "%s"
// 		properties {
// 			identifier = "bool"
// 			type = "boolean"
// 			title = "boolean"
// 			default_value = {"value": true}
// 		}
// 		properties {
// 			identifier = "number"
// 			type = "number"
// 			title = "number"
// 			default_value = {"value": 1}
// 		}
// 		properties {
// 			identifier = "obj"
// 			type = "object"
// 			title = "object"
// 			default_value = {"value": jsonencode({"a":"b"})}
// 		}
// 		properties {
// 			identifier = "array"
// 			type = "array"
// 			items = {
// 				type = "string"
// 				format = "url"
// 			}
// 			title = "array"
// 			default_items = ["https://getport.io", "https://app.getport.io"]
// 		}
// 		properties {
// 			identifier = "text"
// 			type = "string"
// 			title = "text"
// 			icon = "Terraform"
// 			enum = ["a", "b", "c"]
// 			enum_colors = {
// 				a = "red"
// 				b = "blue"
// 			}
// 			default_value = {"value":"a"}
// 		}
// 	}
// `, identifier)
// 	resource.Test(t, resource.TestCase{
// 		Providers: map[string]*schema.Provider{
// 			"port-labs": Provider(),
// 		},
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccActionConfigCreate,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.0.default_items.0", "https://getport.io"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.0.default_items.#", "2"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.0.items.type", "string"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.0.items.format", "url"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.1.default_value.value", "1"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.2.identifier", "text"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.2.enum.0", "a"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.2.enum_colors.a", "red"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.2.default_value.value", "a"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.3.default_value.value", "true"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.4.default_value.value", "{\"a\":\"b\"}"),
// 				),
// 			},
// 		},
// 	})

// }

func TestAccPortBlueprintWithSpecification(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice" {
		title = "TF Provider Test BP0"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_prop = {
				"text" = {
					title = "text"
					type = "string"
					format = "url"
					spec = "embedded-url"
					spec_authentication = {
						token_url = "https://getport.io"
						client_id = "123"
						authorization_url = "https://getport.io"
					}
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
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.text.spec", "embedded-url"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.text.spec_authentication.authorization_url", "https://getport.io"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.text.spec_authentication.client_id", "123"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice", "properties.string_prop.text.spec_authentication.token_url", "https://getport.io"),
				),
			},
		},
	})
}

// func TestAccPortBlueprintUpdate(t *testing.T) {
// 	identifier := utils.GenID()
// 	var testAccActionConfigCreate = fmt.Sprintf(`
// 	resource "port-labs_blueprint" "microservice1" {
// 		description = "Test Description"
// 		title = "TF Provider Test BP2"
// 		icon = "Terraform"
// 		identifier = "%s"
// 		properties {
// 			required = true
// 			identifier = "text"
// 			type = "string"
// 			icon = "Terraform"
// 			title = "text"
// 			enum = ["a", "b", "c"]
// 			enum_colors = {
// 				a = "red"
// 				b = "blue"
// 			}
// 		}
// 		calculation_properties {
// 			identifier = "calc"
// 			type = "number"
// 			icon = "Terraform"
// 			title = "calc"
// 			calculation = "2"
// 			colorized = true
// 			colors = {
// 				0 = "red"
// 				1 = "blue"
// 			}
// 		}
// 	}
// `, identifier)
// 	var testAccActionConfigUpdate = fmt.Sprintf(`
// 	resource "port-labs_blueprint" "microservice1" {
// 		title = "TF Provider Test BP2"
// 		icon = "Terraform"
// 		identifier = "%s"
// 		properties {
// 			required = false
// 			identifier = "text"
// 			type = "string"
// 			title = "text"
// 		}
// 		properties {
// 			identifier = "number"
// 			type = "number"
// 			title = "num"
// 		}
// 	}
// `, identifier)
// 	var testAccActionConfigUpdateAgain = fmt.Sprintf(`
// 	resource "port-labs_blueprint" "microservice1" {
// 		title = "TF Provider Test BP2"
// 		icon = "Terraform"
// 		identifier = "%s"
// 		properties {
// 			identifier = "number"
// 			type = "number"
// 			title = "num"
// 		}
// 	}
// `, identifier)
// 	resource.Test(t, resource.TestCase{
// 		Providers: map[string]*schema.Provider{
// 			"port-labs": Provider(),
// 		},
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccActionConfigCreate,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "description", "Test Description"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.0.title", "text"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.0.required", "true"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.0.icon", "Terraform"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "calculation_properties.0.identifier", "calc"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "calculation_properties.0.icon", "Terraform"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "calculation_properties.0.type", "number"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "calculation_properties.0.colorized", "true"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "calculation_properties.0.colors.0", "red"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "calculation_properties.0.colors.1", "blue"),
// 				),
// 			},
// 			{
// 				Config: testAccActionConfigUpdate,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "description", ""),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.0.title", "num"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.1.title", "text"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.1.required", "false"),
// 				),
// 			},
// 			{
// 				Config: testAccActionConfigUpdateAgain,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.0.title", "num"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "properties.#", "1"),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccPortBlueprintUpdateRelation(t *testing.T) {
	envID := utils.GenID()
	vmID := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "Environment" {
		title = "Environment"
		icon = "Environment"
		identifier = "%s"
		properties = {
			string_prop = {
				"env_name" = {
					title = "Name"
				}
			}
		}
	}
	resource "port-labs_blueprint" "vm" {
		title = "Virtual Machine"
		icon = "Azure"
		identifier = "%s"
		properties = {
			string_prop = {
				"image" = {
					title = "Image"
				}
			}
		}
		relations = {
			"vm-to-environment" = {
				title = "Related Environment"
				target = port-labs_blueprint.Environment.identifier
			}
		}
	}
`, envID, vmID)
	var testAccActionConfigUpdate = fmt.Sprintf(`
	resource "port-labs_blueprint" "Environment" {
		title = "Environment"
		icon = "Environment"
		identifier = "%s"
		properties = {
			string_prop = {
				"env_name" = {
					title = "Name"
				}
			}
		}
	}
	resource "port-labs_blueprint" "vm" {
		title = "Virtual Machine"
		icon = "Azure"
		identifier = "%s"
		properties = {
			string_prop = {
				"image" = {
					title = "Image"
				}
			}
		}
		relations = {
			"environment" = {
				title = "Related Environment"
				target = port-labs_blueprint.Environment.identifier
			}
		}
	}
`, envID, vmID)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.vm-to-environment.title", "Related Environment"),
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.vm-to-environment.target", envID),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.environment.title", "Related Environment"),
					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.environment.target", envID),
				),
			},
		},
	})
}

// 	resource.Test(t, resource.TestCase{
// 		Providers: map[string]*schema.Provider{
// 			"port-labs": Provider(),
// 		},
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccActionConfigCreate,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.#", "1"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.0.title", "Related Environment"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.0.target", envID),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.0.identifier", "vm-to-environment"),
// 				),
// 			},
// 			{
// 				Config: testAccActionConfigUpdate,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.#", "1"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.0.title", "Related Environment"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.0.target", envID),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "relations.0.identifier", "environment"),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccPortBlueprintWithMirrorProperty(t *testing.T) {
	identifier1 := utils.GenID()
	identifier2 := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice1" {
		title = "TF Provider Test BP2"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_prop = {
				"text" = {
					title = "text"
				}
			}
		}
	}
	resource "port-labs_blueprint" "microservice2" {
		title = "TF Provider Test BP3"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_prop = {
				"text" = {
					title = "text"
				}
			}
		}
		mirror_properties = {
			"mirror-for-microservice1" = {
				title = "Mirror for microservice1"
				path = "test-rel.$identifier"
			}
		}
		relations = {
			"test-rel" = {
				title = "Test Relation"
				target = port-labs_blueprint.microservice1.identifier
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
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice2", "relations.test-rel.title", "Test Relation"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice2", "relations.test-rel.target", identifier1),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice2", "mirror_properties.mirror-for-microservice1.title", "Mirror for microservice1"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice2", "mirror_properties.mirror-for-microservice1.path", "test-rel.$identifier"),
				),
			},
		},
	})
}

func TestAccPortBlueprintWithCalculationProperty(t *testing.T) {
	identifier1 := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port-labs_blueprint" "microservice1" {
		title = "TF Provider Test BP2"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_prop = {
				"text" = {
					title = "text"
				}
			}
		}
		calculation_properties = {
			"calculation-for-microservice1" = {
				title = "Calculation for microservice1"
				calculation = "test-rel.$identifier"
				type = "string"
			}
		}
	}`, identifier1)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "calculation_properties.calculation-for-microservice1.title", "Calculation for microservice1"),
					resource.TestCheckResourceAttr("port-labs_blueprint.microservice1", "calculation_properties.calculation-for-microservice1.calculation", "test-rel.$identifier"),
				),
			},
		},
	})
}

// 		Providers: map[string]*schema.Provider{
// 			"port-labs": Provider(),
// 		},
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccActionConfigCreate,
// 			},
// 		},
// 	})
// }

// func TestAccPortBlueprintUpdateMirrorProperty(t *testing.T) {
// 	envID := utils.GenID()
// 	vmID := utils.GenID()
// 	var testAccActionConfigCreate = fmt.Sprintf(`
// 	resource "port-labs_blueprint" "Environment" {
// 		title = "Environment"
// 		icon = "Environment"
// 		identifier = "%s"
// 		properties {
// 			identifier = "env_name"
// 			type = "string"
// 			title = "Name"
// 		}
// 	}
// 	resource "port-labs_blueprint" "vm" {
// 		title = "Virtual Machine"
// 		icon = "Azure"
// 		identifier = "%s"
// 		properties {
// 			identifier = "image"
// 			type = "string"
// 			title = "Image"
// 		}
// 		mirror_properties {
// 			identifier = "mirror-for-environment"
// 			title = "Mirror for environment"
// 			path = "vm-to-environment.$identifier"
// 		}
// 		relations {
// 			identifier = "vm-to-environment"
// 			title = "Related Environment"
// 			target = port-labs_blueprint.Environment.identifier
// 		}
// 	}
// `, envID, vmID)
// 	var testAccActionConfigUpdate = fmt.Sprintf(`
// 	resource "port-labs_blueprint" "Environment" {
// 		title = "Environment"
// 		icon = "Environment"
// 		identifier = "%s"
// 		properties {
// 			identifier = "env_name"
// 			type = "string"
// 			title = "Name"
// 		}
// 	}
// 	resource "port-labs_blueprint" "vm" {
// 		title = "Virtual Machine"
// 		icon = "Azure"
// 		identifier = "%s"
// 		properties {
// 			identifier = "image"
// 			type = "string"
// 			title = "Image"
// 		}
// 		mirror_properties {
// 			identifier = "mirror-for-environment"
// 			title = "Mirror for environment2"
// 			path = "environment.$identifier"
// 		}
// 		relations {
// 			identifier = "environment"
// 			title = "Related Environment"
// 			target = port-labs_blueprint.Environment.identifier
// 		}
// 	}
// `, envID, vmID)
// 	resource.Test(t, resource.TestCase{
// 		Providers: map[string]*schema.Provider{
// 			"port-labs": Provider(),
// 		},
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccActionConfigCreate,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "mirror_properties.#", "1"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "mirror_properties.0.title", "Mirror for environment"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "mirror_properties.0.identifier", "mirror-for-environment"),
// 				),
// 			},
// 			{
// 				Config: testAccActionConfigUpdate,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "mirror_properties.#", "1"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "mirror_properties.0.title", "Mirror for environment2"),
// 					resource.TestCheckResourceAttr("port-labs_blueprint.vm", "mirror_properties.0.identifier", "mirror-for-environment"),
// 				),
// 			},
// 		},
// 	})
// }
