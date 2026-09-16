variable "github_app_private_key" {
  type      = string
  sensitive = true
}

locals {
  installation_id  = "github-prod-app"
  integration_type = "github-ocean"
  github_app_id    = "123456"
}

resource "port_organization_secret" "github_app_private_key" {
  secret_name  = "github_app_private_key"
  secret_value = var.github_app_private_key
}

resource "port_integration" "github" {
  depends_on = [port_organization_secret.github_app_private_key]

  installation_id       = local.installation_id
  installation_app_type = local.integration_type
  installation_type     = "Saas"
  title                 = "GitHub Production (App)"

  spec = jsonencode({
    integrationSpec = {
      authenticationMode    = "Github App"
      githubAppId           = local.github_app_id
      githubAppPrivateKey   = port_organization_secret.github_app_private_key.secret_name
      githubOrganization    = "my-github-org"
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
