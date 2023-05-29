resource "iotcentral_organization" "example" {
  id = "example"
  display_name = "Example"
}

resource "iotcentral_service_principal_user" "example" {
  object_id = "<object_id>"
  tenant_id = "<tenant_id>"
  roles = [ 
    {
      role = "c495eb57-eb18-489e-9802-62c474e5645c",
      organization = iotcentral_organization.example.id 
    }
  ]
}
