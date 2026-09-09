# Azure DevOps integration mapping snippets aligned with Port PR #17668.
# Import an existing Azure DevOps integration and apply updated mappings.
# https://github.com/port-labs/Port/pull/17668

resource "port_integration" "azure_devops" {
  installation_id       = "my-azure-devops-installation-id"
  title                 = "Azure DevOps"
  installation_app_type = "AZURE_DEVOPS"

  config = jsonencode({
    createMissingRelatedEntities = true
    deleteDependentEntities      = true
    resources = [
      {
        kind = "environment"
        selector = {
          query = "true"
        }
        port = {
          entity = {
            mappings = [{
              identifier = ".project.id + \"/\" + (.id | tostring) | gsub(\" \"; \"\")"
              title      = ".name | tostring"
              blueprint  = "'azureDevopsEnvironment'"
              properties = {
                description    = ".description"
                createdOn      = ".createdOn"
                lastModifiedOn = ".lastModifiedOn"
              }
              relations = {
                project = ".project.id | gsub(\" \"; \"\")"
              }
            }]
          }
        }
      },
      {
        kind = "build"
        selector = {
          query = ".finishTime != null"
        }
        port = {
          entity = {
            mappings = [{
              identifier = ".__project.id + \"/\" + (.repository.id | tostring) + \"/\" + (.id | tostring) | gsub(\" \"; \"\")"
              title      = ".buildNumber"
              blueprint  = "'azureDevopsBuild'"
              properties = {
                status         = ".status"
                result         = ".result"
                queueTime      = ".queueTime"
                startTime      = ".startTime"
                finishTime     = ".finishTime"
                definitionName = ".definition.name"
                requestedFor   = ".requestedFor.displayName"
                sourceBranch   = ".sourceBranch | gsub(\"refs/heads/\"; \"\")"
                link           = "._links.web.href"
              }
              relations = {
                project    = ".__project.id | gsub(\" \"; \"\")"
                repository = ".repository.id"
              }
            }]
          }
        }
      },
      {
        kind = "pipeline-stage"
        selector = {
          query = ".finishTime != null"
        }
        port = {
          entity = {
            mappings = [{
              identifier = ".__project.id + \"/\" + (.__build.repository.id | tostring) + \"/\" + (.__build.id | tostring) + \"/\" + (.id | tostring) | gsub(\" \"; \"\")"
              title      = ".name"
              blueprint  = "'azureDevopsPipelineStage'"
              properties = {
                state      = ".state"
                result     = ".result"
                startTime  = ".startTime"
                finishTime = ".finishTime"
                stageType  = ".type"
              }
              relations = {
                project = ".__project.id | gsub(\" \"; \"\")"
                build   = "(.__project.id + \"/\" + (.__build.repository.id | tostring) + \"/\" + (.__build.id | tostring)) | gsub(\" \"; \"\")"
              }
            }]
          }
        }
      },
      {
        kind = "pipeline-run"
        selector = {
          query = ".finishedDate != null"
        }
        port = {
          entity = {
            mappings = [{
              identifier = ".__project.id + \"/\" + (.__pipeline.id | tostring) + \"/\" + (.id | tostring) | gsub(\" \"; \"\")"
              title      = ".name"
              blueprint  = "'azureDevopsPipelineRun'"
              properties = {
                state        = ".state"
                result       = ".result"
                createdDate  = ".createdDate"
                finishedDate = ".finishedDate"
                pipelineName = ".pipeline.name"
              }
              relations = {
                project  = ".__project.id | gsub(\" \"; \"\")"
                pipeline = ".__project.id + \"/\" + (.__pipeline.id | tostring)"
              }
            }]
          }
        }
      },
      {
        kind = "release-deployment"
        selector = {
          query          = ".completedOn != null and .completedOn != \"0001-01-01T00:00:00\""
          includeRelease = true
        }
        port = {
          entity = {
            mappings = [{
              identifier = ".releaseDefinition.projectReference.id + \"/\" + (.release.id | tostring) + \"/\" + (.id | tostring) | gsub(\" \"; \"\")"
              title      = ".release.name + \" - \" + .releaseEnvironment.name"
              blueprint  = "'azureDevopsReleaseDeployment'"
              properties = {
                status          = ".deploymentStatus"
                url             = ".release.webAccessUri"
                reason          = ".reason"
                startedOn       = ".startedOn | if . == \"0001-01-01T00:00:00\" then null else . end"
                completedOn     = ".completedOn"
                requestedBy     = ".requestedBy.displayName"
                operationStatus = ".operationStatus"
                environment     = ".releaseEnvironment.name"
              }
              relations = {
                project = ".releaseDefinition.projectReference.id | gsub(\" \"; \"\")"
                release = ".releaseDefinition.projectReference.id + \"/\" + (.release.id | tostring) | gsub(\" \"; \"\")"
              }
            }]
          }
        }
      },
      {
        kind = "pipeline-deployment"
        selector = {
          query = ".finishTime != null"
        }
        port = {
          entity = {
            mappings = [{
              identifier = ".__project.id + \"/\" + (.environmentId | tostring) + \"/\" + (.id | tostring) | gsub(\" \"; \"\")"
              title      = ".owner.name // (.id | tostring)"
              blueprint  = "'azureDevopsPipelineDeployment'"
              properties = {
                planType    = ".planType"
                stageName   = ".stageName"
                jobName     = ".jobName"
                result      = ".result"
                startTime   = ".startTime"
                finishTime  = ".finishTime"
              }
              relations = {
                project     = ".__project.id | gsub(\" \"; \"\")"
                environment = ".__project.id + \"/\" + (.environmentId | tostring) | gsub(\" \"; \"\")"
              }
            }]
          }
        }
      },
    ]
  })
}
