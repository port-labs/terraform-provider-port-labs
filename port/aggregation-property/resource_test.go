package aggregation_property_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
)

func TestAccPortAggregationPropertyWithCycleRelation(t *testing.T) {
	// Test checks that a cycle aggregation property works.
	// The cycle is created by creating a parent blueprint and a child blueprint.
	// The child blueprint has a relation to the parent blueprint.
	// The parent blueprint has an aggregation property that counts the children of the parent.
	// The aggregation property is created with a cycle relation to the child blueprint, which is allowed.
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	aggregationPropIdentifier := utils.GenID()
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

	resource "port_aggregation_property" "count_entities" {
		aggregation_identifier = "%s"
		blueprint_identifier = port_blueprint.parent_blueprint.identifier
		target_blueprint_identifier = port_blueprint.child_blueprint.identifier
		title = "Count Childrens"
		icon = "Terraform"
		description = "Count Childrens"
		method = {
			count_entities = true
		}
	}
`, parentBlueprintIdentifier, childBlueprintIdentifier, aggregationPropIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "title", "Count Childrens"),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "aggregation_identifier", aggregationPropIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "description", "Count Childrens"),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "target_blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "method.count_entities", "true"),
				),
			},
		},
	})
}

func TestAccCreateAggregationPropertyAverageEntities(t *testing.T) {
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	aggregationPropIdentifier := utils.GenID()
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

	resource "port_aggregation_property" "count_entities" {
		aggregation_identifier = "%s"
		blueprint_identifier = port_blueprint.parent_blueprint.identifier
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
`, parentBlueprintIdentifier, childBlueprintIdentifier, aggregationPropIdentifier)

	var testAccActionConfigUpdate = fmt.Sprintf(`
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

	resource "port_aggregation_property" "count_entities" {
		aggregation_identifier = "%s"
		blueprint_identifier = port_blueprint.parent_blueprint.identifier
		target_blueprint_identifier = port_blueprint.child_blueprint.identifier
		title = "Count Childrens"
		icon = "Terraform"
		description = "Count Childrens"
		method = {
			average_entities = {}
		}
	}
`, parentBlueprintIdentifier, childBlueprintIdentifier, aggregationPropIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "title", "Count Childrens"),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "aggregation_identifier", aggregationPropIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "description", "Count Childrens"),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "target_blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "method.average_entities.average_of", "month"),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "method.average_entities.measure_time_by", "$updatedAt"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "title", "Count Childrens"),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "aggregation_identifier", aggregationPropIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "description", "Count Childrens"),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "target_blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "method.average_entities.average_of", "day"),
					resource.TestCheckResourceAttr("port_aggregation_property.count_entities", "method.average_entities.measure_time_by", "$createdAt"),
				),
			},
		},
	})
}

