terraform {
  required_providers {
    port = {
      source  = "port-labs/port-labs"
      version = "~> 2.0"
    }
  }
}

provider "port" {
  # client_id and secret are read from PORT_CLIENT_ID and PORT_CLIENT_SECRET
}
