resource "port_llm_provider" "openai_compatible" {
  provider_type         = "openai-compatible"
  enabled               = true
  validate_connection   = false

  openai_compatible {
    base_url            = "https://litellm.corp.example.com/openai"
    api_key_secret_name = "LITELLM_API_KEY"

    models {
      name           = "gpt-4o"
      display_name   = "GPT-4o"
      context_window = 128000

      supported_features {
        use_legacy_max_tokens = false
      }
    }
  }
}
