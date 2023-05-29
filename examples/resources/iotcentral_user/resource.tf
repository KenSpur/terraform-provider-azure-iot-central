# Manage example user.
resource "iotcentral_organization" "example" {
  id = "example"
  display_name = "Example"
}

resource "iotcentral_user" "example" {
  email = "example@example.net"
  roles = [ 
    {
      role = "c495eb57-eb18-489e-9802-62c474e5645c",
      organization = iotcentral_organization.example.id 
    }
  ]
}
