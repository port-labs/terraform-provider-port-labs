# This example demonstrates different ways to use pattern_jq_query

# Example 1: Using pattern_jq_query to dynamically generate a regex pattern
resource "port_action" "regex_pattern_example" {
  title      = "Create Service"
  identifier = "create_service"
  self_service_trigger = {
    operation           = "CREATE"
    blueprint_identifier = "service"
    user_properties = {
      string_props = {
        service_name = {
          title      = "Service Name"
          # Dynamic regex pattern based on context
          pattern_jq_query = "if .environment == \"production\" then \"^[a-z][a-z0-9-]{3,20}$\" else \"^[a-z][a-z0-9-]{2,10}$\" end"
          description = "Name of the service (lowercase, numbers, hyphens)"
        }
      }
    }
  }
  description = "Create a new service"
  webhook_method = {
    url = "https://api.example.com/create-service"
  }
}

# Example 2: Using pattern_jq_query to dynamically generate a list of allowed values
resource "port_action" "allowed_values_example" {
  title      = "Deploy To Environment"
  identifier = "deploy_to_env"
  self_service_trigger = {
    operation           = "DAY-2"
    blueprint_identifier = "microservice"
    user_properties = {
      string_props = {
        target_environment = {
          title      = "Target Environment"
          # Dynamic allowed values based on context
          pattern_jq_query = "if .team == \"platform\" then [\"dev\", \"staging\", \"production\"] else [\"dev\", \"staging\"] end"
          description = "Environment to deploy to (platform team can deploy to production)"
        }
      }
    }
  }
  description = "Deploy service to environment"
  webhook_method = {
    url = "https://api.example.com/deploy"
  }
}

# Example 3: Using pattern_jq_query with a direct JSON array of allowed values
resource "port_action" "direct_array_example" {
  title      = "Select Region"
  identifier = "select_region"
  self_service_trigger = {
    operation           = "CREATE"
    blueprint_identifier = "deployment"
    user_properties = {
      string_props = {
        region = {
          title      = "Region"
          # Direct JSON array of allowed values
          pattern_jq_query = "[\"us-east-1\", \"us-west-1\", \"eu-west-1\", \"ap-northeast-1\"]"
          description = "AWS region for deployment"
        }
      }
    }
  }
  description = "Select deployment region"
  webhook_method = {
    url = "https://api.example.com/select-region"
  }
} 