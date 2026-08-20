resource "port_oauth_app" "mcp_connector" {
  name         = "My MCP Connector"
  redirect_uri = "https://api.port.io/v1/mcp/oauth2/callback"
}
