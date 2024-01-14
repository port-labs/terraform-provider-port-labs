package aggregation_properties

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func AggregationPropertySchema() schema.Attribute {
	return schema.MapNestedAttribute{
		MarkdownDescription: "The aggregation property of the blueprint",
		Required:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"target_blueprint_identifier": schema.StringAttribute{
					MarkdownDescription: "The identifier of the blueprint to perform the aggregation on",
					Required:            true,
				},
				"title": schema.StringAttribute{
					MarkdownDescription: "The title of the aggregation property",
					Optional:            true,
				},
				"icon": schema.StringAttribute{
					MarkdownDescription: "The icon of the aggregation property",
					Optional:            true,
				},
				"description": schema.StringAttribute{
					MarkdownDescription: "The description of the aggregation property",
					Optional:            true,
				},
				"method": schema.SingleNestedAttribute{
					MarkdownDescription: "The aggregation method to perform on the target blueprint, one of count_entities, average_entities, average_by_property, aggregate_by_property",
					Required:            true,
					Attributes: map[string]schema.Attribute{
						"count_entities": schema.BoolAttribute{
							MarkdownDescription: "Function to count the entities of the target entities",
							Optional:            true,
							Validators: []validator.Bool{
								boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("average_entities")),
								boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("average_by_property")),
								boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("aggregate_by_property")),
							},
						},
						"average_entities": schema.SingleNestedAttribute{
							MarkdownDescription: "Function to average the entities of the target entities",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"average_of": schema.StringAttribute{
									MarkdownDescription: "The time periods to calculate the average of, e.g. hour, day, week, month",
									Optional:            true,
									Computed:            true,
									Default:             stringdefault.StaticString("day"),
									Validators: []validator.String{
										stringvalidator.OneOf("hour", "day", "week", "month"),
									},
								},
								"measure_time_by": schema.StringAttribute{
									MarkdownDescription: "The property name on which to calculate the the time periods, e.g. $createdAt, $updated_at or any other date property",
									Optional:            true,
									Computed:            true,
									Default:             stringdefault.StaticString("$createdAt"),
								},
							},
							Validators: []validator.Object{
								objectvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("count_entities"),
									path.MatchRelative().AtParent().AtName("average_by_property"),
									path.MatchRelative().AtParent().AtName("aggregate_by_property"),
								),
							},
						},
						"average_by_property": schema.SingleNestedAttribute{
							MarkdownDescription: "Function to calculate the average by property value of the target entities",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"average_of": schema.StringAttribute{
									MarkdownDescription: "The time periods to calculate the average by, e.g. hour, day, week, month",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("hour", "day", "week", "month", "total"),
									},
								},
								"measure_time_by": schema.StringAttribute{
									MarkdownDescription: "The property name on which to calculate the the time periods, e.g. $createdAt, $updated_at or any other date property",
									Required:            true,
								},
								"property": schema.StringAttribute{
									MarkdownDescription: "The property name on which to calculate the average by",
									Required:            true,
								},
							},
							Validators: []validator.Object{
								objectvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("count_entities"),
									path.MatchRelative().AtParent().AtName("average_entities"),
									path.MatchRelative().AtParent().AtName("aggregate_by_property"),
								),
							},
						},
						"aggregate_by_property": schema.SingleNestedAttribute{
							MarkdownDescription: "Function to calculate the aggregate by property value of the target entities, such as sum, min, max, median",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"func": schema.StringAttribute{
									MarkdownDescription: "The func of the aggregate by property",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("sum", "min", "max", "median"),
									},
								},
								"property": schema.StringAttribute{
									MarkdownDescription: "The property of the aggregate by property",
									Required:            true,
								},
							},
							Validators: []validator.Object{
								objectvalidator.ConflictsWith(
									path.MatchRelative().AtParent().AtName("count_entities"),
									path.MatchRelative().AtParent().AtName("average_entities"),
									path.MatchRelative().AtParent().AtName("average_by_property"),
								),
							},
						},
					},
				},
				"query": schema.StringAttribute{
					MarkdownDescription: "Query to filter the target entities",
					Optional:            true,
				},
			},
		},
	}
}

func AggregationPropertiesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"blueprint_identifier": schema.StringAttribute{
			Description: "The identifier of the blueprint the aggregation property will be added to",
			Required:    true,
		},
		"properties": AggregationPropertySchema(),
	}
}

func (r *AggregationPropertiesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: AggregationPropertyResourceMarkdownDescription,
		Attributes:          AggregationPropertiesSchema(),
	}
}

var AggregationPropertyResourceMarkdownDescription = `

# Aggregation Property

This resource allows you to manage an aggregation property.

See the [Port documentation](https://docs.getport.io/build-your-software-catalog/define-your-data-model/setup-blueprint/properties/aggregation-properties/) for more information about aggregation properties.


Supported Methods:

- count_entities - Count the entities of the target blueprint
- average_entities - Average the entities of the target blueprint by time periods
- average_by_property - Calculate the average by property value of the target entities
- aggregate_by_property - Calculate the aggregate by property value of the target entities, such as sum, min, max, median

## Example Usage

Create a parent blueprint with a child blueprint and an aggregation property to count the parent kids:

` + "```hcl" + `

resource "port_blueprint" "parent_blueprint" {
  title       = "Parent Blueprint"
  icon        = "Terraform"
  identifier  = "parent"
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
  title       = "Child Blueprint"
  icon        = "Terraform"
  identifier  = "child"
  description = ""
  properties = {
    number_props = {
      "age" = {
        title = "Age"
      }
    }
  }
  relations = {
    "parent" = {
      title  = "Parent"
      target = port_blueprint.parent_blueprint.identifier
    }
  }
}

resource "port_aggregation_properties" "parent_aggregation_properties" {
  blueprint_identifier        = port_blueprint.parent_blueprint.identifier
  properties = {
    "count_kids" = {
      target_blueprint_identifier = port_blueprint.child_blueprint.identifier
      title                       = "Count Kids"
      icon                        = "Terraform"
      description                 = "Count Kids"
      method                      = {
        count_entities = true
      }
    }
  }
}

` + "```" + `

Create a parent blueprint with a child blueprint and an aggregation property to calculate the average avg of the parent kids age:

` + "```hcl" + `

resource "port_blueprint" "parent_blueprint" {
  title       = "Parent Blueprint"
  icon        = "Terraform"
  identifier  = "parent"
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
  title       = "Child Blueprint"
  icon        = "Terraform"
  identifier  = "child"
  description = ""
  properties = {
    number_props = {
      "age" = {
        title = "Age"
      }
    }
  }
  relations = {
    "parent" = {
      title  = "Parent"
      target = port_blueprint.parent_blueprint.identifier
    }
  }
}

resource "port_aggregation_properties" "parent_aggregation_properties" {
  blueprint_identifier = port_blueprint.parent_blueprint.identifier
  properties           = {
    average_kids_age = {
      target_blueprint_identifier = port_blueprint.child_blueprint.identifier
      title                       = "Average Kids Age"
      icon                        = "Terraform"
      description                 = "Average Kids Age"
      method                      = {
        average_by_property = {
          average_of      = "total"
          measure_time_by = "$createdAt"
          property        = "age"
        }
      }
    }
  }
}


` + "```" + `

Create a repository blueprint and a pull request blueprint and an aggregation property to calculate the average of pull requests created per day:

` + "```hcl" + `

resource "port_blueprint" "repository_blueprint" {
  title       = "Repository Blueprint"
  icon        = "Terraform"
  identifier  = "repository"
  description = ""
}

resource "port_blueprint" "pull_request_blueprint" {
  title       = "Pull Request Blueprint"
  icon        = "Terraform"
  identifier  = "pull_request"
  description = ""
  properties = {
    string_props = {
      "status" = {
        title = "Status"
      }
    }
  }
  relations = {
    "repository" = {
      title  = "Repository"
      target = port_blueprint.repository_blueprint.identifier
    }
  }
}

resource "port_aggregation_properties" "repository_aggregation_properties" {
  blueprint_identifier = port_blueprint.repository_blueprint.identifier
  properties           = {
    "pull_requests_per_day" = {
      target_blueprint_identifier = port_blueprint.pull_request_blueprint.identifier
      title                       = "Pull Requests Per Day"
      icon                        = "Terraform"
      description                 = "Pull Requests Per Day"
      method                      = {
        average_entities = {
          average_of      = "day"
          measure_time_by = "$createdAt"
        }
      }
    }
  }
}
  
` + "```" + `

Create a repository blueprint and a pull request blueprint and an aggregation property to calculate the average of fix pull request per month:

To do that we will add a query to the aggregation property to filter only pull requests with fixed title:

` + "```hcl" + `

resource "port_blueprint" "repository_blueprint" {
  title       = "Repository Blueprint"
  icon        = "Terraform"
  identifier  = "repository"
  description = ""
}

resource "port_blueprint" "pull_request_blueprint" {
  title       = "Pull Request Blueprint"
  icon        = "Terraform"
  identifier  = "pull_request"
  description = ""
  properties = {
    string_props = {
      "status" = {
        title = "Status"
      }
    }
  }
  relations = {
    "repository" = {
      title  = "Repository"
      target = port_blueprint.repository_blueprint.identifier
    }
  }
}

resource "port_aggregation_properties" "repository_aggregation_properties" {
  blueprint_identifier = port_blueprint.repository_blueprint.identifier
  properties           = {
    "fix_pull_requests_count" = {
      target_blueprint_identifier = port_blueprint.pull_request_blueprint.identifier
      title                       = "Pull Requests Per Day"
      icon                        = "Terraform"
      description                 = "Pull Requests Per Day"
      method                      = {
        average_entities = {
          average_of      = "month"
          measure_time_by = "$createdAt"
        }
      }
      query = jsonencode(
        {
          "combinator" : "and",
          "rules" : [
            {
              "property" : "$title",
              "operator" : "ContainsAny",
              "value" : ["fix", "fixed", "fixing", "Fix"]
            }
          ]
        }
      )
    }
  }
}

` + "```" + `


Create multiple aggregation properties in one resource:

` + "```hcl" + `

resource "port_blueprint" "repository_blueprint" {
  title       = "Repository Blueprint"
  icon        = "Terraform"
  identifier  = "repository"
  description = ""
}

resource "port_blueprint" "pull_request_blueprint" {
  title       = "Pull Request Blueprint"
  icon        = "Terraform"
  identifier  = "pull_request"
  description = ""
  properties = {
    string_props = {
      "status" = {
        title = "Status"
      }
    }
  }
  relations = {
    "repository" = {
      title  = "Repository"
      target = port_blueprint.repository_blueprint.identifier
    }
  }
}

resource "port_aggregation_properties" "repository_aggregation_properties" {
  blueprint_identifier = port_blueprint.repository_blueprint.identifier
  properties           = {
    "pull_requests_per_day" = {
      target_blueprint_identifier = port_blueprint.pull_request_blueprint.identifier
      title                       = "Pull Requests Per Day"
      icon                        = "Terraform"
      description                 = "Pull Requests Per Day"
      method                      = {
        average_entities = {
          average_of      = "day"
          measure_time_by = "$createdAt"
        }
      }
    }
    "overall_pull_requests_count" = {
      target_blueprint_identifier = port_blueprint.pull_request_blueprint.identifier
      title                       = "Overall Pull Requests Count"
      icon                        = "Terraform"
      description                 = "Overall Pull Requests Count"
      method                      = {
        count_entities = true
      }
    }
  }
}

` + "```" + ``
