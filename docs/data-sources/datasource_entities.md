---
page_title: "port_datasource_entities Data Source - port"
subcategory: ""
description: |-
  Datasource Entities Data Source
  The datasource entities data source looks up entity identifiers by matching their datasource path against a prefix and suffix.
---

# port_datasource_entities (Data Source)

The datasource entities data source looks up entity identifiers by matching their datasource path against a prefix and suffix.
This mirrors the Port API endpoint updated in [port-labs/Port#18099](https://github.com/port-labs/Port/pull/18099), which no longer accepts an `extra_conditions.not_modified_since_creation` filter.

## Example Usage

```hcl
data "port_datasource_entities" "integration_entities" {
  datasource_prefix = "port-ocean/github/"
  datasource_suffix = "/my-github-integration/resync"
}

output "integration_entity_identifiers" {
  value = [for entity in data.port_datasource_entities.integration_entities.entities : entity.identifier]
}
```

## Schema

### Required

- `datasource_prefix` (String) The datasource prefix to match entities against.
- `datasource_suffix` (String) The datasource suffix to match entities against.

### Optional

- `before` (String) Return only entities updated before this RFC3339 timestamp.
- `limit` (Number) The maximum number of entities to return per request. When omitted, Port uses its default page size.

### Read-Only

- `entities` (Attributes List) Entities whose datasource matches the prefix and suffix. (see [below for nested schema](#nestedatt--entities))
- `id` (String) The ID of this resource.

<a id="nestedatt--entities"></a>
### Nested Schema for `entities`

Read-Only:

- `blueprint` (String) The blueprint identifier of the entity.
- `identifier` (String) The identifier of the entity.
