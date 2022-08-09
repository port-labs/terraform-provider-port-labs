terraform {
  required_providers {
    port = {
      source  = "port-labs/port"
      version = "~> 0.0.1"
    }
  }
}
provider "port" {}

resource "port-labs_entity" "microservice" {
  title     = "monolith"
  blueprint = "microservice_blueprint"
  properties {
    name  = "microservice_name"
    value = "golang_monolith"
    type  = "string"
  }
}
