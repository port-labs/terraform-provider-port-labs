# Contributing to the Port Terraform Provider

## Setting up your developer workspace:

* Have [golang](https://go.dev/doc/install) installed
* Have `wget` installed
* Run `make dev-setup`

## Verifying your contribution:

Be sure to run:

```sh
make lint
```

**NOTE**: Should be installed with `make dev-setup`, if you prefer manually, have a look at the [installation guide](https://golangci-lint.run/welcome/install/#local-installation).


In addition, when changing documentation, run `make gen-docs`.

You can preview how the documentation will look in the Terraform registry with [this tool](https://registry.terraform.io/tools/doc-preview).

## Running your tests

Expose the following environment variables:

`PORT_CLIENT_ID` 

`PORT_CLIENT_SECRET`

`PORT_BASE_URL` - Optional, Port API url

Then run:

```sh
make acctest

# or filtered for your specific test:

TEST_FILTER=.*MyCustomResource.* make acctest
```
## Running your code as the actual terraform provider

```sh
make dev-run-integration
```

Then export the printed `TF_REATTACH_PROVIDERS` environment variable, then your `terraform` will use your running code.

## Debugging your code with `dlv`

Install [Delve](https://github.com/go-delve/delve)

```sh
make dev-debug
```
