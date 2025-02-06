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
	sidebarIdentifier := utils.GenID()
	folderIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortFolderResourceBasic = fmt.Sprintf(`

resource "port_folder" "example_folder" {
  sidebar_identifier    = "%s"
  folder_identifier     = "%s"
  title                 = "Example Folder"
  description           = "This is an example folder"
}
`, sidebarIdentifier, folderIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccPortFolderResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_folder.example_folder", "sidebar_identifier", sidebarIdentifier),
					resource.TestCheckResourceAttr("port_folder.example_folder", "folder_identifier", folderIdentifier),
					resource.TestCheckResourceAttr("port_folder.example_folder", "title", "Example Folder"),
					resource.TestCheckResourceAttr("port_folder.example_folder", "description", "This is an example folder"),
				),
			},
		},
	})
}

func TestAccPortFolderResourceBasicBetaDisabled(t *testing.T) {
	sidebarIdentifier := utils.GenID()
	folderIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "false")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortFolderResourceBasic = fmt.Sprintf(`

resource "port_folder" "example_folder" {
  sidebar_identifier    = "%s"
  folder_identifier     = "%s"
  title                 = "Example Folder"
  description           = "This is an example folder"
}
`, sidebarIdentifier, folderIdentifier)

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
	sidebarIdentifier := utils.GenID()
	parentFolderIdentifier := utils.GenID()
	childFolderIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortFolderResourceParent = fmt.Sprintf(`

resource "port_folder" "parent_folder" {
  sidebar_identifier    = "%s"
  folder_identifier     = "%s"
  title                 = "Parent Folder"
  description           = "This is a parent folder"
}
`, sidebarIdentifier, parentFolderIdentifier)

	var testAccPortFolderResourceChild = fmt.Sprintf(`

resource "port_folder" "child_folder" {
  sidebar_identifier    = "%s"
  folder_identifier     = "%s"
  parent                = port_folder.parent_folder.folder_identifier
  title                 = "Child Folder"
  description           = "This is a child folder" 
}
`, sidebarIdentifier, childFolderIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccPortFolderResourceParent + testAccPortFolderResourceChild,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_folder.parent_folder", "sidebar_identifier", sidebarIdentifier),
					resource.TestCheckResourceAttr("port_folder.parent_folder", "folder_identifier", parentFolderIdentifier),
					resource.TestCheckResourceAttr("port_folder.parent_folder", "title", "Parent Folder"),
					resource.TestCheckResourceAttr("port_folder.parent_folder", "description", "This is a parent folder"),
					resource.TestCheckResourceAttr("port_folder.child_folder", "sidebar_identifier", sidebarIdentifier),
					resource.TestCheckResourceAttr("port_folder.child_folder", "folder_identifier", childFolderIdentifier),
					resource.TestCheckResourceAttr("port_folder.child_folder", "parent", parentFolderIdentifier),
					resource.TestCheckResourceAttr("port_folder.child_folder", "title", "Child Folder"),
					resource.TestCheckResourceAttr("port_folder.child_folder", "description", "This is a child folder"),
				),
			},
		},
	})
}
