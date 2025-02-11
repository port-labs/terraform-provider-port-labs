resource "port_system_blueprint" "user" {
  identifier = "_user"
  properties = {
    number_props = {
      summ = {
        type = "number"
      }
    }
  }
}

resource "port_system_blueprint" "team" {
  identifier = "_team"
  relations = {
    "manager" = {
      title    = "Manager"
      required = "false"
      target   = port_system_blueprint.user.identifier
    }
  }
}
