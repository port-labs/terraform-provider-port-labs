package organization_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func TestAccPortOrganization(t *testing.T) {
	orgName := utils.GenID()
	var testAccOrganizationConfigCreate = fmt.Sprintf(`
	resource "port_organization" "test" {
		name = "%s"
	}`, orgName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccOrganizationConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization.test", "name", orgName),
				),
			},
		},
	})
}

func TestAccPortOrganizationUpdate(t *testing.T) {
	orgName := utils.GenID()
	updatedName := utils.GenID()
	var testAccOrganizationConfigCreate = fmt.Sprintf(`
	resource "port_organization" "test" {
		name = "%s"
	}`, orgName)

	var testAccOrganizationConfigUpdate = fmt.Sprintf(`
	resource "port_organization" "test" {
		name = "%s"
	}`, updatedName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccOrganizationConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization.test", "name", orgName),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccOrganizationConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccPortOrganizationImport(t *testing.T) {
	orgName := utils.GenID()
	var testAccOrganizationConfigCreate = fmt.Sprintf(`
	resource "port_organization" "test" {
		name = "%s"
	}`, orgName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccOrganizationConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization.test", "name", orgName),
				),
			},
			{
				ResourceName:      "port_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     orgName,
			},
		},
	})
}
