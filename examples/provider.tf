terraform {
  required_providers {
    port = {
      source  = "port-labs/port-labs"
      version = "~> 2.0.0"
    }
  }
}
provider "port" {
  client_id = "60EsooJtOqimlekxrNh7nfr2iOgTcyLZ"                                 # or set the environment variable PORT_CLIENT_ID
  secret    = "35D7Hw4ZpjdHW0u1lNS0cE5UXvevhlGQWeXuwkIX91s6UjgLzO44GSBG9yNBdehr" # or set the environment variable PORT_CLIENT_SECRET
  base_url  = "http://localhost:3000"
}
