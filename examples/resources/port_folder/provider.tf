terraform {
  required_providers {
    port = {
      source  = "port-labs/port-labs"
      version = "~> 2.0.0"
    }
  }
}

provider "port" {
  client_id = "" # or set the environment variable PORT_CLIENT_ID
  secret    = "" # or set the environment variable PORT_CLIENT_SECRET
  base_url  = "https://api.getport.io"
}