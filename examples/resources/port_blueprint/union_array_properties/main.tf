resource "port_blueprint" "vulnerability_tracker" {
  title      = "Vulnerability Tracker"
  icon       = "Terraform"
  identifier = "vulnerability-tracker"
  properties = {
    array_props = {
      vulnerabilities = {
        title        = "Vulnerabilities"
        description  = "CVEs reported by multiple scanners"
        union        = true
        string_items = {}
      }
      scanner_scores = {
        title        = "Scanner Scores"
        union        = true
        number_items = {}
      }
    }
  }
}

resource "port_entity" "service" {
  blueprint = port_blueprint.vulnerability_tracker.identifier
  title     = "Example Service"
  properties = {
    array_props = {
      union_string_slices = {
        vulnerabilities = {
          source_key = "terraform"
          items      = ["CVE-2024-1", "CVE-2024-2"]
        }
      }
      union_number_slices = {
        scanner_scores = {
          source_key = "terraform"
          items      = [9.8, 7.5]
        }
      }
    }
  }
}
