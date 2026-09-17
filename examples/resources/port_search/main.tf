data "port_search" "microservices_by_text" {
  query = jsonencode({
    combinator = "and"
    rules = [
      { operator = "=", property = "$blueprint", value = "microservice" },
      { operator = "allPropertiesSearch", value = "payments", properties = ["title", "description"] },
    ]
  })
}

data "port_search" "microservices_with_full_text" {
  query = jsonencode({
    combinator = "and"
    rules = [
      { operator = "=", property = "$blueprint", value = "microservice" },
    ]
  })

  full_text_search = {
    term          = "payments"
    target_fields = ["title"]
  }
}
