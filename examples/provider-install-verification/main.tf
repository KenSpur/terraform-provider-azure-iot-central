terraform {
  required_providers {
    iotcentral = {
      source = "local/develop/azure-iot-central"
    }
  }
}

provider "iotcentral" {}

data "iotcentral_devices" "example" {}