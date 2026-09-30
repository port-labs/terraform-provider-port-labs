# GitHub Ocean pull request integration example (Port PR #18572)

This example documents updated GitHub Ocean `pull-request` mappings and `githubPullRequest` blueprint properties introduced in [port-labs/Port#18572](https://github.com/port-labs/Port/pull/18572).

## Provider impact

No Terraform provider schema changes are required. The `port_integration` resource already accepts arbitrary integration `config` JSON (including new pull-request selector fields and jq mappings) via `jsonencode`. The `port_system_blueprint` resource already supports the property types added to the default `githubPullRequest` blueprint (`string`, `boolean`, `number`, and string enums).

## Integration mapping changes

| Area | Change |
| --- | --- |
| Open `pull-request` selector | Uses GraphQL (`api: graphql`), enables `enrichWithFirstCommit`, and excludes heavy GraphQL fields via `excludeGraphqlFields` |
| Open PR identifier | `.fullDatabaseId\|tostring` (GraphQL database ID) |
| Timestamp / link fields | REST snake_case fields (`created_at`, `html_url`, …) replaced with GraphQL camelCase (`createdAt`, `url`, …) |
| DORA / review metrics | Maps `isDraft`, `reviewDecision`, review timestamps, and derived hour metrics (`codingTimeHours`, `timeToFirstReview`, `timeFromApprovalToMerge`) |
| Closed `pull-request` identifier | `.id\|tostring` for the closed-PR sync resource (unchanged in Port #18572) |

## System blueprint changes

The `githubPullRequest` system blueprint gains properties for draft state, review decision, review timestamps, and DORA-style hour metrics. See `examples/resources/port_system_blueprint/github_pull_request/README.md` for the property list.
