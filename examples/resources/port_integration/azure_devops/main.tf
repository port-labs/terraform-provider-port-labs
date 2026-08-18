<<<<<<< HEAD
# Example: Azure DevOps integration mappings updated in port-labs/Port#16979
#
# Port PR #16979 extends the default Azure DevOps catalog with a pipeline blueprint,
# enriches pull request/build/pipeline-run mappings, and fixes link mappings to use
# _links.web.href. The port_integration resource accepts these mappings via jsonencode.

resource "port_integration" "azure_devops" {
  installation_id       = "my-azure-devops-integration"
  title                 = "Azure DevOps"
  installation_app_type = "AZURE DEVOPS"

  config = jsonencode({
    deleteDependentEntities = true
    resources = [
      {
        kind = "pipeline"
        selector = {
          query       = "true"
          includeRepo = true
        }
        port = {
          entity = {
            mappings = [{
              identifier = ".__projectId + \"/\" + (.id | tostring)"
              title      = ".name"
              blueprint  = "'azureDevopsPipeline'"
              properties = {
                url      = ".url"
                revision = ".revision"
                folder   = ".folder"
              }
              relations = {
                project    = ".__projectId | gsub(\" \"; \"\")"
                repository = "if .__repository then .__repository.id else null end"
              }
            }]
          }
        }
      },
      {
=======
# Azure DevOps integration mapping snippets aligned with Port PR #16979.
# Import an existing Azure DevOps integration and apply updated mappings.
# https://github.com/port-labs/Port/pull/16979

resource "port_integration" "azure_devops" {
  installation_id       = "my-azure-devops-installation-id"
  title                 = "Azure DevOps"
  installation_app_type = "AZURE_DEVOPS"

  config = jsonencode({
    createMissingRelatedEntities = true
    deleteDependentEntities      = true
    resources = [
      {
>>>>>>> 6fd0342 (Add Azure DevOps blueprint and integration examples (Port PR #16979))
        kind = "pull-request"
        selector = {
          query = "true"
        }
        port = {
          entity = {
<<<<<<< HEAD
            mappings = [{
              identifier = ".repository.id + \"/\" + (.pullRequestId | tostring)"
              title      = ".title"
              blueprint  = "'azureDevopsPullRequest'"
=======
            mappings = {
              identifier = ".repository.id + \"/\" + (.pullRequestId | tostring)"
              blueprint  = "\"azureDevopsPullRequest\""
>>>>>>> 6fd0342 (Add Azure DevOps blueprint and integration examples (Port PR #16979))
              properties = {
                status        = ".status"
                createdAt     = ".creationDate"
                closedAt      = ".closedDate"
                description   = ".description"
                creator       = ".createdBy.displayName"
                reviewers     = "[.reviewers[].displayName]"
                sourceBranch  = ".sourceRefName | gsub(\"refs/heads/\"; \"\")"
                targetBranch  = ".targetRefName | gsub(\"refs/heads/\"; \"\")"
                leadTimeHours = "(.creationDate as $createdAt | .status as $status | .closedDate as $closedAt | ($createdAt | sub(\"\\\\..*Z$\"; \"Z\") | strptime(\"%Y-%m-%dT%H:%M:%SZ\") | mktime) as $createdTimestamp | ($closedAt | if . == null then null else sub(\"\\\\..*Z$\"; \"Z\") | strptime(\"%Y-%m-%dT%H:%M:%SZ\") | mktime end) as $closedTimestamp | if $status == \"completed\" and $closedTimestamp != null then (((($closedTimestamp - $createdTimestamp) / 3600) * 100 | floor) / 100) else null end)"
                link          = "._links.web.href"
              }
              relations = {
                repository = ".repository.id"
              }
<<<<<<< HEAD
            }]
=======
            }
>>>>>>> 6fd0342 (Add Azure DevOps blueprint and integration examples (Port PR #16979))
          }
        }
      },
      {
        kind = "build"
        selector = {
          query = "true"
        }
        port = {
          entity = {
<<<<<<< HEAD
            mappings = [{
              identifier = ".__project.id + \"/\" + (.repository.id | tostring) + \"/\" + (.id | tostring) | gsub(\" \"; \"\")"
              title      = ".buildNumber"
              blueprint  = "'azureDevopsBuild'"
=======
            mappings = {
              identifier = ".__project.id + \"/\" + (.repository.id | tostring) + \"/\" + (.id | tostring) | gsub(\" \"; \"\")"
              title      = ".buildNumber"
              blueprint  = "\"azureDevopsBuild\""
>>>>>>> 6fd0342 (Add Azure DevOps blueprint and integration examples (Port PR #16979))
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
<<<<<<< HEAD
            }]
=======
            }
          }
        }
      },
      {
        kind = "pipeline"
        selector = {
          query       = "true"
          includeRepo = true
        }
        port = {
          entity = {
            mappings = {
              identifier = ".__projectId + \"/\" + (.id | tostring)"
              title      = ".name"
              blueprint  = "\"azureDevopsPipeline\""
              properties = {
                url      = ".url"
                revision = ".revision"
                folder   = ".folder"
              }
              relations = {
                project    = ".__projectId | gsub(\" \"; \"\")"
                repository = "if .__repository then .__repository.id else null end"
              }
            }
>>>>>>> 6fd0342 (Add Azure DevOps blueprint and integration examples (Port PR #16979))
          }
        }
      },
      {
        kind = "pipeline-run"
        selector = {
          query = "true"
        }
        port = {
          entity = {
<<<<<<< HEAD
            mappings = [{
              identifier = ".__project.id + \"/\" + (.__pipeline.id | tostring) + \"/\" + (.id | tostring) | gsub(\" \"; \"\")"
              title      = ".name"
              blueprint  = "'azureDevopsPipelineRun'"
=======
            mappings = {
              identifier = ".__project.id + \"/\" + (.__pipeline.id | tostring) + \"/\" + (.id | tostring) | gsub(\" \"; \"\")"
              title      = ".name"
              blueprint  = "\"azureDevopsPipelineRun\""
>>>>>>> 6fd0342 (Add Azure DevOps blueprint and integration examples (Port PR #16979))
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
<<<<<<< HEAD
            }]
          }
        }
      },
=======
            }
          }
        }
      }
>>>>>>> 6fd0342 (Add Azure DevOps blueprint and integration examples (Port PR #16979))
    ]
  })
}
