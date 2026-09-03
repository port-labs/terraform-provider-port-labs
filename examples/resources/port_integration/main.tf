resource "port_integration" "my_onprem_integration" {
  installation_id       = "my-onprem-integration-id"
  installation_app_type = "github"
  version               = "1.0.0"
  title                 = "My OnPrem Integration"
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
  # NOTE: This property is by default not used, since it can change outside of terraform
  # Include this only if you explicitly want to control the version with Terraform
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
