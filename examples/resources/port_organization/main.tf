# Example with sensitive variable
variable "api_secret_name" {
  type      = string
}

variable "home_secret_name" {
  type      = string
}
variable "guest_secret_name" {
  type      = string
}
variable "api_key" {
  type      = string
  sensitive = true
}
variable "api_description" {
  type      = string
}

variable "home_url" {
  type      = string
}

variable "guest_url" {
  type      = string
}

variable "home_description" {
  type      = string
}

variable "guest_description" {
  type      = string
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

