terraform {
  required_providers {
    port = {
      source  = "port-labs/port-labs"
      version = "~> 2.0.0"
    }
  }
}
provider "port" {
  client_id = "60EsooJtOqimlekxrNh7nfr2iOgTcyLZ" # or set the environment variable PORT_CLIENT_ID
  secret    = "4jkgT8GySA3neg2JRqPCf4FboZmVD62Tp3uWwBQRetZ8pohPwRNRfcSVshaefUaa" # or set the environment variable PORT_CLIENT_SECRET
  base_url  = "http://localhost:3000"

}
