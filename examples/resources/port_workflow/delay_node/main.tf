resource "port_workflow" "delayed_action" {
  identifier  = "delayed-action"
  title       = "Workflow with Delay"
  description = "Pause the workflow for a fixed number of seconds before continuing (DELAY node)"

  nodes = jsonencode([
    {
      identifier = "trigger"
      title      = "Start"
      config = {
        type = "SELF_SERVE_TRIGGER"
        userInputs = {
          properties = {
            wait_seconds = {
              type    = "number"
              title   = "Wait Duration (seconds)"
              minimum = 1
              maximum = 86400
            }
          }
          required = ["wait_seconds"]
        }
      }
    },
    {
      identifier = "wait"
      title      = "Wait Before Next Step"
      config = {
        type    = "DELAY"
        seconds = "{{ .outputs.trigger.wait_seconds }}"
      }
    },
    {
      identifier = "next_step"
      title      = "Continue After Delay"
      config = {
        type = "WEBHOOK"
        url  = "https://example.com/webhook"
      }
    }
  ])

  connections = jsonencode([
    {
      sourceIdentifier = "trigger"
      targetIdentifier = "wait"
    },
    {
      sourceIdentifier = "wait"
      targetIdentifier = "next_step"
    }
  ])
}
