resource "port_blueprint" "microservice" {
  title      = "Microservice"
  icon       = "Microservice"
  identifier = "examples-client-side-encryption-microservice"
  properties = {
    string_props = {
      name = {
        type  = "string"
        title = "Name"
      }
    }
  }
}

# Example: Action with client-side encryption for secret inputs
# The value will be encrypted on the client side using the provided public key
# before being sent to Port. This is useful for sensitive data that should
# never be visible to Port's servers in plaintext.
resource "port_action" "deploy_with_secrets" {
  title      = "Deploy with Secrets"
  icon       = "Terraform"
  identifier = "examples-action-deploy-with-secrets"
  self_service_trigger = {
    operation            = "DAY-2"
    blueprint_identifier = port_blueprint.microservice.identifier
    user_properties = {
      string_props = {
        # Regular string property (not encrypted)
        environment = {
          title       = "Environment"
          description = "Target environment for deployment"
          enum        = ["dev", "staging", "prod"]
          required    = true
        }

        # Client-side encrypted string property
        # The value entered by the user will be encrypted using RSA-OAEP-SHA256
        # with the provided public key before being sent to Port
        api_key = {
          title       = "API Key"
          description = "API key for the deployment (encrypted client-side)"
          required    = true
          client_side_encryption = {
            algorithm = "client-side"
            key       = <<-EOT
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0Z3VS5JJcds3xfn/ygWy
f8sK8pPgbJPbL3pvFkOk9vXwB1QsE0p2LXRJ8ABkOC8fAKCVlCcHWoF7AXxEm+FK
xqMJJO7vOxYXf4cF3bHPzR3pHJnFgAtY3aN/VBMAnTvvvfoUBGhLf0oEGoXmCQbZ
zP3zJIzX/O0G8u0L+wMw9e3CnGWMFYVbq3zOdmGBYVDMnR4lqJMfT3+Qr+w/F6Vf
0jG3x8OXrHVCiNxNv0xHp5zRJvQMW7jDk9frYmOxFvACP/yLMDx/PA/kJxZ0IqSy
hBZ0zfjA3bjZkQfT5NrJmzY5C1t5F6x4sFdHb1e5Kv5VDFMQw5tSMPNhLo3qYVVn
TwIDAQAB
-----END PUBLIC KEY-----
EOT
          }
        }

        # Another client-side encrypted property for database credentials
        db_password = {
          title       = "Database Password"
          description = "Database password (encrypted client-side)"
          required    = true
          client_side_encryption = {
            algorithm = "client-side"
            key       = <<-EOT
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0Z3VS5JJcds3xfn/ygWy
f8sK8pPgbJPbL3pvFkOk9vXwB1QsE0p2LXRJ8ABkOC8fAKCVlCcHWoF7AXxEm+FK
xqMJJO7vOxYXf4cF3bHPzR3pHJnFgAtY3aN/VBMAnTvvvfoUBGhLf0oEGoXmCQbZ
zP3zJIzX/O0G8u0L+wMw9e3CnGWMFYVbq3zOdmGBYVDMnR4lqJMfT3+Qr+w/F6Vf
0jG3x8OXrHVCiNxNv0xHp5zRJvQMW7jDk9frYmOxFvACP/yLMDx/PA/kJxZ0IqSy
hBZ0zfjA3bjZkQfT5NrJmzY5C1t5F6x4sFdHb1e5Kv5VDFMQw5tSMPNhLo3qYVVn
TwIDAQAB
-----END PUBLIC KEY-----
EOT
          }
        }
      }

      object_props = {
        # Client-side encrypted object property for complex secrets
        credentials = {
          title       = "Service Credentials"
          description = "JSON object containing service credentials (encrypted client-side)"
          client_side_encryption = {
            algorithm = "client-side"
            key       = <<-EOT
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0Z3VS5JJcds3xfn/ygWy
f8sK8pPgbJPbL3pvFkOk9vXwB1QsE0p2LXRJ8ABkOC8fAKCVlCcHWoF7AXxEm+FK
xqMJJO7vOxYXf4cF3bHPzR3pHJnFgAtY3aN/VBMAnTvvvfoUBGhLf0oEGoXmCQbZ
zP3zJIzX/O0G8u0L+wMw9e3CnGWMFYVbq3zOdmGBYVDMnR4lqJMfT3+Qr+w/F6Vf
0jG3x8OXrHVCiNxNv0xHp5zRJvQMW7jDk9frYmOxFvACP/yLMDx/PA/kJxZ0IqSy
hBZ0zfjA3bjZkQfT5NrJmzY5C1t5F6x4sFdHb1e5Kv5VDFMQw5tSMPNhLo3qYVVn
TwIDAQAB
-----END PUBLIC KEY-----
EOT
          }
        }
      }
    }
  }
  webhook_method = {
    url = "https://api.example.com/deploy"
  }
}

# Example: Comparing server-side vs client-side encryption
resource "port_action" "encryption_comparison" {
  title      = "Encryption Comparison Example"
  icon       = "Lock"
  identifier = "examples-action-encryption-comparison"
  self_service_trigger = {
    operation            = "CREATE"
    blueprint_identifier = port_blueprint.microservice.identifier
    user_properties = {
      string_props = {
        # Server-side encryption (aes256-gcm)
        # The value is encrypted by Port's servers
        server_encrypted_secret = {
          title       = "Server-Side Encrypted Secret"
          description = "This secret is encrypted by Port's servers using AES-256-GCM"
          encryption  = "aes256-gcm"
        }

        # Client-side encryption (rsa-oaep-sha256)
        # The value is encrypted in the browser before being sent to Port
        client_encrypted_secret = {
          title       = "Client-Side Encrypted Secret"
          description = "This secret is encrypted in the browser using RSA-OAEP-SHA256"
          client_side_encryption = {
            algorithm = "client-side"
            key       = <<-EOT
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0Z3VS5JJcds3xfn/ygWy
f8sK8pPgbJPbL3pvFkOk9vXwB1QsE0p2LXRJ8ABkOC8fAKCVlCcHWoF7AXxEm+FK
xqMJJO7vOxYXf4cF3bHPzR3pHJnFgAtY3aN/VBMAnTvvvfoUBGhLf0oEGoXmCQbZ
zP3zJIzX/O0G8u0L+wMw9e3CnGWMFYVbq3zOdmGBYVDMnR4lqJMfT3+Qr+w/F6Vf
0jG3x8OXrHVCiNxNv0xHp5zRJvQMW7jDk9frYmOxFvACP/yLMDx/PA/kJxZ0IqSy
hBZ0zfjA3bjZkQfT5NrJmzY5C1t5F6x4sFdHb1e5Kv5VDFMQw5tSMPNhLo3qYVVn
TwIDAQAB
-----END PUBLIC KEY-----
EOT
          }
        }
      }
    }
  }
  kafka_method = {
    payload = jsonencode({
      runId = "{{.run.id}}"
    })
  }
}
