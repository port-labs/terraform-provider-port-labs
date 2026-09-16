---
page_title: "port_scorecard_group_scorecard Resource - port"
subcategory: ""
description: |-
  Associates an existing scorecard with a scorecard group via the Port membership API.
---

# port_scorecard_group_scorecard (Resource)

Associates an existing scorecard with a scorecard group using `POST /v1/scorecard-groups/{group}/scorecards/{scorecard}`. Destroying the resource calls the remove-membership API.

## Example Usage

```hcl
resource "port_scorecard_group_scorecard" "example" {
  group_identifier     = "production-readiness"
  scorecard_identifier = "legacy-readiness"
  override_group       = true
}
```

## Schema

### Required

- `group_identifier` (String) Scorecard group identifier.
- `scorecard_identifier` (String) Scorecard identifier to add to the group.

### Optional

- `override_group` (Boolean) When true, moves the scorecard to this group even if it already belongs to another group. Defaults to `false`.

### Read-Only

- `id` (String) Composite ID `<group_identifier>/<scorecard_identifier>`.

## Import

Import using `terraform import port_scorecard_group_scorecard.example <group_identifier>/<scorecard_identifier>`.
