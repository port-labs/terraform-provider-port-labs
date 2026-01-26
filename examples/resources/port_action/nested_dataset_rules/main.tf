# This example demonstrates different ways to use nested dataset rules with combinator groups
# Note: Change the identifier prefix if these conflict with existing blueprints in your org

locals {
  # Change this prefix if you have naming conflicts
  prefix = "nested_example"
}

resource "port_blueprint" "service" {
  title      = "Service (Nested Rules Example)"
  icon       = "Microservice"
  identifier = "${local.prefix}_service"
  properties = {
    string_props = {
      "environment" = {
        title = "Environment"
        enum  = ["development", "staging", "production"]
      }
      "team" = {
        title = "Team"
      }
      "region" = {
        title = "Region"
        enum  = ["us-east-1", "us-west-2", "eu-west-1"]
      }
    }
    boolean_props = {
      "user_assignable" = {
        title = "User Assignable"
      }
    }
  }
}

resource "port_blueprint" "permission_set" {
  title      = "Permission Set (Nested Rules Example)"
  icon       = "Lock"
  identifier = "${local.prefix}_permission_set"
  properties = {
    string_props = {
      "name" = {
        title = "Name"
      }
      "level" = {
        title = "Level"
        enum  = ["read", "write", "admin"]
      }
    }
    boolean_props = {
      "user_assignable" = {
        title = "User Assignable"
      }
    }
  }
  relations = {
    "business_domain" = {
      title  = "Business Domain"
      target = port_blueprint.service.identifier
    }
  }
}

# Example 1: Simple nested dataset with OR combinator
# Use case: Filter entities where property matches one of multiple values
resource "port_action" "or_combinator_example" {
  title      = "Assign Permission"
  identifier = "${local.prefix}_assign_permission_or"
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.service.identifier
    user_properties = {
      string_props = {
        permission = {
          title     = "Permission"
          format    = "entity"
          blueprint = port_blueprint.permission_set.identifier
          dataset = {
            combinator = "and"
            rules = [
              # First rule: must be user assignable
              {
                property = "user_assignable"
                operator = "="
                value = {
                  jq_query = "true"
                }
              },
              # Second rule: level must be read OR write (not admin)
              {
                combinator = "or"
                rules = [
                  {
                    property = "level"
                    operator = "="
                    value = {
                      jq_query = "\"read\""
                    }
                  },
                  {
                    property = "level"
                    operator = "="
                    value = {
                      jq_query = "\"write\""
                    }
                  }
                ]
              }
            ]
          }
        }
      }
    }
  }
  description = "Assign a read or write permission to a service"
  webhook_method = {
    url = "https://api.example.com/assign-permission"
  }
}

# Example 2: Nested dataset with dynamic values from form
# Use case: Filter based on user selection plus additional conditions
resource "port_action" "dynamic_filter_example" {
  title      = "Select Service"
  identifier = "${local.prefix}_select_service_dynamic"
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.service.identifier
    user_properties = {
      string_props = {
        environment = {
          title = "Environment"
          enum  = ["development", "staging", "production"]
        }
        target_service = {
          title     = "Target Service"
          format    = "entity"
          blueprint = port_blueprint.service.identifier
          depends_on = ["environment"]
          dataset = {
            combinator = "and"
            rules = [
              # Match the selected environment
              {
                property = "environment"
                operator = "="
                value = {
                  jq_query = ".form.environment"
                }
              },
              # AND (team is platform OR team is devops)
              {
                combinator = "or"
                rules = [
                  {
                    property = "team"
                    operator = "="
                    value = {
                      jq_query = "\"platform\""
                    }
                  },
                  {
                    property = "team"
                    operator = "="
                    value = {
                      jq_query = "\"devops\""
                    }
                  }
                ]
              }
            ]
          }
        }
      }
    }
  }
  description = "Select a service from platform or devops team in the chosen environment"
  webhook_method = {
    url = "https://api.example.com/select-service"
  }
}

# Example 3: Complex nested rules with multiple conditions
# Use case: Filter with complex business logic (AWS SSO style permissions)
resource "port_action" "complex_nested_example" {
  title      = "Assign AWS Permission Set"
  identifier = "${local.prefix}_assign_aws_permission"
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.service.identifier
    user_properties = {
      string_props = {
        business_domain = {
          title = "Business Domain"
          enum  = ["engineering", "finance", "operations"]
        }
        permission_set = {
          title     = "Permission Set"
          required  = true
          format    = "entity"
          blueprint = port_blueprint.permission_set.identifier
          depends_on = ["business_domain"]
          dataset = {
            combinator = "and"
            rules = [
              # Rule 1: Permission must be user assignable
              {
                property = "user_assignable"
                operator = "!="
                value = {
                  jq_query = "false"
                }
              },
              # Rule 2: OR condition for business domain matching
              {
                combinator = "or"
                rules = [
                  # Option A: Related to "engineering" domain
                  {
                    property = "$identifier"
                    operator = "contains"
                    value = {
                      jq_query = "\"engineering\""
                    }
                  },
                  # Option B: Related to user-selected business domain
                  {
                    property = "$identifier"
                    operator = "contains"
                    value = {
                      jq_query = ".form.business_domain"
                    }
                  }
                ]
              }
            ]
          }
        }
      }
    }
  }
  description = "Assign AWS permission set based on business domain"
  webhook_method = {
    url = "https://api.example.com/assign-aws-permission"
  }
}

# Example 4: Two levels of nesting
# Use case: (A AND B) OR (C AND D)
resource "port_action" "two_level_nesting_example" {
  title      = "Select Eligible Service"
  identifier = "${local.prefix}_select_eligible_service"
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.service.identifier
    user_properties = {
      string_props = {
        eligible_service = {
          title     = "Eligible Service"
          format    = "entity"
          blueprint = port_blueprint.service.identifier
          dataset = {
            # Top level: OR - match either condition group
            combinator = "or"
            rules = [
              # Group 1: Production services in us-east-1
              {
                combinator = "and"
                rules = [
                  {
                    property = "environment"
                    operator = "="
                    value = {
                      jq_query = "\"production\""
                    }
                  },
                  {
                    property = "region"
                    operator = "="
                    value = {
                      jq_query = "\"us-east-1\""
                    }
                  }
                ]
              },
              # Group 2: Staging services in eu-west-1
              {
                combinator = "and"
                rules = [
                  {
                    property = "environment"
                    operator = "="
                    value = {
                      jq_query = "\"staging\""
                    }
                  },
                  {
                    property = "region"
                    operator = "="
                    value = {
                      jq_query = "\"eu-west-1\""
                    }
                  }
                ]
              }
            ]
          }
        }
      }
    }
  }
  description = "Select production in us-east-1 OR staging in eu-west-1"
  webhook_method = {
    url = "https://api.example.com/select-eligible"
  }
}
