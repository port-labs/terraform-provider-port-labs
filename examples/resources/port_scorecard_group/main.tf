resource "port_blueprint" "microservice" {
  title      = "VM"
  icon       = "GPU"
  identifier = "examples-scorecard-group-svc"
  ownership = {
    type = "Direct"
  }
  properties = {
    string_props = {
      author = {
        type  = "string"
        title = "Author"
      }
    }
  }
}

resource "port_blueprint" "database" {
  title      = "Database"
  icon       = "Database"
  identifier = "examples-scorecard-group-db"
  ownership = {
    type = "Direct"
  }
  properties = {
    string_props = {
      author = {
        type  = "string"
        title = "Author"
      }
    }
  }
}

resource "port_team" "platform" {
  name        = "Scorecard Group Platform"
  description = "Team for scorecard group example entities"
  users       = ["test-admin-user@test.com"]
}

resource "port_entity" "vm_alpha" {
  title     = "VM Alpha"
  blueprint = port_blueprint.microservice.identifier
  teams     = [port_team.platform.identifier]
  properties = {
    string_props = {
      author = "Platform Team"
    }
  }
}

resource "port_entity" "vm_unowned" {
  title     = "VM Unowned"
  blueprint = port_blueprint.microservice.identifier
  properties = {
    string_props = {
      author = "Nobody"
    }
  }
}

resource "port_entity" "db_primary" {
  title     = "Primary DB"
  blueprint = port_blueprint.database.identifier
  teams     = [port_team.platform.identifier]
  properties = {
    string_props = {
      author = "Platform Team"
    }
  }
}

resource "port_scorecard_group" "production_readiness" {
  identifier = "production-readiness1"
  title      = "Production Readiness"
  blueprints = [
    port_blueprint.microservice.identifier,
    port_blueprint.database.identifier,
  ]
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
  depends_on = [
    port_blueprint.microservice,
    port_blueprint.database,
    port_entity.vm_alpha,
    port_entity.vm_unowned,
    port_entity.db_primary,
  ]
}
