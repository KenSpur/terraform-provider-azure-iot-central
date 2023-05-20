terraform {
  required_providers {
    iotcentral = {
      source = "local/develop/azure-iot-central"
    }
  }
}

provider "iotcentral" {}

resource "iotcentral_device" "example" {
  id = "example"
  display_name = "example"
}

data "iotcentral_devices" "example" {}

output "example_devices" {
  value = data.iotcentral_devices.example
}