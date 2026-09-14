# GitHub Ocean pull-request mapping snippets aligned with Port PR #18572.
# Import an existing GitHub Ocean integration and apply updated mappings.
# https://github.com/port-labs/Port/pull/18572

resource "port_integration" "github_ocean" {
  installation_id       = "my-github-ocean-installation-id"
  title                 = "GitHub"
  installation_app_type = "GitHub"

  config = jsonencode({
    createMissingRelatedEntities = true
    deleteDependentEntities      = true
    resources = [
      {
        kind = "pull-request"
        selector = {
          query                   = "true"
          states                  = ["open"]
          api                     = "graphql"
          enrichWithFirstCommit   = true
          excludeGraphqlFields    = ["additions", "deletions", "changedFiles"]
        }
        port = {
          entity = {
            mappings = [{
              identifier = ".head.repo.name + (.fullDatabaseId|tostring)"
              title      = ".title"
              blueprint  = "'githubPullRequest'"
              properties = {
                status                  = ".state"
                closedAt                = ".closedAt"
                updatedAt               = ".updatedAt"
                mergedAt                = ".mergedAt"
                createdAt               = ".createdAt"
                prNumber                = ".number"
                link                    = ".url"
                branch                  = ".head.ref"
                leadTimeHours           = "(.createdAt as $createdAt | .mergedAt as $mergedAt | if $createdAt == null or $mergedAt == null then null else ((($mergedAt | strptime(\"%Y-%m-%dT%H:%M:%SZ\") | mktime) - ($createdAt | strptime(\"%Y-%m-%dT%H:%M:%SZ\") | mktime)) / 3600 * 100 | floor) / 100 end)"
                isDraft                 = ".isDraft"
                reviewDecision          = ".reviewDecision"
                firstCommitAt           = ".firstCommit.committedDate"
                firstReviewAt           = "(.reviews.nodes | map(select(.createdAt != null)) | sort_by(.createdAt) | first.createdAt) // null"
                approvedAt              = "(.reviews.nodes | map(select(.state == \"APPROVED\" and .createdAt != null)) | sort_by(.createdAt) | last.createdAt) // null"
                readyForReviewAt        = "(if .isDraft == true then null else .createdAt end)"
                codingTimeHours         = "(.firstCommit.committedDate as $first | (if .isDraft == true then null else .createdAt end) as $ready | if $first == null or $ready == null then null else ((($ready | strptime(\"%Y-%m-%dT%H:%M:%SZ\") | mktime) - ($first | strptime(\"%Y-%m-%dT%H:%M:%SZ\") | mktime)) / 3600 * 100 | floor) / 100 end)"
                timeToFirstReview       = "(.createdAt as $created | (.reviews.nodes | map(select(.createdAt != null)) | sort_by(.createdAt) | first.createdAt) as $firstReview | if $created == null or $firstReview == null then null else ((($firstReview | strptime(\"%Y-%m-%dT%H:%M:%SZ\") | mktime) - ($created | strptime(\"%Y-%m-%dT%H:%M:%SZ\") | mktime)) / 3600 * 100 | floor) / 100 end)"
                timeFromApprovalToMerge = "((.reviews.nodes | map(select(.state == \"APPROVED\" and .createdAt != null)) | sort_by(.createdAt) | last.createdAt) as $approved | .mergedAt as $merged | if $approved == null or $merged == null then null else ((($merged | strptime(\"%Y-%m-%dT%H:%M:%SZ\") | mktime) - ($approved | strptime(\"%Y-%m-%dT%H:%M:%SZ\") | mktime)) / 3600 * 100 | floor) / 100 end)"
              }
              relations = {
                repository = ".__repository"
              }
            }]
          }
        }
      },
      {
        kind = "pull-request"
        selector = {
          query  = "true"
          states = ["closed"]
        }
        port = {
          entity = {
            mappings = [{
              identifier = ".head.repo.name + (.id|tostring)"
              blueprint  = "'githubPullRequest'"
              properties = {}
              relations = {
                repository = ".__repository"
              }
            }]
          }
        }
      },
    ]
  })
}
