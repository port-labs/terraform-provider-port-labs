# This example demonstrates how to create an automation action with an empty
# jq_condition expressions array. An empty expressions array means no filtering
# will be applied - the automation will trigger on any matching event.

resource "port_blueprint" "microservice" {
  title      = "Microservice"
  icon       = "Microservice"
  identifier = "examples-automation-empty-jq-microservice"
  properties = {
    string_props = {
      "name" = {
        title = "Name"
      }
      "environment" = {
        title = "Environment"
        enum  = ["production", "staging", "development"]
      }
    }
  }
}

# Automation with empty jq_condition expressions
# This will trigger on ANY entity change without filtering
resource "port_action" "notify_on_any_change" {
  title       = "Notify On Any Change"
  identifier  = "examples-automation-notify-on-any-change"
  icon        = "Notification"
  description = "Triggers on any entity change without condition filtering"
  publish     = true

  automation_trigger = {
    any_entity_change_event = {
      blueprint_identifier = port_blueprint.microservice.identifier
    }
    # Empty expressions array = no filtering, triggers on all events
    jq_condition = {
      expressions = []
      combinator  = "and"
    }
  }

  webhook_method = {
    url = "https://example.com/webhook"
  }
}

# For comparison: automation with actual jq_condition expressions
# This will only trigger when the environment is "production"
resource "port_action" "notify_on_production_change" {
  title       = "Notify On Production Change"
  identifier  = "examples-automation-notify-on-production-change"
  icon        = "Notification"
  description = "Triggers only when production entities change"
  publish     = true

  automation_trigger = {
    any_entity_change_event = {
      blueprint_identifier = port_blueprint.microservice.identifier
    }
    jq_condition = {
      expressions = [".diff.after.properties.environment == \"production\""]
      combinator  = "and"
    }
  }

  webhook_method = {
    url = "https://example.com/webhook"
  }
}

# Automation without jq_condition at all (equivalent to empty expressions)
resource "port_action" "notify_on_creation" {
  title       = "Notify On Creation"
  identifier  = "examples-automation-notify-on-creation"
  icon        = "Notification"
  description = "Triggers on entity creation without any condition"
  publish     = true

  automation_trigger = {
    entity_created_event = {
      blueprint_identifier = port_blueprint.microservice.identifier
    }
    # No jq_condition block - same effect as empty expressions
  }

  webhook_method = {
    url = "https://example.com/webhook"
  }
}
