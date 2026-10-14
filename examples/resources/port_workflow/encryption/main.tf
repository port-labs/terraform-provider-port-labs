resource "port_workflow" "deploy_with_secrets" {
  identifier = "examples-deploy-with-secrets"
  title      = "Deploy with secrets"

  node {
    identifier = "trigger"

    self_serve_trigger {
      action_card_button_text    = "Deploy"
      execute_action_button_text = "Deploy"

      user_inputs {
        user_properties = {
          string_props = {
            "api_key" = {
              title      = "API key"
              required   = true
              encryption = "aes256-gcm"
            }
            "client_secret" = {
              title = "Client secret"
              client_side_encryption = {
                algorithm = "client-side"
                key       = file("${path.module}/public_key.pem")
              }
            }
          }
          object_props = {
            "kubeconfig" = {
              title      = "Kubeconfig"
              format     = "multi-line"
              encryption = "aes256-gcm"
            }
          }
          array_props = {
            "tags" = {
              title = "Tags"
              string_items = {
                pattern    = "^[a-z][a-z0-9-]*$"
                min_length = 2
                max_length = 32
              }
            }
          }
        }
      }

      permissions {
        roles = ["Member"]
      }
    }
  }

  node {
    identifier = "deploy"

    webhook {
      url = "https://ci.example.com/deploy"
    }
  }

  connections {
    source_identifier = "trigger"
    target_identifier = "deploy"
  }
}
