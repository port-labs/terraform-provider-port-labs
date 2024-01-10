resource "port_blueprint" "repository_blueprint" {
  title       = "Repository Blueprint"
  icon        = "Terraform"
  identifier  = "repository"
  description = ""
}

resource "port_blueprint" "pull_request_blueprint" {
  title       = "Pull Request Blueprint"
  icon        = "Terraform"
  identifier  = "pull_request"
  description = ""
  properties = {
    string_props = {
      "status" = {
        title = "Status"
      }
    }
  }
  relations = {
    "repository" = {
      title  = "Repository"
      target = port_blueprint.repository_blueprint.identifier
    }
  }
}

resource "port_aggregation_property" "fix_pull_requests_per_day" {
  aggregation_identifier      = "fix_pull_requests_count"
  blueprint_identifier        = port_blueprint.repository_blueprint.identifier
  target_blueprint_identifier = port_blueprint.pull_request_blueprint.identifier
  title                       = "Pull Requests Per Day"
  icon                        = "Terraform"
  description                 = "Pull Requests Per Day"
  method = {
    average_entities = {
      average_of      = "month"
      measure_time_by = "$createdAt"
    }
  }
  query = jsonencode(
    {
      "combinator" : "and",
      "rules" : [
        {
          "property" : "$title",
          "operator" : "ContainsAny",
          "value" : ["fix", "fixed", "fixing", "Fix"]
        }
      ]
    }
  )
}

