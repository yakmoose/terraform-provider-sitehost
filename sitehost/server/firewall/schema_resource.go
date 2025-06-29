package firewall

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceSchema is the schema with values for a Server Firewall resource.
var resourceSchema = map[string]*schema.Schema{
	"server": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the server to manage firewall rules for.",
	},
	"groups": {
		Type:        schema.TypeList,
		Required:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "List of security group names to apply to the server's firewall (in order).",
	},
}
