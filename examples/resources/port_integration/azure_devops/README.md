# Azure DevOps integration example (Port PR #17668)

This example documents updated Azure DevOps integration mappings introduced in [port-labs/Port#17668](https://github.com/port-labs/Port/pull/17668).

## Provider impact

No Terraform provider schema changes are required. The `port_integration` resource already accepts arbitrary jq expressions in `config` via `jsonencode`, and the `port_system_blueprint` resource already supports relations generically.

## Mapping changes

| Resource kind | Change |
| --- | --- |
| `build` | Selector filters to finished builds (`.finishTime != null`) |
| `pipeline-stage` | Selector filters to finished stages (`.finishTime != null`) |
| `pipeline-run` | Selector filters to finished runs (`.finishedDate != null`) |
| `release-deployment` | Selector excludes unset `completedOn`; title uses release + environment name; `startedOn` nulls out ADO min-date; `url` uses `.release.webAccessUri`; adds `project` relation |
| `pipeline-deployment` | Selector filters to finished deployments; title uses `.owner.name` with id fallback |
| `project` relations | Project IDs strip spaces via `gsub(" "; "")` |

## System blueprint change

The `azureDevopsReleaseDeployment` system blueprint now includes a `project` relation targeting `azureDevopsProject`. This relation is provisioned by Port and appears automatically when refreshing a `port_system_blueprint` resource for that identifier.
