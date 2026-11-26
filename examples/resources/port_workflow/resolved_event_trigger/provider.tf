terraform {
  required_providers {
    port = {
      source = "port-labs/port-labs"
    }
  }
}

provider "port" {
  # client_id     = var.port_client_id
  # client_secret = var.port_client_secret
}