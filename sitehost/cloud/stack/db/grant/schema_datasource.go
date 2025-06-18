package grant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// databaseGrantDataSourceSchema returns a schema for the database user datasource.
var databaseGrantDataSourceSchema = map[string]*schema.Schema{
	"server_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The server id/name",
	},
	"mysql_host": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The mysql host",
	},
	"username": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The username",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database name",
	},
	"grants": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The the assigned grants",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
}
