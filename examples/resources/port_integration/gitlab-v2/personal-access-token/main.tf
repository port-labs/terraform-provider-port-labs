variable "gitlab_token" {
  type      = string
  sensitive = true
}

locals {
  installation_id  = "gitlab-prod"
  integration_type = "gitlab-v2"
}

resource "port_organization_secret" "gitlab_token" {
  secret_name  = "_${upper(replace("${local.installation_id}_${local.integration_type}_gitlab_token", "-", "_"))}"
  secret_value = var.gitlab_token
}

resource "port_integration" "gitlab" {
  depends_on = [port_organization_secret.gitlab_token]

  installation_id        = local.installation_id
  installation_app_type  = local.integration_type
  installation_type     = "Saas"
  title                 = "GitLab Production"

  spec = jsonencode({
    integrationSpec = {
      gitlabToken = port_organization_secret.gitlab_token.secret_name
      gitlabHost  = "https://gitlab.com"
    }
    appSpec = {
      scheduledResyncInterval = "12h"
      liveEventsEnabled       = false
    }
  })
}
