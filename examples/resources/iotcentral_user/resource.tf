data "iotcentral_role" "example" {
  display_name = "Org Administrator"
}

resource "iotcentral_organization" "example" {
  id = "example"
  display_name = "Example"
}

resource "iotcentral_user" "example" {
  email = "example@example.net"
  roles = [ 
    {
      role = data.iotcentral_role.example.id
      organization = iotcentral_organization.example.id 
    }
  ]
}
