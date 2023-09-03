resource "port_scorecard" "production_readiness" {
  identifier = "production-readiness"
  title      = "Production Readiness"
  blueprint  = "microservice"
  rules = [{
    identifier = "high-avalability"
    title      = "High Availability"
    level      = "Gold"
    query = {
      combinator = "and"
      conditions = [{
        property = "replicaCount"
        operator = ">="
        "value"  = "4"
      }]
    }
  }]
}
