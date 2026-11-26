terraform {
  required_providers {
    port = {
      source = "port-labs/port-labs"
    }
  }
}

provider "port" {}

resource "port_factory_line" "idea_to_prod" {
  identifier  = "idea-to-production"
  title       = "Idea to production"
  description = "Ship work through a Factory.ai build stage."

  stage {
    identifier   = "build"
    title        = "Build"
    instructions = "Implement the work item and open a pull request when done."
    executor     = "factory_ai_managed_agent"
    timeout      = "2h"

    factory_ai_managed_agent {
      environment_id   = "my-factory-environment"
      reasoning_effort = "medium"
    }
  }
}