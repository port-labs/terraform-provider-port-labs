resource "port_webhook" "github" {
  identifier = "github"
  title      = "Github"
  icon       = "Terraform"
  enabled    = true
  mappings = [
    {
      "blueprint" : "pullRequest",
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
