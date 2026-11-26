---
page_title: "port_factory_line Resource - port"
subcategory: ""
description: |-
  Manages a Port software factory line via `v1/factory/lines`.
---

# port_factory_line (Resource)

Manages a Port software factory line, including stage definitions and executor settings. Stages can use the `factory_ai_managed_agent` executor (Port PR #23963) to run Factory.ai Droids.

## Example Usage

```terraform
resource "port_factory_line" "example" {
  identifier = "idea-to-production"
  title      = "Idea to production"

  stage {
    identifier   = "build"
    executor     = "factory_ai_managed_agent"
    instructions = "Implement the work item."
    timeout      = "2h"

    factory_ai_managed_agent {
      environment_id   = "my-environment"
      reasoning_effort = "medium"
    }
  }
}
```

## Schema

### Required

- `identifier` (String) Unique line identifier.
- `title` (String) Line title.

### Optional

- `description` (String) Line description.
- `team` (List of String) Team identifiers.
- `max_running` (Number) Maximum concurrent work items (1–1000).
- `status` (String) `active` or `paused`. Defaults to `active`.
- `status_reason` (String) Status explanation.
- `stage` (Block List) Stage definitions (see below).

### Read-Only

- `id` (String) Line identifier.
- `stage_order` (List of String) Ordered stage identifiers from the API.
- `created_at` (String) Creation timestamp.
- `updated_at` (String) Last update timestamp.

### Nested Schema for `stage`

Required:

- `identifier` (String) Stage identifier.
- `executor` (String) One of `claude_managed_agent`, `devin_managed_agent`, `factory_ai_managed_agent`, `custom_workflow`, `human`, `none`.

Optional:

- `title`, `description`, `instructions` (String)
- `team` (List of String)
- `max_tries` (Number, 1–10, default `3`)
- `allow_take_over` (Boolean, default `true`)
- `timeout` (String) Duration such as `2h`, `30m`.
- `checkpoints` (String) JSON array string (use `jsonencode`).

Exactly one executor settings block must match `executor` (not required for `none`).

#### Nested Schema for `stage.factory_ai_managed_agent`

- `installation_id` (String) Factory.ai integration installation ID.
- `environment_id` (String) Factory environment for a new Droid session (xor `computer_id`).
- `computer_id` (String) Existing Factory computer for a new session (xor `environment_id`).
- `model` (String) Model for a new session.
- `reasoning_effort` (String) One of `none`, `dynamic`, `off`, `minimal`, `low`, `medium`, `high`, `xhigh`, `max`.

## Import

Import using the line identifier:

```shell
terraform import port_factory_line.example idea-to-production
```