terraform {
  required_providers {
    port-labs = {
      source  = "port-labs/port-labs"
      version = "~> 0.2.0"
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
  }
}
