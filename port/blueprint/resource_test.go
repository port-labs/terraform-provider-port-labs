package blueprint_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/version"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
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
			string_props = {
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
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.myStringIdentifier.description", "This is a string property"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.myStringIdentifier.title", "text"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.myStringIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.myStringIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.myStringIdentifier.min_length", "1"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.myStringIdentifier.max_length", "10"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.myStringIdentifier.default", "default"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.myStringIdentifier.format", "user"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.myStringIdentifier.enum.0", "default"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.myStringIdentifier.enum.1", "default2"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.myStringIdentifier.pattern", "^[a-zA-Z0-9]*$"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.myStringIdentifier.enum_colors.default", "red"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.myStringIdentifier.enum_colors.default2", "green"),
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
			number_props = {
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
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_props.myNumberIdentifier.description", "This is a number property"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_props.myNumberIdentifier.title", "number"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_props.myNumberIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_props.myNumberIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_props.myNumberIdentifier.minimum", "1"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_props.myNumberIdentifier.maximum", "10"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_props.myNumberIdentifier.default", "3"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_props.myNumberIdentifier.enum.0", "1"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_props.myNumberIdentifier.enum.1", "2"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_props.myNumberIdentifier.enum.2", "3"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_props.myNumberIdentifier.enum_colors.1", "red"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_props.myNumberIdentifier.enum_colors.2", "green"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.number_props.myNumberIdentifier.enum_colors.3", "blue"),
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
			boolean_props = {
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
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.boolean_props.myBooleanIdentifier.description", "This is a boolean property"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.boolean_props.myBooleanIdentifier.title", "boolean"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.boolean_props.myBooleanIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.boolean_props.myBooleanIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.boolean_props.myBooleanIdentifier.default", "true"),
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
			array_props = {
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
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myStringArrayIdentifier.description", "This is a string array property"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myStringArrayIdentifier.title", "array"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myStringArrayIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myStringArrayIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myStringArrayIdentifier.min_items", "1"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myStringArrayIdentifier.max_items", "10"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myStringArrayIdentifier.string_items.default.0", "a"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myStringArrayIdentifier.string_items.default.1", "b"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myStringArrayIdentifier.string_items.default.2", "c"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myNumberArrayIdentifier.number_items.default.0", "1"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myNumberArrayIdentifier.number_items.default.1", "2"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myNumberArrayIdentifier.number_items.default.2", "3"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myBooleanArrayIdentifier.boolean_items.default.0", "false"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myBooleanArrayIdentifier.boolean_items.default.1", "true"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myObjectArrayIdentifier.object_items.default.0", "{\"a\":\"b\"}"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.array_props.myObjectArrayIdentifier.object_items.default.1", "{\"c\":\"d\"}"),
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
			object_props = {
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
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.object_props.myObjectIdentifier.description", "This is an object property"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.object_props.myObjectIdentifier.title", "object"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.object_props.myObjectIdentifier.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.object_props.myObjectIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.object_props.myObjectIdentifier.default", "{\"key\":\"value\"}"),
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
			string_props = {
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
			string_props = {
			   "text" = {
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
					resource.TestCheckNoResourceAttr("port_blueprint.microservice1", "description"),
					resource.TestCheckResourceAttr("port_blueprint.microservice1", "identifier", identifier1),
					resource.TestCheckResourceAttr("port_blueprint.microservice1", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "title", "TF Provider Test BP3"),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "identifier", identifier2),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "relations.test-rel.title", "Test Relation"),
					resource.TestCheckNoResourceAttr("port_blueprint.microservice2", "description"),
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
			string_props = {
			   "text" = {
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
			string_props = {
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
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.text.spec", "embedded-url"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.text.spec_authentication.authorization_url", "https://getport.io"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.text.spec_authentication.client_id", "123"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "properties.string_props.text.spec_authentication.token_url", "https://getport.io"),
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
			string_props = {
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
			string_props = {
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
			string_props = {
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
			string_props = {
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
			string_props = {
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
			string_props = {
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
				description = "Test Relation"
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
					resource.TestCheckResourceAttr("port_blueprint.microservice2", "relations.test-rel.description", "Test Relation"),
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
			string_props = {
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
				icon = "Terraform"
				colors = {
					"test1" = "red"
					"test2" = "blue"
					"test3" = "green"
					"test4" = "yellow"
				}
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
					resource.TestCheckResourceAttr("port_blueprint.microservice1", "calculation_properties.calculation-for-microservice1.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_blueprint.microservice1", "calculation_properties.calculation-for-microservice1.colors.test2", "blue"),
				),
			},
		},
	})
}

func TestAccPortUpdateBlueprintIdentifier(t *testing.T) {
	identifier := utils.GenID()
	updatedIdentifier := utils.GenID()
	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test"
		icon = "Terraform"
		identifier = "%s"
		description = ""
	}
`, identifier)

	var testAccActionConfigUpdate = fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF Provider Test"
		icon = "Terraform"
		identifier = "%s"
		description = ""
	}
`, updatedIdentifier)

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
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "identifier", updatedIdentifier),
					resource.TestCheckResourceAttr("port_blueprint.microservice", "icon", "Terraform"),
				),
			},
		},
	})
}

func TestAccPortDestroyDeleteAllEntities(t *testing.T) {
	identifier := utils.GenID()
	title := "Blueprint with entities1"
	icon := "Terraform"
	var testAccBlueprintConfigImport = fmt.Sprintf(`
		resource "port_blueprint" "microservice" {
			identifier = "%s"
			icon = "%s"
			title = "%s"
		}
	`, identifier, icon, title)

	var testAccBlueprintConfigForceDeleteEntitiesTrue = fmt.Sprintf(`
		resource "port_blueprint" "microservice" {
			identifier = "%s"
			icon = "%s"
			title = "%s"
			force_delete_entities = true
		}
	`, identifier, icon, title)

	portClient, ctx, err := initializePortTestClient(t)
	if err != nil {
		t.Fatalf("Failed to initialize port client: %s", err.Error())
		return
	}

	blueprint := &cli.Blueprint{
		Identifier: identifier,
		Icon:       &icon,
		Title:      title,
		Schema: cli.BlueprintSchema{
			Properties: map[string]cli.BlueprintProperty{},
		},
		CalculationProperties: map[string]cli.BlueprintCalculationProperty{},
		AggregationProperties: map[string]cli.BlueprintAggregationProperty{},
		MirrorProperties:      map[string]cli.BlueprintMirrorProperty{},
		Relations:             map[string]cli.Relation{},
	}

	_, err = portClient.CreateBlueprint(ctx, blueprint, nil)

	if err != nil {
		t.Fatalf("Failed to create blueprint: %s", err.Error())
		return
	}

	entity := &cli.Entity{
		Blueprint:  identifier,
		Properties: map[string]interface{}{},
		Relations:  map[string]any{},
	}

	// create entity
	_, err = portClient.CreateEntity(ctx, entity, "")
	if err != nil {
		t.Fatalf("Failed to create entity: %s", err.Error())
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             acctest.ProviderConfig + testAccBlueprintConfigImport,
				ResourceName:       "port_blueprint.microservice",
				ImportStateId:      identifier,
				ImportState:        true,
				ImportStatePersist: true,
			},
			{
				Config:      acctest.ProviderConfig + testAccBlueprintConfigImport,
				Destroy:     true,
				ExpectError: regexp.MustCompile(".* has dependant entities*"),
			},
			{
				Config: acctest.ProviderConfig + testAccBlueprintConfigForceDeleteEntitiesTrue,
				Check:  resource.TestCheckResourceAttr("port_blueprint.microservice", "force_delete_entities", "true"),
			},
			{
				Config:  acctest.ProviderConfig + testAccBlueprintConfigForceDeleteEntitiesTrue,
				Destroy: true,
			},
		},
	})
}

