# Blueprint definitions

resource "port_blueprint" "Microservice" {
  identifier  = "Microservice"
  title       = "Microservice"
  icon        = "Microservice"
  description = "Microservice entities"

  properties = {
    string_props = {
      language = {
        type = "string"
      }
    }
  }
}

resource "port_blueprint" "cluster" {
  identifier  = "Cluster"
  title       = "Cluster"
  icon        = "Cluster"
  description = "Cluster entities"

  properties = {
    string_props = {
      studio = {
        type = "string"
      }
    }
  }
}

# Entity definitions

resource "port_entity" "service_ruby" {
  title     = "Ruby Microservice"
  blueprint = port_blueprint.Microservice.identifier

  properties = {
    string_props = {
      language = "Ruby"
    }
  }

  depends_on = [port_blueprint.Microservice]
}

resource "port_entity" "service_python" {
  title     = "Python Microservice"
  blueprint = port_blueprint.Microservice.identifier

  properties = {
    string_props = {
      language = "Python"
    }
  }

  depends_on = [port_blueprint.Microservice]
}

resource "port_entity" "cluster_python" {
  title     = "Python Cluster"
  blueprint = port_blueprint.cluster.identifier

  properties = {
    string_props = {
      studio = "Cool studio"
    }
  }

  depends_on = [port_blueprint.cluster]
}

resource "port_entity" "cluster_java" {
  title     = "Java Cluster"
  blueprint = port_blueprint.cluster.identifier

  properties = {
    string_props = {
      studio = "Even cooler studio"
    }
  }

  depends_on = [port_blueprint.cluster]
}

resource "port_entity" "cluster_java_new" {
  title     = "Another Cluster"
  blueprint = port_blueprint.cluster.identifier

  properties = {
    string_props = {
      studio = "Even more cooler studio"
    }
  }

  depends_on = [port_blueprint.cluster]
}

# Page with filters example

resource "port_page" "page_with_filters" {
  identifier = "page_with_filters"
  title      = "Page Filters example"
  icon       = "Apps"
  type       = "dashboard"

  depends_on = [
    port_blueprint.Microservice,
    port_blueprint.cluster,
    port_entity.service_ruby,
    port_entity.service_python,
    port_entity.cluster_python,
    port_entity.cluster_java
  ]

  page_filters = [
    jsonencode(
      {
        "identifier" = "584d867a-a0bc-4880-bcce-f0e62eca4905"
        "title"      = "Microservice: language = Ruby"
        "query" = {
          "combinator" = "and"
          "rules" = [
            {
              "value"    = "Ruby"
              "property" = "language"
              "operator" = "="
            }
          ]
          "blueprint" = "Microservice"
        }
      }
    ),
    jsonencode(
      {
        "identifier" = "3dbeb1e9-fbdb-408e-8db3-e4adf2572a7d"
        "title"      = "Cluster: studio != Cool studio"
        "query" = {
          "combinator" = "and"
          "rules" = [
            {
              "value"    = "Cool studio"
              "property" = "studio"
              "operator" = "!="
            }
          ]
          "blueprint" = "Cluster"
        }
      }
    )
  ]

  widgets = [
    jsonencode(
      {
        "id"   = "6cb29856-9547-44a8-a5a4-1640fda12745"
        "type" = "dashboard-widget"
        "layout" = [
          {
            "height" = 616
            "columns" = [
              {
                "id"   = "2Q3L3IrtuiPQJUyN"
                "size" = 6
              },
              {
                "id"   = "QocSYJUnJdQLNRik"
                "size" = 3
              },
              {
                "id"   = "ssZ3dPYBdlT3rO5z"
                "size" = 3
              }
            ]
          },

        ]
        "widgets" = [
          {
            "id"             = "2Q3L3IrtuiPQJUyN"
            "type"           = "table-entities-explorer"
            "displayMode"    = "widget"
            "title"          = "Clusters table"
            "excludedFields" = []
            "description"    = ""
            "emptyStateText" = ""
            "icon"           = "Table"
            "dataset" = {
              "combinator" = "and"
              "rules"      = []
            }
            "blueprintConfig" = {
              "Cluster" = {
                "filterSettings" = {
                  "filterBy" = {
                    "combinator" = "and"
                    "rules"      = []
                  }
                }
                "groupSettings" = {
                  "groupBy" = []
                }
                "sortSettings" = {
                  "sortBy" = []
                }
                "propertiesSettings" = {
                  "order" = []
                  "shown" = [
                    "$title",
                    "studio",
                    "$team"
                  ]
                }
              }
            }
            "blueprint" = "Cluster"
          },
          {
            "id"             = "QocSYJUnJdQLNRik"
            "type"           = "entities-pie-chart"
            "blueprint"      = "Cluster"
            "title"          = "Clusters by Studio"
            "description"    = ""
            "emptyStateText" = ""
            "icon"           = "Pie"
            "property"       = "property#studio"
            "dataset" = {
              "combinator" = "and"
              "rules"      = []
            }
          },
          {
            "id"             = "ssZ3dPYBdlT3rO5z"
            "type"           = "entities-pie-chart"
            "blueprint"      = "Microservice"
            "title"          = "Services by language"
            "description"    = ""
            "emptyStateText" = ""
            "icon"           = "Pie"
            "property"       = "property#language"
            "dataset" = {
              "combinator" = "and"
              "rules"      = []
            }
          }
        ]
      }
    )
  ]
}
