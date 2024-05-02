resource "port_blueprint" "environment" {
  title      = "Env from Port TF examples"
  icon       = "Environment"
  identifier = "fenrir-env"
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
  identifier  = "fenrir-microservice"
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
  relations = {
    "environment" = {
      title    = "Test Relation"
      required = "true"
      target   = port_blueprint.environment.identifier
    }
  }
}

resource "port_blueprint_permissions" "microservice_permissions" {
  blueprint_identifier = port_blueprint.microservice.identifier
  entities = {
    "register" = {
      "roles" : [
        "Admin",
        "Member",
      ],
      "users" : [],
      "teams" : []
    },
    "unregister" = {
      "roles" : [
        "Admin",
        "Member",
      ],
      "users" : [],
      "teams" : [],
    },
    "update" = {
      "roles" : [
        "Admin",
        "Member",
      ],
      "users" : [
      ],
      "teams" : []
    },
    "update_metadata_properties" = {
      "icon" = {
        "roles" = [
          "Admin",
          "fenrir-microservice-moderator",
        ],
        "users" = [],
        "teams" = [],
        "ownedByTeam" : false
      },
      "identifier" = {
        "roles" = [
          "Admin",
          "Member",
          "fenrir-microservice-moderator",
        ],
        "users" = [],
        "teams" = [],
        "ownedByTeam" : false
      },
      "title" = {
        "roles" = [
          "Admin",
          "Member",
          "fenrir-microservice-moderator",
        ],
        "users" = [],
        "teams" = [],
        "ownedByTeam" : false
      },
      "team" = {
        "roles" = [
          "Admin",
          "Member",
          "fenrir-microservice-moderator",
        ],
        "users" = [],
        "teams" = [],
        "ownedByTeam" : false
      },
    },
    "update_properties" = {
      "myStringIdentifier" : {
        "roles" = [
          "Admin",
          "Member",
          "fenrir-microservice-moderator",
        ],
        "users" = [],
        "teams" = [],
        "ownedByTeam" : false,
      }
    },
    "update_relations" = {
      "environment" = {
        "roles" = [
          "Admin",
          "Member",
          "fenrir-microservice-moderator",
        ],
        "users" = [],
        "teams" = [],
        "ownedByTeam" : false,
      }
    }
  }
}
