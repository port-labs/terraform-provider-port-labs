resource "port_workflow" "review_pr" {
  identifier = "review-pr"
  title      = "Review PR"
  category   = "engineering"

  node {
    identifier = "trigger"
    title      = "PR Trigger"
    event_trigger {
      type                 = "ENTITY_UPDATED"
      blueprint_identifier = "githubPullRequest"
    }
  }

  node {
    identifier = "decision"
    title      = "Review"
    cursor_agent {
      api_key = "{{ .secrets[\"CURSOR_API_KEY\"] }}"
      prompt {
        text = "Review this PR."
      }
      source {
        pr_url = "{{ .outputs.trigger.diff.after.properties.link }}"
      }
    }
  }

  connections {
    source_identifier = "trigger"
    target_identifier = "decision"
  }
}