func TestAccPortBlueprintOwnership(t *testing.T) {
	var testAccConfigDirect = `
	resource "port_blueprint" "parent_service" {
		title = "Parent Blueprint"
		icon = "Terraform"
		identifier = "parent-service"
		properties = {
			string_props = {
				team = {
					type = "string"
					format = "team"
					title = "Team"
				}
			}
		}
		ownership = {
			type = "Direct"
		}
	}
`

	var testAccConfigInherited = `
	resource "port_blueprint" "parent_service" {
		title = "Parent Blueprint"
		icon = "Terraform"
		identifier = "parent-service"
		properties = {
			string_props = {
				team = {
					type = "string"
					format = "team"
					title = "Team"
				}
			}
		}
		ownership = {
			type = "Direct"
		}
	}

	resource "port_blueprint" "child_service" {
		title = "Child Blueprint"
		icon = "Terraform"
		identifier = "child-service"
		properties = {
			string_props = {
				team = {
					type = "string"
					format = "team"
					title = "Team"
				}
			}
		}
		relations = {
			parent = {
				target = port_blueprint.parent_service.identifier
				title = "Parent Service"
				required = false
				many = false
			}
		}
		ownership = {
			type = "Inherited"
			path = "$relations.parent"
		}
	}
`

	var testAccConfigInvalidPath = `
	resource "port_blueprint" "invalid_path_service" {
		title = "Invalid Path Blueprint"
		icon = "Terraform"
		identifier = "invalid-path-service"
		ownership = {
			type = "Inherited"
			path = "invalid_path_format"
		}
	}
`

	var testAccConfigMissingPath = `
	resource "port_blueprint" "missing_path_service" {
		title = "Missing Path Blueprint"
		icon = "Terraform"
		identifier = "missing-path-service"
		ownership = {
			type = "Inherited"
		}
	}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test Direct ownership
			{
				Config: acctest.ProviderConfig + testAccConfigDirect,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.parent_service", "title", "Parent Blueprint"),
					resource.TestCheckResourceAttr("port_blueprint.parent_service", "identifier", "parent-service"),
					resource.TestCheckResourceAttr("port_blueprint.parent_service", "ownership.type", "Direct"),
					resource.TestCheckNoResourceAttr("port_blueprint.parent_service", "ownership.path"),
				),
			},
			// Test Inherited ownership
			{
				Config: acctest.ProviderConfig + testAccConfigInherited,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_blueprint.child_service", "title", "Child Blueprint"),
					resource.TestCheckResourceAttr("port_blueprint.child_service", "identifier", "child-service"),
					resource.TestCheckResourceAttr("port_blueprint.child_service", "ownership.type", "Inherited"),
					resource.TestCheckResourceAttr("port_blueprint.child_service", "ownership.path", "$relations.parent"),
				),
			},
			// Test invalid path format
			{
				Config:      acctest.ProviderConfig + testAccConfigInvalidPath,
				ExpectError: regexp.MustCompile(`path must be a valid relation identifier starting with '\$relations\.' followed by the relation name`),
			},
			// Test missing path for inherited ownership
			{
				Config:      acctest.ProviderConfig + testAccConfigMissingPath,
				ExpectError: regexp.MustCompile(`path is required when type is 'Inherited'`),
			},
		},
	})
}

func TestAccPortBlueprintCatalogPageCreation(t *testing.T) {
	testCases := []struct {
		name               string
		createCatalogPage  bool
		expectedPageStatus int
	}{
		{
			name:               "CatalogPageCreationTrue",
			createCatalogPage:  true,
			expectedPageStatus: http.StatusOK,
		},
		{
			name:               "CatalogPageCreationFalse",
			createCatalogPage:  false,
			expectedPageStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// add `s` to handle the plural default page that is being created by port
			identifier := fmt.Sprintf("test-%ss", utils.GenID()[:10])
			title := "Microservices"
			icon := "Terraform"
			testAccBlueprintConfig := fmt.Sprintf(`
                resource "port_blueprint" "microservices" {
                    identifier = "%s"
                    icon = "%s"
                    title = "%s"
                    create_catalog_page = %t
                }
            `, identifier, icon, title, tc.createCatalogPage)

			portClient, ctx, err := initializePortTestClient(t)
			if err != nil {
				t.Fatalf("Failed to initialize port client: %s", err.Error())
				return
			}

			blueprint := &cli.Blueprint{
				Identifier: identifier,
				Icon:       &icon,
				Title:      title,
				Schema: cli.BlueprintSchema{
					Properties: map[string]cli.BlueprintProperty{},
				},
				CalculationProperties: map[string]cli.BlueprintCalculationProperty{},
				AggregationProperties: map[string]cli.BlueprintAggregationProperty{},
				MirrorProperties:      map[string]cli.BlueprintMirrorProperty{},
				Relations:             map[string]cli.Relation{},
			}

			_, err = portClient.CreateBlueprint(ctx, blueprint, &tc.createCatalogPage)

			if err != nil {
				t.Fatalf("Failed to create blueprint: %s", err.Error())
				return
			}

			// give grace time for page creation
			time.Sleep(10 * time.Second)

			_, statusCode, err := portClient.GetPage(ctx, identifier)
			if err != nil {
				if statusCode != tc.expectedPageStatus {
					t.Fatalf("Unexpected status code: got %v want %v", statusCode, tc.expectedPageStatus)
				}
			}

			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { acctest.TestAccPreCheck(t) },
				ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:             acctest.ProviderConfig + testAccBlueprintConfig,
						ResourceName:       "port_blueprint.microservices",
						ImportStateId:      identifier,
						ImportState:        true,
						ImportStatePersist: true,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("port_blueprint.microservice", "create_catalog_page", strconv.FormatBool(tc.createCatalogPage)),
						),
					},
				},
			})
		})
	}
}

func initializePortTestClient(t *testing.T) (*cli.PortClient, context.Context, error) {
	baseUrl := os.Getenv("PORT_BASE_URL")
	clientId := os.Getenv("PORT_CLIENT_ID")
	clientSecret := os.Getenv("PORT_CLIENT_SECRET")

	if baseUrl == "" {
		baseUrl = consts.DefaultBaseUrl
	}
	c, err := cli.New(baseUrl, cli.WithHeader("User-Agent", version.ProviderVersion))
	if err != nil {
		t.Fatalf("Failed to create Port-labs client: %s", err.Error())
	}
	ctx := context.Background()
	_, err = c.Authenticate(ctx, clientId, clientSecret)
	if err != nil {
		t.Fatalf("Failed to authenticate with Port-labs: %s", err.Error())
		return nil, ctx, err
	}
	return c, ctx, nil
}
