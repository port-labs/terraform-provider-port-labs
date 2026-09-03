resource "port_llm_provider" "anthropic_compatible" {
  provider_type         = "anthropic-compatible"
  enabled             = true
  validate_connection = false

  anthropic_compatible {
    base_url            = "https://litellm.corp.example.com/anthropic"
    api_key_secret_name = "LITELLM_API_KEY"

    models {
      name           = "claude-sonnet"
      display_name   = "Claude Sonnet"
      context_window = 200000

      supported_features {
        caching                  = true
        native_structured_output = false
      }
    }
  }
}
