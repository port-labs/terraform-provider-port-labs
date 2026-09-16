---
page_title: "port_scorecard_group Data Source - port"
subcategory: ""
description: |-
  Reads an existing scorecard group from the Port `/v1/scorecard-groups/{identifier}` API.
---

# port_scorecard_group (Data Source)

Reads an existing scorecard group by identifier. Use this data source to reference scorecard groups that are not managed by Terraform.

## Example Usage

```hcl
data "port_scorecard_group" "production_readiness" {
  identifier = "production-readiness"
}

output "title" {
  value = data.port_scorecard_group.production_readiness.title
}
```

## Schema

### Required

- `identifier` (String) The identifier of the scorecard group to read.

### Read-Only

- `id` (String) The scorecard group identifier.
- `title` (String) The scorecard group title.
- `levels` (Attributes List) Custom levels shared by member scorecards.
- `properties` (String) JSON-encoded `_scorecard` blueprint properties on member scorecards.
- `blueprints` (Set of String) Blueprint identifiers in shared-rules mode.
- `rules` (Attributes List) Shared rules in shared-rules mode.
- `filters` (Attributes Map) Per-blueprint filters in shared-rules mode.
- `scorecards` (Attributes Map) Per-blueprint member configuration.
- `created_at`, `created_by`, `updated_at`, `updated_by` (String) Audit metadata.
