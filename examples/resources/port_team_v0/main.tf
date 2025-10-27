resource "port_team" "Example-Team" {
  name        = "Example Team"
  description = "An example teammmmmmm"
  # Note, this will need real users to work!
  users = [
  ]
}

# resource "port_blueprint" "some-blueprint" {
#   identifier = "some-blueprint"
#   title      = "Some Blueprint"
# }

# resource "port_entity" "some-entity" {
#   blueprint = port_blueprint.some-blueprint.identifier
#   title = "Some Entity"
#   teams = [
#     port_team.Example-Team.name,
#   ]
# }