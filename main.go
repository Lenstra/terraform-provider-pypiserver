package main

import (
	"gitlab.com/lenstra/pypi-server/terraform-provider-pypiserver/pypiserver"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return pypiserver.Provider()
		},
	})
}
