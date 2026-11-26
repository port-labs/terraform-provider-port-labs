# GitHub Pull Request system blueprint (Port PR #18572)

[port-labs/Port#18572](https://github.com/port-labs/Port/pull/18572) extends the default `githubPullRequest` system blueprint with DORA and review lifecycle properties. New GitHub Ocean installations receive these properties automatically from Port provisioning.

No Terraform changes are required to consume the new fields on ingested entities. Use `port_system_blueprint` only when you need to add **additional** custom properties or relations on top of the system blueprint.

## Properties added in Port #18572

| Property | Type | Notes |
| --- | --- | --- |
| `firstCommitAt` | `string` (`date-time`) | First commit timestamp on the PR branch (requires `enrichWithFirstCommit` in integration mapping) |
| `isDraft` | `boolean` | Whether the PR is in draft state |
| `reviewDecision` | `string` enum | `APPROVED`, `CHANGES_REQUESTED`, or `REVIEW_REQUIRED` |
| `firstReviewAt` | `string` (`date-time`) | Earliest review activity timestamp |
| `approvedAt` | `string` (`date-time`) | Most recent standing approval (excludes dismissed approvals) |
| `readyForReviewAt` | `string` (`date-time`) | When the PR left draft; falls back to `createdAt` when never drafted |
| `codingTimeHours` | `number` | Hours from first commit to ready-for-review |
| `timeToFirstReview` | `number` | Hours from PR creation to first review |
| `timeFromApprovalToMerge` | `number` | Hours from latest standing approval to merge |

See `examples/resources/port_integration/github_ocean_pull_request/` for matching GitHub Ocean `pull-request` integration mappings.
