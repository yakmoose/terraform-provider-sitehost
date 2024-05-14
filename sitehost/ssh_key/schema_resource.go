package sshkey

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// resourceSchema is the schema with values for a SSH Key resource.
var resourceSchema = map[string]*schema.Schema{
	"label": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The `label` is the name of the SSH Key, and is displayed in CP.",
	},
	"content": {
		Type:        schema.TypeString,
		Sensitive:   true,
		Required:    true,
		ForceNew:    true,
		Description: "The `content` is the contents of the public key.",
	},
	"custom_image_access": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "0",
	},
	"id": {
		Type:        schema.TypeString,
		Required:    false,
		Optional:    false,
		Computed:    true,
		Description: "The `id` is the ID of the SSH Key within SiteHost's systems.",
	},
	"date_added": {
		Type:        schema.TypeString,
		Required:    false,
		Optional:    false,
		Computed:    true,
		Description: "The `date_added` is the date/time when the SSH Key was added.",
	},
	"date_updated": {
		Type:        schema.TypeString,
		Required:    false,
		Optional:    false,
		Computed:    true,
		Description: "The `date_updated` is the date/time when the SSH Key was updated.",
	},
}
