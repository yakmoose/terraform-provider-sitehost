package stack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// stackDataSourceSchema returns a schema with the function to read Server resource.
func stackDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The stack name",
		},

		"label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The stack label",
		},

		"server_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The server id",
		},

		"server_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The server name",
		},

		"server_label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The server label",
		},

		"docker_file": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The docker file used by the stack/container",
		},

		"ip_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The server ip address",
		},

		// containers...
		// is this even needed?
		// add the rest of the schema here... whheee
	}
}
