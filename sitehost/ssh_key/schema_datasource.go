package sshkey

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// sshKeyDataSourceSchema is the schema with values for a SSH Key DataSource.
func sshKeyDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Required:    false,
			Optional:    false,
			Computed:    true,
			Description: "The `id` is the ID of the SSH Key within SiteHost's systems.",
		},
		"label": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The `label` is the name of the SSH Key, and is displayed in CP.",
		},
		"content": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The `content` is the contents of the public key.",
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
		"custom_image_access": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "`custom_image_access` determines whether the key can be used to access custom images.",
		},
	}
}
