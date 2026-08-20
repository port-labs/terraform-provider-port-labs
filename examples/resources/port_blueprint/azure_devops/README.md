# Azure DevOps blueprint examples (Port PR #16979)

These examples define the Azure DevOps blueprint schema updates from [port-labs/Port#16979](https://github.com/port-labs/Port/pull/16979):

- **azureDevopsBuild** — `result` enum colors, `sourceBranch` property, and `repository` relation
- **azureDevopsPipeline** — new blueprint for pipeline definitions
- **azureDevopsPullRequest** — `closedAt`, `description`, `creator`, `reviewers`, branch properties
- **azureDevopsPipelineRun** — `pipeline` relation to `azureDevopsPipeline`

See `examples/resources/port_integration/azure_devops/` for matching integration mapping configuration.
