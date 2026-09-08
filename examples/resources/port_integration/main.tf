variable "github_token" {
  type      = string
  sensitive = true
}

# =============================================================================
# Self-hosted (OnPrem) — Step 1: Create without config
# =============================================================================

resource "port_integration" "my_custom_integration" {
  installation_id = "my-custom-integration-id"
  title           = "My Custom Integration"

  # config is set on the next apply, after provisioning creates default mappings.
}

# After the first apply, add the config block below and run 'terraform apply' again.
#
# resource "port_integration" "my_custom_integration" {
#   installation_id = "my-custom-integration-id"
#   title           = "My Custom Integration"
#
#   config = jsonencode({
#     createMissingRelatedEntities = true
#     deleteDependentEntities      = true
#     resources = [{
#       kind = "my-custom-kind"
#       selector = {
#         query = ".title"
#       }
#       port = {
#         entity = {
#           mappings = [{
#             identifier = "'my-identifier'"
#             title      = ".title"
#             blueprint  = "'my-blueprint'"
#             properties = {
#               my_property = 123
#             }
#             relations = {}
#           }]
#         }
#       }
#     }]
#   })
# }

# =============================================================================
# K8s exporter with version pinning
# =============================================================================

resource "port_integration" "my_k8s_exporter" {
  installation_id       = "my-k8s-exporter"
  title                 = "My K8S Exporter with version managed by Terraform"
  installation_app_type = "K8S EXPORTER"
  # NOTE: version can change outside of Terraform.
  # Include this only if you explicitly want to control the version with Terraform.
  version = "1.33.7"

  # config is set on the next apply, after provisioning creates default mappings.
}

# =============================================================================
# SaaS (Ocean) — Step 1: Create without config
# =============================================================================

# Secret naming convention: _{INSTALLATION_ID}_{INTEGRATION_TYPE}_{PROPERTY} in SCREAMING_SNAKE_CASE.
# This matches the frontend's generateSecretName() so Ocean can resolve the secret at runtime.
#
# Each integration type has its own spec fields — check the integration's .port/spec.json
# for the list of configurations (name, type, sensitive, dependencies, etc.).

locals {
  github_installation_id = "github-prod"
  github_type            = "github-ocean"
}

resource "port_organization_secret" "github_token" {
  secret_name  = "_${upper(replace("${local.github_installation_id}_${local.github_type}_github_token", "-", "_"))}"
  secret_value = var.github_token
}

resource "port_integration" "github" {
  depends_on = [port_organization_secret.github_token]

  installation_id       = local.github_installation_id
  installation_app_type = local.github_type
  installation_type     = "Saas"
  title                 = "GitHub Production"

  spec = jsonencode({
    integrationSpec = {
      authenticationMode = "Personal Access Token"
      githubToken        = port_organization_secret.github_token.secret_name
    }

    # appSpec controls feature toggles. If omitted, the server applies its own
    # defaults which may differ from the Port UI defaults. Declare explicitly
    # to match the UI behaviour.
    appSpec = {
      scheduledResyncInterval  = "12h"
      sendRawDataExamples      = true
      liveEventsEnabled        = false
      actionsProcessingEnabled = false
      incrementalSyncEnabled   = false
    }
  })

  # Do NOT set config here — provisioning will populate default mappings.
  # Add config on a subsequent apply to override them.
}

# After the first apply, add config to override mappings:
#
# resource "port_integration" "github" {
#   ...same as above...
#
#   config = jsonencode({
#     resources = [
#       {
#         kind = "repository"
#         selector = { query = "true" }
#         port = {
#           entity = {
#             mappings = {
#               identifier = ".name"
#               blueprint  = "\"githubRepository\""
#               title      = ".name"
#               properties = {
#                 url           = ".html_url"
#                 defaultBranch = ".default_branch"
#               }
#             }
#           }
#         }
#       }
#     ]
#   })
# }
