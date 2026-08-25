resource "port_blueprint" "microservice" {
  title      = "VM"
  icon       = "GPU"
  identifier = "examples-scorecard-group-svc"
  properties = {
    string_props = {
      author = {
        type  = "string"
        title = "Author"
      }
    }
  }
}

resource "port_scorecard_group" "production_readiness" {
  identifier = "production-readiness"
  title      = "Production Readiness"
  properties = {
    owner = "platform-team"
  }
  blueprints = [port_blueprint.microservice.identifier]
  rules = [{
    identifier = "has-owner"
    title      = "Has Owner"
    level      = "Gold"
    query = {
      combinator = "and"
      conditions = [jsonencode({
        property = "$team"
        operator = "isNotEmpty"
      })]
    }
  }]
  depends_on = [port_blueprint.microservice]
}
