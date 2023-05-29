package iotcentral

// import (
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
// )

// func TestAccIotCentralADGroupUserResource(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Create and Read testing
// 			{
// 				Config: providerConfig + `
// 				resource "iotcentral_ad_group_user" "test" {
// 					object_id = "<object_id>"
//   					tenant_id = "<tenant_id>"
// 					roles = [
// 					  {
// 						role = "ca310b8d-2f4a-44e0-a36e-957c202cd8d4"
// 					  }
// 					]
// 				  }
// `,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					// Verify object id is set
// 					resource.TestCheckResourceAttr("iotcentral_ad_group_user.test", "object_id", "<object_id>"),
// 					// Verify tenant id is set
// 					resource.TestCheckResourceAttr("iotcentral_ad_group_user.test", "tenant_id", "<tenant_id>"),
// 					// Verify roles are set
// 					resource.TestCheckResourceAttr("iotcentral_ad_group_user.test", "roles.#", "1"),
// 					// Verify role is set
// 					resource.TestCheckResourceAttr("iotcentral_ad_group_user.test", "roles.0.role", "ca310b8d-2f4a-44e0-a36e-957c202cd8d4"),
// 				),
// 			},
// 			// ImportState testing
// 			{
// 				ResourceName:      "iotcentral_ad_group_user.test",
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			// Update and Read testing
// 			{
// 				Config: providerConfig + `
// 				resource "iotcentral_organization" "user_test_org" {
// 					id = "usertestorg"
// 					display_name = "User Test Org"
// 				}

// 				resource "iotcentral_ad_group_user" "test" {
// 					object_id = "<object_id>"
//   					tenant_id = "<tenant_id>"
// 					roles = [
// 					  {
// 						role = "ca310b8d-2f4a-44e0-a36e-957c202cd8d4"
// 					  },
// 					  {
// 						role = "c495eb57-eb18-489e-9802-62c474e5645c",
// 						organization = iotcentral_organization.user_test_org.id
// 					  }
// 					]
// 				  }
// `,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					// Verify object id is set
// 					resource.TestCheckResourceAttr("iotcentral_ad_group_user.test", "object_id", "<object_id>"),
// 					// Verify tenant id is set
// 					resource.TestCheckResourceAttr("iotcentral_ad_group_user.test", "tenant_id", "<tenant_id>"),
// 					// Verify roles are set
// 					resource.TestCheckResourceAttr("iotcentral_ad_group_user.test", "roles.#", "2"),
// 				),
// 			},
// 			// Update and Read testing
// 			{
// 				Config: providerConfig + `
// 				resource "iotcentral_organization" "user_test_org" {
// 					id = "usertestorg"
// 					display_name = "User Test Org"
// 				}

// 				resource "iotcentral_ad_group_user" "test" {
// 					object_id = "<object_id>"
//   					tenant_id = "<tenant_id>"
// 					roles = [
// 					  {
// 						role = "c495eb57-eb18-489e-9802-62c474e5645c",
// 						organization = iotcentral_organization.user_test_org.id
// 					  }
// 					]
// 				  }
// `,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					// Verify object id is set
// 					resource.TestCheckResourceAttr("iotcentral_ad_group_user.test", "object_id", "<object_id>"),
// 					// Verify tenant id is set
// 					resource.TestCheckResourceAttr("iotcentral_ad_group_user.test", "tenant_id", "<tenant_id>"),
// 					// Verify roles are set
// 					resource.TestCheckResourceAttr("iotcentral_ad_group_user.test", "roles.#", "1"),
// 					// Verify role is set
// 					resource.TestCheckResourceAttr("iotcentral_ad_group_user.test", "roles.0.role", "c495eb57-eb18-489e-9802-62c474e5645c"),
// 					resource.TestCheckResourceAttr("iotcentral_ad_group_user.test", "roles.0.organization", "usertestorg"),
// 				),
// 			},
// 			// Delete testing automatically occurs in TestCase
// 		},
// 	})
// }
