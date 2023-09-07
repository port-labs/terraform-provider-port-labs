package team_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/internal/acctest"
)

func TestAccPortTeam(t *testing.T) {
	var testAccTeamConfigCreate = `
	resource "port_team" "team" {
		name = "Test"
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
					resource.TestCheckResourceAttr("port_team.team", "name", "Test"),
					resource.TestCheckResourceAttr("port_team.team", "description", "Test description"),
					resource.TestCheckResourceAttr("port_team.team", "users.#", "0"),
				),
			},
		},
	})
}
