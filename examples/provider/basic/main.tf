provider "port" {}

resource "port_entity" "microservice" {
  title = "monolith"
  blueprint = "microservice_blueprint"
  properties {
    name = "microservice_name"
    value = "golang_monolith"
    type = "string"
  }
}