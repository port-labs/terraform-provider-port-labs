terraform {
  required_providers {
    port = {
      source  = "port-labs/port-labs"
      version = "~> 2"
    }
  }
}
provider "port" {
  client_id = "60EsooJtOqimlekxrNh7nfr2iOgTcyLZ" # or set the environment variable PORT_CLIENT_ID
  secret    = "oCn21vichy9IXn0KQCG03kUDkcsqKSowfsKkn9YhlcsiOiexMDttpyETZmXO9s8m" # or set the environment variable PORT_CLIENT_SECRET
  base_url  = "http://localhost:3000"

}
