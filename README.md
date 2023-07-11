<img align="right" src="https://user-images.githubusercontent.com/8277210/183290078-f38cdfd2-e5da-4562-82e6-f274d0330825.svg#gh-dark-mode-only" width="100" height="74" /> <img align="right" width="100" height="74" src="https://user-images.githubusercontent.com/8277210/183290025-d7b24277-dfb4-4ce1-bece-7fe0ecd5efd4.svg#gh-light-mode-only" />

# Port Terraform Provider

[![Slack](https://img.shields.io/badge/Slack-4A154B?style=for-the-badge&logo=slack&logoColor=white)](https://join.slack.com/t/devex-community/shared_invite/zt-1bmf5621e-GGfuJdMPK2D8UN58qL4E_g)

Port is the Developer Platform meant to supercharge your DevOps and Developers, and allow you to regain control of your environment.

## Documentation

- [Terraform registry docs](https://registry.terraform.io/providers/port-labs/port/latest/docs)
- [Port docs](https://docs.getport.io/build-your-software-catalog/sync-data-to-catalog/iac/terraform)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.19 (to build the provider plugin)
- [Port Credentials](https://docs.getport.io/build-your-software-catalog/sync-data-to-catalog/api/#find-your-port-credentials)

## Installation

Terraform utilizes the Terraform Registry to download and install providers. To install the `port` provider, copy and paste the following code into your Terraform file:

```terraform
terraform {
  required_providers {
    port = {
      source  = "port-labs/port-labs"
      version = "~> 1.0.0"
    }
  }
}

provider "port" {
  client_id = "{YOUR CLIENT ID}"     # or set the environment variable PORT_CLIENT_ID
  secret    = "{YOUR CLIENT SECRET}" # or set the environment variable PORT_CLIENT_SECRET
}
```

After you have added the code above, run the following command:

```bash
terraform init
```

## Examples

please refer to the [examples](./examples) directory
