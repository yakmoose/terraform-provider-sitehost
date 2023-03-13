package db

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// databaseDataSourceSchema returns a schema with the function to read a stack resource.
func databaseDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"server_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The server id/name",
		},
		"mysql_host": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The mysqlhost",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The database nume",
		},
		"backup_container": {
			Type:        schema.TypeString,
			Description: "The container where backups are stored",
			Computed:    true,
		},
		//"grants": {
		//
		//},
	}
}