func TestAccPortCreateAggregationAverageProperties(t *testing.T) {
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	aggregationPropIdentifier := utils.GenID()
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

	resource "port_aggregation_property" "average_age" {
		aggregation_identifier = "%s"
		blueprint_identifier = port_blueprint.child_blueprint.identifier
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
`, parentBlueprintIdentifier, childBlueprintIdentifier, aggregationPropIdentifier)

	var testAccActionConfigUpdate = fmt.Sprintf(`
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

	resource "port_aggregation_property" "average_age" {
		aggregation_identifier = "%s"
		blueprint_identifier = port_blueprint.child_blueprint.identifier
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
`, parentBlueprintIdentifier, childBlueprintIdentifier, aggregationPropIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "title", "Average Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "aggregation_identifier", aggregationPropIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "description", "Average Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "method.average_by_property.average_of", "month"),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "method.average_by_property.measure_time_by", "$updatedAt"),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "method.average_by_property.property", "age"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "title", "Average Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "aggregation_identifier", aggregationPropIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "description", "Average Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "method.average_by_property.average_of", "day"),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "method.average_by_property.measure_time_by", "$createdAt"),
					resource.TestCheckResourceAttr("port_aggregation_property.average_age", "method.average_by_property.property", "age"),
				),
			},
		},
	})
}

func TestAccPortCreateAggregationPropertyAggregateByProperty(t *testing.T) {
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	aggregationPropIdentifier := utils.GenID()
	var testAccActionConfigCreateAggrByPropMin = fmt.Sprintf(`
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

	resource "port_aggregation_property" "aggr" {
		aggregation_identifier = "%s"
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
		title = "Min Age"
		icon = "Terraform"
		description = "Min Age"
		method = {
			aggregate_by_property = {
				"func" = "min"
				"property" = "age"
			}
		}
	}
`, parentBlueprintIdentifier, childBlueprintIdentifier, aggregationPropIdentifier)

	var testAccAggregatePropertyUpdateMax = fmt.Sprintf(`
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

	resource "port_aggregation_property" "aggr" {
		aggregation_identifier = "%s"
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
		title = "Max Age"
		icon = "Terraform"
		description = "Max Age"
		method = {
			aggregate_by_property = {
				"func" = "max"
				"property" = "age"
			}
		}
	}
`, parentBlueprintIdentifier, childBlueprintIdentifier, aggregationPropIdentifier)

	var testAccAggregatePropertyUpdateSum = fmt.Sprintf(`
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

	resource "port_aggregation_property" "aggr" {
		aggregation_identifier = "%s"
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
		title = "Sum Age"	
		icon = "Terraform"
		description = "Sum Age"
		method = {
			aggregate_by_property = {
				"func" = "sum"
				"property" = "age"
			}
		}
	}
`, parentBlueprintIdentifier, childBlueprintIdentifier, aggregationPropIdentifier)

	var testAccAggregatePropertyUpdateMedian = fmt.Sprintf(`
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

	resource "port_aggregation_property" "aggr" {	
		aggregation_identifier = "%s"
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
		title = "Median Age"
		icon = "Terraform"
		description = "Median Age"
		method = {
			aggregate_by_property = {
				"func" = "median"
				"property" = "age"
			}
		}
	}
`, parentBlueprintIdentifier, childBlueprintIdentifier, aggregationPropIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreateAggrByPropMin,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "title", "Min Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "aggregation_identifier", aggregationPropIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "description", "Min Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "method.aggregate_by_property.func", "min"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "method.aggregate_by_property.property", "age"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccAggregatePropertyUpdateMax,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "title", "Max Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "aggregation_identifier", aggregationPropIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "description", "Max Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "method.aggregate_by_property.func", "max"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "method.aggregate_by_property.property", "age"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccAggregatePropertyUpdateSum,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "title", "Sum Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "aggregation_identifier", aggregationPropIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "description", "Sum Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "method.aggregate_by_property.func", "sum"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "method.aggregate_by_property.property", "age"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccAggregatePropertyUpdateMedian,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "title", "Median Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "aggregation_identifier", aggregationPropIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "description", "Median Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "method.aggregate_by_property.func", "median"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "method.aggregate_by_property.property", "age"),
				),
			},
		},
	})
}

func TestAccPortCreateBlueprintWithAggregationByPropertyWithFilter(t *testing.T) {
	parentBlueprintIdentifier := utils.GenID()
	childBlueprintIdentifier := utils.GenID()
	aggregationPropIdentifier := utils.GenID()
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

	resource "port_aggregation_property" "aggr" {
		aggregation_identifier = "%s"
		blueprint_identifier = port_blueprint.child_blueprint.identifier
		target_blueprint_identifier = port_blueprint.parent_blueprint.identifier
		title = "Min Age"
		icon = "Terraform"
		description = "Min Age"
		method = {
			aggregate_by_property = {
				"func" = "min"
				"property" = "age"
			}
		}
		query = jsonencode(
			{
				"combinator": "and",
				"rules": [
					{
						"property": "age",
					  	"operator": "=",
					  	"value": 10
					}
				]
			}
		)
	}
`, parentBlueprintIdentifier, childBlueprintIdentifier, aggregationPropIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccActionConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "title", "Min Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "aggregation_identifier", aggregationPropIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "icon", "Terraform"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "description", "Min Age"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "blueprint_identifier", childBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "target_blueprint_identifier", parentBlueprintIdentifier),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "method.aggregate_by_property.func", "min"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "method.aggregate_by_property.property", "age"),
					resource.TestCheckResourceAttr("port_aggregation_property.aggr", "query", "{\"combinator\":\"and\",\"rules\":[{\"operator\":\"=\",\"property\":\"age\",\"value\":10}]}"),
				),
			},
		},
	})
}
