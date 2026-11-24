resource "port_virtual_mcp_tool" "list_skills" {
  identifier = "skills_distribution_list_skills"
  title      = "list_skills"
  name       = "list_skills"
  tool       = "list_entities"
  description = "Lists skills visible to the connected user."
  expose_arguments = "listed"
  arguments = jsonencode({
    blueprintIdentifier = { value = "skill" }
    include             = { value = ["$identifier", "$title", "description"] }
    limit               = { value = 1000 }
  })
}