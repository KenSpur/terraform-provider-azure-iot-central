package iotcentral

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDevicesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "iotcentral_devices" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("data.iotcentral_devices.test", "devices.#", "1"),
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttr("data.iotcentral_devices.test", "devices.0.id", "example_2"),
					resource.TestCheckResourceAttr("data.iotcentral_devices.test", "devices.0.display_name", "example_2"),
				),
			},
		},
	})
}
