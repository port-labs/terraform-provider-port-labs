package team_test

import (
	"fmt"
	"github.com/port-labs/terraform-provider-port-labs/internal/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/acctest"
)

func TestAccPortTeam(t *testing.T) {
	teamName := utils.GenID()
	var testAccTeamConfigCreate = fmt.Sprintf(`
	resource "port_team" "team" {
		name = "%s"
		description = "Test description"
		users = []
	}`, teamName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccTeamConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_team.team", "name", teamName),
					resource.TestCheckResourceAttr("port_team.team", "description", "Test description"),
					resource.TestCheckResourceAttr("port_team.team", "users.#", "0"),
				),
			},
		},
	})
}

func TestAccPortTeamUpdate(t *testing.T) {
	teamName := utils.GenID()
	var testAccTeamConfigCreate = fmt.Sprintf(`
	resource "port_team" "team" {
		name = "%s"
		description = "Test description"
		users = []
	}`, teamName)

	var testAccTeamConfigUpdate = fmt.Sprintf(`
	resource "port_team" "team" {
		name = "%s"
		description = "Test description2"
		users = ["pluacbcqsnsqhfvcqs@cazlg.com"]
	}`, teamName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccTeamConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_team.team", "name", teamName),
					resource.TestCheckResourceAttr("port_team.team", "description", "Test description"),
					resource.TestCheckResourceAttr("port_team.team", "users.#", "0"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccTeamConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_team.team", "name", teamName),
					resource.TestCheckResourceAttr("port_team.team", "description", "Test description2"),
					resource.TestCheckResourceAttr("port_team.team", "users.#", "1"),
					resource.TestCheckResourceAttr("port_team.team", "users.0", "pluacbcqsnsqhfvcqs@cazlg.com"),
				),
			},
		},
	})
}

func TestAccPortTeamEmptyDescription(t *testing.T) {
	teamName := utils.GenID()
	var testAccTeamConfigCreate = fmt.Sprintf(`
	resource "port_team" "team" {
		name = "%s"
		description = "abc"
		users = []
	}`, teamName)

	var testAccTeamConfigUpdate = fmt.Sprintf(`
	resource "port_team" "team" {
		name = "%s"
		users = []
	}`, teamName)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccTeamConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_team.team", "name", teamName),
					resource.TestCheckResourceAttr("port_team.team", "description", "abc"),
					resource.TestCheckResourceAttr("port_team.team", "users.#", "0"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccTeamConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_team.team", "name", teamName),
					resource.TestCheckNoResourceAttr("port_team.team", "description"),
					resource.TestCheckResourceAttr("port_team.team", "users.#", "0"),
				),
			},
		},
	})
}

func TestAccPortTeamImport(t *testing.T) {
	teamName := utils.GenID()
	var testAccTeamConfigCreate = fmt.Sprintf(`
	resource "port_team" "team" {
		name = "%s"
		description = "Test description"
		users = ["pluacbcqsnsqhfvcqs@cazlg.com"]
	}`, teamName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccTeamConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_team.team", "name", teamName),
					resource.TestCheckResourceAttr("port_team.team", "description", "Test description"),
					resource.TestCheckResourceAttr("port_team.team", "users.#", "1"),
					resource.TestCheckResourceAttr("port_team.team", "users.0", "pluacbcqsnsqhfvcqs@cazlg.com"),
				),
			},
			{
				ResourceName:            "port_team.team",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           teamName,
				ImportStateVerifyIgnore: []string{"provider_name"},
			},
		},
	})
}
