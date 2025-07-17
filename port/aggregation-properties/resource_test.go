package aggregation_properties_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func baseBlueprintsTemplate(parentBlueprintIdentifier string, childBlueprintIdentifier string) string {
	return fmt.Sprintf(`
	resource "port_blueprint" "parent_blueprint" {
		title = "Parent Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		properties = {
			number_props = {
				"age" = {
					title = "Age"
				}
			}
		}
	}

	resource "port_blueprint" "child_blueprint" {
		title = "Child Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		relations = {
			"parent" = {
				title = "Parent"
				target = port_blueprint.parent_blueprint.identifier
			}
		}
	}
`, parentBlueprintIdentifier, childBlueprintIdentifier)
}

func TestAccPortAggregationPropertyWithCycleRelation(t *testing.T) {
	// Test checks that a cycle aggregation property works.
	// The cycle is created by creating a parent blueprint and a child blueprint.
	// The child blueprint has a relation to the parent blueprint.
	// The parent blueprint has an aggregation property that counts the children of the parent.
	// The aggregation property is created with a cycle relation to the child blueprint, which is allowed.
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	var testAccActionConfigCreate = baseBlueprintsTemplate(parentBlueprintIdentifier, childBlueprintIdentifier) + `
	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.parent_blueprint.identifier
		properties = {
			"count_entities" = {
				target_blueprint_identifier = port_blueprint.child_blueprint.identifier
				title = "Count Childrens"
				icon = "Terraform"
				description = "Count Childrens"
				method = {
					count_entities = true
				}
			}
		}
	}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_entities.title", "Count Childrens"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_entities.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_entities.description", "Count Childrens"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_entities.target_blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_entities.method.count_entities", "true"),
				),
			},
		},
	})
}

func TestAccCreateAggregationPropertyAverageEntities(t *testing.T) {
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	var testAccActionConfigCreate = baseBlueprintsTemplate(parentBlueprintIdentifier, childBlueprintIdentifier) + `
	resource "port_aggregation_properties" "parent_aggregation_properties" {
		blueprint_identifier = port_blueprint.parent_blueprint.identifier
		properties = {
			"count_entities" = {
				target_blueprint_identifier = port_blueprint.child_blueprint.identifier
				title = "Count Childrens"
				icon = "Terraform"
				description = "Count Childrens"
				method = {
					average_entities = {
					"average_of" = "month"
					"measure_time_by" = "$updatedAt"
					}
				}
			}
		}
	}
`

	var testAccActionConfigUpdate = baseBlueprintsTemplate(parentBlueprintIdentifier, childBlueprintIdentifier) + `
	resource "port_aggregation_properties" "parent_aggregation_properties" {
		blueprint_identifier = port_blueprint.parent_blueprint.identifier
		properties = {
			"count_entities" = {
				target_blueprint_identifier = port_blueprint.child_blueprint.identifier
				title = "Count Childrens"
				icon = "Terraform"
				description = "Count Childrens"
				method = {
					average_entities = {}
				}
			}
		}
	}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "properties.count_entities.title", "Count Childrens"),
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "properties.count_entities.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "properties.count_entities.description", "Count Childrens"),
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "properties.count_entities.target_blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "properties.count_entities.method.average_entities.average_of", "month"),
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "properties.count_entities.method.average_entities.measure_time_by", "$updatedAt"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "properties.count_entities.title", "Count Childrens"),
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "properties.count_entities.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "properties.count_entities.description", "Count Childrens"),
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "properties.count_entities.target_blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "properties.count_entities.method.average_entities.average_of", "day"),
					resource.TestCheckResourceAttr("port_aggregation_properties.parent_aggregation_properties", "properties.count_entities.method.average_entities.measure_time_by", "$createdAt"),
				),
			},
		},
	})
}

