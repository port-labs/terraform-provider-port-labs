terraform {
  required_providers {
    port = {
      source  = "port-labs/port-labs"
      version = "0.9.6"
    }
  }
}

provider "port" {
  base_url  = "https://api.getport.io"
}
