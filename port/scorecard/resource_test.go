package scorecard_test

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
				"author" = {
					type = "string"
					title = "text"
				}
				"url" = {
					type = "string"
					title = "text"
				}
			}
		}
	}
	`, identifier)
}

func TestAccPortScorecardBasic(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	scorecardIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_scorecard" "test" {
		identifier = "%s"
		title      = "Scorecard 1"
		blueprint  = "%s"
		rules = [{
		  identifier = "hasTeam"
		  title      = "Has Team"
		  level      = "Gold" 
		  query = {
			combinator = "and"
			conditions = [{
			  property = "$team"
			  operator = "isNotEmpty"
			}]
		  }
		}]

		depends_on = [
		port_blueprint.microservice
		]		
	  }`, scorecardIdentifier, blueprintIdentifier)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard.test", "title", "Scorecard 1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.identifier", "hasTeam"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.title", "Has Team"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.level", "Gold"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.combinator", "and"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.property", "$team"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.operator", "isNotEmpty"),
				),
			},
		},
	})
}

func TestAccPortScorecard(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	scorecardIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_scorecard" "test" {
		identifier = "%s"
		title      = "Scorecard 1"
		blueprint  = "%s"
		rules = [{
		  identifier = "test1"
		  title      = "Test1"
		  level      = "Gold"
		  query = {
			combinator = "and"
			conditions = [{
			  property = "$team"
			  operator = "isNotEmpty"
			  },
			  {
				property = "author",
				"operator" : "=",
				"value" : "myValue"
			}]
		  }
		  },
		  {
			identifier = "test2"
			title      = "Test2"
			level      = "Silver"
			query = {
			  combinator = "and"
			  conditions = [{
				property = "url"
				operator = "isNotEmpty"
			  }]
			}
		}]
	  
		depends_on = [
		  port_blueprint.microservice
		]
	  }`, scorecardIdentifier, blueprintIdentifier)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard.test", "title", "Scorecard 1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.#", "2"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.identifier", "test1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.title", "Test1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.level", "Gold"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.combinator", "and"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.#", "2"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.property", "$team"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.operator", "isNotEmpty"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.1.property", "author"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.1.operator", "="),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.1.value", "myValue"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.1.identifier", "test2"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.1.title", "Test2"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.1.level", "Silver"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.1.query.combinator", "and"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.1.query.conditions.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.1.query.conditions.0.property", "url"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.1.query.conditions.0.operator", "isNotEmpty"),
				),
			},
		},
	})
}

func TestAccPortScorecardUpdate(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	scorecardIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_scorecard" "test" {
		identifier = "%s"
		title      = "Scorecard 1"
		blueprint  = "%s"
		rules = [{
		  identifier = "hasTeam"
		  title      = "Has Team"
		  level      = "Gold" 
		  query = {
			combinator = "and"
			conditions = [{
			  property = "$team"
			  operator = "isNotEmpty"
			}]
		  }
		}]

		depends_on = [
		port_blueprint.microservice
		]		
	  }`, scorecardIdentifier, blueprintIdentifier)

	var testAccActionConfigUpdate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_scorecard" "test" {
		identifier = "%s"
		title      = "Scorecard 2"
		blueprint  = "%s"
		rules = [{
					identifier = "hasTeam"
					title      = "Has Team"
					level      = "Bronze"
					query = {
						combinator = "or"
						conditions = [{
							property = "$team"
							operator = "isNotEmpty"
						}]
					}
				}]
		depends_on = [
		port_blueprint.microservice
		]
	 }`, scorecardIdentifier, blueprintIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard.test", "title", "Scorecard 1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.identifier", "hasTeam"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.title", "Has Team"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.level", "Gold"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.combinator", "and"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.property", "$team"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.operator", "isNotEmpty"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard.test", "title", "Scorecard 2"),
					resource.TestCheckResourceAttr("port_scorecard.test", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.identifier", "hasTeam"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.title", "Has Team"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.level", "Bronze"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.combinator", "or"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.property", "$team"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.operator", "isNotEmpty"),
				),
			},
		},
	})
}

func TestAccPortScorecardImport(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	scorecardIdentifier := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_scorecard" "test" {
		identifier = "%s"
		title      = "Scorecard 1"
		blueprint  = "%s"
		rules = [{
		  identifier = "hasTeam"
		  title      = "Has Team"
		  level      = "Gold" 
		  query = {
			combinator = "and"
			conditions = [{
			  property = "$team"
			  operator = "isNotEmpty"
			}]
		  }
		}]

		depends_on = [
		port_blueprint.microservice
		]		
	  }`, scorecardIdentifier, blueprintIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard.test", "title", "Scorecard 1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.identifier", "hasTeam"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.title", "Has Team"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.level", "Gold"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.combinator", "and"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.property", "$team"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.operator", "isNotEmpty"),
				),
			},
			{
				ResourceName:      "port_scorecard.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s:%s", blueprintIdentifier, scorecardIdentifier),
			},
		},
	})
}

func TestAccPortScorecardUpdateIdentifier(t *testing.T) {
	blueprintIdentifier := utils.GenID()
	scorecardIdentifier := utils.GenID()
	scorecardIdentifierUpdated := utils.GenID()
	var testAccActionConfigCreate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_scorecard" "test" {
		identifier = "%s"
		title      = "Scorecard 1"
		blueprint  = "%s"
		rules = [{
		  identifier = "hasTeam"
		  title      = "Has Team"
		  level      = "Gold" 
		  query = {
			combinator = "and"
			conditions = [{
			  property = "$team"
			  operator = "isNotEmpty"
			}]
		  }
		}]

		depends_on = [
		port_blueprint.microservice
		]		
	  }`, scorecardIdentifier, blueprintIdentifier)

	var testAccActionConfigUpdate = testAccCreateBlueprintConfig(blueprintIdentifier) + fmt.Sprintf(`
	resource "port_scorecard" "test" {
		identifier = "%s"
		title      = "Scorecard 2"
		blueprint  = "%s"
		rules = [{
			identifier = "hasTeam"
			title      = "Has Team"
			level      = "Gold" 
			query = {
			  combinator = "and"
			  conditions = [{
				property = "$team"
				operator = "isNotEmpty"
			  }]
			}
		  }]
		depends_on = [
		port_blueprint.microservice
		]
	 }`, scorecardIdentifierUpdated, blueprintIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard.test", "identifier", scorecardIdentifier),
					resource.TestCheckResourceAttr("port_scorecard.test", "title", "Scorecard 1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.identifier", "hasTeam"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.title", "Has Team"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.level", "Gold"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.combinator", "and"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.property", "$team"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.operator", "isNotEmpty"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_scorecard.test", "identifier", scorecardIdentifierUpdated),
					resource.TestCheckResourceAttr("port_scorecard.test", "title", "Scorecard 1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "blueprint", blueprintIdentifier),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.identifier", "hasTeam"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.title", "Has Team"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.level", "Gold"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.combinator", "and"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.#", "1"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.property", "$team"),
					resource.TestCheckResourceAttr("port_scorecard.test", "rules.0.query.conditions.0.operator", "isNotEmpty"),
				),
			},
		},
	})
}
