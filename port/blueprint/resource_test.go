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
	resource "port_blueprint" "microservice" {
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
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckNoResourceAttr("port_blueprint.microservice", "properties"),
				),
			},
		},
	})
}

func TestAccPortBlueprintStringProperty(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
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
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.myStringIdentifier.description", "This is a string property"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.myStringIdentifier.title", "text"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.myStringIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.myStringIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.myStringIdentifier.min_length", "1"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.myStringIdentifier.max_length", "10"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.myStringIdentifier.default", "default"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.myStringIdentifier.format", "user"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.myStringIdentifier.enum.0", "default"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.myStringIdentifier.enum.1", "default2"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.myStringIdentifier.pattern", "^[a-zA-Z0-9]*$"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.myStringIdentifier.enum_colors.default", "red"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.myStringIdentifier.enum_colors.default2", "green"),
				),
			},
		},
	})
}

func TestAccPortBlueprintNumberProperty(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
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
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_prop.myNumberIdentifier.description", "This is a number property"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_prop.myNumberIdentifier.title", "number"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_prop.myNumberIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_prop.myNumberIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_prop.myNumberIdentifier.minimum", "1"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_prop.myNumberIdentifier.maximum", "10"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_prop.myNumberIdentifier.default", "3"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_prop.myNumberIdentifier.enum.0", "1"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_prop.myNumberIdentifier.enum.1", "2"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_prop.myNumberIdentifier.enum.2", "3"),
					// resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_prop.myNumberIdentifier.enum_colors.1", "red"),
					// resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_prop.myNumberIdentifier.enum_colors.2", "green"),
					// resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_prop.myNumberIdentifier.enum_colors.3", "blue"),
				),
			},
		},
	})
}

func TestAccPortBlueprintBooleanProperty(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
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
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.boolean_prop.myBooleanIdentifier.description", "This is a boolean property"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.boolean_prop.myBooleanIdentifier.title", "boolean"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.boolean_prop.myBooleanIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.boolean_prop.myBooleanIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.boolean_prop.myBooleanIdentifier.default", "true"),
				),
			},
		},
	})
}

func TestAccPortBlueprintArrayProperty(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			array_prop = {
				myStringArrayIdentifier = {
					description = "This is a string array property"
					title = "array"
					icon = "Terraform"
					required = true
					min_items = 1
					max_items = 10
					string_items = {
						default = ["a", "b", "c"]
					}
				}
				myNumberArrayIdentifier = {
					description = "This is a number array property"
					title = "array"
					number_items = {
						default = [1, 2, 3]
					}
				}
				myBooleanArrayIdentifier = {
					description = "This is a boolean array property"
					title = "array"
					boolean_items = {
						default = [false,true]
					}
				}
				myObjectArrayIdentifier = {
					description = "This is a object array property"
					title = "array"
					object_items = {
						default = [jsonencode({"a": "b"}), jsonencode({"c": "d"})]
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
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myStringArrayIdentifier.description", "This is a string array property"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myStringArrayIdentifier.title", "array"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myStringArrayIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myStringArrayIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myStringArrayIdentifier.min_items", "1"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myStringArrayIdentifier.max_items", "10"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myStringArrayIdentifier.string_items.default.0", "a"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myStringArrayIdentifier.string_items.default.1", "b"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myStringArrayIdentifier.string_items.default.2", "c"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myNumberArrayIdentifier.number_items.default.0", "1"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myNumberArrayIdentifier.number_items.default.1", "2"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myNumberArrayIdentifier.number_items.default.2", "3"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myBooleanArrayIdentifier.boolean_items.default.0", "false"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myBooleanArrayIdentifier.boolean_items.default.1", "true"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myObjectArrayIdentifier.object_items.default.0", "{\"a\":\"b\"}"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_prop.myObjectArrayIdentifier.object_items.default.1", "{\"c\":\"d\"}"),
				),
			},
		},
	})
}

func TestAccPortBlueprintObjectProperty(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
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
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.object_prop.myObjectIdentifier.description", "This is an object property"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.object_prop.myObjectIdentifier.title", "object"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.object_prop.myObjectIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.object_prop.myObjectIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.object_prop.myObjectIdentifier.default", "{\"key\":\"value\"}"),
				),
			},
		},
	})
}

