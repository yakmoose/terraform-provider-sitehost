package user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// userDataSourceSchema returns a schema with the function to read a stack resource.
func userDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"server_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Server where the user exists",
		},

		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The username",
		},

		"ssh_key": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The `id` is the ID of the SSH Key within SiteHost's systems.",
					},
					"label": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The `label` is the name of the SSH Key, and is displayed in CP.",
					},
					"content": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The `content` is the contents of the public key.",
					},
				},
			},
		},

		"containers": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"volumes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"read_only_config": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}
