resource "port_integration" "k8s_exporter_with_selective_deletes" {
  installation_id       = "my-k8s-exporter-selective-deletes"
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
        # enableDelete omitted — reconciliation deletes remain enabled for this resource
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
