package user

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceSchema returns a schema with the function to read an ssh user.
var resourceSchema = map[string]*schema.Schema{
	"server_name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The server where the user is configured",
	},

	"username": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The user name",
		ValidateFunc: validation.StringMatch(
			regexp.MustCompile("^[a-z0-9]+$"),
			"Usernames can only be lowercase alphanumeric characters",
		),
	},

	// this is a one way trip.....
	// can't read it back...
	"password": {
		Type:        schema.TypeString,
		Sensitive:   true,
		Optional:    true,
		Description: "The password for the user",

		// DiffSuppressOnRefresh: true,
		// DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
		//	// how to check/handle this...
		//	return true
		// },
	},

	"read_only_config": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},

	"ssh_key": {
		Type:     schema.TypeList,
		Optional: true,
		Required: false,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The `id` is the ID of the SSH Key within SiteHost's systems.",
				},

				// read these back on create??? / known after apply
				"label": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The `label` is the name of the SSH Key, and is displayed in CP.",
				},
				"content": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The `content` is the contents of the public key.",
				},
			},
		},
	},

	// not sure how to validate this correctly
	// since we can only have one or the other...

	"container": {
		Type:     schema.TypeList,
		Optional: true,
		Required: false,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	},

	"volume": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	},
}
