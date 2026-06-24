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

resource "port_entity" "example_microservice" {
  title     = "Example Microservice"
  blueprint = port_blueprint.microservice.identifier

  properties = {
    string_props = {
      url    = "https://example.com"
      author = "default"
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

resource "port_page" "microservice_entity_page" {
  identifier = "microserviceEntity"
  title      = "Microsvc from Port TF Examples"
  icon       = "Terraform"
  type       = "entity"
  blueprint  = port_blueprint.microservice.identifier

  depends_on = [port_blueprint.microservice]

  page_filters = [
    jsonencode(
      {
        "identifier" = "fac6b5aa-272c-4a20-9635-add07d097bb9"
        "title"      = "Entity Creation Date is in the past 30 days"
        "query" = {
          "combinator" = "and"
          "rules" = [
            {
              "property" = "$createdAt"
              "operator" = "between"
              "value" = {
                "preset" = "lastMonth"
              }
            }
          ]
          "blueprint" = "dashboard-filters-meta-blueprint"
        }
      }
    ),
  ]

  widgets = [
    jsonencode(
      {
        "id"                  = "entityPageGrouper",
        "type"                = "grouper",
        "displayMode"         = "tabs",
        "activeGroupUrlParam" = "activeTab",
        "groupsOrder" = [
          "Overview",
          "Related Entities",
          "Runs",
          "Audit Log",
        ],
        "groups" = [
          {
            "title" = "Overview",
            "widgets" = [
              {
                "id"   = "overviewDashboard",
                "type" = "dashboard-widget",
                "layout" = [
                  {
                    "height" = 400,
                    "columns" = [
                      {
                        "id"   = "entityDetails",
                        "size" = 12,
                      },
                    ],
                  },
                ],
                "widgets" = [
                  {
                    "id"        = "entityDetails",
                    "type"      = "entity-info",
                    "title"     = "Details",
                    "blueprint" = "{{blueprint}}",
                    "entity"    = "{{url.identifier}}",
                  },
                ],
              },
            ],
          },
          {
            "title" = "Related Entities",
            "widgets" = [
              {
                "id"          = "relatedEntitiesGrouper",
                "type"        = "grouper",
                "title"       = "Related Entities",
                "displayMode" = "switch",
                "groups" = [
                  {
                    "title" = "Table",
                    "icon"  = "Table",
                    "widgets" = [
                      {
                        "id"   = "relatedTable",
                        "title" = "Related Entities Table",
                        "type" = "table-entities-explorer-by-direction",
                      },
                    ],
                  },
                  {
                    "title" = "Graph",
                    "icon"  = "Relation",
                    "widgets" = [
                      {
                        "id"               = "relatedGraph",
                        "type"             = "graph-entities-explorer",
                        "hiddenBlueprints" = [],
                        "dataset" = {
                          "combinator" = "or",
                          "rules" = [
                            {
                              "operator"  = "relatedTo",
                              "value"     = "{{url.identifier}}",
                              "blueprint" = "{{blueprint}}",
                            },
                            {
                              "combinator" = "and",
                              "rules" = [
                                {
                                  "operator" = "=",
                                  "value"    = "{{url.identifier}}",
                                  "property" = "$identifier",
                                },
                                {
                                  "operator" = "=",
                                  "value"    = "{{blueprint}}",
                                  "property" = "$blueprint",
                                },
                              ],
                            },
                          ],
                        },
                      },
                    ],
                  },
                ],
              },
            ],
          },
          {
            "title" = "Runs",
            "widgets" = [
              {
                "id"    = "runsTable",
                "type"  = "runs-table",
                "title" = "Run Log",
                "query" = {
                  "entity"    = "{{url.identifier}}",
                  "blueprint" = "{{blueprint}}",
                },
              },
            ],
          },
          {
            "title" = "Audit Log",
            "widgets" = [
              {
                "id"   = "auditLogTable",
                "type" = "table-audit-log",
                "query" = {
                  "entity"    = "{{url.identifier}}",
                  "blueprint" = "{{blueprint}}",
                },
              },
            ],
          },
        ],
      }
    ),
  ]
}
