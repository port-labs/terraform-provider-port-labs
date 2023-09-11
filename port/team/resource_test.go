package team_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/acctest"
)

func TestAccPortTeam(t *testing.T) {
	var testAccTeamConfigCreate = `
	resource "port_team" "team" {
		name = "Tf-Test"
		description = "Test description"
		users = []
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccTeamConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_team.team", "name", "Tf-Test"),
					resource.TestCheckResourceAttr("port_team.team", "description", "Test description"),
					resource.TestCheckResourceAttr("port_team.team", "users.#", "0"),
				),
			},
		},
	})
}

func TestAccPortTeamUpdate(t *testing.T) {
	var testAccTeamConfigCreate = `
	resource "port_team" "team" {
		name = "Tf-Test"
		description = "Test description"
		users = []
	}`

	var testAccTeamConfigUpdate = `
	resource "port_team" "team" {
		name = "Test"
		description = "Test description2"
		users = ["devops-port@port-test.io"]
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccTeamConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_team.team", "name", "Tf-Test"),
					resource.TestCheckResourceAttr("port_team.team", "description", "Test description"),
					resource.TestCheckResourceAttr("port_team.team", "users.#", "0"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccTeamConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_team.team", "name", "Tf-Test"),
					resource.TestCheckResourceAttr("port_team.team", "description", "Test description2"),
					resource.TestCheckResourceAttr("port_team.team", "users.#", "1"),
					resource.TestCheckResourceAttr("port_team.team", "users.0", "devops-port@port-test.io"),
				),
			},
		},
	})
}

func TestAccPortTeamImport(t *testing.T) {
	var testAccTeamConfigCreate = `
	resource "port_team" "team" {
		name = "Tf-Test"
		description = "Test description"
		users = ["devops-port@port-test.io"]
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccTeamConfigCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_team.team", "name", "Tf-Test"),
					resource.TestCheckResourceAttr("port_team.team", "description", "Test description"),
					resource.TestCheckResourceAttr("port_team.team", "users.#", "1"),
					resource.TestCheckResourceAttr("port_team.team", "users.0", "devops-port@port-test.io"),
				),
			},
			{
				ResourceName:            "port_team.team",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateId:           "Test",
				ImportStateVerifyIgnore: []string{"provider_name"},
			},
		},
	})
}
