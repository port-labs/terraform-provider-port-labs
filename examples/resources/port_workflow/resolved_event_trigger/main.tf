# System workflows use a resolved event trigger to watch entity changes via a resolver URL.
resource "port_workflow" "watch_entities" {
  identifier = "system.watch"
  title      = "Watch entities"

  node {
    identifier = "entity_changed"

    resolved_event_trigger {
      resolved_triggers {
        url = "https://example.com/resolve-entity-events"
      }
    }
  }
}