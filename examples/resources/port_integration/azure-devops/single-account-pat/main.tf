variable "ado_pat" {
  type      = string
  sensitive = true
}

locals {
  installation_id  = "ado-single-prod"
  integration_type = "azure-devops"
}

resource "port_organization_secret" "ado_pat" {
  secret_name  = "_${upper(replace("${local.installation_id}_${local.integration_type}_personal_access_token", "-", "_"))}"
  secret_value = var.ado_pat
}

resource "port_integration" "ado" {
  depends_on = [port_organization_secret.ado_pat]

  installation_id       = local.installation_id
  installation_app_type = local.integration_type
  installation_type     = "Saas"
  title                 = "Azure DevOps (single account)"

  spec = jsonencode({
    integrationSpec = {
      accountMode         = "Single Account"
      organizationUrl     = "https://dev.azure.com/my-organization"
      personalAccessToken = port_organization_secret.ado_pat.secret_name
      isProjectsLimited   = true
    }
    appSpec = {
      scheduledResyncInterval = "12h"
      liveEventsEnabled       = false
    }
  })
}
