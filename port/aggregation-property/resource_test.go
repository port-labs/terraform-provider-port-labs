package aggregation_property

//
//import (
//	"fmt"
//	"github.com/port-labs/terraform-provider-port-labs/internal/acctest"
//	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
//)
//
//func TestAccPortCreateBlueprintWithAggregationCountEntitiesProperties(t *testing.T) {
//	baseIdentifier := utils.GenID()
//	aggrBlueprintIdentifier := utils.GenID()
//	var testAccActionConfigCreate = fmt.Sprintf(`
//	resource "port_blueprint" "base_blueprint" {
//		title = "Base Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		properties = {
//			number_props = {
//				"age" = {
//					title = "Age"
//				}
//			}
//		}
//	}
//
//	resource "port_blueprint" "aggr_blueprint" {
//		title = "Aggregation Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		aggregation_properties = {
//			"count_entities" = {
//				title = "Count Entities"
//				icon = "Terraform"
//				description = "Count Entities"
//                target = port_blueprint.base_blueprint.identifier
//			    method = {
//					count_entities = true
//				}
//			}
//		}
//        relations = {
//             "base_blueprint" = {
//					 title = "Base Blueprint"
//					 target = port_blueprint.base_blueprint.identifier
//			 }
//		}
//	}
//`, baseIdentifier, aggrBlueprintIdentifier)
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
//		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: acctest.ProviderConfig + testAccActionConfigCreate,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "title", "Aggregation Blueprint"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "identifier", aggrBlueprintIdentifier),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.count_entities.title", "Count Entities"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.count_entities.icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.count_entities.description", "Count Entities"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.count_entities.target", baseIdentifier),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.count_entities.method.count_entities", "true"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "relations.base_blueprint.title", "Base Blueprint"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "relations.base_blueprint.target", baseIdentifier),
//				),
//			},
//		},
//	})
//}
//
//func TestAccPortCreateBlueprintWithAggregationAverageEntities(t *testing.T) {
//	baseIdentifier := utils.GenID()
//	aggrBlueprintIdentifier := utils.GenID()
//	var testAccActionConfigCreate = fmt.Sprintf(`
//	resource "port_blueprint" "base_blueprint" {
//		title = "Base Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		properties = {
//			number_props = {
//				"age" = {
//					title = "Age"
//				}
//			}
//		}
//	}
//
//	resource "port_blueprint" "aggr_blueprint" {
//		title = "Aggregation Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		aggregation_properties = {
//			"average_entities" = {
//				title = "Average Entities"
//				icon = "Terraform"
//				description = "Average Entities"
//				target = port_blueprint.base_blueprint.identifier
//			    method = {
//					average_entities = {}
//				}
//			}
//		}
//		relations = {
//			 "base_blueprint" = {
//					 title = "Base Blueprint"
//					 target = port_blueprint.base_blueprint.identifier
//			 }
//		}
//	}
//`, baseIdentifier, aggrBlueprintIdentifier)
//
//	var testAccActionConfigUpdate = fmt.Sprintf(`
//	resource "port_blueprint" "base_blueprint" {
//		title = "Base Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		properties = {
//			number_props = {
//				"age" = {
//					title = "Age"
//					type = "number"
//				}
//			}
//		}
//	}
//
//	resource "port_blueprint" "aggr_blueprint" {
//		title = "Aggregation Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		aggregation_properties = {
//			"average_entities" = {
//				title = "Average Entities"
//				icon = "Terraform"
//				description = "Average Entities"
//				target = port_blueprint.base_blueprint.identifier
//			    method = {
//					average_entities = {
//                        "average_of" = "month"
//                        "measure_time_by" = "$updatedAt"
//                    }
//				}
//			}
//		}
//		relations = {
//			 "base_blueprint" = {
//					 title = "Base Blueprint"
//					 target = port_blueprint.base_blueprint.identifier
//			 }
//		}
//	}
//`, baseIdentifier, aggrBlueprintIdentifier)
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
//		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: acctest.ProviderConfig + testAccActionConfigCreate,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "title", "Aggregation Blueprint"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "identifier", aggrBlueprintIdentifier),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average_entities.title", "Average Entities"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average_entities.icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average_entities.description", "Average Entities"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average_entities.target", baseIdentifier),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average_entities.method.average_entities.average_of", "day"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average_entities.method.average_entities.measure_time_by", "$createdAt"),
//				),
//			},
//			{
//				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "title", "Aggregation Blueprint"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "identifier", aggrBlueprintIdentifier),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average_entities.title", "Average Entities"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average_entities.icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average_entities.description", "Average Entities"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average_entities.target", baseIdentifier),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average_entities.method.average_entities.average_of", "month"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average_entities.method.average_entities.measure_time_by", "$updatedAt"),
//				),
//			},
//		},
//	})
//}
//
//func TestAccPortCreateBlueprintWithAggregationAverageProperties(t *testing.T) {
//	baseIdentifier := utils.GenID()
//	aggrBlueprintIdentifier := utils.GenID()
//	var testAccActionConfigCreate = fmt.Sprintf(`
//	resource "port_blueprint" "base_blueprint" {
//		title = "Base Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		properties = {
//			number_props = {
//				"age" = {
//					title = "Age"
//					type = "number"
//				}
//			}
//		}
//	}
//
//	resource "port_blueprint" "aggr_blueprint" {
//		title = "Aggregation Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		aggregation_properties = {
//			"average" = {
//				title = "Average"
//				icon = "Terraform"
//				description = "Average"
//                target = port_blueprint.base_blueprint.identifier
//				method = {
//					average_by_property = {
//						"average_of" = "month"
//						"measure_time_by" = "$updatedAt"
//                        "property" = "age"
//					}
//				}
//			}
//		}
//		relations = {
//			 "base_blueprint" = {
//					 title = "Base Blueprint"
//					 target = port_blueprint.base_blueprint.identifier
//			 }
//		}
//	}
//`, baseIdentifier, aggrBlueprintIdentifier)
//
//	var testAccActionConfigUpdate = fmt.Sprintf(`
//	resource "port_blueprint" "base_blueprint" {
//		title = "Base Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		properties = {
//			number_props = {
//				"age" = {
//					title = "Age"
//					type = "number"
//				}
//			}
//		}
//	}
//
//	resource "port_blueprint" "aggr_blueprint" {
//		title = "Aggregation Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		aggregation_properties = {
//			"average" = {
//				title = "Average"
//				icon = "Terraform"
//				description = "Average"
//				target = port_blueprint.base_blueprint.identifier
//				method = {
//					average_by_property = {
//						"average_of" = "day"
//						"measure_time_by" = "$createdAt"
//                        "property" = "age"
//					}
//				}
//			}
//		}
//		relations = {
//			 "base_blueprint" = {
//					 title = "Base Blueprint"
//					 target = port_blueprint.base_blueprint.identifier
//			 }
//		}
//	}
//`, baseIdentifier, aggrBlueprintIdentifier)
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
//		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: acctest.ProviderConfig + testAccActionConfigCreate,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "title", "Aggregation Blueprint"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "identifier", aggrBlueprintIdentifier),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average.title", "Average"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average.icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average.description", "Average"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average.method.average_by_property.average_of", "month"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average.method.average_by_property.measure_time_by", "$updatedAt"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average.method.average_by_property.property", "age"),
//				),
//			},
//			{
//				Config: acctest.ProviderConfig + testAccActionConfigUpdate,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "title", "Aggregation Blueprint"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "identifier", aggrBlueprintIdentifier),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average.title", "Average"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average.icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average.description", "Average"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average.method.average_by_property.average_of", "day"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average.method.average_by_property.measure_time_by", "$createdAt"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.average.method.average_by_property.property", "age"),
//				),
//			},
//		},
//	})
//}
//func TestAccPortCreateBlueprintWithAggregationByProperty(t *testing.T) {
//	baseIdentifier := utils.GenID()
//	aggrBlueprintIdentifier := utils.GenID()
//	var testAccActionConfigCreate = fmt.Sprintf(`
//	resource "port_blueprint" "base_blueprint" {
//		title = "Base Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		properties = {
//			number_props = {
//				"age" = {
//					title = "Age"
//					type = "number"
//				}
//			}
//		}
//	}
//
//	resource "port_blueprint" "aggr_blueprint" {
//		title = "Aggregation Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		aggregation_properties = {
//			"max_age" = {
//				title = "max age"
//				icon = "Terraform"
//				description = "max age"
//				target = port_blueprint.base_blueprint.identifier
//				method = {
//					aggregate_by_property = {
//						"property" = "age"
//                        "func" = "max"
//					}
//   				}
//			}
//		}
//		relations = {
//			 "base_blueprint" = {
//					 title = "Base Blueprint"
//					 target = port_blueprint.base_blueprint.identifier
//			 }
//		}
//	}
//`, baseIdentifier, aggrBlueprintIdentifier)
//
//	var testAccActionConfigUpdateOverallAge = fmt.Sprintf(`
//	resource "port_blueprint" "base_blueprint" {
//		title = "Base Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		properties = {
//			number_props = {
//				"age" = {
//					title = "Age"
//					type = "number"
//				}
//			}
//		}
//	}
//
//	resource "port_blueprint" "aggr_blueprint" {
//		title = "Aggregation Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		aggregation_properties = {
//			"overall_age" = {
//				title = "overall age"
//				icon = "Terraform"
//				description = "overall age"
//				target = port_blueprint.base_blueprint.identifier
//				method = {
//					aggregate_by_property = {
//						"property" = "age"
//                        "func" = "sum"
//					}
//				}
//			}
//		}
//		relations = {
//			 "base_blueprint" = {
//					 title = "Base Blueprint"
//					 target = port_blueprint.base_blueprint.identifier
//			 }
//		}
//	}
//`, baseIdentifier, aggrBlueprintIdentifier)
//
//	var testAccActionConfigUpdateMedianAge = fmt.Sprintf(`
//resource "port_blueprint" "base_blueprint" {
//		title = "Base Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		properties = {
//			number_props = {
//				"age" = {
//					title = "Age"
//					type = "number"
//				}
//			}
//		}
//	}
//
//	resource "port_blueprint" "aggr_blueprint" {
//		title = "Aggregation Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		aggregation_properties = {
//			"median_age" = {
//				title = "median_age"
//				icon = "Terraform"
//				description = "median age"
//				target = port_blueprint.base_blueprint.identifier
//				method = {
//					aggregate_by_property = {
//						"property" = "age"
//                        "func" = "median"
//					}
//				}
//			}
//		}
//		relations = {
//			 "base_blueprint" = {
//					 title = "Base Blueprint"
//					 target = port_blueprint.base_blueprint.identifier
//			 }
//		}
//	}
//`, baseIdentifier, aggrBlueprintIdentifier)
//
//	var testAccActionConfigUpdateMinAge = fmt.Sprintf(`
//	resource "port_blueprint" "base_blueprint" {
//		title = "Base Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		properties = {
//			number_props = {
//				"age" = {
//					title = "Age"
//					type = "number"
//				}
//			}
//		}
//	}
//
//	resource "port_blueprint" "aggr_blueprint" {
//		title = "Aggregation Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		aggregation_properties = {
//			"min_age" = {
//				title = "min age"
//				icon = "Terraform"
//				description = "min age"
//				target = port_blueprint.base_blueprint.identifier
//				method = {
//					aggregate_by_property = {
//						"property" = "age"
//                        "func" = "min"
//					}
//				}
//			}
//		}
//		relations = {
//			 "base_blueprint" = {
//					 title = "Base Blueprint"
//					 target = port_blueprint.base_blueprint.identifier
//			 }
//		}
//	}
//`, baseIdentifier, aggrBlueprintIdentifier)
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
//		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: acctest.ProviderConfig + testAccActionConfigCreate,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "title", "Aggregation Blueprint"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "identifier", aggrBlueprintIdentifier),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.max_age.title", "max age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.max_age.icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.max_age.description", "max age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.max_age.method.aggregate_by_property.property", "age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.max_age.method.aggregate_by_property.func", "max"),
//				),
//			},
//			{
//				Config: acctest.ProviderConfig + testAccActionConfigUpdateOverallAge,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "title", "Aggregation Blueprint"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "identifier", aggrBlueprintIdentifier),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.overall_age.title", "overall age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.overall_age.icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.overall_age.description", "overall age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.overall_age.method.aggregate_by_property.property", "age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.overall_age.method.aggregate_by_property.func", "sum"),
//				),
//			},
//			{
//				Config: acctest.ProviderConfig + testAccActionConfigUpdateMedianAge,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "title", "Aggregation Blueprint"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "identifier", aggrBlueprintIdentifier),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.median_age.title", "median_age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.median_age.icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.median_age.description", "median age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.median_age.method.aggregate_by_property.property", "age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.median_age.method.aggregate_by_property.func", "median"),
//				),
//			},
//			{
//				Config: acctest.ProviderConfig + testAccActionConfigUpdateMinAge,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "title", "Aggregation Blueprint"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "identifier", aggrBlueprintIdentifier),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.min_age.title", "min age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.min_age.icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.min_age.description", "min age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.min_age.method.aggregate_by_property.property", "age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.min_age.method.aggregate_by_property.func", "min"),
//				),
//			},
//		},
//	})
//}
//
//func TestAccPortCreateBlueprintWithAggregationByPropertyWithFilter(t *testing.T) {
//	baseIdentifier := utils.GenID()
//	aggrBlueprintIdentifier := utils.GenID()
//
//	var testAccActionConfigCreate = fmt.Sprintf(`
//	resource "port_blueprint" "base_blueprint" {
//		title = "Base Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		properties = {
//			number_props = {
//				"age" = {
//					title = "Age"
//					type = "number"
//				}
//				"height" = {
//					title = "Height"
//					type = "number"
//				}
//			}
//		}
//	}
//
//	resource "port_blueprint" "aggr_blueprint" {
//		title = "Aggregation Blueprint"
//		icon = "Terraform"
//		identifier = "%s"
//		description = ""
//		aggregation_properties = {
//			"overall_age" = {
//				title = "overall age"
//				icon = "Terraform"
//				description = "overall age"
//				target = port_blueprint.base_blueprint.identifier
//				method = {
//					aggregate_by_property = {
//						"property" = "age"
//						"func" = "sum"
//					}
//				}
//			  	query = jsonencode(
//					{
//					  "combinator": "and",
//					  "rules": [
//						{
//						  "property": "age",
//						  "operator": "=",
//						  "value": 10
//						}
//					  ]
//					}
//			  	)
//			}
//		}
//		relations = {
//			 "base_blueprint" = {
//					 title = "Base Blueprint"
//					 target = port_blueprint.base_blueprint.identifier
//			 }
//		}
//	}
//`, baseIdentifier, aggrBlueprintIdentifier)
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
//		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: acctest.ProviderConfig + testAccActionConfigCreate,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "title", "Aggregation Blueprint"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "identifier", aggrBlueprintIdentifier),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.overall_age.title", "overall age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.overall_age.icon", "Terraform"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.overall_age.description", "overall age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.overall_age.method.aggregate_by_property.property", "age"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.overall_age.method.aggregate_by_property.func", "sum"),
//					resource.TestCheckResourceAttr("port_blueprint.aggr_blueprint", "aggregation_properties.overall_age.query", "{\"combinator\":\"and\",\"rules\":[{\"operator\":\"=\",\"property\":\"age\",\"value\":10}]}"),
//				),
//			},
//		},
//	})
//}
//
