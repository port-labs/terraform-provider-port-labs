resource "port_blueprint" "microservice" {
  title      = "VM"
  icon       = "GPU"
  identifier = "examples-scorecard-svc"
  properties = {
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
      },
      replicaCount = {
        type = "number"
      }
    }
  }
}

resource "port_system_blueprint" "scorecard" {
  identifier = "_scorecard"
  properties = {
    string_props = {
      owner = {
        type  = "string"
        title = "Owner"
      }
    }
    number_props = {
      priority = {
        type  = "number"
        title = "Priority"
      }
    }
  }
  relations = {
    owner_team = {
      title  = "Owner Team"
      target = "_team"
    }
  }
}

resource "port_scorecard" "production_readiness" {
  identifier = "production-readiness"
  title      = "Production Readiness"
  blueprint  = port_blueprint.microservice.identifier
  properties = jsonencode({
    owner    = "platform-team"
    priority = 1
  })
  relations = jsonencode({
    owner_team = "platform-team"
  })
  rules = [{
    identifier = "high-avalability"
    title      = "High Availability"
    level      = "Gold"
    query = {
      combinator = "and"
      conditions = [jsonencode({
        property = "replicaCount"
        operator = ">="
        value    = "4"
      })]
    }
  }]
  depends_on = [
    port_system_blueprint.scorecard,
  ]
}
