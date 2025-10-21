# Example with sensitive variable
variable "api_secret_name" {
  type      = string
  default   = "api_key"
}

variable "home_secret_name" {
  type      = string
  default   = "home_url"
}
variable "guest_secret_name" {
  type      = string
  default   = "guest_url"
}
variable "api_key" {
  type      = string
  sensitive = true
  default   = "api_key"
}
variable "api_description" {
  type      = string
  default   = "api description"
}

variable "home_url" {
  type      = string
  default   = "home url"
}

variable "guest_url" {
  type      = string
  default   = "guest url"
}

variable "home_description" {
  type      = string
  default   = "home description"
}

variable "guest_description" {
  type      = string
  default   = "guest description"
}

resource "port_organization_secret" "api_key" {
  secret_name  = var.api_secret_name
  secret_value = var.api_key
  description  = var.api_description
}

resource "port_organization_secret" "home_url" {
  secret_name  = var.home_secret_name
  secret_value = var.home_url
  description  = var.home_description
}

resource "port_organization_secret" "guest_url" {
  secret_name  = var.guest_secret_name
  secret_value = var.guest_url
  description  = var.guest_description
}

