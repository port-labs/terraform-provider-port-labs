# This example demonstrates different ways to use max_items_jq_query and min_items_jq_query

resource "port_blueprint" "microservice" {
  title      = "Microservice"
  icon       = "Microservice"
  identifier = "microservice"
  properties = {
    string_props = {
      "environment" = {
        title = "Environment"
        enum  = ["development", "staging", "production"]
      }
    }
  }
}

resource "port_blueprint" "project" {
  title      = "Project"
  icon       = "Blueprint"
  identifier = "project"
  properties = {
    string_props = {
      "name" = {
        title = "Project Name"
      }
    }
  }
}

resource "port_blueprint" "deployment" {
  title      = "Deployment"
  icon       = "Deployment"
  identifier = "deployment"
  properties = {
    string_props = {
      "name" = {
        title = "Deployment Name"
      }
    }
  }
}

resource "port_action" "dynamic_max_items_example" {
  title      = "Add Reviewers"
  identifier = "add_reviewers"
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.microservice.identifier
    user_properties = {
      array_props = {
        reviewers = {
          title = "Reviewers"
          string_items = {
            format = "user"
          }
          # Dynamically calculate remaining reviewer slots (max 6 total):
          # 1. (.form.reviewers // []) - Get current reviewers array, default to empty array if null
          # 2. | length - Count how many reviewers are already selected
          # 3. (6 - ...) as $n - Subtract from 6 to get remaining slots, store in variable $n
          # 4. if $n < 0 then 0 else $n end - Ensure result is never negative (min 0)
          # Result: As users select reviewers, the max limit decreases accordingly
          max_items_jq_query = "(6 - ((.form.reviewers // []) | length)) as $n | if $n < 0 then 0 else $n end"
          description = "Select up to 6 reviewers (limit decreases as you select)"
        }
      }
    }
  }
  description = "Add reviewers to a service"
  webhook_method = {
    url = "https://api.example.com/add-reviewers"
  }
}

resource "port_action" "conditional_min_items_example" {
  title      = "Update Service Tags"
  identifier = "update_service_tags"
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.microservice.identifier
    user_properties = {
      array_props = {
        tags = {
          title = "Tags"
          string_items = {}
          min_items_jq_query = "if .entity.properties.environment == \"production\" then 2 else 1 end"
          max_items = 10
          description = "Service tags (production requires at least 2)"
        }
      }
    }
  }
  description = "Update service tags"
  webhook_method = {
    url = "https://api.example.com/update-tags"
  }
}

resource "port_action" "priority_based_max_items_example" {
  title      = "Request Approval"
  identifier = "request_approval"
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.microservice.identifier
    user_properties = {
      string_props = {
        priority = {
          title = "Priority"
          enum  = ["low", "medium", "high"]
        }
      }
      array_props = {
        approvers = {
          title = "Approvers"
          string_items = {
            format = "user"
          }
          max_items_jq_query = "if .form.priority == \"high\" then 5 elif .form.priority == \"medium\" then 3 else 1 end"
          description = "Select approvers (max depends on priority)"
        }
      }
    }
  }
  description = "Request approval for changes"
  webhook_method = {
    url = "https://api.example.com/request-approval"
  }
}

resource "port_action" "both_jq_queries_example" {
  title      = "Assign Team Members"
  identifier = "assign_team_members"
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.project.identifier
    user_properties = {
      string_props = {
        project_size = {
          title = "Project Size"
          enum  = ["small", "medium", "large"]
        }
      }
      array_props = {
        team_members = {
          title = "Team Members"
          string_items = {
            format = "user"
          }
          # Dynamically set minimum team members based on project size:
          # - Large projects: require at least 3 members
          # - Medium projects: require at least 2 members
          # - Small projects: require at least 1 member
          min_items_jq_query = "if .form.project_size == \"large\" then 3 elif .form.project_size == \"medium\" then 2 else 1 end"
          # Dynamically set maximum team members based on project size:
          # - Large projects: allow up to 10 members
          # - Medium projects: allow up to 5 members
          # - Small projects: allow up to 3 members
          max_items_jq_query = "if .form.project_size == \"large\" then 10 elif .form.project_size == \"medium\" then 5 else 3 end"
          description = "Assign team members (range depends on project size)"
        }
      }
    }
  }
  description = "Assign team members to project"
  webhook_method = {
    url = "https://api.example.com/assign-team"
  }
}

resource "port_action" "static_values_example" {
  title      = "Add Labels"
  identifier = "add_labels"
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.microservice.identifier
    user_properties = {
      array_props = {
        labels = {
          title = "Labels"
          string_items = {}
          min_items = 1 ## test with integer
          max_items = 5 ## test with integer
          description = "Service labels (1-5 required)"
        }
      }
    }
  }
  description = "Add labels to service"
  webhook_method = {
    url = "https://api.example.com/add-labels"
  }
}

resource "port_action" "simple_jq_query_example" {
  title      = "Configure Replicas"
  identifier = "configure_replicas"
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.deployment.identifier
    user_properties = {
      array_props = {
        availability_zones = {
          title = "Availability Zones"
          string_items = {}
          min_items_jq_query = "1" ##test with string
          max_items_jq_query = "3" ##test with string
          description = "Select availability zones (1-3)"
        }
      }
    }
  }
  description = "Configure deployment replicas"
  webhook_method = {
    url = "https://api.example.com/configure-replicas"
  }
}
