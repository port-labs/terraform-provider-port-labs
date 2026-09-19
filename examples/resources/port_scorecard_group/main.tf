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

resource "port_blueprint" "team" {
  title      = "Team"
  icon       = "Team"
  identifier = "examples-scorecard-group-team"
}

resource "port_entity" "platform_team" {
  identifier = "platform-team"
  title      = "Platform Team"
  blueprint  = port_blueprint.team.identifier
}

resource "port_system_blueprint" "scorecard_group" {
  identifier = "_scorecard_group"
  properties = {
    string_props = {
      category = {
        type  = "string"
        title = "Category"
      }
    }
  }
  relations = {
    owning_team = {
      title  = "Owning Team"
      target = port_blueprint.team.identifier
    }
  }
}

resource "port_system_blueprint" "scorecard" {
  identifier = "_scorecard"
  properties = {
    string_props = {
      owner = {
        type  = "string"
        title = "Owner"
      }
    }
    number_props = {
      priority = {
        type  = "number"
        title = "Priority"
      }
    }
  }
  relations = {
    owner_team = {
      title  = "Owner Team"
      target = port_blueprint.team.identifier
    }
  }
}

resource "port_scorecard_group" "production_readiness" {
  identifier = "production-readiness"
  title      = "Production Readiness"
  blueprints = [
    port_blueprint.microservice.identifier,
    port_blueprint.database.identifier,
  ]
  group_properties = jsonencode({
    category = "governance"
  })
  group_relations = jsonencode({
    owning_team = port_entity.platform_team.identifier
  })
  scorecard_properties = jsonencode({
    owner    = "platform-team"
    priority = 1
  })
  scorecard_relations = jsonencode({
    owner_team = port_entity.platform_team.identifier
  })
  rules = [
    {
      identifier = "zebra-production-environment"
      title      = "Production Environment"
      level      = "Gold"
      query = {
        combinator = "and"
        conditions = [jsonencode({
          property = "environment"
          operator = "="
          value    = "production"
        })]
      }
    },
    {
      identifier = "alpha-has-platform-author"
      title      = "Has Platform Author"
      level      = "Silver"
      query = {
        combinator = "and"
        conditions = [jsonencode({
          property = "author"
          operator = "="
          value    = "Platform Team"
        })]
      }
    },
    {
      identifier = "beta-has-author"
      title      = "Has Author"
      level      = "Bronze"
      query = {
        combinator = "and"
        conditions = [jsonencode({
          property = "author"
          operator = "isNotEmpty"
        })]
      }
    },
  ]
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
    port_system_blueprint.scorecard_group,
    port_system_blueprint.scorecard,
    port_entity.platform_team,
  ]
}

# Per-blueprint mode: each blueprint gets its own rules (and optional filter).
resource "port_scorecard_group" "blueprint_specific_readiness" {
  identifier = "blueprint-specific-readiness"
  title      = "Blueprint-Specific Readiness"
  group_properties = jsonencode({
    category = "governance"
  })
  group_relations = jsonencode({
    owning_team = port_entity.platform_team.identifier
  })
  scorecard_properties = jsonencode({
    owner    = "platform-team"
    priority = 2
  })
  scorecard_relations = jsonencode({
    owner_team = port_entity.platform_team.identifier
  })
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
      rules = [
        {
          identifier = "zebra-staging-environment"
          title      = "Staging Environment"
          level      = "Gold"
          query = {
            combinator = "and"
            conditions = [jsonencode({
              property = "environment"
              operator = "="
              value    = "staging"
            })]
          }
        },
        {
          identifier = "alpha-has-platform-author"
          title      = "Has Platform Author"
          level      = "Silver"
          query = {
            combinator = "and"
            conditions = [jsonencode({
              property = "author"
              operator = "="
              value    = "Platform Team"
            })]
          }
        },
      ]
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
      rules = [
        {
          identifier = "zebra-production-environment"
          title      = "Production Environment"
          level      = "Gold"
          query = {
            combinator = "and"
            conditions = [jsonencode({
              property = "environment"
              operator = "="
              value    = "production"
            })]
          }
        },
        {
          identifier = "beta-has-author"
          title      = "Has Author"
          level      = "Bronze"
          query = {
            combinator = "and"
            conditions = [jsonencode({
              property = "author"
              operator = "isNotEmpty"
            })]
          }
        },
      ]
    }
  }
  depends_on = [
    port_blueprint.microservice,
    port_blueprint.database,
    port_entity.vm_alpha,
    port_entity.vm_unowned,
    port_entity.db_primary,
    port_system_blueprint.scorecard_group,
    port_system_blueprint.scorecard,
    port_entity.platform_team,
  ]
}
