resource "port_blueprint" "environment" {
  title      = "Environment"
  icon       = "Environment"
  identifier = "examples-action-env"
  properties = {
    string_props = {
      "name" = {
        type  = "string"
        title = "name"
      }
      "docs-url" = {
        title  = "Docs URL"
        format = "url"
      }
    }
  }
}

resource "port_blueprint" "microservice" {
  title      = "VM"
  icon       = "GPU"
  identifier = "examples-action-microservice"
  properties = {
    string_props = {
      name = {
        type  = "string"
        title = "Name"
      },
      author = {
        type  = "string"
        title = "Author"
      },
      url = {
        type  = "string"
        title = "URL"
      },
      myStringIdentifier = {
        description = "This is a string property"
        title       = "text"
        icon        = "Terraform"
        required    = true
        min_length  = 1
        max_length  = 10
        default     = "default"
        enum        = ["default", "default2"]
        pattern     = "^[a-zA-Z0-9]*$"
        format      = "user"
        enum_colors = {
          default  = "red"
          default2 = "green"
        }
      }
    }
    boolean_props = {
      required = {
        type = "boolean"
      }
    }
    number_props = {
      sum = {
        type = "number"
      }
    }
  }

  relations = {
    "environment" = {
      title    = "Test Relation"
      required = "true"
      target   = port_blueprint.environment.identifier
    }
  }
}

resource "port_action" "restart_microservice" {
  title      = "Restart microservice"
  icon       = "Terraform"
  identifier = "examples-action-restart-microservice"
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.microservice.identifier
    user_properties = {
      string_props = {
        "webhook_url" = {
          title       = "Webhook URL"
          description = "Webhook URL to send the request to"
          format      = "url"
          default     = "https://example.com"
          pattern     = "^https://.*"
        }
      }
    }
  }
  webhook_method = {
    type = "WEBHOOK"
    url  = "https://app.getport.io"
  }
}

resource "port_action" "notifiy_on_mocrosiervice_creation" {
  title      = "Notify On Microservice Creation"
  icon       = "Terraform"
  identifier = "examples-automation-notify-on-microservice-creation"
  automation_trigger = {
    entity_created_event = {
      blueprint_identifier = port_blueprint.microservice.identifier
    }
  }
  webhook_method = {
    type = "WEBHOOK"
    url  = "https://example.com"
  }
  publish = true
}

resource "port_action" "notifiy_on_microservice_restart_failed" {
  title      = "Notify On Microservice Restart Failed"
  icon       = "Terraform"
  identifier = "examples-automation-notify-on-microservice-restart-failed"
  automation_trigger = {
    run_updated_event = {
      action_identifier = port_action.restart_microservice.identifier
    }
    jq_condition = {
      combinator = "and"
      expressions = [
        ".diff.after.status == \"FAILURE\""
      ],
    }
  }
  webhook_method = {
    type = "WEBHOOK"
    url  = "https://example.com"
  }
  publish = true
}
