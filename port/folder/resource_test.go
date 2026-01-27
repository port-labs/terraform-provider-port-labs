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
		identifier  = "%s"
        title              = "Example Folder"
    }
    `, folderIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccPortFolderResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_folder.example_folder", "title", "Example Folder"),
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
        identifier  = "%s"
        title              = "Example Folder"
    }
    `, folderIdentifier)

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
        identifier  = "%s"
        title              = "Parent Folder"
    }
    `, parentFolderIdentifier)

	var testAccPortFolderResourceChild = fmt.Sprintf(`
    resource "port_folder" "child_folder" {
        identifier  = "%s"
        parent             = port_folder.parent_folder.identifier
        title              = "Child Folder"
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
					resource.TestCheckResourceAttr("port_folder.parent_folder", "title", "Parent Folder"),
					resource.TestCheckResourceAttr("port_folder.child_folder", "identifier", childFolderIdentifier),
					resource.TestCheckResourceAttr("port_folder.child_folder", "parent", parentFolderIdentifier),
					resource.TestCheckResourceAttr("port_folder.child_folder", "title", "Child Folder"),
				),
			},
		},
	})
}

func TestAccPortFolderResourceUpdateFolder(t *testing.T) {
	folderIdentifier := utils.GenID()
	updatedTitle := "Updated Folder Title"
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortFolderResource = fmt.Sprintf(`
    resource "port_folder" "example_folder" {
        identifier  = "%s"
        title              = "Example Folder"
    }
    `, folderIdentifier)

	var testAccPortFolderResourceUpdated = fmt.Sprintf(`
    resource "port_folder" "example_folder" {
        identifier  = "%s"
        title              = "%s"
    }
    `, folderIdentifier, updatedTitle)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccPortFolderResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_folder.example_folder", "identifier", folderIdentifier),
					resource.TestCheckResourceAttr("port_folder.example_folder", "title", "Example Folder"),
				),
			},
			{
				Config: acctest.ProviderConfig + testAccPortFolderResourceUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_folder.example_folder", "identifier", folderIdentifier),
					resource.TestCheckResourceAttr("port_folder.example_folder", "title", updatedTitle),
				),
			},
		},
	})
}

func TestAccPortFolderResourceImport(t *testing.T) {
	folderIdentifier := utils.GenID()
	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}
	var testAccPortFolderResource = fmt.Sprintf(`
    resource "port_folder" "example_folder" {
        identifier  = "%s"
        title              = "Example Folder"
    }
    `, folderIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccPortFolderResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_folder.example_folder", "identifier", folderIdentifier),
					resource.TestCheckResourceAttr("port_folder.example_folder", "title", "Example Folder"),
				),
			},
			{
				ResourceName:      "port_folder.example_folder",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccPortFolderResourceAutomaticParentInheritance tests the fix for the bug where
// folders with 'after' set but 'parent' null would inherit the parent from the 'after' folder.
// This test ensures the provider correctly handles the Port API's automatic parent assignment.
func TestAccPortFolderResourceAutomaticParentInheritance(t *testing.T) {
	parentFolderIdentifier := utils.GenID()
	firstChildIdentifier := utils.GenID()
	secondChildIdentifier := utils.GenID()

	err := os.Setenv("PORT_BETA_FEATURES_ENABLED", "true")
	if err != nil {
		t.Fatal(err)
	}

	// Create a parent folder, a first child with explicit parent,
	// and a second child with 'after' pointing to first child but no parent specified.
	// The second child should automatically inherit the parent from the first child.
	var testAccPortFolderResourceConfig = fmt.Sprintf(`
    resource "port_folder" "parent_folder" {
        identifier = "%s"
        title      = "Parent Folder"
    }

    resource "port_folder" "first_child" {
        identifier = "%s"
        parent     = port_folder.parent_folder.identifier
        title      = "First Child"
    }

    resource "port_folder" "second_child" {
        identifier = "%s"
        after      = port_folder.first_child.identifier
        title      = "Second Child"
        # Note: parent is not explicitly set, but should be inherited from first_child
    }
    `, parentFolderIdentifier, firstChildIdentifier, secondChildIdentifier)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.ProviderConfig + testAccPortFolderResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("port_folder.parent_folder", "identifier", parentFolderIdentifier),
					resource.TestCheckResourceAttr("port_folder.first_child", "identifier", firstChildIdentifier),
					resource.TestCheckResourceAttr("port_folder.first_child", "parent", parentFolderIdentifier),
					resource.TestCheckResourceAttr("port_folder.second_child", "identifier", secondChildIdentifier),
					// This is the key check: second_child should have inherited parent from first_child
					resource.TestCheckResourceAttr("port_folder.second_child", "parent", parentFolderIdentifier),
					resource.TestCheckResourceAttr("port_folder.second_child", "after", firstChildIdentifier),
				),
			},
		},
	})
}
