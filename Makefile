HOSTNAME=registry.terraform.io
NAMESPACE=port-labs
NAME=port-labs
BINARY=terraform-provider-${NAME}
VERSION=0.9.6
OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)
OS_ARCH=${OS}_${ARCH}
TEST_FILTER?=.*

default: install

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

clean:
	rm -rf examples/.terraform examples/.terraform.lock.hcl examples/terraform*

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

setup:
	cd tools && go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

acctest:
	# TEST_FILTER can be any regex, E.g: .*PageResource.*
	# TEST_FILTER='TestAccPortPageResource*' make acctest
	TF_ACC=1 PORT_CLIENT_ID=$(PORT_CLIENT_ID) PORT_CLIENT_SECRET=$(PORT_CLIENT_SECRET) PORT_BASE_URL=$(PORT_BASE_URL) go test -timeout 40m -p 1 ./... -run "$(TEST_FILTER)"

gen-docs:
	@echo "Generating documentation..."
	@if [ ! -f "./terraform-provider-port-labs" ]; then \
		echo "Provider binary not found. Building..."; \
		go build -o terraform-provider-port-labs; \
	fi
	tfplugindocs generate --provider-name port
	@echo "Documentation generated successfully!"

lint: build
	# https://golangci-lint.run/welcome/install/#local-installation
	golangci-lint run

format: build
	# https://golangci-lint.run/welcome/install/#local-installation
	golangci-lint fmt

dev-run-integration: build
	PORT_BETA_FEATURES_ENABLED=true go run . --debug

dev-setup: setup
	wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v2.8.0

dev-debug: build
	PORT_BETA_FEATURES_ENABLED=true dlv exec --accept-multiclient --continue --headless ./terraform-provider-port-labs -- --debug

