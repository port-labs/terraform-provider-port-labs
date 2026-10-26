package integration_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/cli"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/consts"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

const (
	integrationResourceName = "port_integration.kafkush"
	defaultTitle            = "my-kafka-cluster"
	defaultVersion          = "1.33.7"
)

const integrationMappingConfig = `
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
		})`

const integrationMappingConfigNonAlphabetical = `
		config = jsonencode({
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
								zebra = 123
								alpha = 456
							}
							relations  = {}
						}]
					}
				}
			}]
			deleteDependentEntities = true
		})`

const integrationWebhook = `
		webhook_changelog_destination = {
			type  = "WEBHOOK"
			url   = "https://google.com"
			agent = true
		}`

func integrationHCL(installationID, appType string, blocks ...string) string {
	return integrationHCLWithInstallationType(installationID, appType, "", blocks...)
}

func integrationHCLWithInstallationType(installationID, appType, installationType string, blocks ...string) string {
	extras := strings.Join(blocks, "\n")
	if extras != "" {
		extras = "\n" + extras
	}

	installationTypeLine := ""
	if installationType != "" {
		installationTypeLine = fmt.Sprintf("\n\t\tinstallation_type     = \"%s\"", installationType)
	}

	return fmt.Sprintf(`
	resource "port_integration" "kafkush" {
		installation_id       = "%s"
		installation_app_type = "%s"
		title                 = "%s"
		version               = "%s"%s%s
	}
`, installationID, appType, defaultTitle, defaultVersion, installationTypeLine, extras)
}

func integrationDefaultChecks(installationID, appType string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(integrationResourceName, "installation_id", installationID),
		resource.TestCheckResourceAttr(integrationResourceName, "installation_app_type", appType),
		resource.TestCheckResourceAttr(integrationResourceName, "title", defaultTitle),
		resource.TestCheckResourceAttr(integrationResourceName, "version", defaultVersion),
		resource.TestCheckResourceAttr(integrationResourceName, "webhook_changelog_destination.%", "0"),
	)
}

func integrationCreateThenConfigSteps(installationID, appType, configBlock string, checks ...resource.TestCheckFunc) []resource.TestStep {
	withConfig := integrationHCL(installationID, appType, configBlock)
	allChecks := []resource.TestCheckFunc{integrationDefaultChecks(installationID, appType)}
	allChecks = append(allChecks, checks...)

	return []resource.TestStep{
		{Config: integrationHCL(installationID, appType)},
		{Config: withConfig, Check: resource.ComposeTestCheckFunc(allChecks...)},
	}
}

func enableIntegrationBetaFeatures(t *testing.T) {
	if err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true"); err != nil {
		t.Fatal(err)
	}
}

