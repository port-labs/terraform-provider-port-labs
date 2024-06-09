resource "port_blueprint" "environment" {
  title      = "Environment"
  icon       = "Environment"
  identifier = "examples-entity-env"
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

resource "port_blueprint" "microservice" {
  identifier  = "examples-entity-srvc"
  title       = "Microsvc from Port TF Examples"
  icon        = "Terraform"
  description = ""
  properties = {
    string_props = {
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
  }

}


resource "port_entity" "microservice" {
  title     = "monolith"
  blueprint = port_blueprint.microservice.identifier
  relations = {
    "tfRelation" = {
      "title"  = "Test Relation"
      "target" = port_blueprint.environment.identifier
    }
  }
  properties = {
    string_props = {
      "microservice_name" = "golang_monolith"
    }
  }
}

resource "port_entity" "prod_env" {
  title     = "production"
  blueprint = port_blueprint.environment.identifier
  properties = {
    string_props = {
      "name" = "production-env"
    }
  }
}
