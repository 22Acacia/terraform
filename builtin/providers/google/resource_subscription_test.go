package google

import (
	"os"
	"fmt"
	"testing"

	"github.com/22acacia/terraform-gcloud"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSubscriptionCreate(t *testing.T) {
	if os.Getenv("GOOGLE_GCLOUD_TESTS") == "TRUE" {

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSubscriptionDestroy,
			Steps: []resource.TestStep{
				resource.TestStep{
					Config: testAccSubscription,
					Check: resource.ComposeTestCheckFunc(
						testAccSubscriptionExists(
							"google_subscription.foobar"),
					),
				},
			},
		})
	} 
}

func testAccCheckSubscriptionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_subscription" {
			continue
		}

		fullid, err := terraformGcloud.ReadSubscription(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Failed to read subscription list")
		}

		if fullid == rs.Primary.ID {
			return fmt.Errorf("Subscription still exists")
		}
	}

	return nil
}

func testAccSubscriptionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		found, err := terraformGcloud.ReadSubscription(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found != rs.Primary.ID {
			return fmt.Errorf("Subscription not found.  Found: %q, id: %q", found, rs.Primary.ID)
		}

		return nil
	}
}

const testAccSubscription = `
resource "google_pubsub" "foobar" {
	name = "foobar"
}
resource "google_subscription" "foobar" {
	name = "foobar"
	topic = "foobar"
}`
