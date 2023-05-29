resource "iotcentral_organization" "example" {
  id = "example"
  display_name = "Example"
}

resource "iotcentral_organization" "example_child" {
  id = "examplechild"
  display_name = "Example Child"
  parent = iotcentral_organization.example.id
}