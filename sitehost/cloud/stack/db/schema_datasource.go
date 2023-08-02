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
	}
}

// listDatabaseDataSourceSchema is the datasource for a listing of databases.
func listDatabaseDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"databases": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "The list of databases",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"server_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The server name",
					},
					"server_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The server id",
					},
					"server_label": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The server label",
					},
					"mysql_host": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The mysqlhost",
					},
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The database name",
					},
					"backup_container": {
						Type:        schema.TypeString,
						Description: "The container where backups are stored",
						Computed:    true,
					},
				},
			},
		},
	}
}
