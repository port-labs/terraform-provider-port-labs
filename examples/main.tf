terraform {
  required_providers {
    port-labs = {
      source  = "port-labs/port-labs"
      version = "~> 0.10.3"
    }
  }
}
provider "port-labs" {
  client_id = "{YOUR CLIENT ID}"     # or set the environment variable PORT_CLIENT_ID
  secret    = "{YOUR CLIENT SECRET}" # or set the environment variable PORT_CLIENT_SECRET
}

resource "port-labs_entity" "microservice" {
  title     = "monolith"
  blueprint = "microservice_blueprint"
  properties {
    name  = "microservice_name"
    value = "golang_monolith"
  }
}
