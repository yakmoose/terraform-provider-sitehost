package db

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// databaseResourceSchema returns a schema with the function to read stack resource.
var databaseResourceSchema = map[string]*schema.Schema{
	"server_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The server id/name",
		ForceNew:    true,
	},
	"mysql_host": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The mysqlhost",
		ForceNew:    true,
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database name",
		ForceNew:    true,
	},
	"backup_container": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The container where backups are stored",
	},

	// grants here is problematic, since we can't create grants until after we have database
}
