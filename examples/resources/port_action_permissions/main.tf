resource "port_action_permissions" "restart_microservice_permissions_without_policy" {
  action_identifier = port_action.restart_microservice.identifier
  blueprint_identifier = port_blueprint.microservice.identifier
  permissions = {
    "execute": {
      "roles": [
        "Admin"
      ],
      "users": [],
      "teams": [],
      "owned_by_team": true
    },
    "approve": {
      "roles": ["Member", "Admin"],
      "users": [],
      "teams": []
    }
  }
}

resource "port_action_permissions" "restart_microservice_permissions_with_policy" {
  action_identifier = port_action.restart_microservice.identifier
  blueprint_identifier = port_blueprint.microservice.identifier
  permissions = {
    "execute": {
      "roles": [
        "Admin"
      ],
      "users": [],
      "teams": [],
      "owned_by_team": true
    },
    "approve": {
      "roles": ["Member", "Admin"],
      "users": [],
      "teams": []
      "policy": jsonencode(
        {
          queries: {
            executingUser: {
              rules: [
                {
                  value: "user",
                  operator: "=",
                  property: "$blueprint"
                },
                {
                  value: "{{.trigger.user.email}}",
                  operator: "=",
                  property: "$identifier"
                },
                {
                  value: "true",
                  operator: "=",
                  property: "$owned_by_team"

                }
              ],
              combinator: "or"
            }
          },
          conditions: [
            "true"]
        }
      )
    }
  }
}