package iotcentral

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "iotcentral_device" "test_1" {
					id = "test_id_1"
					display_name = "test_1"
				}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify id is set
					resource.TestCheckResourceAttr("iotcentral_device.test_1", "id", "test_id_1"),
					// Verify display_name is set
					resource.TestCheckResourceAttr("iotcentral_device.test_1", "display_name", "test_1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "iotcentral_device.test_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "iotcentral_device" "test_1" {
					id = "test_id_1"
					display_name = "test_1_name_updated"
				}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify id is same
					resource.TestCheckResourceAttr("iotcentral_device.test_1", "id", "test_id_1"),
					// Verify display_name is updated
					resource.TestCheckResourceAttr("iotcentral_device.test_1", "display_name", "test_1_name_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
