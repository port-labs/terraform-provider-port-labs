resource "port-labs_action" "restart_microservice" {
		title = "Restart microservice"
		icon = "Terraform"
		identifier = "restart-micrservice"
		blueprint_identifier = port-labs_blueprint.microservice.identifier
		trigger = "DAY-2"
		invocation_method = "KAFKA"
		user_properties {
			identifier = "webhook_url"
			type = "string"
			title = "Webhook URL"
			description = "Webhook URL to send the request to"
			format = "url"
			default = "https://example.com"
			pattern = "^https://.*"
		}
	}