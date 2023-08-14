package action_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func testAccCreateBlueprintConfig(identifier string) string {
	return fmt.Sprintf(`
	resource "port_blueprint" "microservice" {
		title = "TF test microservice"
		icon = "Terraform"
		identifier = "%s"
		properties = {
			string_props = {
			"text" = {
				type = "string"
				title = "text"
				}
			}
		}
	}
	`, identifier)
}
func TestAccPortActionBasic(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		kafka_method = {}
	}`, actionIdentifier)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
				),
			},
		},
	})
}

func TestAccPortAction(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		kafka_method = {}
		user_properties = {
			"string_props" = {
				"myStringIdentifier" = {
					"title" = "My String Identifier"
					"required" = true

				}
			}
			"number_props" = {
				"myNumberIdentifier" = {
					"title" = "My Number Identifier"
					"required" = true
					maximum = 100
					minimum = 0
				}
			}
			"boolean_props" = {
				"myBooleanIdentifier" = {
					"title" = "My Boolean Identifier"
					"required" = true
				}
			}
			"object_props" = {
				"myObjectIdentifier" = {
					"title" = "My Object Identifier"
					"required" = true
				}
			}
			"array_props" = {
				"myArrayIdentifier" = {
					"title" = "My Array Identifier"
					"required" = true
					string_items = {
						format = "email"
					}
				}
			}
		}

	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier.title", "My String Identifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.number_props.myNumberIdentifier.title", "My Number Identifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.number_props.myNumberIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.number_props.myNumberIdentifier.maximum", "100"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.number_props.myNumberIdentifier.minimum", "0"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.boolean_props.myBooleanIdentifier.title", "My Boolean Identifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.boolean_props.myBooleanIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.object_props.myObjectIdentifier.title", "My Object Identifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.object_props.myObjectIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.array_props.myArrayIdentifier.title", "My Array Identifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.array_props.myArrayIdentifier.required", "true"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.array_props.myArrayIdentifier.string_items.format", "email"),
				),
			},
		},
	})
}

func TestAccPortActionWebhookInvocation(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		webhook_method = {
			url = "https://example.com"
			agent = true
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "webhook_method.url", "https://example.com"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "webhook_method.agent", "true"),
				),
			},
		},
	})
}

func TestAccPortActionGitlabInvocation(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		gitlab_method = {
			project_name = "terraform-provider-port"
			group_name = "port"
			omit_payload = true
			omit_user_inputs = true
			default_ref = "main"
			agent = true
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "gitlab_method.project_name", "terraform-provider-port"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "gitlab_method.group_name", "port"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "gitlab_method.omit_payload", "true"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "gitlab_method.omit_user_inputs", "true"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "gitlab_method.default_ref", "main"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "gitlab_method.agent", "true"),
				),
			},
		},
	})
}
func TestAccPortActionAzureInvocation(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		azure_method = {
			org = "port",
			webhook = "https://getport.io"
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "azure_method.org", "port"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "azure_method.webhook", "https://getport.io"),
				),
			},
		},
	})
}

func TestAccPortActionGithubInvocation(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		github_method = {
			org = "port",
			repo = "terraform-provider-port",
			workflow = "main.yml"
			omit_payload = true
			omit_user_inputs = true
			report_workflow_status = false
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "github_method.org", "port"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "github_method.repo", "terraform-provider-port"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "github_method.workflow", "main.yml"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "github_method.omit_payload", "true"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "github_method.omit_user_inputs", "true"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "github_method.report_workflow_status", "false"),
				),
			},
		},
	})
}

func TestAccPortActionImport(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		webhook_method = {
			url = "https://getport.io"
		}
		user_properties = {
			"string_props" = {
				"myStringIdentifier" = {
					"title" = "My String Identifier"
					"required" = true
				}
			}
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "webhook_method.url", "https://getport.io"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier.title", "My String Identifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier.required", "true"),
				),
			},
			{
				ResourceName:      "port_action.create_microservice",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s:%s", blueprintIdentifier, actionIdentifier),
			},
		},
	})
}

func TestAccPortActionUpdate(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		webhook_method = {
			url = "https://getport.io"
		}
		user_properties = {
			"string_props" = {
				"myStringIdentifier" = {
					"title" = "My String Identifier"
					"required" = true
				}
			}
		}
	}`, actionIdentifier)

	var testAccActionConfigUpdate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title = "TF Provider Test"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		webhook_method = {
			url = "https://getport.io"
		}
		user_properties = {
			"string_props" = {
				"myStringIdentifier2" = {
					"title" = "My String Identifier"
					"required" = false
				}
			}
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "webhook_method.url", "https://getport.io"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier.title", "My String Identifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier.required", "true"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "TF Provider Test"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "blueprint", identifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "webhook_method.url", "https://getport.io"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier2.title", "My String Identifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier2.required", "false"),
				),
			},
		},
	})
}

