resource "port_system_blueprint" "user" {
  identifier = "_user"
  properties = {
    string_props = {
      "age" = {
        type  = "number"
        title = "Age"
      }
    }
  }
}

resource "port_system_blueprint" "team" {
  identifier = "_team"
  relations = {
    "manager" = {
      title    = "Manager"
      required = "true"
      target   = port_system_blueprint.user.identifier
    }
  }
}
