resource "port_workflow" "tagged_workflow" {
  identifier  = "tagged-workflow"
  title       = "Tagged Workflow"
  description = "Example workflow demonstrating free-form tags for filtering and grouping"
  tags        = ["terraform", "example", "lorekeeper"]

  nodes = jsonencode([
    {
      identifier = "trigger"
      title      = "Manual Trigger"
      config = {
        type = "SELF_SERVE_TRIGGER"
      }
    }
  ])

  connections = jsonencode([])
}
