# Self-hosted — Step 1: Create without config.
# After the first apply, add config and run 'terraform apply' again to override default mappings.

resource "port_integration" "my_custom_integration" {
  installation_id = "my-custom-integration-id"
  title           = "My Custom Integration"
}

# Step 2 example (uncomment after first apply):
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
