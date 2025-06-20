<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://port-graphical-assets.s3.eu-west-1.amazonaws.com/Logo+Typo+%2B+Logo+Symbol+-+white.png">
  <source media="(prefers-color-scheme: light)" srcset="https://port-graphical-assets.s3.eu-west-1.amazonaws.com/Logo+Typo+%2B+Logo+Symbol.svg">
  <img align="right" height="54" src="https://port-graphical-assets.s3.eu-west-1.amazonaws.com/Logo+Typo+%2B+Logo+Symbol.svg">
</picture>

# Port Terraform Provider

[![Slack](https://img.shields.io/badge/Slack-4A154B?style=for-the-badge&logo=slack&logoColor=white)](https://join.slack.com/t/devex-community/shared_invite/zt-1bmf5621e-GGfuJdMPK2D8UN58qL4E_g)
[![Built with Devbox](https://www.jetify.com/img/devbox/shield_galaxy.svg)](https://www.jetify.com/devbox/docs/contributor-quickstart/)

Port is the Developer Platform meant to supercharge your DevOps and Developers, and allow you to regain control of your environment.

> **⚠️ Deprecation Notice**  
> Version 1.x of the Port Terraform Provider is now deprecated.  
> Please upgrade to version 2.x for continued support, new features, and improvements.  
> Refer to the [upgrade instructions](https://github.com/port-labs/terraform-provider-port-labs/releases/tag/v2.0.0) to upgrade to version 2.x.

## Documentation

- [Terraform registry docs](https://registry.terraform.io/providers/port-labs/port-labs/latest/docs)
- [Port docs](https://docs.getport.io/build-your-software-catalog/custom-integration/iac/terraform/)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html)
- [Go](https://golang.org/doc/install) >= 1.24 (to build the provider plugin)
- [Port Credentials](https://docs.getport.io/build-your-software-catalog/sync-data-to-catalog/api/#find-your-port-credentials)

## Installation

Terraform utilizes the Terraform Registry to download and install providers. To install the `port` provider, copy and paste the following code into your Terraform file:

```terraform
terraform {
  required_providers {
    port = {
      source  = "port-labs/port-labs"
      version = "~> 2.0.0"
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

Please refer to the [examples](./examples) directory

## Contributing

Please refer to [contributing.md](./CONTRIBUTING.md)
