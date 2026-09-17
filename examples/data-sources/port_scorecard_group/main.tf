data "port_scorecard_groups" "all" {}

data "port_scorecard_group" "production_readiness" {
  identifier = "production-readiness"
}

output "scorecard_group_titles" {
  value = [for group in data.port_scorecard_groups.all.groups : group.title]
}

output "production_readiness_title" {
  value = data.port_scorecard_group.production_readiness.title
}
