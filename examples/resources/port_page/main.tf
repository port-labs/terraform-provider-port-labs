resource "port_blueprint" "microservice" {
  identifier  = "microservice"
  title       = "Microsvc from Port TF Examples"
  icon        = "Terraform"
  description = ""
  properties = {
    string_props = {
      url = {
        type = "string"
      }
      author = {
        icon       = "github"
        required   = true
        min_length = 1
        max_length = 10
        default    = "default"
        enum       = ["default", "default2"]
        pattern    = "^[a-zA-Z0-9]*$"
        format     = "user"
        enum_colors = {
          default  = "red"
          default2 = "green"
        }
      }
    }
  }
}

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
                "id"   = "microservice-table-entities",
                "size" = 12
              }
            ]
          }
        ],
        "type" = "dashboard-widget",
        "widgets" = [
          {
            "id" : "microservice-table-entities",
            "title": "Table test",
            "type" : "table-entities-explorer",
            "blueprint" : port_blueprint.microservice.identifier,
            "displayMode": "widget",
            "dataset" : {
              "combinator" : "and",
              "rules" : [
              ]
            }
          }
        ],
      }
    )
  ]
}
