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

	resource.Test(t, resource.TestCase{
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

	resource.Test(t, resource.TestCase{
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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

func TestPortIntegrationWithWebhook(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()

	resource.Test(t, resource.TestCase{
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

	resource.Test(t, resource.TestCase{
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

	resource.Test(t, resource.TestCase{
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

	resource.Test(t, resource.TestCase{
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

func TestPortIntegrationImmutableInstallationType(t *testing.T) {
	enableIntegrationBetaFeatures(t)

	installationID := utils.GenID()
	appType := "kafka"
	saasSpec := `
		spec = jsonencode({
			integrationSpec = { token = "my-token" }
		})`

	resource.Test(t, resource.TestCase{
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
				Config: integrationHCLWithInstallationType(installationID, appType, "Saas", saasSpec),
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

	resource.Test(t, resource.TestCase{
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
			resource.Test(t, resource.TestCase{
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

			resource.Test(t, resource.TestCase{
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
