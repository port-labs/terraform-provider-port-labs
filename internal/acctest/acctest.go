package acctest

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/provider"
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		consts.ProviderName: providerserver.NewProtocol6WithError(provider.New()),
	}
)

var ProviderConfig = fmt.Sprintf(`provider "port" {
	client_id = "%s"
	secret = "%s"
	base_url = "%s"
	}
`, os.Getenv("PORT_CLIENT_ID"), os.Getenv("PORT_CLIENT_SECRET"), os.Getenv("PORT_BASE_URL"))

var ProviderConfigNoPropertyTypeProtection = fmt.Sprintf(`provider "port" {
	client_id = "%s"
	secret = "%s"
	base_url = "%s"
	blueprint_property_type_change_protection = false
	}
`, os.Getenv("PORT_CLIENT_ID"), os.Getenv("PORT_CLIENT_SECRET"), os.Getenv("PORT_BASE_URL"))

var ProviderConfigNoEscapeHTML = fmt.Sprintf(`provider "port" {
	client_id = "%s"
	secret = "%s"
	base_url = "%s"
	json_escape_html = false
	}
`, os.Getenv("PORT_CLIENT_ID"), os.Getenv("PORT_CLIENT_SECRET"), os.Getenv("PORT_BASE_URL"))

func TestAccPreCheck(t *testing.T) {
	if v := os.Getenv("PORT_CLIENT_ID"); v == "" {
		t.Fatal("PORT_CLIENT_ID must be set for acceptance tests")
	}

	if v := os.Getenv("PORT_CLIENT_SECRET"); v == "" {
		t.Fatal("PORT_CLIENT_SECRET must be set for acceptance tests")
	}
}
