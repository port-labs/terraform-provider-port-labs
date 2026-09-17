# Port Hosted `appSpec` and create options

Examples for `spec.appSpec` feature toggles and `create_port_resources_origin`. Uses the Linear integration as a minimal credential setup; the same `appSpec` fields apply to other Port Hosted integration types.

| Case | Path |
|------|------|
| Feature toggles and resync interval | [`feature-toggles/`](./feature-toggles/) |
| Skip default blueprints and mappings | [`skip-default-resources/`](./skip-default-resources/) |

`appSpec` fields demonstrated:

- `scheduledResyncInterval` — periodic full resync (e.g. `"6h"`, `"12h"`).
- `liveEventsEnabled` — real-time webhook-driven sync.
- `actionsProcessingEnabled` — allow the integration to process Port actions.
- `incrementalSyncEnabled` — scheduled incremental syncs.
- `incrementalSyncInterval` — incremental sync frequency when enabled (`"15m"`, `"30m"`, or `"1h"`; defaults to `"15m"`).

If `appSpec` is omitted, Port applies its own defaults, which may differ from the Port UI.
