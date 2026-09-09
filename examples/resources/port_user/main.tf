resource "port_user" "example" {
  email              = "user@example.com"
  roles              = ["Member"]
  teams              = ["engineering"]
  inactivity_timeout = 30
}