func TestAccPortActionAdvancedFormConfigurations(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "action1" {
		title             = "Action 1"
		blueprint         =  port_blueprint.microservice.id
		identifier        = "%s"
		trigger           = "DAY-2"
		description       = "This is a test action"
		required_approval = true
		github_method = {
		  org      = "port-labs"
		  repo     = "Port"
		  workflow = "lint"
		}
		user_properties = {
		  string_props = {
			myStringIdentifier = {
			  title      = "myStringIdentifier"
			  default    = "default"
			  required   = false
			}
			myStringIdentifier2 = {
			  title      = "myStringIdentifier2"
			  default    = "default"
			  required   = false
			  depends_on = ["myStringIdentifier"]
			}
			myStringIdentifier3 = {
			  title      = "myStringIdentifier3"
			  required   = false
			  dataset = {
				"combinator" : "and",
				"rules" : [
				  {
					"property" : "$team",
					"operator" : "containsAny",
					"value" : {
					  "jq_query" : "Test"
					}
				  }
				]
			  }
			}
		  }
		}
	  }`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.action1", "title", "Action 1"),
					resource.TestCheckResourceAttr("port_action.action1", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.action1", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.action1", "description", "This is a test action"),
					resource.TestCheckResourceAttr("port_action.action1", "required_approval", "true"),
					resource.TestCheckResourceAttr("port_action.action1", "github_method.org", "port-labs"),
					resource.TestCheckResourceAttr("port_action.action1", "github_method.repo", "Port"),
					resource.TestCheckResourceAttr("port_action.action1", "github_method.workflow", "lint"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier.title", "myStringIdentifier"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier.default", "default"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier.required", "false"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier2.title", "myStringIdentifier2"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier2.default", "default"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier2.required", "false"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier2.depends_on.0", "myStringIdentifier"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier3.title", "myStringIdentifier3"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier3.required", "false"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier3.dataset.combinator", "and"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier3.dataset.rules.0.property", "$team"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier3.dataset.rules.0.operator", "containsAny"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier3.dataset.rules.0.value.jq_query", "Test"),
				),
			},
		},
	})
}

func TestAccPortActionJqDefault(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title             = "Action 1"
		blueprint         =  port_blueprint.microservice.id
		identifier        = "%s"
		trigger           = "DAY-2"
		description       = "This is a test action"
		kafka_method = {} 
		user_properties = {
			string_props = {
				myStringIdentifier = {
					title      = "myStringIdentifier"
					default_jq_query = "'Test'"
					required   = false
				}
			}
			number_props = {
				myNumberIdentifier = {
					title      = "myNumberIdentifier"
					default_jq_query = "1"
				}
			}
			boolean_props = {
				myBooleanIdentifier = {
					title      = "myBooleanIdentifier"
					default_jq_query = "true"
				}
			}
			object_props = {
				myObjectIdentifier = {
					title      = "myObjectIdentifier"
					default_jq_query = "{ \"test\": \"test\" }"
				}
			}
			array_props = {
				myArrayIdentifier = {
					title      = "myArrayIdentifier"
					default_jq_query = "[ \"test\" ]"
				}
			}
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "Action 1"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "description", "This is a test action"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier.title", "myStringIdentifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier.default_jq_query", "'Test'"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.number_props.myNumberIdentifier.title", "myNumberIdentifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.number_props.myNumberIdentifier.default_jq_query", "1"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.boolean_props.myBooleanIdentifier.title", "myBooleanIdentifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.boolean_props.myBooleanIdentifier.default_jq_query", "true"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.object_props.myObjectIdentifier.title", "myObjectIdentifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.object_props.myObjectIdentifier.default_jq_query", "{ \"test\": \"test\" }"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.array_props.myArrayIdentifier.title", "myArrayIdentifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.array_props.myArrayIdentifier.default_jq_query", "[ \"test\" ]"),
				),
			},
		},
	})

}

func TestAccPortActionEnumJqQuery(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title             = "Action 1"
		blueprint         =  port_blueprint.microservice.id
		identifier        = "%s"
		trigger           = "DAY-2"
		description       = "This is a test action"
		kafka_method = {}
		user_properties = {
			string_props = {
				myStringIdentifier = {
					title      = "myStringIdentifier"
					enum_jq_query = "[\"test1\", \"test2\"]"
				}
			}
			number_props = {
				myNumberIdentifier = {
					title 	= "myNumberIdentifier"
					enum_jq_query = "[1, 2]"
				}
			}
			array_props = {
				myStringArrayIdentifier = {
					title 	= "myStringArrayIdentifier"
					string_items = {
						enum_jq_query = "'example' | [ . ]"
					}
				}
				myNumberArrayIdentifier = {
					title 	= "myNumberArrayIdentifier"
					number_items = {
						enum_jq_query = "[1, 2]"
					}
				}

			}
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "Action 1"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "description", "This is a test action"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier.title", "myStringIdentifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier.enum_jq_query", "[\"test1\", \"test2\"]"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.number_props.myNumberIdentifier.title", "myNumberIdentifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.number_props.myNumberIdentifier.enum_jq_query", "[1, 2]"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.array_props.myStringArrayIdentifier.title", "myStringArrayIdentifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.array_props.myStringArrayIdentifier.string_items.enum_jq_query", "'example' | [ . ]"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.array_props.myNumberArrayIdentifier.title", "myNumberArrayIdentifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.array_props.myNumberArrayIdentifier.number_items.enum_jq_query", "[1, 2]"),
				),
			},
		},
	})
}

func TestAccPortActionEnum(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "create_microservice" {
		title             = "Action 1"
		blueprint         =  port_blueprint.microservice.id
		identifier        = "%s"
		trigger           = "DAY-2"
		description       = "This is a test action"
		kafka_method = {}
		user_properties = {
			string_props = {
				myStringIdentifier = {
					title      = "myStringIdentifier"
					enum = ["test1", "test2"]
				}
			}
			number_props = {
				myNumberIdentifier = {
					title 	= "myNumberIdentifier"
					enum = [1, 2]
				}
			}
			array_props = {
				myStringArrayIdentifier = {
					title 	= "myStringArrayIdentifier"
					string_items = {
						enum = ["example"]
					}
				}
			}
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.create_microservice", "title", "Action 1"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.create_microservice", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "description", "This is a test action"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier.title", "myStringIdentifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier.enum.0", "test1"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.string_props.myStringIdentifier.enum.1", "test2"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.number_props.myNumberIdentifier.title", "myNumberIdentifier"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.number_props.myNumberIdentifier.enum.0", "1"),
					resource.TestCheckResourceAttr("port_action.create_microservice", "user_properties.number_props.myNumberIdentifier.enum.1", "2"),
				),
			},
		},
	})
}
func TestAccPortActionOrderProperties(t *testing.T) {
	identifier := utils.GenID()
	actionIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(identifier) + fmt.Sprintf(`
	resource "port_action" "action1" {
		title = "Action 1"
		identifier = "%s"
		icon = "Terraform"
		blueprint = port_blueprint.microservice.id
		trigger = "DAY-2"
		kafka_method = {}
		order_properties = ["myStringIdentifier2", "myStringIdentifier1"]
		user_properties = {
			string_props = {
				myStringIdentifier1 = {
					title      = "myStringIdentifier1"
					required   = false
				}
				myStringIdentifier2 = {
					title      = "myStringIdentifier2"
					required   = false
				}
			}
		}
	}`, actionIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_action.action1", "title", "Action 1"),
					resource.TestCheckResourceAttr("port_action.action1", "identifier", actionIdentifier),
					resource.TestCheckResourceAttr("port_action.action1", "trigger", "DAY-2"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier1.title", "myStringIdentifier1"),
					resource.TestCheckResourceAttr("port_action.action1", "user_properties.string_props.myStringIdentifier2.title", "myStringIdentifier2"),
					resource.TestCheckResourceAttr("port_action.action1", "order_properties.0", "myStringIdentifier2"),
					resource.TestCheckResourceAttr("port_action.action1", "order_properties.1", "myStringIdentifier1"),
				),
			},
		},
	})
}
