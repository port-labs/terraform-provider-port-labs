resource "port_workflow" "deploy_service" {
  identifier  = "deploy-service"
  title       = "Deploy service"
  description = "Collects deployment inputs, asks for approval and deploys"
  category    = "engineering"

  node {
    identifier = "trigger"
    title      = "Deploy request"

    self_serve_trigger {
      action_card_button_text    = "Deploy"
      execute_action_button_text = "Deploy"

      user_inputs {
        user_properties = {
          string_props = {
            "service" = {
              title    = "Service"
              required = true
            }
          }
          number_props = {
            "min_replicas" = {
              title   = "Minimum replicas"
              default = 1
            }
            "max_replicas" = {
              title   = "Maximum replicas"
              default = 3
            }
          }
        }

        # Checked when the form is submitted. Move the rules into the individual
        # steps when the form is split with `steps`.
        validations = [
          {
            constraint = ".form.max_replicas >= .form.min_replicas"
            message    = "Maximum replicas must be greater than or equal to minimum replicas"
          },
        ]
      }

      permissions {
        roles = ["Member"]
      }
    }
  }

  node {
    identifier = "approval"

    input {
      description = "Approve this deployment?"

      user_inputs {
        buttons = [
          {
            identifier = "approve"
            label      = "Approve"
            variant    = "PRIMARY"
          },
          {
            identifier = "reject"
            label      = "Reject"
            variant    = "DANGER"
          },
        ]
      }

      outlets {
        identifier        = "approve"
        title             = "Approved"
        num_of_responders = 1
      }

      outlets {
        identifier        = "reject"
        title             = "Rejected"
        num_of_responders = 1
      }

      notifications {
        target = "slack"
      }

      responders {
        roles = ["Admin"]
      }
    }
  }

  node {
    identifier = "deploy"
    verbose    = true
    links      = ["https://ci.example.com/runs/{{ .result.runId }}"]

    webhook {
      url    = "https://ci.example.com/deploy"
      method = "POST"
      body = jsonencode({
        service      = "{{ .outputs.trigger.inputs.service }}"
        min_replicas = "{{ .outputs.trigger.inputs.min_replicas }}"
        max_replicas = "{{ .outputs.trigger.inputs.max_replicas }}"
      })
    }
  }

  node {
    identifier = "record"

    upsert_entity {
      blueprint_identifier = "service"
      mapping {
        identifier = "{{ .outputs.trigger.inputs.service }}"
        title      = "{{ .outputs.trigger.inputs.service }}"
        properties = jsonencode({ last_deployed_at = "{{ .run.completedAt }}" })
      }
    }
  }

  connections {
    source_identifier = "trigger"
    target_identifier = "approval"
  }

  connections {
    source_identifier        = "approval"
    target_identifier        = "deploy"
    source_outlet_identifier = "approve"
  }

  connections {
    source_identifier = "deploy"
    target_identifier = "record"
  }
}

# An event driven workflow that branches on a JQ condition.
resource "port_workflow" "audit_service_changes" {
  identifier = "audit-service-changes"
  title      = "Audit service changes"

  node {
    identifier = "trigger"

    event_trigger {
      type                 = "ENTITY_UPDATED"
      blueprint_identifier = "service"

      condition {
        expressions = [".diff.after.properties.tier == \"production\""]
        combinator  = "and"
      }
    }
  }

  node {
    identifier = "branch"

    condition {
      outlets {
        identifier = "owner_missing"
        title      = "Owner missing"
        expression = ".outputs.trigger.diff.after.properties.owner == null"

        status_label {
          text    = "Missing owner"
          variant = "alert"
        }
      }
    }
  }

  node {
    identifier = "alert"

    webhook {
      url        = "https://alerts.example.com/service-owner-missing"
      on_failure = "continue"
    }
  }

  node {
    identifier = "summarize"

    ai {
      user_prompt   = "Summarize the change to {{ .outputs.trigger.diff.after.identifier }}"
      system_prompt = "You are a concise release auditor."
      tools         = ["get_.*"]
    }
  }

  connections {
    source_identifier = "trigger"
    target_identifier = "branch"
  }

  connections {
    source_identifier        = "branch"
    target_identifier        = "alert"
    source_outlet_identifier = "owner_missing"
  }

  connections {
    source_identifier = "branch"
    target_identifier = "summarize"
    fallback          = true
  }
}
