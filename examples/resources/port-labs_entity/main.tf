resource "port_entity" "microservice" {
  title     = "monolith"
  blueprint = "microservice_blueprint"
  relations = {
    "tfRelation" = {
      "title"  = "Test Relation"
      "target" = port_entity.prod_env.id
    }
  }
  properties {
    string_prop = {
      "microservice_name" = "golang_monolith"
    }
  }
}

resource "port_entity" "prod_env" {
  title     = "production"
  blueprint = "environments"
  properties {
    string_prop = {
      "name" = "production-env"
    }
  }
}