func TestAccPortCreateAggregationAverageProperties(t *testing.T) {
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	var testAccActionConfigCreate = baseBlueprintsTemplate(parentBlueprintIdentifier, childBlueprintIdentifier) + `
	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		properties = {
			"average_age" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title = "Average Age"
				icon = "Terraform"
				description = "Average Age"
				method = {
					average_by_property = {
						"average_of" = "month"
						"measure_time_by" = "$updatedAt"
						"property" = "age"
					}
				}
			}
		}
	}
`

	var testAccActionConfigUpdate = baseBlueprintsTemplate(parentBlueprintIdentifier, childBlueprintIdentifier) + `
	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		properties = {
			"average_age" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title = "Average Age"
				icon = "Terraform"
				description = "Average Age"
				method = {
					average_by_property = {
						"average_of" = "day"
						"measure_time_by" = "$createdAt"
						"property" = "age"
					}
				}
			}
		}
	}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.title", "Average Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.description", "Average Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.method.average_by_property.average_of", "month"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.method.average_by_property.measure_time_by", "$updatedAt"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.method.average_by_property.property", "age"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.title", "Average Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.description", "Average Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.method.average_by_property.average_of", "day"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.method.average_by_property.measure_time_by", "$createdAt"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.average_age.method.average_by_property.property", "age"),
				),
			},
		},
	})
}

func TestAccPortCreateAggregationPropertyAggregateByProperty(t *testing.T) {
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	var testAccActionConfigCreateAggrByPropMin = baseBlueprintsTemplate(parentBlueprintIdentifier, childBlueprintIdentifier) + `
	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		properties           = {
			"aggr" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title                       = "Min Age"
				icon                        = "Terraform"
				description                 = "Min Age"
				method                      = {
					aggregate_by_property = {
						"func"     = "min"
						"property" = "age"
					}
				}
			}
		}
	}
`

	var testAccAggregatePropertyUpdateMax = baseBlueprintsTemplate(parentBlueprintIdentifier, childBlueprintIdentifier) + `
	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		properties           = {
			"aggr" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title                       = "Max Age"
				icon                        = "Terraform"
				description                 = "Max Age"
				method                      = {
					aggregate_by_property = {
						"func"     = "max"
						"property" = "age"
					}
				}
			}
		}
	}
`

	var testAccAggregatePropertyUpdateSum = baseBlueprintsTemplate(parentBlueprintIdentifier, childBlueprintIdentifier) + `
	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		properties           = {
			"aggr" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title                       = "Sum Age"
				icon                        = "Terraform"
				description                 = "Sum Age"
				method                      = {
					aggregate_by_property = {
						"func"     = "sum"
						"property" = "age"
					}
				}
			}
		}
	}
`

	var testAccAggregatePropertyUpdateMedian = baseBlueprintsTemplate(parentBlueprintIdentifier, childBlueprintIdentifier) + `
	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		properties           = {
			"aggr" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title                       = "Median Age"
				icon                        = "Terraform"
				description                 = "Median Age"
				method                      = {
					aggregate_by_property = {
						"func"     = "median"
						"property" = "age"
					}
				}
			}
		}
	}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreateAggrByPropMin,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.title", "Min Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.description", "Min Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.method.aggregate_by_property.func", "min"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.method.aggregate_by_property.property", "age"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccAggregatePropertyUpdateMax,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.title", "Max Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.description", "Max Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.method.aggregate_by_property.func", "max"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.method.aggregate_by_property.property", "age"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccAggregatePropertyUpdateSum,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.title", "Sum Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.description", "Sum Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.method.aggregate_by_property.func", "sum"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.method.aggregate_by_property.property", "age"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccAggregatePropertyUpdateMedian,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.title", "Median Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.description", "Median Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.method.aggregate_by_property.func", "median"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.method.aggregate_by_property.property", "age"),
				),
			},
		},
	})
}

