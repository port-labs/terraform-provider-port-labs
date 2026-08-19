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
      environment = {
        type  = "string"
        title = "Environment"
        enum  = ["production", "staging"]
      }
    }
  }
}

resource "port_blueprint" "database" {
  title      = "Database"
  icon       = "Database"
  identifier = "examples-scorecard-group-db"
  properties = {
    string_props = {
      author = {
        type  = "string"
        title = "Author"
      }
      environment = {
        type  = "string"
        title = "Environment"
        enum  = ["production", "staging"]
      }
    }
  }
}

resource "port_entity" "vm_alpha" {
  title     = "VM Alpha"
  blueprint = port_blueprint.microservice.identifier
  properties = {
    string_props = {
      author      = "Platform Team"
      environment = "production"
    }
  }
}

resource "port_entity" "vm_unowned" {
  title     = "VM Unowned"
  blueprint = port_blueprint.microservice.identifier
  properties = {
    string_props = {
      author      = "Nobody"
      environment = "staging"
    }
  }
}

resource "port_entity" "db_primary" {
  title     = "Primary DB"
  blueprint = port_blueprint.database.identifier
  properties = {
    string_props = {
      author      = "Platform Team"
      environment = "production"
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
    identifier = "has-platform-author"
    title      = "Has Platform Author"
    level      = "Gold"
    query = {
      combinator = "and"
      conditions = [jsonencode({
        property = "author"
        operator = "="
        value    = "Platform Team"
      })]
    }
  }]
  filters = {
    (port_blueprint.microservice.identifier) = {
      combinator = "and"
      conditions = [jsonencode({
        property = "environment"
        operator = "="
        value    = "production"
      })]
    }
    (port_blueprint.database.identifier) = {
      combinator = "and"
      conditions = [jsonencode({
        property = "author"
        operator = "isNotEmpty"
      })]
    }
  }
  depends_on = [
    port_blueprint.microservice,
    port_blueprint.database,
    port_entity.vm_alpha,
    port_entity.vm_unowned,
    port_entity.db_primary,
  ]
}

# Per-blueprint mode: each blueprint gets its own rules (and optional filter).
resource "port_scorecard_group" "blueprint_specific_readiness" {
  identifier = "blueprint-specific-readiness"
  title      = "Blueprint-Specific Readiness"
  scorecards = {
    (port_blueprint.microservice.identifier) = {
      filter = {
        combinator = "and"
        conditions = [jsonencode({
          property = "environment"
          operator = "="
          value    = "staging"
        })]
      }
      rules = [{
        identifier = "has-platform-author"
        title      = "Has Platform Author"
        level      = "Gold"
        query = {
          combinator = "and"
          conditions = [jsonencode({
            property = "author"
            operator = "="
            value    = "Platform Team"
          })]
        }
      }]
    }
    (port_blueprint.database.identifier) = {
      filter = {
        combinator = "and"
        conditions = [jsonencode({
          property = "environment"
          operator = "="
          value    = "production"
        })]
      }
      rules = [{
        identifier = "has-author"
        title      = "Has Author"
        level      = "Gold"
        query = {
          combinator = "and"
          conditions = [jsonencode({
            property = "author"
            operator = "isNotEmpty"
          })]
        }
      }]
    }
  }
  depends_on = [
    port_blueprint.microservice,
    port_blueprint.database,
    port_entity.vm_alpha,
    port_entity.vm_unowned,
    port_entity.db_primary,
  ]
}
