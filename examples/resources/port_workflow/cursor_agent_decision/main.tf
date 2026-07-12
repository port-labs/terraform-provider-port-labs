resource "port_workflow" "lorekeeper_decision" {
  identifier  = "lorekeeper-decision"
  title       = "LoreKeeper Decision"
  description = "Review a PR using a Cursor cloud agent (CURSOR_AGENT node)"
  tags        = ["lorekeeper", "cursor-agent"]

  nodes = jsonencode([
    {
      identifier = "trigger"
      title      = "PR Updated"
      config = {
        type = "EVENT_TRIGGER"
        event = {
          type                = "ENTITY_UPDATED"
          blueprintIdentifier = "githubPullRequest"
        }
      }
    },
    {
      identifier = "decision"
      title      = "LoreKeeper Decision"
      config = {
        type   = "CURSOR_AGENT"
        apiKey = "{{ .secrets[\"CURSOR_API_KEY\"] }}"
        prompt = {
          text = <<-EOT
            Review the PR changes and decide whether companion docs, skills, or terraform PRs are required.
            Changed files: {{ .outputs.trigger.diff.after.properties.changedFiles | tostring }}
          EOT
        }
        source = {
          prUrl = "{{ .outputs.trigger.diff.after.properties.link }}"
        }
      }
    }
  ])

  connections = jsonencode([
    {
      sourceIdentifier = "trigger"
      targetIdentifier = "decision"
    }
  ])
}
