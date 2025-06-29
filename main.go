// Package main is the starting point of the project to initialise the Terraform provider.
package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var version = "1.0.0"

// main is a function to initialise the Terraform provider.
func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        debugMode,
		ProviderFunc: sitehost.New(version),
		ProviderAddr: "sitehostnz/sitehost",
	}

	plugin.Serve(opts)
}
