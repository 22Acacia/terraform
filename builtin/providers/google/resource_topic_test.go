package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTopicCreate(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTopicDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccTopic,
				Check: resource.ComposeTestCheckFunc(
					testAccTopicExists(
						"google_topic.foobar"),
				),
			},
		},
	})
}

func testAccCheckTopicDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_topic" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		_, err := config.clientPubsub.Projects.Topics.Get(rs.Primary.ID).Do()
		if err != nil {
			fmt.Errorf("Topic still present")
		}
	}

	return nil
}

func testAccTopicExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		_, err := config.clientPubsub.Projects.Topics.Get(rs.Primary.ID).Do()
		if err != nil {
			fmt.Errorf("Topic still present")
		}

		return nil
	}
}

const testAccTopic = `
resource "google_topic" "foobar" {
	name = "foobar"
}`
