resource "port_page" "microservice_dashboard_page" {
  identifier = "microservice_dashboard_page"
  title      = "Microservices"
  icon       = "GitHub"
  type       = "dashboard"
  widgets = [
    jsonencode(
      {
        "id" = "dashboardWidget",
        "layout" = [
          {
            "height" = 400,
            "columns" = [
              {
                "id"   = "microserviceGuide",
                "size" = 12
              }
            ]
          }
        ],
        "type" = "dashboard-widget",
        "widgets" = [
          {
            "title"       = "Microservices Guide",
            "icon"        = "BlankPage",
            "markdown"    = "# This is the new Microservice Dashboard",
            "type"        = "markdown",
            "description" = "",
            "id"          = "microserviceGuide"
          }
        ],
      }
    )
  ]
}
