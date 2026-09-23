# port_integration examples

Each integration type has its own directory. Subdirectories represent supported installation cases.

| Integration | Case | Path |
|-------------|------|------|
| Self-hosted (custom) | Basic create (step 1 — no config) | [`onprem/basic/`](./onprem/basic/) |
| K8s exporter | Version pinning | [`k8s-exporter/version-pinning/`](./k8s-exporter/version-pinning/) |
| K8s exporter | Per-resource `enableDelete` | [`enable_delete/`](./enable_delete/) |
| GitHub Ocean | Personal Access Token | [`github-ocean/pat/`](./github-ocean/pat/) |
| GitHub Ocean | Personal Access Token (skip default resources) | [`github-ocean/pat-skip-default-resources/`](./github-ocean/pat-skip-default-resources/) |
| GitHub Ocean | GitHub App | [`github-ocean/github-app/`](./github-ocean/github-app/) |
| Azure DevOps | Single account (PAT) | [`azure-devops/single-account-pat/`](./azure-devops/single-account-pat/) |
| Azure DevOps | Multiple accounts (service principal) | [`azure-devops/multi-account-service-principal/`](./azure-devops/multi-account-service-principal/) |
| Jira | API token | [`jira/api-token/`](./jira/api-token/) |
| Linear | API key | [`linear/api-key/`](./linear/api-key/) |
| GitLab v2 | Personal access token | [`gitlab-v2/personal-access-token/`](./gitlab-v2/personal-access-token/) |

### Port Hosted `appSpec` and create options

See [`app-spec/`](./app-spec/) for `liveEventsEnabled`, `actionsProcessingEnabled`, `incrementalSyncEnabled`, `incrementalSyncInterval`, `scheduledResyncInterval`, and `create_port_resources_origin` examples.

Port Hosted integrations: create without `config` first, then add `config` on a subsequent apply to override default mappings.
