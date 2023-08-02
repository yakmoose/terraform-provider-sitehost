package grant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// databaseGrantDataSourceSchema returns a schema for the database user datasource.
func databaseGrantDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		//"server_name": {
		//	Type:        schema.TypeString,
		//	Required:    true,
		//	Description: "The server id/name",
		//	ForceNew:    true,
		//},
		//"mysql_host": {
		//	Type:        schema.TypeString,
		//	Required:    true,
		//	Description: "The mysqlhost",
		//	ForceNew:    true,
		//},
		//"username": {
		//	Type:        schema.TypeString,
		//	Required:    true,
		//	Description: "The username",
		//	ForceNew:    true,
		//},
		//"password": {
		//	Type:        schema.TypeString,
		//	Required:    true,
		//	Description: "The users password",
		//	ForceNew:    true,
		//},
		//"grants": {
		//	Type:     schema.TypeList,
		//	Optional: true,
		//	Required: false,
		//	Elem: &schema.Schema{
		//		Type: schema.TypeString,
		//	},
		//},
	}
}
