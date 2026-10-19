package user_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
)

func TestAccPortUserInactivityTimeout(t *testing.T) {
	userEmail := os.Getenv("CI_USER_NAME")
	if userEmail == "" {
		t.Skip("CI_USER_NAME must be set for user acceptance tests")
	}

	inactivityTimeout := 30
	var testAccUserConfig = fmt.Sprintf(`
	resource "port_user" "test" {
		email              = "%s"
		inactivity_timeout = %d
	}`, userEmail, inactivityTimeout)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_user.test", "email", userEmail),
					resource.TestCheckResourceAttr("port_user.test", "inactivity_timeout", fmt.Sprintf("%d", inactivityTimeout)),
				),
			},
		},
	})
}

func TestAccPortUserImport(t *testing.T) {
	userEmail := os.Getenv("CI_USER_NAME")
	if userEmail == "" {
		t.Skip("CI_USER_NAME must be set for user acceptance tests")
	}

	var testAccUserConfig = fmt.Sprintf(`
	resource "port_user" "test" {
		email = "%s"
	}`, userEmail)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_user.test", "email", userEmail),
				),
			},
			{
				ResourceName:      "port_user.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     userEmail,
			},
		},
	})
}
