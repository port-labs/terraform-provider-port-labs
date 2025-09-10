terraform {
  required_providers {
    port = {
        source  = "port-labs/port-labs"
        version = "0.9.6"
    }
  }
}

provider "port" {
  client_id = "60EsooJtOqimlekxrNh7nfr2iOgTcyLZ"                                 # or set the env var PORT_CLIENT_ID
  secret    = "1yp6DQM6OwW4svTEIbrEFthMQqVVrIQ3fcJK0aJUpMwdc07ZuA6v27C3G08oj8Kx" # or set the env var PORT_CLIENT_SECRET
  base_url  = "http://api.localhost:9080"                                            # or set the env var PORT_BASE_URL
}

resource "port_blueprint" "author" {
    title = "Author"
    icon = "User"
    identifier = "author"
    properties = {
      string_props = {
        "name" = {
          type = "string"
          title = "Name"
        }
      }
    }
  }

  resource "port_blueprint" "team" {
    title = "Team"
    icon = "Team"
    identifier = "team"
    properties = {
      string_props = {
        "name" = {
          type = "string"
          title = "Team Name"
        }
      }
    }
  }

  resource "port_blueprint" "microservice" {
    title = "TF test microservice"
    icon = "Terraform"
    identifier = "microservice"
    properties = {
      string_props = {
        "url" = {
          type = "string"
          title = "URL"
        }
      }
    }
    relations = {
      "author" = {
        title = "Author"
        target = port_blueprint.author.identifier
      }
      "team" = {
        title = "Team"
        target = port_blueprint.team.identifier
      }
    }
  }


#   resource "port_webhook" "create_pr" {
#     identifier = "pr_webhook"
#     title      = "Webhook with mixed relations"
#     icon       = "Terraform"
#     enabled    = true
    
#     mappings = [
#       {
#         blueprint = port_blueprint.microservice.identifier
#         operation = { "type" = "create" }
#         filter    = ".headers.\"x-github-event\" == \"pull_request\""
#         entity = {
#           identifier = ".body.pull_request.id | tostring"
#           title      = ".body.pull_request.title"
#           properties = {
#             url = ".body.pull_request.html_url"
#           }
#           relations = {
#             # Complex object relation with search query
#             author = jsonencode({
#               combinator = "'and'",
#               rules = [
#                 {
#                   property = "'$identifier'"
#                   operator = "'='"
#                   value    = ".body.pull_request.user.login | tostring"
#                 }
#               ]
#             })
            
#             # Simple string relation
#             team = ".body.repository.owner.login | tostring"
#           }
#         }
#       }
#     ]
    
#     depends_on = [
#       port_blueprint.microservice,
#       port_blueprint.author,
#       port_blueprint.team
#     ]
#   }
  
# resource "port_webhook" "test_webhook_1" {
#   identifier = "testWebhook1"
#   title      = "Webhook with string relation"
#   enabled    = true
#   mappings = [
#     {
#       blueprint = "githubPullRequest"
#       operation = { "type" = "create" }
#       filter    = ".headers.\"x-github-event\" == \"pull_request\""
#       entity = {
#         identifier = ".body.pull_request.id | tostring"
#         title      = ".body.pull_request.title"
#         url        = ".body.pull_request.html_url"
#         relations = {
#           author = jsonencode({
#             rules = [
#               {
#                 operator = "'='"
#                 property = "'$identifier'"
#                 value    = ".body.pull_request.user.login | tostring"
#               }
#             ]
#           })
#         }
#       }
#     }
#   ]
# }

# resource "port_webhook" "test_webhook_2" {
#   identifier = "testWebhook2"
#   title      = "Webhook with json relation"
#   enabled    = true
#   mappings = [
#     {
#       blueprint = "githubPullRequest"
#       operation = { "type" = "create" }
#       filter    = ".headers.\"x-github-event\" == \"pull_request\""
#       entity = {
#         identifier = ".body.pull_request.id | tostring"
#         title      = ".body.pull_request.title"
#         url        = ".body.pull_request.html_url"
#         relations = {
#           author = jsonencode({
#             combinator = "'and'",
#             rules = [
#               {
#                 operator = "'='"
#                 property = "'$identifier'"
#                 value    = ".body.pull_request.user.login | tostring"
#               }
#             ]
#           })
#         }
#       }
#     }
#   ]
# }

resource "port_webhook" "test_webhook_3" {
  identifier = "testWebhook32"
  title      = "wh with both relations"
  enabled    = true
  mappings = [
    {
      blueprint = "githubPullRequest"
      operation = { "type" = "create" }
      filter    = ".headers.\"x-github-event\" == \"pull_request\""
      entity = {
        identifier = ".body.pull_request.id | tostring"
        title      = ".body.pull_request.title"
        url        = ".body.pull_request.html_url"
        # relations = {
        #   author = jsonencode({
        #     combinator = "'and'",
        #     rules = [
        #       {
        #         operator = "'='"
        #         property = "'$identifier'"
        #         value    = ".body.pull_request.user.login | tostring"
        #       }
        #     ]
        #   })
        #   team = ".body.repository.owner.login | tostring"
        # }
      }
    }
  ]
}

resource "port_webhook" "test_webhook_4" {
  identifier = "testWebhook3"
  title      = "wh with both relations"
  enabled    = true
  mappings = [
    {
      blueprint = "githubPullRequest"
      operation = { "type" = "create" }
      filter    = ".headers.\"x-github-event\" == \"pull_request\""
      entity = {
        # identifier = ".body.pull_request.id | tostring"
        identifier = jsonencode({
          combinator = "'and'",
          rules = [
            {
              property = "'$identifier'"
              operator = "'='"
              value    = ".body.pull_request.user.login | tostring"
            }
          ]
        })
        title      = ".body.pull_request.title"
        url        = ".body.pull_request.html_url"
        # relations = {
        #   author = jsonencode({
        #     combinator = "'and'",
        #     rules = [
        #       {
        #         operator = "'='"
        #         property = "'$identifier'"
        #         value    = ".body.pull_request.user.login | tostring"
        #       }
        #     ]
        #   })
        #   team = ".body.repository.owner.login | tostring"
        # }
      }
    }
  ]
}