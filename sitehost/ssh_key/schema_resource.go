package sshkey

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

// resourceSchema is the schema with values for a SSH Key resource.
var resourceSchema = map[string]*schema.Schema{
	"label": {
		Type:     schema.TypeString,
		Required: true,
		// don't want things re-keying if they are updated
		ForceNew:    false,
		Description: "The `label` is the name of the SSH Key, and is displayed in CP.",
	},
	"content": {
		Type: schema.TypeString,
		// It's a public key, it's right there in the name, public keys are not sensitive.
		Sensitive: false,
		Required:  true,
		// don't want things re-keying if they are updated
		ForceNew:    false,
		Description: "The `content` is the contents of the public key.",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// we need this in place, since the keys may end up with whitespace and crap, so make sure we're
			// just looking at the key
			return strings.TrimSpace(old) == strings.TrimSpace(new)
		},
	},
	"custom_image_access": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
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
