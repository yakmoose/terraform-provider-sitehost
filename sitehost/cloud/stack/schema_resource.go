package stack

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// resourceSchema returns a schema with the function to read Server resource.
var resourceSchema = map[string]*schema.Schema{
	"server": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Server name",
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Stack name",
	},
	"label": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Stack label",
	},
	"enable_ssl": {
		Type:        schema.TypeBool,
		Required:    true,
		Description: "enable or disable SSL",
	},
	"docker_compose": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The docker compose file",
	},
	"ip_address": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The server IP address",
	},
	//"containers": {
	//	Type:        schema.Array of some sort or another,
	//	Computed:    true,
	//	Description: "The container list",
	//},

	// not sure about this
	"settings": {
		Type:     schema.TypeMap,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
}
