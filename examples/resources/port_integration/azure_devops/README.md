# Azure DevOps integration examples (Port PR #16979)

These examples show integration mapping updates from [port-labs/Port#16979](https://github.com/port-labs/Port/pull/16979):

- **pull-request** — enriched properties (`closedAt`, `creator`, `reviewers`, branch names) and a constructed PR link
- **build** — `sourceBranch` property and `repository` relation
- **pipeline** — new resource kind with `includeRepo` selector
- **pipeline-run** — `title` mapping and `pipeline` relation to `azureDevopsPipeline`

See `examples/resources/port_blueprint/azure_devops/` for matching blueprint schema definitions.
