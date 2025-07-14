# Contributing to the Port Terraform Provider

## Setting up your developer workspace:

Our project uses [Devbox](https://github.com/jetpack-io/devbox) to manage its development environment.

Using Devbox will get your dev environment up and running easily and make sure we're all using the same tools with the same versions.

- Install `devbox`

```sh
curl -fsSL https://get.jetpack.io/devbox | bash
```

- Start the `devbox` shell (will take a while on first time)

```sh
devbox shell
```

_This will create a shell where all required tools are installed._

### Optional

- Set up `direnv` so `devbox shell` runs automatically
  - [direnv](https://www.jetify.com/devbox/docs/ide_configuration/direnv/) is a tool that allows setting unique environment variables per directory in your file system.
    - Install `direnv` with: `brew install direnv`
    - Add the following line at the end of the `~/.bashrc` file: `eval "$(direnv hook bash)"`
      - See [direnv's installation instructions](https://direnv.net/docs/hook.html) for other shells.
    - Enable `direnv` by running `direnv allow`
- Install the VSCode Extension
  - Follow [this guide](https://www.jetify.com/devbox/docs/ide_configuration/vscode/) to set up VSCode to automatically run `devbox shell`.

### Manual Setup

If you don't want to use `devbox`:

- Have [golang](https://go.dev/doc/install) installed
- Have `wget` installed
- Run `make dev-setup`

## Verifying your contribution:

Be sure to run:

```sh
make lint
```

**NOTE**: Should be installed with `make dev-setup`, if you prefer manually, have a look at the [installation guide](https://golangci-lint.run/welcome/install/#local-installation).

## Documentation

The resource examples are generated from the schema. This means that in order to update the examples, you will need to update the `ResourceMarkdownDescription` in the resource schema.go file.

After adding a new resource, make sure adding it to provider/provider.go/Resources

In addition, when changing schema, run `make gen-docs`.

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
