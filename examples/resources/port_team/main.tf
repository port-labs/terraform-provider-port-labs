resource "port_team" "example" {
  name        = "example"
  description = "example"
  # Note, this will need real users to work!
  users = [
    "user1@test.com",
    "user2@test.com",
    "user3@test.com",
  ]
}
