resource "port_blueprint" "service" {
  title      = "Service"
  icon       = "Microservice"
  identifier = "service"

  properties = {
    array_props = {
      vulnerabilities = {
        title              = "Vulnerabilities"
        union              = true
        include_duplicates = false
        string_items       = {}
      }
    }
  }
}

resource "port_entity" "service" {
  title     = "Example Service"
  blueprint = port_blueprint.service.identifier

  properties = {
    union_array_props = {
      string_items = {
        vulnerabilities = {
          source_key = "terraform-sync"
          items      = ["CVE-2024-1", "CVE-2024-2"]
        }
      }
    }
  }
}
