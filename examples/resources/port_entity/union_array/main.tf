resource "port_blueprint" "service" {
  title      = "Service"
  icon       = "Terraform"
  identifier = "examples-entity-union-service"
  properties = {
    array_props = {
      tags = {
        title              = "Tags"
        union              = true
        include_duplicates = false
        string_items       = {}
      }
    }
  }
}

resource "port_entity" "service" {
  title                 = "example-service"
  blueprint             = port_blueprint.service.identifier
  union_array_source_key = "terraform"
  properties = {
    array_props = {
      string_items = {
        tags = ["alpha", "beta"]
      }
    }
  }
}
