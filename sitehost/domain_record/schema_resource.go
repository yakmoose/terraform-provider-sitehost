package domain_record

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
)

// resourceSchema is the schema with values for a Server resource.
var resourceSchema = map[string]*schema.Schema{
	"domain": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.NoZeroValues,
		Description:  "The base domain",
	},

	"name": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.NoZeroValues,
		Description:  "The subdomain",
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {

			domain := d.Get("domain").(string)

			return (oldValue == "@" && newValue == domain) || (oldValue+"."+domain == newValue)
		},
	},

	"type": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice([]string{
			"A",
			"AAAA",
			"CAA",
			"CNAME",
			"MX",
			"NS",
			"TXT",
			"SRV",
		}, false),
		Description: "The record type",
	},

	"priority": {
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 65535),
		Description:  "The priority type",
		Default:      0,
	},

	"ttl": {
		Type:         schema.TypeInt,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntAtLeast(1),
	},

	"record": {
		Type:     schema.TypeString,
		Optional: true,
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			// bloody dots at the end of records...
			return strings.TrimSuffix(oldValue, ".") == strings.TrimSuffix(newValue, ".")
		},
	},

	"fqdn": {
		Type:     schema.TypeString,
		Computed: true,
	},

	"change_date": {
		Type:     schema.TypeString,
		Computed: true,
	},
}