func TestAccPortCreateBlueprintWithAggregationByPropertyWithFilter(t *testing.T) {
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	var testAccActionConfigCreate = baseBlueprintsTemplate(parentBlueprintIdentifier, childBlueprintIdentifier) + `

	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		properties           = {
			"aggr" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title                       = "Min Age"
				icon                        = "Terraform"
				description                 = "Min Age"
				method                      = {
					aggregate_by_property = {
						"func"     = "min"
						"property" = "age"
					}
				}
				query = jsonencode(
					{
						"combinator" : "and",
						"rules" : [
							{
								"property" : "age",
								"operator" : "=",
								"value" : 10
							}
						]
					}
				)
			}
		}
	}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.title", "Min Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.description", "Min Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.method.aggregate_by_property.func", "min"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.method.aggregate_by_property.property", "age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.query", "{\"combinator\":\"and\",\"rules\":[{\"operator\":\"=\",\"property\":\"age\",\"value\":10}]}"),
				),
			},
		},
	})
}

func TestAccMultipleAggregationPropertiesForBlueprintCreate(t *testing.T) {
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	var testAccActionConfigCreate = baseBlueprintsTemplate(parentBlueprintIdentifier, childBlueprintIdentifier) + `
	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		properties           = {
			"aggr" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title                       = "Min Age"
				icon                        = "Terraform"
				description                 = "Min Age"
				method                      = {
					aggregate_by_property = {
						"func"     = "min"
						"property" = "age"
					}
				}
				query = jsonencode(
					{
						"combinator" : "and",
						"rules" : [
							{
								"property" : "age",
								"operator" : "=",
								"value" : 10
							}
						]
					}
				)
			}
			"aggr2" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title                       = "Max age"
				icon                        = "Terraform"
				description                 = "Max age"
				method                      = {
					aggregate_by_property = {
						"func"     = "max"
						"property" = "age"
					}
				}
				query = jsonencode(
					{
						"combinator" : "and",
						"rules" : [
							{
								"property" : "age",
								"operator" : "=",
								"value" : 10
							}
						]
					}
				)
			}
		}
	}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.title", "Min Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.description", "Min Age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.method.aggregate_by_property.func", "min"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.method.aggregate_by_property.property", "age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr.query", "{\"combinator\":\"and\",\"rules\":[{\"operator\":\"=\",\"property\":\"age\",\"value\":10}]}"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr2.title", "Max age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr2.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr2.description", "Max age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr2.target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr2.method.aggregate_by_property.func", "max"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr2.method.aggregate_by_property.property", "age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.aggr2.query", "{\"combinator\":\"and\",\"rules\":[{\"operator\":\"=\",\"property\":\"age\",\"value\":10}]}"),
				),
			},
		},
	})
}

func TestAccPortCreateAggregationPropertyWithPathFilter(t *testing.T) {
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	intermediateBlueprintIdentifier := utils.GenID()

	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "parent_blueprint" {
		title = "Parent Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		properties = {
			number_props = {
				"age" = {
					title = "Age"
				}
			}
		}
	}

	resource "port_blueprint" "intermediate_blueprint" {
		title = "Intermediate Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		relations = {
			"parent" = {
				title = "Parent"
				target = port_blueprint.parent_blueprint.identifier
			}
		}
	}

	resource "port_blueprint" "child_blueprint" {
		title = "Child Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		relations = {
			"intermediate" = {
				title = "Intermediate"
				target = port_blueprint.intermediate_blueprint.identifier
			}
		}
	}

	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		properties = {
			"count_parents" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title = "Count Parents"
				icon = "Terraform"
				description = "Count Parents via Path Filter"
				method = {
					count_entities = true
				}
				path_filter = [
					{
						path = ["intermediate", "parent"]
					}
				]
			}
		}
	}
`, parentBlueprintIdentifier, intermediateBlueprintIdentifier, childBlueprintIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.title", "Count Parents"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.description", "Count Parents via Path Filter"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.method.count_entities", "true"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.path_filter.0.path.0", "intermediate"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.path_filter.0.path.1", "parent"),
				),
			},
		},
	})
}

func TestAccPortAddAndRemovePathFilterFromAggregationProperty(t *testing.T) {
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	intermediateBlueprintIdentifier := utils.GenID()

	var testAccActionConfigWithoutPathFilter = fmt.Sprintf(`
	resource "port_blueprint" "parent_blueprint" {
		title = "Parent Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		properties = {
			number_props = {
				"age" = {
					title = "Age"
				}
			}
		}
	}

	resource "port_blueprint" "intermediate_blueprint" {
		title = "Intermediate Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		relations = {
			"parent" = {
				title = "Parent"
				target = port_blueprint.parent_blueprint.identifier
			}
		}
	}

	resource "port_blueprint" "child_blueprint" {
		title = "Child Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		relations = {
			"intermediate" = {
				title = "Intermediate"
				target = port_blueprint.intermediate_blueprint.identifier
			}
		}
	}

	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		properties = {
			"count_parents" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title = "Count Parents"
				icon = "Terraform"
				description = "Count Parents"
				method = {
					count_entities = true
				}
			}
		}
	}
`, parentBlueprintIdentifier, intermediateBlueprintIdentifier, childBlueprintIdentifier)

	var testAccActionConfigWithPathFilter = fmt.Sprintf(`
	resource "port_blueprint" "parent_blueprint" {
		title = "Parent Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		properties = {
			number_props = {
				"age" = {
					title = "Age"
				}
			}
		}
	}

	resource "port_blueprint" "intermediate_blueprint" {
		title = "Intermediate Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		relations = {
			"parent" = {
				title = "Parent"
				target = port_blueprint.parent_blueprint.identifier
			}
		}
	}

	resource "port_blueprint" "child_blueprint" {
		title = "Child Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		relations = {
			"intermediate" = {
				title = "Intermediate"
				target = port_blueprint.intermediate_blueprint.identifier
			}
		}
	}

	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		properties = {
			"count_parents" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title = "Count Parents"
				icon = "Terraform"
				description = "Count Parents with Path Filter"
				method = {
					count_entities = true
				}
				path_filter = [
					{
						path = ["intermediate", "parent"]
					}
				]
			}
		}
	}
