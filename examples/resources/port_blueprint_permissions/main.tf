resource "port_blueprint" "environment" {
  title      = "Environment"
  icon       = "Environment"
  identifier = "examples-blueprint-perms-env"
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
  identifier  = "examples-blueprint-perms-srvc"
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
          "${port_blueprint.microservice.identifier}-moderator",
        ],
        "users" = [],
        "teams" = [],
        "ownedByTeam" : false
      },
      "identifier" = {
        "roles" = [
          "Admin",
          "Member",
          "${port_blueprint.microservice.identifier}-moderator",
        ],
        "users" = [],
        "teams" = [],
        "ownedByTeam" : false
      },
      "title" = {
        "roles" = [
          "Admin",
          "Member",
          "${port_blueprint.microservice.identifier}-moderator",
        ],
        "users" = [],
        "teams" = [],
        "ownedByTeam" : false
      },
      "team" = {
        "roles" = [
          "Admin",
          "Member",
          "${port_blueprint.microservice.identifier}-moderator",
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
          "${port_blueprint.microservice.identifier}-moderator",
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
          "${port_blueprint.microservice.identifier}-moderator",
        ],
        "users" = [],
        "teams" = [],
        "ownedByTeam" : false,
      }
    }
  }
}
