package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"testing"
)

// test suites
func TestAccSmartFolder(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSmartFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartFolderConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSmartFolderExists("turbot_smart_folder.test"),
					resource.TestCheckResourceAttr("turbot_smart_folder.test", "title", "smart_folder"),
					resource.TestCheckResourceAttr("turbot_smart_folder.test", "description", "Smart Folder Testing"),
					resource.TestCheckResourceAttr("turbot_smart_folder.test", "parent", "tmod:@turbot/turbot#/"),
					resource.TestCheckResourceAttr("turbot_smart_folder.test", "filters", "arn:aws:iam::013122550996:user/pratik/accesskey/AKIAQGDRKHTKBON32K3J"),
				),
			},
		},
	})
}

// configs
func testAccSmartFolderConfig() string {
	return `
	resource "turbot_smart_folder" "test" {
		parent  = "tmod:@turbot/turbot#/"
		filter = "arn:aws:iam::013122550996:user/pratik/accesskey/AKIAQGDRKHTKBON32K3J"
		description = "Smart Folder Testing"
		title = "smart_folder"
	}
`
}

func testAccSmartFolderUpdateDescConfig() string {
	return `
	resource "turbot_smart_folder" "test" {
		parent  = "tmod:@turbot/turbot#/"
		filter = "arn:aws:iam::013122550996:user/pratik/accesskey/AKIAQGDRKHTKBON32K3J"
		description = "Smart Folder updated"
		title ="smart_folder"
	}
`
}

// helper functions
func testAccCheckSmartFolderExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiclient.Client)
		_, err := client.ReadSmartFolder(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckSmartFolderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "smartFolder" {
			continue
		}
		_, err := client.ReadSmartFolder(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		if !apiclient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}

	return nil
}
