resource "port_integration" "my_custom_integration" {
  installation_id       = "my-custom-integration-id"
  title                 = "My Custom Integration"
  version               = "1.33.7"
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
