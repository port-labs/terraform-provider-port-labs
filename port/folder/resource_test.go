package folder_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/acctest"
	"github.com/port-labs/terraform-provider-port-labs/v2/internal/utils"
)

func TestAccPortFolderResourceBasicBetaEnabled(t *testing.T) {
	folderIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortFolderResourceBasic = fmt.Sprintf(`

resource "port_folder" "example_folder" {
  identifier            = "%s"
  name                  = "Example Folder"
  description           = "This is an example folder"
}
`, folderIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccPortFolderResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_folder.example_folder", "identifier", folderIdentifier),
					resource.TestCheckResourceAttr("port_folder.example_folder", "name", "Example Folder"),
					resource.TestCheckResourceAttr("port_folder.example_folder", "description", "This is an example folder"),
				),
			},
		},
	})
}

func TestAccPortFolderResourceBasicBetaDisabled(t *testing.T) {
	folderIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "false")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortFolderResourceBasic = fmt.Sprintf(`

resource "port_folder" "example_folder" {
  identifier            = "%s"
  name                  = "Example Folder"
  description           = "This is an example folder"
}
`, folderIdentifier)

	// expect to fail on beta feature not enabled
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      acctest.ProviderConfig + testAccPortFolderResourceBasic,
				ExpectError: regexp.MustCompile("Beta features are not enabled"),
			},
		},
	})
}

func TestAccPortFolderResourceCreateFolderWithParent(t *testing.T) {
	parentFolderIdentifier := utils.GenID()
	childFolderIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortFolderResourceParent = fmt.Sprintf(`

resource "port_folder" "parent_folder" {
  identifier            = "%s"
  name                  = "Parent Folder"
  description           = "This is a parent folder"
}
`, parentFolderIdentifier)

	var testAccPortFolderResourceChild = fmt.Sprintf(`

resource "port_folder" "child_folder" {
  identifier            = "%s"
  name                  = "Child Folder"
  parent                = port_folder.parent_folder.identifier
  description           = "This is a child folder"
}
`, childFolderIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccPortFolderResourceParent + testAccPortFolderResourceChild,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_folder.parent_folder", "identifier", parentFolderIdentifier),
					resource.TestCheckResourceAttr("port_folder.parent_folder", "name", "Parent Folder"),
					resource.TestCheckResourceAttr("port_folder.parent_folder", "description", "This is a parent folder"),
					resource.TestCheckResourceAttr("port_folder.child_folder", "identifier", childFolderIdentifier),
					resource.TestCheckResourceAttr("port_folder.child_folder", "name", "Child Folder"),
					resource.TestCheckResourceAttr("port_folder.child_folder", "parent", parentFolderIdentifier),
					resource.TestCheckResourceAttr("port_folder.child_folder", "description", "This is a child folder"),
				),
			},
		},
	})
}
