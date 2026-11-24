resource "port_virtual_mcp_tool" "list_skills" {
  identifier = "skills_distribution_list_skills"
  title      = "list_skills"
  name       = "list_skills"
  tool       = "list_entities"
  expose_arguments = "listed"
  arguments = jsonencode({
    blueprintIdentifier = { value = "skill" }
    include             = { value = ["$identifier", "$title", "description"] }
  })
}

resource "port_virtual_mcp_server" "skills_distribution" {
  identifier = "skills_distribution"
  title      = "Skills Distribution"
  description = "Distributes organizational skills to MCP clients."
  instructions  = <<-EOT
    Call `list_skills` at the beginning of every session.
    When a task matches a skill, call `load_skill` with that skill's identifier.
  EOT
  passthrough_tools = ["load_skill"]
  virtual_mcp_tools = [port_virtual_mcp_tool.list_skills.identifier]

  depends_on = [port_virtual_mcp_tool.list_skills]
}

output "skills_distribution_url" {
  value       = port_virtual_mcp_server.skills_distribution.url
  description = "OAuth MCP endpoint for the virtual server"
}