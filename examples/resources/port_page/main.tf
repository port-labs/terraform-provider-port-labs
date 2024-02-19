resource "port_page" "microservice_blueprint_page" {
  identifier            = "microservice_blueprint_page"
  title                 = "Microservices"
  type                  = "blueprint-entities"
  icon                  = "Microservice"
  blueprint             = port_blueprint.base_blueprint.identifier
  widgets               = [
    jsonencode(
      {
        "id" : "microservice-table-entities",
        "type" : "table-entities-explorer",
        "dataset" : {
          "combinator" : "and",
          "rules" : [
            {
              "operator" : "=",
              "property" : "$blueprint",
              "value" : "{{`\"{{blueprint}}\"`}}"
            }
          ]
        }
      }
    )
  ]
}


resource "port_page" "microservice_dashboard_page" {
  identifier            = "microservice_dashboard_page"
  title                 = "Microservices"
  icon                  = "GitHub"
  type                  = "dashboard"
  widgets               = [
    jsonencode(
      {
        "id" : "dashboardWidget",
        "layout" : [
          {
            "height" : 400,
            "columns" : [
              {
                "id" : "microserviceGuide",
                "size" : 12
              }
            ]
          }
        ],
        "type" : "dashboard-widget",
        "widgets" : [
          {
            "title" : "Microservices Guide",
            "icon" : "BlankPage",
            "markdown" : "# This is the new Microservice Dashboard",
            "type" : "markdown",
            "description" : "",
            "id" : "microserviceGuide"
          }
        ],
      }
    )
  ]
}
