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

func TestAccPortOrganizationWithHiddenBlueprints(t *testing.T) {
	orgName := utils.GenID()
	bp1 := utils.GenID()
	bp2 := utils.GenID()
	var testAccOrganizationConfigCreate = fmt.Sprintf(`
	resource "port_blueprint" "bp1" {
		identifier = "%s"
		title      = "BP1"
		icon       = "Microservice"
	}
	resource "port_blueprint" "bp2" {
		identifier = "%s"
		title      = "BP2"
		icon       = "Microservice"
	}
	resource "port_organization" "test" {
		name              = "%s"
		hidden_blueprints = [port_blueprint.bp1.identifier, port_blueprint.bp2.identifier]
	}`, bp1, bp2, orgName)

	var testAccOrganizationConfigUpdate = fmt.Sprintf(`
	resource "port_blueprint" "bp1" {
		identifier = "%s"
		title      = "BP1"
		icon       = "Microservice"
	}
	resource "port_blueprint" "bp2" {
		identifier = "%s"
		title      = "BP2"
		icon       = "Microservice"
	}
	resource "port_organization" "test" {
		name              = "%s"
		hidden_blueprints = [port_blueprint.bp1.identifier]
	}`, bp1, bp2, orgName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccOrganizationConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization.test", "name", orgName),
					resource.TestCheckResourceAttr("port_organization.test", "hidden_blueprints.#", "2"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccOrganizationConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization.test", "name", orgName),
					resource.TestCheckResourceAttr("port_organization.test", "hidden_blueprints.#", "1"),
				),
			},
		},
	})
}

func TestAccPortOrganizationWithSettings(t *testing.T) {
	orgName := utils.GenID()
	var testAccOrganizationConfigCreate = fmt.Sprintf(`
	resource "port_organization" "test" {
		name          = "%s"
		portal_title  = "My Dev Portal"
		portal_icon   = "https://example.com/icon.png"
		include_blueprints_in_global_search_by_default = true
	}`, orgName)

	var testAccOrganizationConfigUpdate = fmt.Sprintf(`
	resource "port_organization" "test" {
		name          = "%s"
		portal_title  = "Updated Portal"
		portal_icon   = "https://example.com/icon2.png"
		include_blueprints_in_global_search_by_default = false
	}`, orgName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccOrganizationConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization.test", "name", orgName),
					resource.TestCheckResourceAttr("port_organization.test", "portal_title", "My Dev Portal"),
					resource.TestCheckResourceAttr("port_organization.test", "portal_icon", "https://example.com/icon.png"),
					resource.TestCheckResourceAttr("port_organization.test", "include_blueprints_in_global_search_by_default", "true"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccOrganizationConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization.test", "name", orgName),
					resource.TestCheckResourceAttr("port_organization.test", "portal_title", "Updated Portal"),
					resource.TestCheckResourceAttr("port_organization.test", "portal_icon", "https://example.com/icon2.png"),
					resource.TestCheckResourceAttr("port_organization.test", "include_blueprints_in_global_search_by_default", "false"),
				),
			},
		},
	})
}

func TestAccPortOrganizationWithAnnouncement(t *testing.T) {
	orgName := utils.GenID()
	var testAccOrganizationConfigCreate = fmt.Sprintf(`
	resource "port_organization" "test" {
		name                  = "%s"
		announcement_enabled  = true
		announcement_content  = "Welcome to the portal"
		announcement_link     = "https://example.com"
		announcement_color    = "blue"
	}`, orgName)

	var testAccOrganizationConfigUpdate = fmt.Sprintf(`
	resource "port_organization" "test" {
		name                  = "%s"
		announcement_enabled  = false
		announcement_content  = "Portal under maintenance"
		announcement_color    = "red"
	}`, orgName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccOrganizationConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization.test", "name", orgName),
					resource.TestCheckResourceAttr("port_organization.test", "announcement_enabled", "true"),
					resource.TestCheckResourceAttr("port_organization.test", "announcement_content", "Welcome to the portal"),
					resource.TestCheckResourceAttr("port_organization.test", "announcement_link", "https://example.com"),
					resource.TestCheckResourceAttr("port_organization.test", "announcement_color", "blue"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccOrganizationConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_organization.test", "name", orgName),
					resource.TestCheckResourceAttr("port_organization.test", "announcement_enabled", "false"),
					resource.TestCheckResourceAttr("port_organization.test", "announcement_content", "Portal under maintenance"),
					resource.TestCheckResourceAttr("port_organization.test", "announcement_color", "red"),
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
