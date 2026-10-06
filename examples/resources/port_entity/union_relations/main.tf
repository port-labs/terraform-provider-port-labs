resource "port_blueprint" "dependency" {
  title      = "Dependency"
  identifier = "dependency"
  icon       = "Package"
}

resource "port_blueprint" "service" {
  title      = "Service"
  identifier = "service"
  icon       = "Microservice"

  relations = {
    dependencies = {
      title  = "Dependencies"
      target = port_blueprint.dependency.identifier
      many   = true
      union  = true
    }
  }
}

resource "port_entity" "dependency_a" {
  blueprint  = port_blueprint.dependency.identifier
  identifier = "dependency-a"
  title      = "Dependency A"
}

resource "port_entity" "service_a" {
  blueprint  = port_blueprint.service.identifier
  identifier = "service-a"
  title      = "Service A"

  relations = {
    union_many_relations = {
      dependencies = {
        source_key = "terraform"
        items      = [port_entity.dependency_a.identifier]
      }
    }
  }
}
