package main

import (
	"context"

	"terraform-provider-azure-iot-central/iotcentral"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name iotcentral

func main() {
	providerserver.Serve(context.Background(), iotcentral.New, providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as registry.terraform.io/hashicorp/azure-iot-central. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address: "local/develop/azure-iot-central",
	})
}
