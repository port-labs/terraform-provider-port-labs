variable "ado_client_secret" {
  type      = string
  sensitive = true
}

locals {
  installation_id  = "ado-multi-prod"
  integration_type = "azure-devops"
}

resource "port_organization_secret" "ado_client_secret" {
  secret_name  = "ado_client_secret"
  secret_value = var.ado_client_secret
}

resource "port_integration" "ado" {
  depends_on = [port_organization_secret.ado_client_secret]

  installation_id       = local.installation_id
  installation_app_type = local.integration_type
  installation_type     = "Saas"
  title                 = "Azure DevOps (multiple accounts)"

  spec = jsonencode({
    integrationSpec = {
      accountMode       = "Multiple Accounts"
      clientId          = "00000000-0000-0000-0000-000000000000"
      clientSecret      = port_organization_secret.ado_client_secret.secret_name
      tenantId          = "00000000-0000-0000-0000-000000000000"
      organizationUrls  = ["https://dev.azure.com/org-one", "https://dev.azure.com/org-two"]
      isProjectsLimited = true
    }
    appSpec = {
      scheduledResyncInterval = "12h"
      liveEventsEnabled       = false
    }
  })
}