func TestPortIntegrationConfigKeyOrder(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()
	appType := "kafka"
	withConfig := integrationHCL(installationID, appType, integrationMappingConfigNonAlphabetical)

	steps := integrationCreateThenConfigSteps(installationID, appType, integrationMappingConfigNonAlphabetical)
	steps = append(steps, resource.TestStep{Config: withConfig, PlanOnly: true})

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

func TestPortIntegrationBasic(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()
	appType := "kafka"
	withConfig := integrationHCL(installationID, appType, integrationMappingConfig)

	steps := integrationCreateThenConfigSteps(installationID, appType, integrationMappingConfig)
	steps = append(steps, resource.TestStep{
		Config: strings.ReplaceAll(withConfig, defaultVersion, "1.33.8"),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(integrationResourceName, "installation_id", installationID),
			resource.TestCheckResourceAttr(integrationResourceName, "title", defaultTitle),
			resource.TestCheckResourceAttr(integrationResourceName, "version", "1.33.8"),
			resource.TestCheckResourceAttr(integrationResourceName, "webhook_changelog_destination.%", "0"),
		),
	})

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

func TestPortIntegrationPatchTitleNull(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()
	appType := "kafka"
	withConfig := integrationHCL(installationID, appType, integrationMappingConfig)

	steps := integrationCreateThenConfigSteps(installationID, appType, integrationMappingConfig)
	steps = append(steps, resource.TestStep{
		Config: strings.ReplaceAll(withConfig, fmt.Sprintf(`"%s"`, defaultTitle), "null"),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(integrationResourceName, "installation_id", installationID),
			resource.TestCheckResourceAttr(integrationResourceName, "installation_app_type", appType),
			resource.TestCheckNoResourceAttr(integrationResourceName, "title"),
			resource.TestCheckResourceAttr(integrationResourceName, "version", defaultVersion),
			resource.TestCheckResourceAttr(integrationResourceName, "webhook_changelog_destination.%", "0"),
		),
	})

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

const integrationNotificationChannels = `
		notification_channels = {
			email = {
				enabled = true
			}
			slack = {
				enabled = false
			}
		}`

const integrationNotificationChannelsUpdated = `
		notification_channels = {
			email = {
				enabled = false
			}
			slack = {
				enabled = true
			}
		}`

func TestPortIntegrationNotificationChannels(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()
	secretName := utils.GenID()
	title := "Linear notification channels acctest"
	resourceName := "port_integration.linear"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: saasLinearIntegrationHCL(installationID, secretName, title, integrationNotificationChannels),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "installation_id", installationID),
					resource.TestCheckResourceAttr(resourceName, "notification_channels.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "notification_channels.email.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification_channels.slack.enabled", "false"),
				),
			},
			{
				Config: saasLinearIntegrationHCL(installationID, secretName, title, integrationNotificationChannelsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "notification_channels.email.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "notification_channels.slack.enabled", "true"),
				),
			},
		},
	})
}

func TestPortIntegrationOnPremRejectsNotificationChannels(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()
	config := integrationHCLWithInstallationType(installationID, "kafka", "OnPrem", integrationNotificationChannels)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`notification_channels is only supported when installation_type is`),
			},
		},
	})
}

func TestPortIntegrationWithWebhook(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: integrationHCL(installationID, "kafka", integrationWebhook),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(integrationResourceName, "installation_id", installationID),
					resource.TestCheckResourceAttr(integrationResourceName, "title", defaultTitle),
					resource.TestCheckResourceAttr(integrationResourceName, "version", defaultVersion),
					resource.TestCheckResourceAttr(integrationResourceName, "installation_app_type", "kafka"),
					resource.TestCheckResourceAttr(integrationResourceName, "webhook_changelog_destination.%", "2"),
					resource.TestCheckResourceAttr(integrationResourceName, "webhook_changelog_destination.url", "https://google.com"),
					resource.TestCheckResourceAttr(integrationResourceName, "webhook_changelog_destination.agent", "true"),
				),
			},
		},
	})
}

func TestPortIntegrationImport(t *testing.T) {
	installationID := utils.GenID()
	config := integrationHCL(installationID, "kafka", integrationWebhook)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(integrationResourceName, "installation_app_type", "kafka"),
					resource.TestCheckResourceAttr(integrationResourceName, "installation_id", installationID),
				),
			},
			{
				ResourceName:      integrationResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     installationID,
			},
		},
	})
}

func TestPortIntegrationImmutableInstallationId(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()
	newInstallationID := utils.GenID()
	appType := "kafka"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: integrationHCL(installationID, appType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(integrationResourceName, "installation_id", installationID),
				),
			},
			{
				Config:      integrationHCL(newInstallationID, appType),
				ExpectError: regexp.MustCompile(`cannot change installation_id`),
			},
		},
	})
}

func TestPortIntegrationImmutableInstallationAppType(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: integrationHCL(installationID, "kafka"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(integrationResourceName, "installation_id", installationID),
					resource.TestCheckResourceAttr(integrationResourceName, "installation_app_type", "kafka"),
				),
			},
			{
				Config:      integrationHCL(installationID, "pagerduty"),
				ExpectError: regexp.MustCompile(`cannot change installation_app_type`),
			},
		},
	})
}

