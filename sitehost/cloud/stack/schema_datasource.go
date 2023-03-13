package stack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// stackDataSourceSchema returns a schema with the function to read a stack resource.
func stackDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Stack name",
		},

		"label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The Stack label",
		},

		// this is a create option, and is reflected in the docker file
		"enable_ssl": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Enable or disable SSL",
		},

		// These options are used to create the docker_file
		"monitored": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Enable or disable SSL",
		},
		"type": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"backup_disable": {
			Computed: true,
			Type:     schema.TypeBool,
		},

		// this likely has rules in the main sh api around custom vs sh containers
		"image_update": {
			Computed: true,
			Type:     schema.TypeBool,
		},

		// virtual hosts are called aliases in the sh UI
		"aliases": {
			Computed: true,
			Type:     schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"image": {
			Computed: true,
			Type:     schema.TypeString,
		},

		"docker_file": {
			Computed:    true,
			Type:        schema.TypeString,
			Description: "The docker compose file as returned from the server, that we have generated on create and bundles thigns together",
		},

		"expose": {
			Computed: true,
			Type:     schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		// questions, is the volumes standard? or comes from elsewhere? ports etc...
		// this will be an abstraction on top of the other schema stuff
		// since we can add them when creating, but updating does not change these
		"settings": {
			Computed: true,
			Type:     schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		// server properties, can't change these, informational only
		"server_name": {
			Required:    true,
			Type:        schema.TypeString,
			Description: "The Server name where the stack lives",
		},
		"server_label": {
			Computed:    true,
			Type:        schema.TypeString,
			Description: "The Server label",
		},
		"server_ip_address": {
			Computed:    true,
			Type:        schema.TypeString,
			Description: "The server IP address",
		},
	}
}
