# Example: Create organization secrets for integration credentials
variable "github_token" {
  type      = string
  sensitive = true
}

variable "slack_webhook_url" {
  type      = string
  sensitive = true
}

resource "port_organization_secret" "github_token" {
  secret_name  = "github_token"
  secret_value = var.github_token
  description  = "GitHub personal access token for the GitHub integration"
}

resource "port_organization_secret" "slack_webhook" {
  secret_name  = "slack_webhook_url"
  secret_value = var.slack_webhook_url
  description  = "Slack webhook URL for notifications"
}

# Example: Use secrets in an integration
resource "port_integration" "github" {
  depends_on = [port_organization_secret.github_token]

  installation_id       = "github-prod"
  installation_app_type = "github-ocean"
  installation_type     = "Saas"
  title                 = "GitHub Production"

  spec = jsonencode({
    integrationSpec = {
      githubToken = port_organization_secret.github_token.secret_name
    }
  })
}