func TestPortIntegrationCreatePortResourcesOriginEmpty(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()
	appType := "kafka"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: integrationHCL(installationID, appType, `create_port_resources_origin = "Empty"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(integrationResourceName, "installation_id", installationID),
					resource.TestCheckResourceAttr(integrationResourceName, "create_port_resources_origin", "Empty"),
				),
			},
		},
	})
}

func TestPortIntegrationImmutableCreatePortResourcesOrigin(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()
	appType := "kafka"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: integrationHCL(installationID, appType, `create_port_resources_origin = "Empty"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(integrationResourceName, "create_port_resources_origin", "Empty"),
				),
			},
			{
				Config:      integrationHCL(installationID, appType, `create_port_resources_origin = "Port"`),
				ExpectError: regexp.MustCompile(`cannot change create_port_resources_origin`),
			},
		},
	})
}

func TestPortIntegrationInvalidCreatePortResourcesOrigin(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()
	appType := "kafka"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      integrationHCL(installationID, appType, `create_port_resources_origin = "Ocean"`),
				ExpectError: regexp.MustCompile(`create_port_resources_origin`),
			},
		},
	})
}

func TestPortIntegrationImmutableInstallationType(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()
	appType := "kafka"
	saasSpec := `
		spec = jsonencode({
			integrationSpec = { token = "my-token" }
		})`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: integrationHCLWithInstallationType(installationID, appType, "OnPrem"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(integrationResourceName, "installation_id", installationID),
					resource.TestCheckResourceAttr(integrationResourceName, "installation_type", "OnPrem"),
				),
			},
			{
				Config:      integrationHCLWithInstallationType(installationID, appType, "Saas", saasSpec),
				ExpectError: regexp.MustCompile(`cannot change installation_type`),
			},
		},
	})
}

func TestPortIntegrationOnPremRejectsSpec(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()
	config := integrationHCLWithInstallationType(
		installationID,
		"kafka",
		"OnPrem",
		`spec = jsonencode({ integrationSpec = { token = "my-token" } })`,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`spec is only supported when installation_type is`),
			},
		},
	})
}

func TestPortIntegrationInvalidIdentifier(t *testing.T) {
	appType := "kafka"

	testCases := []struct {
		name         string
		identifier   string
		errorPattern string
	}{
		{name: "spaces", identifier: "my integration with spaces", errorPattern: `installation_id must match the pattern`},
		{name: "uppercase letters", identifier: "MyIntegration", errorPattern: `installation_id must match the pattern`},
		{name: "special characters underscore", identifier: "my_integration", errorPattern: `installation_id must match the pattern`},
		{name: "special characters exclamation", identifier: "my-integration!", errorPattern: `installation_id must match the pattern`},
		{name: "special characters at", identifier: "my@integration", errorPattern: `installation_id must match the pattern`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:                 func() { acctest.TestAccPreCheck(t) },
				ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      integrationHCL(tc.identifier, appType),
						ExpectError: regexp.MustCompile(tc.errorPattern),
					},
				},
			})
		})
	}
}

func saasLinearIntegrationHCL(installationID, secretName, title string, blocks ...string) string {
	extras := strings.Join(blocks, "\n")
	if extras != "" {
		extras = "\n" + extras
	}

	return fmt.Sprintf(`
%s
resource "port_organization_secret" "linear_api_key" {
	secret_name  = "%s"
	secret_value = "lin_api_acctest_dummy"
}

resource "port_integration" "linear" {
	depends_on = [port_organization_secret.linear_api_key]

	installation_id       = "%s"
	installation_app_type = "linear"
	installation_type     = "Saas"
	title                 = "%s"

	spec = jsonencode({
		integrationSpec = {
			linearApiKey = port_organization_secret.linear_api_key.secret_name
		}
		appSpec = {
			liveEventsEnabled       = false
			scheduledResyncInterval = "12h"
		}
	})%s
}`, acctest.ProviderConfig, secretName, installationID, title, extras)
}