`, parentBlueprintIdentifier, intermediateBlueprintIdentifier, childBlueprintIdentifier)

	var testAccActionConfigRemovePathFilter = fmt.Sprintf(`
	resource "port_blueprint" "parent_blueprint" {
		title = "Parent Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		properties = {
			number_props = {
				"age" = {
					title = "Age"
				}
			}
		}
	}

	resource "port_blueprint" "intermediate_blueprint" {
		title = "Intermediate Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		relations = {
			"parent" = {
				title = "Parent"
				target = port_blueprint.parent_blueprint.identifier
			}
		}
	}

	resource "port_blueprint" "child_blueprint" {
		title = "Child Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		relations = {
			"intermediate" = {
				title = "Intermediate"
				target = port_blueprint.intermediate_blueprint.identifier
			}
		}
	}

	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		properties = {
			"count_parents" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title = "Count Parents"
				icon = "Terraform"
				description = "Count Parents without Path Filter"
				method = {
					count_entities = true
				}
			}
		}
	}
`, parentBlueprintIdentifier, intermediateBlueprintIdentifier, childBlueprintIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigWithoutPathFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.title", "Count Parents"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.description", "Count Parents"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.method.count_entities", "true"),
					resource.TestCheckNoResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.path_filter"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigWithPathFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.title", "Count Parents"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.description", "Count Parents with Path Filter"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.method.count_entities", "true"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.path_filter.0.path.0", "intermediate"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.path_filter.0.path.1", "parent"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigRemovePathFilter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.title", "Count Parents"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.description", "Count Parents without Path Filter"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.method.count_entities", "true"),
					resource.TestCheckNoResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.path_filter"),
				),
			},
		},
	})
}

func TestAccPortCreateAggregationPropertyWithMultiplePathFilters(t *testing.T) {
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	intermediateBlueprintIdentifier := utils.GenID()
	otherBlueprintIdentifier := utils.GenID()

	var testAccActionConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "parent_blueprint" {
		title = "Parent Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		properties = {
			number_props = {
				"age" = {
					title = "Age"
				}
			}
		}
	}

	resource "port_blueprint" "intermediate_blueprint" {
		title = "Intermediate Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		relations = {
			"parent" = {
				title = "Parent"
				target = port_blueprint.parent_blueprint.identifier
			}
		}
	}

	resource "port_blueprint" "other_blueprint" {
		title = "Other Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		relations = {
			"intermediate" = {
				title = "Intermediate"
				target = port_blueprint.intermediate_blueprint.identifier
			}
		}
	}

	resource "port_blueprint" "child_blueprint" {
		title = "Child Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = ""
		relations = {
			"intermediate" = {
				title = "Intermediate"
				target = port_blueprint.intermediate_blueprint.identifier
			}
			"other" = {
				title = "Other"
				target = port_blueprint.other_blueprint.identifier
			}
		}
	}

	resource "port_aggregation_properties" "child_aggregation_properties" {
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		properties = {
			"count_parents" = {
				target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
				title = "Count Parents"
				icon = "Terraform"
				description = "Count Parents via Multiple Path Filters"
				method = {
					aggregate_by_property = {
						"func"     = "sum"
						"property" = "age"
					}
				}
				path_filter = [
					{
						path = ["intermediate", "parent"]
					},
					{
						path = ["other", "intermediate", "parent"]
					}
				]
			}
		}
	}
`, parentBlueprintIdentifier, intermediateBlueprintIdentifier, otherBlueprintIdentifier, childBlueprintIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.title", "Count Parents"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.description", "Count Parents via Multiple Path Filters"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.method.aggregate_by_property.func", "sum"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.method.aggregate_by_property.property", "age"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.path_filter.0.path.0", "intermediate"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.path_filter.0.path.1", "parent"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.path_filter.1.path.0", "other"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.path_filter.1.path.1", "intermediate"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.path_filter.1.path.2", "parent"),
				),
			},
		},
	})
}
