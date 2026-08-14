resource "port_blueprint" "service" {
  title      = "Service"
  icon       = "Microservice"
  identifier = "service"

  properties = {
    array_props = {
      vulnerabilities = {
        title               = "Vulnerabilities"
        union               = true
        include_duplicates  = false
        string_items        = {}
      }
    }
  }
}