func TestPortIntegrationSaasLinear(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()
	secretName := utils.GenID()
	title := "Linear acctest"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: saasLinearIntegrationHCL(installationID, secretName, title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_integration.linear", "installation_id", installationID),
					resource.TestCheckResourceAttr("port_integration.linear", "installation_app_type", "linear"),
					resource.TestCheckResourceAttr("port_integration.linear", "installation_type", "Saas"),
					resource.TestCheckResourceAttr("port_integration.linear", "title", title),
					resource.TestCheckResourceAttrSet("port_integration.linear", "spec"),
				),
			},
			{
				Config: saasLinearIntegrationHCL(installationID, secretName, title+" updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_integration.linear", "installation_id", installationID),
					resource.TestCheckResourceAttr("port_integration.linear", "title", title+" updated"),
				),
			},
		},
	})
}

func TestPortIntegrationValidIdentifier(t *testing.T) {
	appType := "kafka"

	testCases := []struct {
		name    string
		buildID func(uniqueID string) string
	}{
		{"simple lowercase", func(id string) string { return "myintegration-" + id }},
		{"with dashes", func(id string) string { return "my-integration-" + id }},
		{"with numbers", func(id string) string { return "my-integration-123-" + id }},
		{"starts with letter ends with number", func(id string) string { return "integration1-" + id }},
		{"multiple dashes", func(id string) string { return "my-custom-integration-v2-" + id }},
		{"starts with number", func(id string) string { return "123-integration-" + id }},
		{"ends with dash", func(id string) string { return "my-integration-" + id + "-" }},
		{"starts with dash", func(id string) string { return "-" + id + "-my-integration" }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			installationID := tc.buildID(utils.GenID())

			resource.ParallelTest(t, resource.TestCase{
				PreCheck:                 func() { acctest.TestAccPreCheck(t) },
				ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: integrationHCL(installationID, appType),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(integrationResourceName, "installation_id", installationID),
							resource.TestCheckResourceAttr(integrationResourceName, "installation_app_type", appType),
						),
					},
				},
			})
		})
	}
}

// createIntegrationOutsideTerraform registers an integration as a self-hosted Ocean integration does when it
// starts, and deletes it when the test ends
func createIntegrationOutsideTerraform(t *testing.T, installationID, appType string) {
	t.Helper()
	acctest.SkipUnlessAcceptance(t)
	portClient := acctest.PortClient(t)
	title, version := "Registered outside Terraform", defaultVersion
	installationType := consts.InstallationTypeOnPrem
	_, err := portClient.CreateIntegration(context.Background(), &cli.Integration{
		InstallationId:      installationID,
		Title:               &title,
		Version:             &version,
		InstallationAppType: &appType,
		InstallationType:    &installationType,
	})
	if err != nil {
		t.Fatalf("failed to create integration: %s", err)
	}
	t.Cleanup(func() { _, _ = portClient.DeleteIntegration(context.Background(), installationID) })
}

func TestAccPortIntegrationAdoptsExistingIntegration(t *testing.T) {
	enableIntegrationBetaFeatures(t)
	installationID := utils.GenID()
	appType := "kafka"
	createIntegrationOutsideTerraform(t, installationID, appType)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      acctest.ProviderConfig + integrationHCL(installationID, appType),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("integration already exists"),
			},
			{
				Config: acctest.ProviderConfigWithAdoption("port_integration") + integrationHCL(installationID, appType),
				Check:  integrationDefaultChecks(installationID, appType),
			},
		},
	})
}

func TestAccPortIntegrationAdoptionRejectsDifferentAppType(t *testing.T) {
	enableIntegrationBetaFeatures(t)
	installationID := utils.GenID()
	createIntegrationOutsideTerraform(t, installationID, "kafka")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      acctest.ProviderConfigWithAdoption("port_integration") + integrationHCL(installationID, "gitlab"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("cannot take over existing integration"),
			},
		},
	})
}
