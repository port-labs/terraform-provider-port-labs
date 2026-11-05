terraform {
  required_providers {
    port = {
      source  = "port-labs/port-labs"
      version = "~> 2"
    }
  }
}

provider "port" {
  # Set credentials via environment variables:
  # PORT_CLIENT_ID, PORT_CLIENT_SECRET
  # Optionally set PORT_BASE_URL (defaults to https://api.getport.io)
}
