data "iotcentral_role" "example" {
  display_name = "App Administrator"
}

resource "iotcentral_user" "example" {
  email = "example@example.net"
  roles = [ 
    {
      role = data.iotcentral_role.example.id
    }
  ]
}