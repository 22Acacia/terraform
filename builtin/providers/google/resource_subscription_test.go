package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSubscriptionCreate(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubscriptionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSubscription,
				Check: resource.ComposeTestCheckFunc(
					testAccSubscriptionExists(
						"google_subscription.foobar_sub"),
				),
			},
		},
	})
}

func testAccCheckSubscriptionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_subscription" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		_, err := config.clientPubsub.Projects.Subscriptions.Get(rs.Primary.ID).Do()
		if err != nil {
			fmt.Errorf("Subscription still present")
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
		config := testAccProvider.Meta().(*Config)
		_, err := config.clientPubsub.Projects.Subscriptions.Get(rs.Primary.ID).Do()
		if err != nil {
			fmt.Errorf("Subscription still present")
		}

		return nil
	}
}

const testAccSubscription = `
resource "google_topic" "foobar_sub" {
	name = "foobar_sub"
}

resource "google_subscription" "foobar_sub" {
	name = "foobar_sub"
	topic = "${google_topic.foobar_sub.name}"
}`

