# Associates an existing scorecard (managed outside the group) with a scorecard group.
resource "port_scorecard_group_scorecard" "example" {
  group_identifier     = "production-readiness"
  scorecard_identifier = "legacy-readiness"
  override_group       = true
}
