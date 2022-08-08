terraform {
  required_providers {
    port = {
      source  = "port-labs/port"
      version = "~> 0.0.1"
    }
  }
}
provider "port" {}

resource "port_entity" "microservice" {
  title = "monolith"
  blueprint = "microservice_blueprint"
  relations = {}
  properties {
    name = "microservice_name"
    value = "golang_monolith"
    type = "string"
  }
}