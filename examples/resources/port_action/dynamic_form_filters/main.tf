# Example: Dynamic Form Filters with JQ Expressions
#
# This example demonstrates how to use dynamic JQ expressions in dataset rules
# to filter entity selectors based on form input values.

terraform {
  required_providers {
    port = {
      source  = "port-labs/port-labs"
      version = "~> 2.0"
    }
  }
}

provider "port" {
  client_id = var.port_client_id
  secret    = var.port_client_secret
}

variable "port_client_id" {
  type        = string
  description = "Port client ID"
}

variable "port_client_secret" {
  type        = string
  sensitive   = true
  description = "Port client secret"
}

# Create a service blueprint
resource "port_blueprint" "service" {
  title      = "Service"
  icon       = "Microservice"
  identifier = "dynamic_filter_example_service"
  properties = {
    string_props = {
      "environment" = {
        title = "Environment"
        enum  = ["development", "staging", "production"]
      }
      "team" = {
        title = "Team"
      }
    }
  }
}

# Create test entities - mix of environments and teams
resource "port_entity" "service_dev_engineering" {
  identifier = "service-dev-api"
  title      = "dev-api-service"
  blueprint  = port_blueprint.service.identifier
  
  properties = {
    string_props = {
      "environment" = "development"
      "team"        = "engineering"
    }
  }
}

resource "port_entity" "service_staging_platform" {
  identifier = "service-staging-web"
  title      = "staging-web-frontend"
  blueprint  = port_blueprint.service.identifier
  
  properties = {
    string_props = {
      "environment" = "staging"
      "team"        = "platform"
    }
  }
}

resource "port_entity" "service_prod_engineering" {
  identifier = "service-prod-api"
  title      = "prod-api-service"
  blueprint  = port_blueprint.service.identifier
  
  properties = {
    string_props = {
      "environment" = "production"
      "team"        = "engineering"
    }
  }
}

resource "port_entity" "service_dev_platform" {
  identifier = "service-dev-web"
  title      = "dev-web-frontend"
  blueprint  = port_blueprint.service.identifier
  
  properties = {
    string_props = {
      "environment" = "development"
      "team"        = "platform"
    }
  }
}

# Action with dynamic form filtering
# The entity selector filters based on the team selected in the form AND environment
resource "port_action" "select_service" {
  title      = "Select Service by Team"
  identifier = "dynamic_filter_example_select_service"
  
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.service.identifier
    
    order_properties = ["selected_team", "target_service"]
    
    user_properties = {
      string_props = {
        # First field: Team selector dropdown
        selected_team = {
          title       = "Select Team"
          description = "Choose a team to filter services"
          enum        = ["engineering", "platform"]
        }
        
        # Second field: Entity selector filtered by the selected team AND non-prod environment
        target_service = {
          title      = "Target Service"
          format     = "entity"
          blueprint  = port_blueprint.service.identifier
          depends_on = ["selected_team"]
          
          # Dynamic filtering using JQ expressions
          dataset = {
            combinator = "and"
            rules = [
              # Rule 1: Dynamic filter - team matches form selection
              # .form.selected_team evaluates to the user's selection at runtime
              {
                property = "team"
                operator = "="
                value = {
                  jq_query = ".form.selected_team"
                }
              },
              # Rule 2: Literal string filter - only non-production environments
              # The value is wrapped in quotes to indicate it's a literal string in JQ
              {
                property = "environment"
                operator = "!="
                value = {
                  jq_query = "\"production\""
                }
              }
            ]
          }
        }
      }
    }
  }
  
  webhook_method = {
    url = "https://example.com/webhook"
  }
}

output "action_identifier" {
  value = port_action.select_service.identifier
}
