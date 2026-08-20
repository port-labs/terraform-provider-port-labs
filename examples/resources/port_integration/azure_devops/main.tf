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
        kind = "pull-request"
        selector = {
          query = "true"
        }
        port = {
          entity = {
            mappings = [{
              identifier = ".repository.id + \"/\" + (.pullRequestId | tostring)"
              title      = ".title"
              blueprint  = "'azureDevopsPullRequest'"
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
                link          = "(.url | split(\"/_apis/\")[0]) + \"/_git/\" + .repository.name + \"/pullrequest/\" + (.pullRequestId | tostring)"
              }
              relations = {
                repository = ".repository.id"
              }
            }]
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
        kind = "pipeline-run"
        selector = {
          query = "true"
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
    ]
  })
}
