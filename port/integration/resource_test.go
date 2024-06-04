package integration_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func createIntegration(
	installationId string,
) string {
	return fmt.Sprintf(`
	resource "port_integration" "kafkush" {
		installation_id       = "%s"
		title                 = "ZOMG"
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
`, installationId)
}

func createIntegrationWithWebHook(
	installationId string,
	installationAppType string,
) string {
	return fmt.Sprintf(`
	resource "port_integration" "kafkush" {
		installation_id       = "%s"
		title                 = "ZOMG"
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

func TestAccPortIntegrationBasic(t *testing.T) {
	integrationIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortIntegrationResourceBasic = createIntegration(integrationIdentifier)

	var testAccBaseIntegrationPermissionsConfigUpdate = strings.Replace(testAccPortIntegrationResourceBasic, "1.33.7", "1.33.8", -1)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPortIntegrationResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_id", integrationIdentifier),
					resource.TestCheckResourceAttr("port_integration.kafkush", "title", "ZOMG"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "version", "1.33.7"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "webhook_changelog_destination.%", "0"),
				),
			},
			{
				Config: testAccBaseIntegrationPermissionsConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_id", integrationIdentifier),
					resource.TestCheckResourceAttr("port_integration.kafkush", "title", "ZOMG"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "version", "1.33.8"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "webhook_changelog_destination.%", "0"),
				),
			},
		},
	})
}

func TestAccPortIntegrationWithWebhook(t *testing.T) {
	integrationIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortIntegrationResourceBasic = createIntegrationWithWebHook(integrationIdentifier, "KAFKA")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPortIntegrationResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_id", integrationIdentifier),
					resource.TestCheckResourceAttr("port_integration.kafkush", "title", "ZOMG"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "version", "1.33.7"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "installation_app_type", "KAFKA"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "webhook_changelog_destination.%", "2"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "webhook_changelog_destination.url", "https://google.com"),
					resource.TestCheckResourceAttr("port_integration.kafkush", "webhook_changelog_destination.agent", "true"),
				),
			},
		},
	})
}
