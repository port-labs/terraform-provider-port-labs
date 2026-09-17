terraform {
  required_providers {
    port = {
      source  = "port-labs/port-labs"
      version = "~> 2"
    }
  }
}

provider "port" {
  client_id = "" # or set the environment variable PORT_CLIENT_ID
  secret    = "" # or set the environment variable PORT_CLIENT_SECRET
  base_url  = "https://api.getport.io"
}

data "port_scorecard_groups" "all" {}

output "scorecard_group_identifiers" {
  value = [for group in data.port_scorecard_groups.all.groups : group.identifier]
}
