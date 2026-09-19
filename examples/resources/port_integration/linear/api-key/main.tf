variable "linear_api_key" {
  type      = string
  sensitive = true
}

locals {
  installation_id  = "linear-prod"
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
  title                 = "Linear Production"

  spec = jsonencode({
    integrationSpec = {
      linearApiKey = port_organization_secret.linear_api_key.secret_name
    }
    appSpec = {
      scheduledResyncInterval = "12h"
      liveEventsEnabled       = false
    }
  })
}
