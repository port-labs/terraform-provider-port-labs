resource "port_blueprint" "parent_blueprint" {
  title       = "Parent Blueprint"
  icon        = "Terraform"
  identifier  = "parent"
  description = ""
  properties = {
    number_props = {
      "age" = {
        title = "Age"
      }
    }
  }
}

resource "port_blueprint" "child_blueprint" {
  title       = "Child Blueprint"
  icon        = "Terraform"
  identifier  = "child"
  description = ""
  properties = {
    number_props = {
      "age" = {
        title = "Age"
      }
    }
  }
  relations = {
    "parent" = {
      title  = "Parent"
      target = port_blueprint.parent_blueprint.identifier
    }
  }
}
resource "port_aggregation_properties" "parent_aggregation_properties" {
  blueprint_identifier = port_blueprint.parent_blueprint.identifier
  properties = {
    "count_kids" = {
      target_blueprint_identifier = port_blueprint.child_blueprint.identifier
      title                       = "Count Kids"
      icon                        = "Terraform"
      description                 = "Count Kids"
      method = {
        count_entities = true
      }
      path_filter = [
        {
          from_blueprint = port_blueprint.child_blueprint.identifier
          path          = ["parent"]
        }
      ]
    }
  }
}
