package iotcentral

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCoffeesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
				data "iotcentral_role" "test" {
					display_name = "Administrator"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name of role returned
					resource.TestCheckResourceAttr("data.iotcentral_role.test", "display_name", "Administrator"),
				),
			},
			{
				Config: providerConfig + `
				data "iotcentral_role" "test" {
					display_name = "App Administrator"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name of role returned
					resource.TestCheckResourceAttr("data.iotcentral_role.test", "display_name", "Administrator"),
				),
			},
			{
				Config: providerConfig + `
				data "iotcentral_role" "test" {
					display_name = "Org Administrator"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name of role returned
					resource.TestCheckResourceAttr("data.iotcentral_role.test", "display_name", "Org Admin"),
				),
			},
			{
				Config: providerConfig + `
				data "iotcentral_role" "test" {
					display_name = "Org Admin"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name of role returned
					resource.TestCheckResourceAttr("data.iotcentral_role.test", "display_name", "Org Admin"),
				),
			},
		},
	})
}
