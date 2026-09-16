variable "github_token" {
  type      = string
  sensitive = true
}

locals {
  installation_id = "github-prod"
  integration_type = "github-ocean"
}

resource "port_organization_secret" "github_token" {
  secret_name  = "_${upper(replace("${local.installation_id}_${local.integration_type}_github_token", "-", "_"))}"
  secret_value = var.github_token
}

resource "port_integration" "github" {
  depends_on = [port_organization_secret.github_token]

  installation_id       = local.installation_id
  installation_app_type = local.integration_type
  installation_type     = "Saas"
  title                 = "GitHub Production (PAT)"

  spec = jsonencode({
    integrationSpec = {
      authenticationMode = "Personal Access Token"
      githubToken        = port_organization_secret.github_token.secret_name
    }
    appSpec = {
      scheduledResyncInterval  = "12h"
      sendRawDataExamples      = true
      liveEventsEnabled        = false
      actionsProcessingEnabled = false
      incrementalSyncEnabled   = false
    }
  })
}
