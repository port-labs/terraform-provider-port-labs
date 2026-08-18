# Azure DevOps blueprint definitions aligned with Port PR #16979.
# https://github.com/port-labs/Port/pull/16979

resource "port_blueprint" "azure_devops_build" {
  identifier = "azureDevopsBuild"
  title      = "Build"
  icon       = "AzureDevops"

  properties = {
    string_props = {
      status = {
        title = "Status"
        type  = "string"
      }
      result = {
        title = "Result"
        type  = "string"
        enum  = ["succeeded", "partiallySucceeded", "failed", "canceled", "none"]
        enum_colors = {
          succeeded          = "green"
          partiallySucceeded = "yellow"
          failed             = "red"
          canceled           = "lightGray"
          none               = "darkGray"
        }
      }
      queueTime = {
        title  = "Queue Time"
        type   = "string"
        format = "date-time"
      }
      startTime = {
        title  = "Start Time"
        type   = "string"
        format = "date-time"
      }
      finishTime = {
        title  = "Finish Time"
        type   = "string"
        format = "date-time"
      }
      definitionName = {
        title = "Definition Name"
        type  = "string"
      }
      requestedFor = {
        title = "Requested For"
        type  = "string"
      }
      link = {
        title  = "Link"
        type   = "string"
        format = "url"
      }
      sourceBranch = {
        title = "Source Branch"
        type  = "string"
      }
    }
  }

  relations = {
    project = {
      title    = "Project"
      target   = "azureDevopsProject"
      required = false
      many     = false
    }
    repository = {
      title    = "Repository"
      target   = "azureDevopsRepository"
      required = false
      many     = false
    }
  }
}

resource "port_blueprint" "azure_devops_pipeline" {
  identifier = "azureDevopsPipeline"
  title      = "Azure Devops Pipeline"
  icon       = "AzureDevops"

  properties = {
    string_props = {
      url = {
        title  = "URL"
        type   = "string"
        format = "url"
      }
      folder = {
        title = "Folder"
        type  = "string"
      }
    }
    number_props = {
      revision = {
        title = "Revision"
        type  = "number"
      }
    }
  }

  relations = {
    project = {
      title    = "Project"
      target   = "azureDevopsProject"
      required = false
      many     = false
    }
    repository = {
      title    = "Repository"
      target   = "azureDevopsRepository"
      required = false
      many     = false
    }
  }
}

resource "port_blueprint" "azure_devops_pull_request" {
  identifier = "azureDevopsPullRequest"
  title      = "Pull Request"
  icon       = "AzureDevops"

  properties = {
    string_props = {
      status = {
        title = "Status"
        type  = "string"
        enum  = ["completed", "abandoned", "active"]
        enum_colors = {
          completed = "yellow"
          abandoned = "red"
          active    = "green"
        }
      }
      createdAt = {
        title  = "Create At"
        type   = "string"
        format = "date-time"
      }
      closedAt = {
        title  = "Closed At"
        type   = "string"
        format = "date-time"
      }
      description = {
        title = "Description"
        type  = "string"
      }
      creator = {
        title = "Creator"
        type  = "string"
      }
      sourceBranch = {
        title = "Source Branch"
        type  = "string"
      }
      targetBranch = {
        title = "Target Branch"
        type  = "string"
      }
      link = {
        title  = "Link"
        type   = "string"
        format = "url"
      }
    }
    number_props = {
      leadTimeHours = {
        title = "Lead Time in hours"
        type  = "number"
      }
    }
    array_props = {
      reviewers = {
        title = "Reviewers"
        type  = "array"
        string_items = {}
      }
    }
  }

  relations = {
    repository = {
      title    = "repository"
      target   = "azureDevopsRepository"
      required = false
      many     = false
    }
  }
}

resource "port_blueprint" "azure_devops_pipeline_run" {
  identifier = "azureDevopsPipelineRun"
  title      = "Pipeline Run"
  icon       = "AzureDevops"

  properties = {
    string_props = {
      state = {
        title = "State"
        type  = "string"
      }
      result = {
        title = "Result"
        type  = "string"
      }
      createdDate = {
        title  = "Created Date"
        type   = "string"
        format = "date-time"
      }
      finishedDate = {
        title  = "Finished Date"
        type   = "string"
        format = "date-time"
      }
      pipelineName = {
        title = "Pipeline Name"
        type  = "string"
      }
    }
  }

  relations = {
    project = {
      title    = "Project"
      target   = "azureDevopsProject"
      required = true
      many     = false
    }
    pipeline = {
      title    = "Pipeline"
      target   = port_blueprint.azure_devops_pipeline.identifier
      required = false
      many     = false
    }
  }
}
