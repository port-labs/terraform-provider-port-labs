# Example: Deeply Nested Dataset Rules (3+ levels)
#
# This example demonstrates a dataset with 3 levels of nested rules,
# which requires datasetRuleSchema depth > 2 to work properly.
#
# Structure:
#   Level 1: AND combinator with leaf rule + group rule
#   Level 2: OR combinator with two AND groups
#   Level 3: AND groups containing leaf rules

# Blueprint for business domains
resource "port_blueprint" "aws_business_domain" {
  identifier  = "deep_nested_business_domain"
  title       = "AWS Business Domain"
  description = "Example blueprint for business domains"
  icon        = "Organization"

  properties = {
    string_props = {
      "domain_code" = {
        title = "Domain Code"
      }
    }
  }
}

# Blueprint for organization units
resource "port_blueprint" "aws_organization_unit" {
  identifier  = "deep_nested_org_unit"
  title       = "AWS Organization Unit"
  description = "Example blueprint for organization units"
  icon        = "Environment"

  properties = {
    string_props = {
      "ou_id" = {
        title = "OU ID"
      }
    }
  }
}

# Blueprint for permission sets (with relations)
resource "port_blueprint" "aws_permission_set" {
  identifier  = "deep_nested_permission_set"
  title       = "AWS Permission Set"
  description = "Example blueprint for testing deeply nested dataset rules"
  icon        = "AWS"

  properties = {
    boolean_props = {
      "user_assignable" = {
        title       = "User Assignable"
        description = "Whether users can be assigned this permission set"
        default     = true
      }
    }
  }

  relations = {
    "business_domain" = {
      title    = "Business Domain"
      target   = port_blueprint.aws_business_domain.identifier
      required = false
    }
    "organization_unit" = {
      title    = "Organization Unit"
      target   = port_blueprint.aws_organization_unit.identifier
      required = false
    }
  }

  depends_on = [
    port_blueprint.aws_business_domain,
    port_blueprint.aws_organization_unit
  ]
}

# Create sample entities
resource "port_entity" "eit_domain" {
  identifier = "eit"
  title      = "EIT Domain"
  blueprint  = port_blueprint.aws_business_domain.identifier

  properties = {
    string_props = {
      "domain_code" = "eit"
    }
  }
}

resource "port_entity" "finance_domain" {
  identifier = "finance"
  title      = "Finance Domain"
  blueprint  = port_blueprint.aws_business_domain.identifier

  properties = {
    string_props = {
      "domain_code" = "finance"
    }
  }
}

resource "port_entity" "sandbox_ou" {
  identifier = "ou-sandbox"
  title      = "Sandbox OU"
  blueprint  = port_blueprint.aws_organization_unit.identifier

  properties = {
    string_props = {
      "ou_id" = "ou-umvq-pmje9xva"
    }
  }
}

resource "port_entity" "prod_ou" {
  identifier = "ou-prod"
  title      = "Production OU"
  blueprint  = port_blueprint.aws_organization_unit.identifier

  properties = {
    string_props = {
      "ou_id" = "ou-umvq-prod1234"
    }
  }
}

# Create permission set entities
resource "port_entity" "admin_permission" {
  identifier = "admin-access"
  title      = "Admin Access"
  blueprint  = port_blueprint.aws_permission_set.identifier

  properties = {
    boolean_props = {
      "user_assignable" = true
    }
  }

  relations = {
    single_relations = {
      "business_domain"   = port_entity.eit_domain.identifier
      "organization_unit" = port_entity.sandbox_ou.identifier
    }
  }

  depends_on = [port_blueprint.aws_permission_set]
}

resource "port_entity" "readonly_permission" {
  identifier = "readonly-access"
  title      = "Read Only Access"
  blueprint  = port_blueprint.aws_permission_set.identifier

  properties = {
    boolean_props = {
      "user_assignable" = false
    }
  }

  relations = {
    single_relations = {
      "business_domain"   = port_entity.finance_domain.identifier
      "organization_unit" = port_entity.prod_ou.identifier
    }
  }

  depends_on = [port_blueprint.aws_permission_set]
}

# Action with DEEPLY NESTED dataset rules (3 levels)
# This is the key test case - it has:
#   Level 1: AND with leaf + group
#   Level 2: OR with two AND groups  
#   Level 3: AND groups with leaf rules (relatedTo)
resource "port_action" "request_sandbox_account" {
  identifier = "deep_nested_request_sandbox"
  title      = "Request Sandbox Account"
  icon       = "AWS"

  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.aws_permission_set.identifier

    user_properties = {
      string_props = {
        # Input for selecting domain
        "domain" = {
          title       = "Business Domain"
          description = "Select your business domain"
          format      = "entity"
          blueprint   = port_blueprint.aws_business_domain.identifier
        }

        # This property has 3 LEVELS of nested dataset rules
        "role_primary" = {
          title       = "Primary Role"
          icon        = "Lock"
          description = "Select a permission set filtered by complex nested rules"
          format      = "entity"
          blueprint   = port_blueprint.aws_permission_set.identifier

          dataset = {
            combinator = "and"
            rules = [
              # Level 1: Leaf rule - filter by user_assignable
              {
                property = "user_assignable"
                operator = "="
                value = {
                  jq_query = "true"
                }
              },
              # Level 1: Group rule with OR combinator
              {
                combinator = "or"
                rules = [
                  # Level 2: First AND group
                  {
                    combinator = "and"
                    rules = [
                      # Level 3: Leaf rules
                      {
                        blueprint = port_blueprint.aws_business_domain.identifier
                        operator  = "relatedTo"
                        value = {
                          jq_query = "\"eit\""
                        }
                      },
                      {
                        blueprint = port_blueprint.aws_organization_unit.identifier
                        operator  = "relatedTo"
                        value = {
                          jq_query = "\"ou-sandbox\""
                        }
                      }
                    ]
                  },
                  # Level 2: Second AND group  
                  {
                    combinator = "and"
                    rules = [
                      # Level 3: Leaf rules with dynamic jqQuery
                      {
                        blueprint = port_blueprint.aws_business_domain.identifier
                        operator  = "relatedTo"
                        value = {
                          jq_query = ".form.domain.identifier"
                        }
                      },
                      {
                        blueprint = port_blueprint.aws_organization_unit.identifier
                        operator  = "relatedTo"
                        value = {
                          jq_query = "\"ou-sandbox\""
                        }
                      }
                    ]
                  }
                ]
              }
            ]
          }

          depends_on = ["domain"]
        }
      }
    }

    order_properties = ["domain", "role_primary"]
  }

  webhook_method = {
    url = "https://example.com/webhook"
  }

  depends_on = [
    port_entity.admin_permission,
    port_entity.readonly_permission
  ]
}

output "action_identifier" {
  value = port_action.request_sandbox_account.identifier
}
