package iotcentral

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIotCentralOrganizationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "iotcentral_organization" "test_1" {
					id = "testid1"
					display_name = "Test 1"
				}

				resource "iotcentral_organization" "test_1_child" {
					id = "testid1child"
					display_name = "Test 1 child"
					parent = iotcentral_organization.test_1.id
				}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify id is set
					resource.TestCheckResourceAttr("iotcentral_organization.test_1", "id", "testid1"),
					// Verify display_name is set
					resource.TestCheckResourceAttr("iotcentral_organization.test_1", "display_name", "Test 1"),

					// Verify id is set
					resource.TestCheckResourceAttr("iotcentral_organization.test_1_child", "id", "testid1child"),
					// Verify display_name is set
					resource.TestCheckResourceAttr("iotcentral_organization.test_1_child", "display_name", "Test 1 child"),
					// Verify parent is set
					resource.TestCheckResourceAttr("iotcentral_organization.test_1_child", "parent", "testid1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "iotcentral_organization.test_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "iotcentral_organization" "test_1" {
					id = "testid1"
					display_name = "test_1_name_updated"
				}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify id is same
					resource.TestCheckResourceAttr("iotcentral_organization.test_1", "id", "testid1"),
					// Verify display_name is updated
					resource.TestCheckResourceAttr("iotcentral_organization.test_1", "display_name", "test_1_name_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
