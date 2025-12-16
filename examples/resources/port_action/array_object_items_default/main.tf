resource "port_blueprint" "repository" {
  title      = "Repository"
  icon       = "Git"
  identifier = "examples-action-repository"
  properties = {
    string_props = {
      name = {
        type  = "string"
        title = "Name"
      }
      url = {
        type  = "string"
        title = "URL"
        format = "url"
      }
    }
  }
}

resource "port_action" "deploy_with_repositories" {
  title      = "Deploy with Repositories"
  icon       = "Terraform"
  identifier = "examples-action-deploy-with-repositories"
  publish    = true

  self_service_trigger = {
    operation            = "CREATE"
    blueprint_identifier = port_blueprint.repository.identifier
    user_properties = {
      array_props = {
        "related_repositories" = {
          title       = "Related Repositories"
          description = "Select related repositories for deployment"
          required    = true
          object_items = {
            default = [
              {
                name = "example-repo-1"
                url  = "https://github.com/example/repo-1"
              },
              {
                name = "example-repo-2"
                url  = "https://github.com/example/repo-2"
              }
            ]
          }
        }
      }
    }
  }

  webhook_method = {
    url = "https://api.example.com/deploy"
  }
}
