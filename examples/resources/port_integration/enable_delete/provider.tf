terraform {
  required_providers {
    port = {
      source = "port-labs/port-labs"
    }
  }
}

provider "port" {
  # client_id and client_secret are read from PORT_CLIENT_ID and PORT_CLIENT_SECRET
}
