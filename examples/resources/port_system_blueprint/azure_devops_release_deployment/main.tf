# Azure DevOps Release Deployment system blueprint extension aligned with Port PR #18246.
# https://github.com/port-labs/Port/pull/18246

resource "port_system_blueprint" "azure_devops_release_deployment" {
  identifier = "azureDevopsReleaseDeployment"

  relations = {
    project = {
      title    = "Project"
      target   = "azureDevopsProject"
      required = false
      many     = false
    }
  }
}
