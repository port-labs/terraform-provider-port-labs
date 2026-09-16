# Port Terraform Provider Examples

### Getting started:

Edit [`provider.tf`](./provider.tf) in this directory (add your credentials).


### Running examples:

When developing unreleased provider features locally, run `make install` from the repository root, then:

```sh
export TF_CLI_CONFIG_FILE=~/.terraform.d/port-labs-dev.rc
```

For editor validation, copy [`.vscode/settings.json.example`](../.vscode/settings.json.example) to `.vscode/settings.json`.

`cd` into any of the `resources/*` then:

```sh
terraform init
terraform plan
terraform apply
```
