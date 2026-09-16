---
page_title: "port_scorecard_groups Data Source - port"
subcategory: ""
description: |-
  Lists all scorecard groups in the organization via `GET /v1/scorecard-groups`.
---

# port_scorecard_groups (Data Source)

Use this data source to list scorecard groups. It mirrors the Port API list endpoint introduced for scorecard group management in the product UI.

~> **Note:** Scorecard groups require the `SCORECARD_GROUPS` organization feature flag.

## Example Usage

```hcl
data "port_scorecard_groups" "all" {}

output "scorecard_group_identifiers" {
  value = [for group in data.port_scorecard_groups.all.groups : group.identifier]
}
```

## Schema

### Read-Only

- `id` (String) Stable identifier for this data source result.
- `groups` (List of Object) All scorecard groups returned by the Port API.
  - `identifier` (String) The scorecard group identifier.
  - `title` (String) The scorecard group title.
