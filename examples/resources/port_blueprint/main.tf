resource "port_blueprint" "environment" {
  title      = "Environment"
  icon       = "Environment"
  identifier = "hedwig-env"
  properties = {
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
  properties = {
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

resource "port_blueprint" "microservice" {
  title      = "VM"
  icon       = "GPU"
  identifier = "hedwig-microservice"
  properties = {
    array_props = {
      enum = {
        type  = "array"
        title = "enum"
        string_items = {
          enum        = ["default", "default2"]
          enum_colors = {
            default  = "red"
            default2 = "green"
          }
        }
      }
    }
    string_props = {
      name = {
        type  = "string"
        title = "Name"
      },
      author = {
        type  = "string"
        title = "Author"
      },
      url = {
        type  = "string"
        title = "URL"
      },
      myStringIdentifier = {
        description = "This is a string property"
        title       = "text"
        icon        = "Terraform"
        required    = true
        min_length  = 1
        max_length  = 10
        default     = "default"
        enum        = ["default", "default2"]
        pattern     = "^[a-zA-Z0-9]*$"
        format      = "user"
        enum_colors = {
          default  = "red"
          default2 = "green"
        }
      }
    }
    boolean_props = {
      required = {
        type = "boolean"
      }
    }
    number_props = {
      sum = {
        type = "number"
      }
    }
  }

  calculation_properties = {
    "spec-calculation" = {
      title       = "Spec Calculation"
      calculation = "@blueprint/environment.properties.name"
      type        = "string"
      spec        = "embedded-url"
      spec_authentication = {
        client_id         = "my-client-id"
        token_url         = "https://example.com/token"
        authorization_url = "https://example.com/auth"
      }
    }
    "simple-calculation" = {
      title       = "Simple Calculation"
      calculation = "1 + 1"
      type        = "number"
      description = "A simple calculation without spec"
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

resource "port_blueprint" "repository_blueprint" {
  title       = "Repository Blueprint"
  icon        = "Terraform"
  identifier  = "repository"
  description = ""
}

resource "port_blueprint" "pull_request_blueprint" {
  title       = "Pull Request Blueprint"
  icon        = "Terraform"
  identifier  = "pull_request"
  description = ""
  properties = {
    string_props = {
      "status" = {
        title = "Status"
      }
    }
  }
  relations = {
    "self" = {
      title  = "Self"
      target = "pull_request"
    }
  }
}
