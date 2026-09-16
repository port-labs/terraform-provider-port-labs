# port_integration examples

Each integration type has its own directory. Subdirectories represent supported installation cases.

| Integration | Case | Path |
|-------------|------|------|
| OnPrem (custom) | Basic create (step 1 — no config) | [`onprem/basic/`](./onprem/basic/) |
| K8s exporter | Version pinning | [`k8s-exporter/version-pinning/`](./k8s-exporter/version-pinning/) |
| GitHub Ocean | Personal Access Token | [`github-ocean/pat/`](./github-ocean/pat/) |
| GitHub Ocean | GitHub App | [`github-ocean/github-app/`](./github-ocean/github-app/) |
| Azure DevOps | Single account (PAT) | [`azure-devops/single-account-pat/`](./azure-devops/single-account-pat/) |
| Azure DevOps | Multiple accounts (service principal) | [`azure-devops/multi-account-service-principal/`](./azure-devops/multi-account-service-principal/) |
| Jira | API token | [`jira/api-token/`](./jira/api-token/) |
| GitLab v2 | Personal access token | [`gitlab-v2/personal-access-token/`](./gitlab-v2/personal-access-token/) |

Secret naming convention (SaaS): `_{INSTALLATION_ID}_{INTEGRATION_TYPE}_{PROPERTY}` in `SCREAMING_SNAKE_CASE` — matches the Port UI `generateSecretName()`.

SaaS integrations: create without `config` first, then add `config` on a subsequent apply to override default mappings.
