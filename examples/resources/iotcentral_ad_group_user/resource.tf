data "iotcentral_role" "example" {
  display_name = "Org Administrator"
}

resource "iotcentral_organization" "example" {
  id = "example"
  display_name = "Example"
}

resource "iotcentral_ad_group_user" "example" {
  object_id = "<object_id>"
  tenant_id = "<tenant_id>"
  roles = [ 
    {
      role = data.iotcentral_role.example.id
      organization = iotcentral_organization.example.id 
    }
  ]
}
