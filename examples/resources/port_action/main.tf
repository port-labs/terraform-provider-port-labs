resource "port_action" "restart_microservice" {
  title      = "Restart microservice"
  icon       = "Terraform"
  identifier = "restart-micrservice"
  blueprint  = port_blueprint.microservice.identifier
  trigger    = "DAY-2"
  webhook_method {
    type = "WEBHOOK"
    url  = "https://app.getport.io"
  }
  user_properties {
    string_props = {
      "webhook_url" = {
        title       = "Webhook URL"
        description = "Webhook URL to send the request to"
        format      = "url"
        default     = "https://example.com"
        pattern     = "^https://.*"
      }
    }
  }
}
