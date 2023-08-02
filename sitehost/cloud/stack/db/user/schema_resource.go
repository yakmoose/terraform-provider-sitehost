package user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// databaseUserResourceSchema returns a schema with the function to create/manipulate a stack database user.
func databaseUserResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

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
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The username",
			ForceNew:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The users password",
		},
		"database": {
			Type:         schema.TypeString,
			Required:     false,
			Optional:     true,
			Description:  "The database name",
			RequiredWith: []string{"grants"},
		},
		"grants": {
			Type:         schema.TypeList,
			Optional:     true,
			Required:     false,
			RequiredWith: []string{"database"},
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
