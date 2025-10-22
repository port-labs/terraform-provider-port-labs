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
    titles = {
      "titleIdentifier" = {
        title = "My String Title"
        description = "My String Description",
        visible_jq_query = "true"
      }
    }
    order_properties = ["titleIdentifier","webhook_url","service","testString","testNumber"]
    user_properties = {
      string_props = {
        "webhook_url" = {
          title       = "Webhook URL"
          description = "Webhook URL to send the request to"
          format      = "url"
          default     = "https://example.com"
          pattern     = "^https://.*"
          disabled    = true
        }
        service = {
          title             = "Service"
          description       = "The service to restart"
          format            = "entity"
          blueprint         = port_blueprint.microservice.identifier
          disabled_jq_query = "1 == 1"

          sort = {
            property = "$updatedAt"
            order    = "DESC"
          }
        }
        testString = {
          type        = "string"
          title       = "String enum"
          icon        = "Terraform"
          default     = "a"
          enum        = ["a","b"]
          enum_colors = {
            a  ="darkGray"
            b = "turquoise"
          }
        }
        testNumber = {
          type        = "number"
          title       = "Number enum"
          icon        = "Terraform"
          default     = 1
          enum        = [1, 2]
          enum_colors = {
            "1"  ="darkGray"
            "2" = "turquoise"
          }
        }
      }
    }
  }
  webhook_method = {
    type = "WEBHOOK"
    url  = "https://app.getport.io"
  }
}

resource "port_action" "restart_microservice_with_steps" {
  title      = "Restart Microservice With Steps"
  icon       = "Terraform"
  identifier = "examples-action-restart-microservice-with-steps"
  publish    = true
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.environment.identifier
    title                = "Restart Microservice Workflow"
    user_properties = {
      string_props = {
        service_name = {
          type  = "string"
          title = "Service Name"
        }
        restart_reason = {
          type  = "string"
          title = "Restart Reason"
        }
        advanced_mode = {
          type  = "string"
          title = "Advanced Options"
        }
        confirm_restart = {
          type  = "boolean"
          title = "Confirm Restart"
        }
      }
    }
    steps = [
      {
        title = "Basic Information"
        order = ["service_name", "restart_reason"]
      },
      {
        title = "Advanced Settings"
        order = ["advanced_mode"]
        visible_jq_query = "1==1"
      },
      {
        title = "Confirmation"
        order = ["confirm_restart"]
        visible = true
      }
    ]
  }
  webhook_method = {
    type = "WEBHOOK"
    url  = "https://api.example.com/restart"
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
