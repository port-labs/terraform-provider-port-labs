# Example: Create a team with members
resource "port_team" "engineering" {
  name        = "Engineering"
  description = "The engineering team"
  users = [
    "developer1@example.com",
    "developer2@example.com",
  ]
}

# Example: Use team identifier to assign entities to a team
resource "port_blueprint" "service" {
  identifier = "service"
  title      = "Service"
  ownership = {
    type = "Direct"
  }
}

resource "port_entity" "my_service" {
  blueprint = port_blueprint.service.identifier
  title     = "My Service"
  teams = [
    port_team.engineering.identifier,
  ]
}
