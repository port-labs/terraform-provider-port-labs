resource "port_integration" "my_k8s_exporter" {
  installation_id       = "my-k8s-exporter"
  title                 = "My K8S Exporter with version managed by Terraform"
  installation_app_type = "K8S EXPORTER"
  # NOTE: version can change outside of Terraform.
  # Include this only if you explicitly want to control the version with Terraform.
  version = "1.33.7"
}
