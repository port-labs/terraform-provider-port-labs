resource "port_mcp_server" "client_credentials" {
  title       = "Example MCP Server"
  url         = "https://mcp.example.com/v1"
  auth_method = "oauth"
  oauth_config = {
    client_id     = "my-client-id"
    client_secret = var.mcp_client_secret
    scope         = "read write"
    grant_type    = "client_credentials"
  }
}

variable "mcp_client_secret" {
  type      = string
  sensitive = true
}
