package environment

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// environmentDatasourceSchema returns a schema with the function to read Server resource.
var environmentDatasourceSchema = map[string]*schema.Schema{
	"server_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The server id/name",
	},
	"project": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The the project id/name",
	},
	"service": {
		Type:        schema.TypeString,
		Required:    false,
		Optional:    true,
		Computed:    true,
		Description: "The service id, this is optional and defaults to the project id/name",
	},

	// key pairs here...
	"settings": {
		Type:     schema.TypeMap,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		//		ValidateFunc: validation.StringIsNotWhiteSpace,
	},
}
