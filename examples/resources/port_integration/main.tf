resource "port_integration" "my_custom_integration" {
  installation_id = "my-custom-integration-id"
  title           = "My Custom Integration"
  config = jsonencode({
    createMissingRelatedEntitiesboolean = true
    deleteDependentEntities             = true
    resources = [{
      kind = "my-custom-kind"
      selector = {
        query = ".title"
      }
      port = {
        entity = {
          mappings = [{
            identifier = "'my-identifier'"
            title      = ".title"
            blueprint  = "'my-blueprint'"
            properties = {
              my_property = 123
            }
            relations = {}
          }]
        }
      }
    }]
  })
}

resource "port_integration" "my_k8s_exporter" {
  installation_id       = "my-k8s-exporter"
  title                 = "My K8S Exporter with version managed by Terraform"
  installation_app_type = "K8S EXPORTER"
  # NOTE: This property is by default not used, since it can change outside of terraform
  # Include this only if you explicitly want to control the version with Terraform
  version = "1.33.7"
  config = jsonencode({
    createMissingRelatedEntitiesboolean = true
    deleteDependentEntities             = true
    resources = [{
      kind = "apps/v1/replicasets"
      selector = {
        query = ".metadata.namespace | startswith(\"kube\") | not"
      }
      port = {
        entity = {
          mappings = [{
            identifier = ".metadata.name"
            title      = ".metadata.name"
            blueprint  = "'deploymentConfig'"
            properties = {
              creationTimestamp = ".metadata.creationTimestamp"
              annotations       = ".metadata.annotations"
              status            = ".status"
            }
          }]
        }
      }
    }]
  })
}

# Port-hosted (SaaS) integration.
# Sensitive integrationSpec fields must reference a Port organization secret by name.
resource "port_organization_secret" "gitlab_token_mapping" {
  secret_name  = "my-gitlab-token-mapping"
  secret_value = jsonencode({ "glpat-example" : ["**"] })
  description  = "GitLab tokens and the group scopes they may ingest"
}

resource "port_integration" "my_gitlab_saas" {
  installation_id       = "my-gitlab-saas"
  installation_type     = "Saas"
  installation_app_type = "gitlab"
  title                 = "My GitLab (Port-hosted)"
  # version is intentionally omitted: Port manages it once the integration is provisioned

  spec = jsonencode({
    integrationSpec = {
      gitlabHost   = "https://gitlab.com"
      tokenMapping = port_organization_secret.gitlab_token_mapping.secret_name
    }
    appSpec = {
      scheduledResyncInterval = "2h"
    }
  })

  config = jsonencode({
    createMissingRelatedEntities = true
    deleteDependentEntities      = true
    resources = [{
      kind = "project"
      selector = {
        query = "true"
      }
      port = {
        entity = {
          mappings = [{
            identifier = ".path_with_namespace"
            title      = ".name"
            blueprint  = "'gitlabProject'"
          }]
        }
      }
    }]
  })
}
