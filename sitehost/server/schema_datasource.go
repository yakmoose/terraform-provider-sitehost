package server

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// serverDataSourceSchema is the schema with values for a Server DataSource.
var serverDataSourceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The `name` is the ID and is provided for a Server.",
	},
	"label": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Each Server is assigned a single public IPv4 address upon creation.",
	},
	"location": {
		Type:     schema.TypeString,
		Required: true,
		Description: "This is the location where the Server was deployed. This cannot be changed without " +
			"opening a support ticket.",
	},
	"product_code": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The product code of the server to be deployed, determining the price and size.",
	},
	"image": {
		Type:     schema.TypeString,
		Required: true,
		Description: "An Image ID to deploy the Disk from. The complete list of images ID you can see " +
			"in our official documentation.",
	},
}
