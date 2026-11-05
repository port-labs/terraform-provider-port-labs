# Example: Environment blueprint with string properties
resource "port_blueprint" "environment" {
  title      = "Environment"
  icon       = "Environment"
  identifier = "examples-entity-env"
  properties = {
    string_props = {
      "name" = {
        title = "Name"
      }
      "docs-url" = {
        title  = "Docs URL"
        format = "url"
      }
    }
  }
}

# Example: Microservice blueprint with various property configurations
resource "port_blueprint" "microservice" {
  identifier  = "examples-entity-srvc"
  title       = "Microservice"
  icon        = "Terraform"
  properties = {
    string_props = {
      "service_name" = {
        description = "The name of the microservice"
        title       = "Service Name"
        required    = true
      }
      "language" = {
        title   = "Language"
        enum    = ["go", "python", "javascript"]
        default = "go"
        enum_colors = {
          go         = "blue"
          python     = "green"
          javascript = "yellow"
        }
      }
    }
    number_props = {
      "version" = {
        title   = "Version"
        default = 1
      }
    }
    boolean_props = {
      "is_active" = {
        title   = "Is Active"
        default = true
      }
    }
  }
  relations = {
    "environment" = {
      title    = "Environment"
      target   = port_blueprint.environment.identifier
      many     = false
      required = false
    }
  }
}

# Example: Basic entity with properties
resource "port_entity" "microservice" {
  title     = "Payments Service"
  blueprint = port_blueprint.microservice.identifier
  properties = {
    string_props = {
      "service_name" = "payments"
      "language"     = "go"
    }
    number_props = {
      "version" = 2
    }
    boolean_props = {
      "is_active" = true
    }
  }
}

# Example: Entity with a single relation
resource "port_entity" "prod_env" {
  title     = "Production"
  blueprint = port_blueprint.environment.identifier
  properties = {
    string_props = {
      "name" = "production"
    }
  }
}

resource "port_entity" "microservice_with_relation" {
  title     = "Orders Service"
  blueprint = port_blueprint.microservice.identifier
  properties = {
    string_props = {
      "service_name" = "orders"
    }
  }
  relations = {
    single_relations = {
      "environment" = port_entity.prod_env.identifier
    }
  }
}

# Example: Entity with create_missing_related_entities
# This will create the related entity if it doesn't exist
resource "port_entity" "service_with_auto_create" {
  title     = "Notifications Service"
  blueprint = port_blueprint.microservice.identifier
  properties = {
    string_props = {
      "service_name" = "notifications"
    }
  }
  relations = {
    single_relations = {
      "environment" = "staging"
    }
  }
  create_missing_related_entities = true
}
