resource "port-labs_blueprint" "environment" {
  title      = "Environment"
  icon       = "Environment"
  identifier = "hedwig-env"
  properties {
    identifier = "name"
    type       = "string"
    title      = "name"
  }
  properties {
    identifier = "docs-url"
    type       = "string"
    title      = "Docs URL"
    format     = "url"
  }
}

resource "port-labs_blueprint" "vm" {
  title      = "VM"
  icon       = "GPU"
  identifier = "hedwig-vm"
  properties {
    identifier = "name"
    type       = "string"
    title      = "Name"
  }
  relations {
    identifier = "environment"
    title      = "Test Relation"
    required   = "true"
    target     = port-labs_blueprint.environment.identifier
  }
}
