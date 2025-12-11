package integration_test

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func createIntegration(
	installationId string,
	installationAppType string,
) string {
	return fmt.Sprintf(`
	resource "port_integration" "kafkush" {
		installation_id       = "%s"
		installation_app_type = "%s"
		title                 = "my-kafka-cluster"
		version               = "1.33.7"
		config = jsonencode({
			deleteDependentEntities = true,
			resources = [{
				kind = "ZOMG"
				selector = {
					query = ".title"
				}
				port = {
					entity = {
						mappings = [{
							identifier = "'my-identifier'"
							title      = ".title"
							blueprint  = "'my-blueprint'"
							properties = {
								bla = 123
							}
							relations  = {}
						}]
					}
				}
			}]
		})
	}
`, installationId, installationAppType)
}

func createIntegrationWithWebHook(
	installationId string,
	installationAppType string,
) string {
	return fmt.Sprintf(`
	resource "port_integration" "kafkush" {
		installation_id       = "%s"
		title                 = "my-kafka-cluster"
		version               = "1.33.7"
		installation_app_type = "%s"
		config = jsonencode({
			deleteDependentEntities = true,
			resources = [{
				kind = "ZOMG"
				selector = {
					query = ".title"
				}
				port = {
					entity = {
						mappings = [{
							identifier = "'my-identifier'"
							title      = ".title"
							blueprint  = "'my-blueprint'"
							properties = {
								bla = 123
							}
							relations  = {}
						}]
					}
				}
			}]
		})
		webhook_changelog_destination = {
			type = "WEBHOOK"
			url = "https://google.com"
			agent = true
		}
	}`, installationId, installationAppType)
}

func TestPortIntegrationBasic(t *testing.T) {
	integrationIdentifier := utils.GenID()
	installationAppType := "kafka"
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testPortIntegrationResourceBasic = createIntegration(integrationIdentifier, installationAppType)

	var testAccBaseIntegrationUpdate = strings.Replace(testPortIntegrationResourceBasic, "1.33.7", "1.33.8", -1)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testPortIntegrationResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_id", integrationIdentifier),
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_app_type", installationAppType),
					resource.TestCheckResourceAttr("port_integration.kafkush", "title", "my-kafka-cluster"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "version", "1.33.7"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "webhook_changelog_destination.%", "0"),
				),
			},
			{
				Config: testAccBaseIntegrationUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_id", integrationIdentifier),
					resource.TestCheckResourceAttr("port_integration.kafkush", "title", "my-kafka-cluster"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "version", "1.33.8"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "webhook_changelog_destination.%", "0"),
				),
			},
		},
	})
}

func TestPortIntegrationPatchTitleNull(t *testing.T) {
	integrationIdentifier := utils.GenID()
	installationAppType := "kafka"
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testPortIntegrationResourceBasic = createIntegration(integrationIdentifier, installationAppType)

	var testAccBaseIntegrationUpdate = strings.Replace(testPortIntegrationResourceBasic, "\"my-kafka-cluster\"", "null", -1)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testPortIntegrationResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_id", integrationIdentifier),
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_app_type", installationAppType),
					resource.TestCheckResourceAttr("port_integration.kafkush", "title", "my-kafka-cluster"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "version", "1.33.7"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "webhook_changelog_destination.%", "0"),
				),
			},
			{
				Config: testAccBaseIntegrationUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_id", integrationIdentifier),
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_app_type", installationAppType),
					resource.TestCheckNoResourceAttr("port_integration.kafkush", "title"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "version", "1.33.7"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "webhook_changelog_destination.%", "0"),
				),
			},
		},
	})
}

func TestPortIntegrationWithWebhook(t *testing.T) {
	integrationIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testPortIntegrationResourceBasic = createIntegrationWithWebHook(integrationIdentifier, "kafka")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testPortIntegrationResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_id", integrationIdentifier),
					resource.TestCheckResourceAttr("port_integration.kafkush", "title", "my-kafka-cluster"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "version", "1.33.7"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_app_type", "kafka"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "webhook_changelog_destination.%", "2"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "webhook_changelog_destination.url", "https://google.com"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "webhook_changelog_destination.agent", "true"),
				),
			},
		},
	})
}

func TestPortIntegrationImport(t *testing.T) {
	integrationIdentifier := utils.GenID()
	var testPortIntegrationResourceBasic = createIntegrationWithWebHook(integrationIdentifier, "kafka")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testPortIntegrationResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_app_type", "kafka"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_id", integrationIdentifier),
				),
			},
			{
				ResourceName:      "port_integration.kafkush",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     integrationIdentifier,
			},
		},
	})
}

func TestPortIntegrationInvalidIdentifier(t *testing.T) {
	installationAppType := "kafka"

	testCases := []struct {
		name         string
		identifier   string
		errorPattern string
	}{
		{
			name:         "spaces",
			identifier:   "my integration with spaces",
			errorPattern: `installation_id must match the pattern`,
		},
		{
			name:         "uppercase letters",
			identifier:   "MyIntegration",
			errorPattern: `installation_id must match the pattern`,
		},
		{
			name:         "starts with number",
			identifier:   "123-integration",
			errorPattern: `installation_id must match the pattern`,
		},
		{
			name:         "ends with dash",
			identifier:   "my-integration-",
			errorPattern: `installation_id must match the pattern`,
		},
		{
			name:         "starts with dash",
			identifier:   "-my-integration",
			errorPattern: `installation_id must match the pattern`,
		},
		{
			name:         "special characters",
			identifier:   "my_integration!",
			errorPattern: `installation_id must match the pattern`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := createIntegration(tc.identifier, installationAppType)
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { acctest.TestAccPreCheck(t) },
				ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      config,
						ExpectError: regexp.MustCompile(tc.errorPattern),
					},
				},
			})
		})
	}
}

func TestPortIntegrationValidIdentifier(t *testing.T) {
	installationAppType := "kafka"

	testCases := []struct {
		name       string
		identifier string
	}{
		{
			name:       "simple lowercase",
			identifier: "myintegration",
		},
		{
			name:       "with dashes",
			identifier: "my-integration",
		},
		{
			name:       "with numbers",
			identifier: "my-integration-123",
		},
		{
			name:       "starts with letter ends with number",
			identifier: "integration1",
		},
		{
			name:       "multiple dashes",
			identifier: "my-custom-integration-v2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := createIntegration(tc.identifier, installationAppType)
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { acctest.TestAccPreCheck(t) },
				ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: config,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("port_integration.kafkush", "installation_id", tc.identifier),
							resource.TestCheckResourceAttr("port_integration.kafkush", "installation_app_type", installationAppType),
						),
					},
				},
			})
		})
	}
}
