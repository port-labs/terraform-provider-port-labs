resource "port_team" "Example-Team" {
  name        = "Example Team"
  description = "An example team"
  # Note, this will need real users to work!
  users = [
    "user1@test.com",
    "user2@test.com",
    "user3@test.com",
  ]
}

resource "port_blueprint" "some-blueprint" {
  identifier = "some-blueprint"
  title      = "Some Blueprint"
  ownership = {
    type = "Direct"
  }
}

resource "port_entity" "some-entity" {
  blueprint = port_blueprint.some-blueprint.identifier
  title = "Some Entity"
  teams = [
    port_team.Example-Team.identifier,
  ]
}