variable "linear_api_key" {
  type      = string
  sensitive = true
}

locals {
  installation_id  = "linear-app-spec-toggles"
  integration_type = "linear"
}

resource "port_organization_secret" "linear_api_key" {
  secret_name  = "linear_api_key"
  secret_value = var.linear_api_key
}

resource "port_integration" "linear" {
  depends_on = [port_organization_secret.linear_api_key]

  installation_id       = local.installation_id
  installation_app_type = local.integration_type
  installation_type     = "Saas"
  title                 = "Linear (appSpec feature toggles)"

  spec = jsonencode({
    integrationSpec = {
      linearApiKey = port_organization_secret.linear_api_key.secret_name
    }
    appSpec = {
      scheduledResyncInterval  = "6h"
      liveEventsEnabled        = true
      actionsProcessingEnabled = true
      incrementalSyncEnabled   = true
      incrementalSyncInterval  = "30m"
      sendRawDataExamples      = true
    }
  })
}
