package user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// databaseUserResourceSchema returns a schema with the function to create/manipulate a stack database user.
var databaseUserResourceSchema = map[string]*schema.Schema{
	//	return map[string]*schema.Schema{
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
	//	}
}
