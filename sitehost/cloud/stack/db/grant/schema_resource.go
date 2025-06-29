package grant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// databaseGrantResourceSchema returns a schema with the function to create/manipulate a stack database user.
var databaseGrantResourceSchema = map[string]*schema.Schema{
	"server_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The server id/name",
		ForceNew:    true,
	},
	"mysql_host": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The mysql host",
		ForceNew:    true,
	},
	"username": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The username",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database name",
		ForceNew:    true,
	},
	"grants": {
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
}
