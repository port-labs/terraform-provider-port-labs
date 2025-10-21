package organization_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func TestAccPortOrganizationSecret(t *testing.T) {
	secretName := utils.GenID()
	var testAccOrganizationSecretConfigCreate = fmt.Sprintf(`
	resource "port_organization_secret" "test" {
		secret_name  = "%s"
		secret_value = "test-secret-value"
		description  = "Test organization secret"
	}`, secretName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccOrganizationSecretConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization_secret.test", "secret_name", secretName),
					resource.TestCheckResourceAttr("port_organization_secret.test", "secret_value", "test-secret-value"),
					resource.TestCheckResourceAttr("port_organization_secret.test", "description", "Test organization secret"),
				),
			},
		},
	})
}

func TestAccPortOrganizationSecretUpdate(t *testing.T) {
	secretName := utils.GenID()
	var testAccOrganizationSecretConfigCreate = fmt.Sprintf(`
	resource "port_organization_secret" "test" {
		secret_name  = "%s"
		secret_value = "initial-value"
		description  = "Initial description"
	}`, secretName)

	var testAccOrganizationSecretConfigUpdate = fmt.Sprintf(`
	resource "port_organization_secret" "test" {
		secret_name  = "%s"
		secret_value = "updated-value"
		description  = "Updated description"
	}`, secretName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccOrganizationSecretConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization_secret.test", "secret_name", secretName),
					resource.TestCheckResourceAttr("port_organization_secret.test", "secret_value", "initial-value"),
					resource.TestCheckResourceAttr("port_organization_secret.test", "description", "Initial description"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccOrganizationSecretConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization_secret.test", "secret_name", secretName),
					resource.TestCheckResourceAttr("port_organization_secret.test", "secret_value", "updated-value"),
					resource.TestCheckResourceAttr("port_organization_secret.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccPortOrganizationSecretEmptyDescription(t *testing.T) {
	secretName := utils.GenID()
	var testAccOrganizationSecretConfigCreate = fmt.Sprintf(`
	resource "port_organization_secret" "test" {
		secret_name  = "%s"
		secret_value = "test-value"
		description  = "Initial description"
	}`, secretName)

	var testAccOrganizationSecretConfigUpdate = fmt.Sprintf(`
	resource "port_organization_secret" "test" {
		secret_name  = "%s"
		secret_value = "test-value"
	}`, secretName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccOrganizationSecretConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization_secret.test", "secret_name", secretName),
					resource.TestCheckResourceAttr("port_organization_secret.test", "description", "Initial description"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccOrganizationSecretConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization_secret.test", "secret_name", secretName),
					resource.TestCheckNoResourceAttr("port_organization_secret.test", "description"),
				),
			},
		},
	})
}

func TestAccPortOrganizationSecretImport(t *testing.T) {
	secretName := utils.GenID()
	var testAccOrganizationSecretConfigCreate = fmt.Sprintf(`
	resource "port_organization_secret" "test" {
		secret_name  = "%s"
		secret_value = "test-value"
		description  = "Test description"
	}`, secretName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccOrganizationSecretConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization_secret.test", "secret_name", secretName),
					resource.TestCheckResourceAttr("port_organization_secret.test", "description", "Test description"),
				),
			},
			{
				ResourceName:      "port_organization_secret.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     secretName,
				// secret_value is not returned by the API, so we ignore it in verification
				ImportStateVerifyIgnore: []string{"secret_value"},
			},
		},
	})
}
