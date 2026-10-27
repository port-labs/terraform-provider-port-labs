package aggregation_properties_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
	"github.com/port-labs/terraform-provider-port-labs/v2/version"
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

	resource.ParallelTest(t, resource.TestCase{
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

	resource.ParallelTest(t, resource.TestCase{
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

	resource.ParallelTest(t, resource.TestCase{
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

	resource.ParallelTest(t, resource.TestCase{
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

	resource.ParallelTest(t, resource.TestCase{
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

	resource.ParallelTest(t, resource.TestCase{
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

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccPortCreateAggregationPropertyWithPathFilterAndFromBlueprint(t *testing.T) {
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
		blueprint_identifier = port_blueprint.parent_blueprint.identifier
		properties = {
			"count_parents" = {
				target_blueprint_identifier = port_blueprint.child_blueprint.identifier
				title = "Count Parents"
				icon = "Terraform"
				description = "Count Parents via Path Filter"
				method = {
					count_entities = true
				}
				path_filter = [
					{
						path = ["intermediate", "parent"]
						from_blueprint = port_blueprint.child_blueprint.identifier
					}
				]
			}
		}
	}
`, parentBlueprintIdentifier, intermediateBlueprintIdentifier, childBlueprintIdentifier)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.title", "Count Parents"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.description", "Count Parents via Path Filter"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.target_blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.method.count_entities", "true"),
					resource.TestCheckResourceAttr("port_aggregation_properties.child_aggregation_properties", "properties.count_parents.path_filter.0.from_blueprint", childBlueprintIdentifier),
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

	resource.ParallelTest(t, resource.TestCase{
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

// testAccPortClient creates an authenticated Port client for acceptance tests.
func testAccPortClient(t *testing.T) (*cli.PortClient, context.Context) {
	t.Helper()

	baseURL := os.Getenv("PORT_BASE_URL")
	if baseURL == "" {
		baseURL = consts.DefaultBaseUrl
	}
	client, err := cli.New(baseURL, cli.WithHeader("User-Agent", version.ProviderVersion))
	if err != nil {
		t.Fatalf("failed to create Port client: %v", err)
	}
	ctx := context.Background()
	if _, err := client.Authenticate(ctx, os.Getenv("PORT_CLIENT_ID"), os.Getenv("PORT_CLIENT_SECRET")); err != nil {
		t.Fatalf("failed to authenticate with Port: %v", err)
	}
	return client, ctx
}

// TestAccPreservesSystemAggregationOnTeamBlueprint tests that the aggregation_properties
// resource preserves system-managed aggregations (like _team.size) when creating, updating,
// and destroying user-defined aggregation properties on a system blueprint.
//
// This test reproduces the issue from ZD 11534 where applying port_aggregation_properties
// on _team would fail with "protected_blueprint_violation" because the provider was
// replacing the entire aggregationProperties map instead of merging.
func TestAccPreservesSystemAggregationOnTeamBlueprint(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}
	acctest.TestAccPreCheck(t)
	client, ctx := testAccPortClient(t)

	// Read the _team blueprint before the test to capture existing aggregations
	teamBpBefore, statusCode, err := client.ReadBlueprint(ctx, "_team")
	if err != nil {
		if statusCode == 404 {
			t.Skip("_team blueprint not available in this organization")
		}
		t.Fatalf("failed to read _team blueprint before test: %v", err)
	}

	// Record which system aggregations exist (e.g., "size")
	systemAggregations := make(map[string]bool)
	for k := range teamBpBefore.AggregationProperties {
		systemAggregations[k] = true
	}
	if len(systemAggregations) == 0 {
		t.Skip("_team blueprint has no existing aggregations to preserve")
	}
	t.Logf("System aggregations on _team before test: %v", systemAggregations)

	// Create a helper blueprint that _team can reference
	helperBlueprintIdentifier := utils.GenID()

	var testAccConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "helper" {
		title = "Helper Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = "Helper for aggregation test"
	}

	resource "port_aggregation_properties" "team_aggregation" {
		blueprint_identifier = "_team"
		properties = {
			"tf_test_count" = {
				target_blueprint_identifier = port_blueprint.helper.identifier
				title = "TF Test Count"
				icon = "Terraform"
				description = "Test aggregation from Terraform"
				method = {
					count_entities = true
				}
			}
		}
	}
`, helperBlueprintIdentifier)

	var testAccConfigUpdate = fmt.Sprintf(`
	resource "port_blueprint" "helper" {
		title = "Helper Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = "Helper for aggregation test"
	}

	resource "port_aggregation_properties" "team_aggregation" {
		blueprint_identifier = "_team"
		properties = {
			"tf_test_count" = {
				target_blueprint_identifier = port_blueprint.helper.identifier
				title = "TF Test Count Updated"
				icon = "Terraform"
				description = "Test aggregation from Terraform - updated"
				method = {
					count_entities = true
				}
			}
		}
	}
`, helperBlueprintIdentifier)

	// Config with just the helper blueprint (no aggregation resource)
	var testAccConfigDestroy = fmt.Sprintf(`
	resource "port_blueprint" "helper" {
		title = "Helper Blueprint"
		icon = "Terraform"
		identifier = "%s"
		description = "Helper for aggregation test"
	}
`, helperBlueprintIdentifier)

	// Helper function to verify system aggregations are preserved
	checkSystemAggregationsPreserved := func(phase string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			teamBp, _, err := client.ReadBlueprint(ctx, "_team")
			if err != nil {
				return fmt.Errorf("failed to read _team blueprint during %s: %w", phase, err)
			}
			for sysAgg := range systemAggregations {
				if _, exists := teamBp.AggregationProperties[sysAgg]; !exists {
					return fmt.Errorf("system aggregation %q was removed from _team during %s", sysAgg, phase)
				}
			}
			t.Logf("System aggregations preserved during %s: %v", phase, systemAggregations)
			return nil
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.team_aggregation", "blueprint_identifier", "_team"),
					resource.TestCheckResourceAttr("port_aggregation_properties.team_aggregation", "properties.tf_test_count.title", "TF Test Count"),
					checkSystemAggregationsPreserved("create"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_properties.team_aggregation", "properties.tf_test_count.title", "TF Test Count Updated"),
					checkSystemAggregationsPreserved("update"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccConfigDestroy,
				Check: resource.ComposeTestCheckFunc(
					checkSystemAggregationsPreserved("destroy"),
				),
			},
		},
	})
}
