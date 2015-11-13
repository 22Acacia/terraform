package google

import (
	"os"
	"fmt"
	"testing"

	"github.com/22acacia/terraform-gcloud"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceControllerCreate(t *testing.T) {

	if os.Getenv("GOOGLE_GCLOUD_TESTS") == "TRUE" {
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckResourceControllerDestroy,
			Steps: []resource.TestStep{
				resource.TestStep{
					Config: testAccResourceController,
					Check: resource.ComposeTestCheckFunc(
						testAccResourceControllerExists(
							"google_container_replica_controller.foobar"),
					),
				},
			},
		})
	}
}

func testAccCheckResourceControllerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_container_replica_controller" {
			continue
		}

		_, err := terraformGcloud.ReadKubeRC(rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Failed to read resource controller list")
		}
		fmt.Println("tests, woo")
	}

	fmt.Println("party")

	return nil
}


func testAccResourceControllerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		count, err := terraformGcloud.ReadKubeRC(rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Command line read has errored: %q with rs.Primary hash: %q", err, rs.Primary)
		}

		if count == 0 {
			return fmt.Errorf("ResourceController pods never started")
		}

		return nil
	}
}

var testAccResourceController = `
resource "google_container_cluster" "foobar" {
    name = "foobar"
    zone = "europe-west1-b"
    initial_node_count = 3

    master_auth {
        username = "foo"
        password = "bar"
    }

    node_config {
        oauth_scopes = [
            "https://www.googleapis.com/auth/compute",
            "https://www.googleapis.com/auth/devstorage.read_only",
            "https://www.googleapis.com/auth/logging.write",
            "https://www.googleapis.com/auth/monitoring"
        ]
    }
}

resource "google_container_replica_controller" "foobar" {
	name = "foobar"
	docker_image = "nginx"
	container_name = "${google_container_cluster.foobar.name}"
	zone = "${google_container_cluster.foobar.zone}"
}`
