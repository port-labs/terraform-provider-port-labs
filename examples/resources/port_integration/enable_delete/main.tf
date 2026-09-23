resource "port_integration" "k8s_with_selective_reconciliation_deletes" {
  installation_id       = "my-k8s-exporter-enable-delete"
  title                 = "K8S exporter with per-resource enableDelete"
  installation_app_type = "K8S EXPORTER"

  config = jsonencode({
    entityDeletionThreshold = 1
    resources = [
      {
        kind         = "namespace"
        enableDelete = false
        selector = {
          query = "true"
        }
        port = {
          entity = {
            mappings = [{
              identifier = ".metadata.uid"
              title      = ".metadata.name"
              blueprint  = "'namespace'"
            }]
          }
        }
      },
      {
        kind = "application"
        selector = {
          query = "true"
        }
        port = {
          entity = {
            mappings = [{
              identifier = ".metadata.uid"
              title      = ".metadata.name"
              blueprint  = "'argocdApplication'"
            }]
          }
        }
      },
    ]
  })
}
