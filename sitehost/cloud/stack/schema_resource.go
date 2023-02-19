package stack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceSchema returns a schema with the function to read Server resource.
var resourceSchema = map[string]*schema.Schema{

	"name": {
		Type:        schema.TypeString,
		Computed:    true,
		Optional:    false,
		Description: "The Stack name",
		// name is computed... if we are creating
		// not if we are importing
	},

	"label": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Stack label",
	},

	// this is a create option, and exists against an individual container
	// when getting, not sure how we update
	"enable_ssl": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Enable or disable SSL, it will default to false, as domain names must be pointing at the server in order to issue",
		Default:     false,
	},

	// These options are used to create the docker_file
	"monitored": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Enable or disable SSL",
		Default:     true,
	},
	"type": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "www",
		ValidateFunc: validation.StringInSlice([]string{
			"www",
			"service",
			"integrated",     // internal sh?
			"infrastructure", // internal sh?
		}, false),
	},
	"backup_disable": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},

	"restart": {
		Type: schema.TypeString,
		ValidateFunc: validation.StringInSlice([]string{
			"always",
			"unless-stopped",
			"on-failure", // this one is tricksy as it can also take a count... let's not for now..
			"no",
		}, false),
	},

	// this likely has rules in the main sh api around custom vs sh containers
	"image_update": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},

	//virtual hosts are called aliases in the sh UI
	"aliases": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"image": {
		Type:     schema.TypeString,
		Required: true,
	},
	"docker_file": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The docker compose file as returned from the server, that we have generated on create and bundles thigns together",
	},

	"expose": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},

	// questions, is the volumes standard? or comes from elsewhere? ports etc...
	// this will be an abstraction on top of the other schema stuff
	// since we can add them when creating, but updating does not change these
	"settings": {
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},

	// server properties, can't change these, informational only
	"server_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Server name where the stack lives",
		ForceNew:    true,
	},

	"server_label": {
		Type:        schema.TypeString,
		Computed:    true,
		Optional:    false,
		Description: "The Server label",
	},
	"server_ip_address": {
		Type:        schema.TypeString,
		Computed:    true,
		Optional:    false,
		Description: "The server IP address",
	},

	"containers": {
		Type:     schema.TypeList,
		Computed: true,
		Optional: false,
	},
}
