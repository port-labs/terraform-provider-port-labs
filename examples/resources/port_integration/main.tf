resource "port_integration" "my_custom_integration" {
  installation_id       = "my-custom-integration-id"
  title                 = "My Custom Integration"
  installation_app_type = "WEBHOOK"
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
  installation_id       = "my-"
  title                 = "My K8S Exporter with version managed by Terraform"
  installation_app_type = "K8S EXPORTER"
  # Optional: pin the integration version in Terraform. Omit to ignore version changes made outside Terraform.
  version               = "1.33.7"
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
              annotations = ".metadata.annotations"
              status = ".status"
            }
          }]
        }
      }
    }]
  })
}
