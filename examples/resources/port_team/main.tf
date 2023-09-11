resource "port_team" "example" {
  name        = "example"
  description = "example"
  users = [
    "user1@test.com",
    "user2@test.com",
    "user3@test.com",
  ]
}
