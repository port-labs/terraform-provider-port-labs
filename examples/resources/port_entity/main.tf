resource "port-labs_entity" "microservice" {
  title     = "monolith"
  blueprint = "microservice_blueprint"
  relations {
    name       = "tf-relation"
    identifier = port-labs_entity.prod_env.id
  }
  properties {
    name  = "microservice_name"
    value = "golang_monolith"
    type  = "string"
  }
}

resource "port-labs_entity" "prod_env" {
  title     = "production"
  blueprint = "environments"
  properties {
    name  = "name"
    value = "production-env"
    type  = "string"
  }
}
