resource "port_blueprint" "microservice" {
  identifier  = "examples-webhook-svc"
  title       = "Microsvc from Port TF Examples"
  icon        = "Terraform"
  description = ""
  properties = {
    string_props = {
      url = {
        type = "string"
      }
      author = {
        icon       = "github"
        required   = true
        min_length = 1
        max_length = 10
        default    = "default"
        enum       = ["default", "default2"]
        pattern    = "^[a-zA-Z0-9]*$"
        format     = "user"
        enum_colors = {
          default  = "red"
          default2 = "green"
        }
      }
    }
  }
}

resource "port_webhook" "github" {
  identifier = "github"
  title      = "Github"
  icon       = "Github"
  enabled    = true
  mappings = [
    {
      "blueprint" : port_blueprint.microservice.identifier,
      "filter" : ".headers.\"X-GitHub-Event\" == \"pull_request\"",
      "entity" : {
        "identifier" : ".body.pull_request.id | tostring",
        "title" : ".body.pull_request.title",
        "properties" : {
          "author" : ".body.pull_request.user.login",
          "url" : ".body.pull_request.html_url"
        }
      }
    }
  ]
}
