data "port_datasource_entities" "integration_entities" {
  datasource_prefix = "port-ocean/github/"
  datasource_suffix = "/my-github-integration/resync"
}

output "integration_entity_identifiers" {
  value = [for entity in data.port_datasource_entities.integration_entities.entities : entity.identifier]
}
