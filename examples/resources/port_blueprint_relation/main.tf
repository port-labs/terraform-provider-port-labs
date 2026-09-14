# Two blueprints that relate to each other.
#
# The relations cannot be declared inline on the blueprints: Port requires a relation's target
# blueprint to already exist, so whichever blueprint is created first would be rejected. Neither
# blueprint here references the other, so both are created first and the relations follow.

resource "port_blueprint" "entra_id_user" {
  title      = "Entra ID User"
  identifier = "entra-id-user"
  icon       = "User"

  properties = {
    string_props = {
      "email" = {
        title = "Email"
      }
    }
  }
}

resource "port_blueprint" "entra_id_group" {
  title      = "Entra ID Group"
  identifier = "entra-id-group"
  icon       = "Group"

  properties = {
    string_props = {
      "displayName" = {
        title = "Display Name"
      }
    }
  }
}

resource "port_blueprint_relation" "user_groups" {
  blueprint  = port_blueprint.entra_id_user.identifier
  identifier = "groups"
  target     = port_blueprint.entra_id_group.identifier
  title      = "Groups"
  many       = true
}

resource "port_blueprint_relation" "group_members" {
  blueprint  = port_blueprint.entra_id_group.identifier
  identifier = "members"
  target     = port_blueprint.entra_id_user.identifier
  title      = "Members"
  many       = true
}
