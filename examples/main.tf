terraform {
  required_providers {
    port-labs = {
      source  = "port-labs/port-labs"
      version = "~> 1.0.0"
    }
  }
}
provider "port-labs" {
  client_id = "{YOUR CLIENT ID}"     # or set the environment variable PORT_CLIENT_ID
  secret    = "{YOUR CLIENT SECRET}" # or set the environment variable PORT_CLIENT_SECRET
}

resource "port_entity" "microservice" {
  title     = "monolith"
  blueprint = "microservice_blueprint"
  properties {
    string_prop = {
      "microservice_name" = "golang_monolith"
    }
  }
}
