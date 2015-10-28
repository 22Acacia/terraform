package google

import (
	"os"
	"fmt"
	"testing"

	"github.com/22acacia/terraform-gcloud"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPubsubCreate(t *testing.T) {

	if os.Getenv("GOOGLE_GCLOUD_TESTS") == "TRUE" {
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckPubsubDestroy,
			Steps: []resource.TestStep{
				resource.TestStep{
					Config: testAccPubsub,
					Check: resource.ComposeTestCheckFunc(
						testAccPubsubExists(
							"google_pubsub.foobar"),
					),
				},
			},
		})
	}
}

func testAccCheckPubsubDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_pubsub" {
			continue
		}

		fullid, _, err := terraformGcloud.ReadPubsub(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Failed to read pubsub list")
		}

		if fullid == rs.Primary.ID {
			return fmt.Errorf("Pubsub still exists")
		}
	}

	return nil
}

func testAccPubsubExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		found, _, err := terraformGcloud.ReadPubsub(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found != rs.Primary.ID {
			return fmt.Errorf("Pubsub not found.  Found: %q, id: %q", found, rs.Primary.ID)
		}

		return nil
	}
}

const testAccPubsub = `
resource "google_pubsub" "foobar" {
	name = "foobar"
}`
