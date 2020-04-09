package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"testing"
)

// test suites
func TestAccGrantActivate_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testAccCheckLocalGrantDestroy, testAccCheckActiveGrantDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantActivateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalGrantExists("turbot_grant.test_grant"),
					testAccCheckActiveGrantExists("turbot_grant_activation.test_activation"),
					resource.TestCheckResourceAttr(
						"turbot_grant.test_grant", "resource", "tmod:@turbot/turbot#/"),
					resource.TestCheckResourceAttr(
						"turbot_grant.test_grant", "type", "tmod:@turbot/turbot-iam#/permission/types/turbot"),
					resource.TestCheckResourceAttr(
						"turbot_grant.test_grant", "level", "tmod:@turbot/turbot-iam#/permission/levels/owner"),
					resource.TestCheckResourceAttr(
						"turbot_grant_activation.test_activation", "resource", "178806508050433"),
				),
			},
		},
	})
}

// configs
func testAccGrantActivateConfig() string {
	return `
resource "turbot_local_directory" "test_dir" {
	parent              = "tmod:@turbot/turbot#/"
	title               = "provider_test_directory"
	description         = "provider_test_directory"
	profile_id_template = "{{profile.email}}"
}

resource "turbot_local_directory_user" "test_user" {
	title        = "Kai Daguerre"
	email        = "kai@turbot.com"
	display_name = "Kai Daguerre"
	parent       = turbot_local_directory.test_dir.id
}

resource "turbot_profile" "test_profile" {
	title             = turbot_local_directory_user.test_user.title
	email             = turbot_local_directory_user.test_user.email
	directory_pool_id = "dpi"
	given_name 		  = "Kai"
	family_name       = "Daguerre"
	display_name      = turbot_local_directory_user.test_user.display_name
	parent            = turbot_local_directory.test_dir.id
	profile_id        = turbot_local_directory_user.test_user.email
}

resource "turbot_grant" "test_grant" {
	resource         = "tmod:@turbot/turbot#/"
	type  = "tmod:@turbot/turbot-iam#/permission/types/turbot"
	level = "tmod:@turbot/turbot-iam#/permission/levels/owner"
	identity          = turbot_profile.test_profile.id
}

resource "turbot_grant_activation" "test_activation" {
	resource = turbot_grant.test_grant.resource
	grant = turbot_grant.test_grant.id
}
`
}

// helper functions
func testAccCheckLocalGrantExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadGrant(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckLocalGrantDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "turbot_grant" {
			_, err := client.ReadGrant(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Alert still exists")
			}
			if !apiClient.NotFoundError(err) {
				return fmt.Errorf("expected 'not found' error, got %s", err)
			}
		}
	}

	return nil
}

func testAccCheckActiveGrantExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := testAccProvider.Meta().(*apiClient.Client)
		_, err := client.ReadGrantActivation(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckActiveGrantDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "Grant" {
			continue
		}
		_, err := client.ReadGrant(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		if !apiClient.NotFoundError(err) {
			return fmt.Errorf("expected 'not found' error, got %s", err)
		}
	}

	return nil
}
