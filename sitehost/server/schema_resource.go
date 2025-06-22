package server

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// resourceSchema is the schema with values for a Server resource.
var resourceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: "The `name` is the ID and is provided for a Server.",
	},
	"password": {
		Type:        schema.TypeString,
		Sensitive:   true,
		Computed:    true,
		Description: "The password that will be assigned to the 'root' user account.",
	},
	"ips": {
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Each Server is assigned a single public IPv4 address upon creation.",
	},
	"label": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The SiteHost's label is for display purposes only.",
	},
	"location": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		Description: "This is the location where the Server was deployed. This cannot be changed without " +
			"opening a support ticket.",
	},
	"product_code": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The product code of the server to be deployed, determining the price and size.",
	},
	"securitygroups": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "The security groups which this server uses.",
	},
	"image": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		Description: "An Image ID to deploy the Disk from. The complete list of images ID you can see " +
			"in our official documentation.",
	},
	"ssh_keys": {
		Type:        schema.TypeList,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "A list of SSH public keys to deploy for the root user on the newly created Server.",
	},
}