func TestAccPortBlueprintWithChangelogDestination(t *testing.T) {
	identifier := utils.GenID()
	identifier2 := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test"
		icon = "Terraform"
		identifier = "%s"
		webhook_changelog_destination = {
						type = "WEBHOOK"
						url = "https://google.com"
						agent = true
					}
	}
	resource "port_blueprint" "microservice2" {
		title = "TF Provider Test"
		icon = "Terraform"
		identifier = "%s"
		kafka_changelog_destination = {}
	}
	`, identifier, identifier2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", identifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "webhook_changelog_destination.url", "https://google.com"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "webhook_changelog_destination.agent", "true"),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "identifier", identifier2),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "icon", "Terraform"),
				),
			},
		},
	})
}

func TestAccPortBlueprintWithRelation(t *testing.T) {
	identifier1 := utils.GenID()
	identifier2 := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice1" {
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

	resource "port_blueprint" "microservice2" {
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
				target = port_blueprint.microservice1.identifier	
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
					resource.TestCheckResourceAttr("port_blueprint.microservice1", "title", "TF Provider Test BP2"),
					resource.TestCheckResourceAttr("port_blueprint.microservice1", "identifier", identifier1),
					resource.TestCheckResourceAttr("port_blueprint.microservice1", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "title", "TF Provider Test BP3"),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "identifier", identifier2),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "relations.test-rel.title", "Test Relation"),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "relations.test-rel.target", identifier1),
				),
			},
		},
	})
}

func TestAccPortBlueprintImport(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice2" {
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
	}`, identifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "title", "TF Provider Test BP3"),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "identifier", identifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "icon", "Terraform"),
				),
			},
			{
				ResourceName:      "port_blueprint.microservice2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccPortBlueprintWithSpecification(t *testing.T) {
	identifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
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
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.text.spec", "embedded-url"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.text.spec_authentication.authorization_url", "https://getport.io"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.text.spec_authentication.client_id", "123"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_prop.text.spec_authentication.token_url", "https://getport.io"),
				),
			},
		},
	})
}

func TestAccPortBlueprintUpdateRelation(t *testing.T) {
	envID := utils.GenID()
	vmID := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "Environment" {
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
	resource "port_blueprint" "vm" {
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
				target = port_blueprint.Environment.identifier
			}
		}
	}
`, envID, vmID)
	var testAccActionConfigUpdate = fmt.Sprintf(`
	resource "port_blueprint" "Environment" {
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
	resource "port_blueprint" "vm" {
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
				target = port_blueprint.Environment.identifier
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
					resource.TestCheckResourceAttr("port_blueprint.vm", "relations.vm-to-environment.title", "Related Environment"),
					resource.TestCheckResourceAttr("port_blueprint.vm", "relations.vm-to-environment.target", envID),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.vm", "relations.environment.title", "Related Environment"),
					resource.TestCheckResourceAttr("port_blueprint.vm", "relations.environment.target", envID),
				),
			},
		},
	})
}

func TestAccPortBlueprintWithMirrorProperty(t *testing.T) {
	identifier1 := utils.GenID()
	identifier2 := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice1" {
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
	resource "port_blueprint" "microservice2" {
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
				target = port_blueprint.microservice1.identifier
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
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "relations.test-rel.title", "Test Relation"),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "relations.test-rel.target", identifier1),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "mirror_properties.mirror-for-microservice1.title", "Mirror for microservice1"),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "mirror_properties.mirror-for-microservice1.path", "test-rel.$identifier"),
				),
			},
		},
	})
}

func TestAccPortBlueprintWithCalculationProperty(t *testing.T) {
	identifier1 := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice1" {
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
					resource.TestCheckResourceAttr("port_blueprint.microservice1", "calculation_properties.calculation-for-microservice1.title", "Calculation for microservice1"),
					resource.TestCheckResourceAttr("port_blueprint.microservice1", "calculation_properties.calculation-for-microservice1.calculation", "test-rel.$identifier"),
				),
			},
		},
	})
}
