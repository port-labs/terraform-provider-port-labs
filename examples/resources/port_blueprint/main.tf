resource "port_blueprint" "environment" {
  title      = "Environment"
  icon       = "Environment"
  identifier = "hedwig-env"
  properties {
    string_props = {
      "name" = {
        type  = "string"
        title = "name"
      }
      "docs-url" = {
        title  = "Docs URL"
        format = "url"
      }
    }
  }
}

resource "port_blueprint" "vm" {
  title      = "VM"
  icon       = "GPU"
  identifier = "hedwig-vm"
  properties {
    string_props = {
      name = {
        type  = "string"
        title = "Name"
      }
    }
  }
  relations = {
    "environment" = {
      title    = "Test Relation"
      required = "true"
      target   = port_blueprint.environment.identifier
    }
  }
}
