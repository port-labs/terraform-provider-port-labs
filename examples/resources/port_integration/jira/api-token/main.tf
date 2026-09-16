variable "jira_token" {
  type      = string
  sensitive = true
}

variable "jira_email" {
  type      = string
  sensitive = true
}

locals {
  installation_id  = "jira-prod"
  integration_type = "jira"
}

resource "port_organization_secret" "jira_token" {
  secret_name  = "jira_token"
  secret_value = var.jira_token
}

resource "port_organization_secret" "jira_email" {
  secret_name  = "jira_email"
  secret_value = var.jira_email
}

resource "port_integration" "jira" {
  depends_on = [
    port_organization_secret.jira_token,
    port_organization_secret.jira_email,
  ]

  installation_id       = local.installation_id
  installation_app_type = local.integration_type
  installation_type     = "Saas"
  title                 = "Jira Production"

  spec = jsonencode({
    integrationSpec = {
      jiraHost                = "https://example.atlassian.net"
      atlassianUserEmail      = port_organization_secret.jira_email.secret_name
      atlassianUserToken      = port_organization_secret.jira_token.secret_name
      atlassianOrganizationId = "00000000-0000-0000-0000-000000000000"
    }
    appSpec = {
      scheduledResyncInterval = "12h"
      liveEventsEnabled       = false
    }
  })
}
