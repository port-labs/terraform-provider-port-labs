terraform {
  required_providers {
    port = {
      source  = "port-labs/port-labs"
      version = "~> 0.0.1"
    }
  }
}
provider "port-labs" {}

resource "port-labs_entity" "microservice" {
  title     = "monolith"
  blueprint = "microservice_blueprint"
  properties {
    name  = "microservice_name"
    value = "golang_monolith"
    type  = "string"
  }
}
