package dns

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sitehostnz/gosh/pkg/utils"
)

// dnsRecordResourceSchema is the schema with values for a Server resource.
var dnsRecordResourceSchema = map[string]*schema.Schema{
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
		ForceNew:     false,
		ValidateFunc: validation.NoZeroValues,
		Description:  "The subdomain",
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			domain := d.Get("domain").(string)

			oldValue = utils.ConstructFqdn(oldValue, domain)
			newValue = utils.ConstructFqdn(newValue, domain)

			return newValue == oldValue
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

	"content": {
		Type:     schema.TypeString,
		Optional: true,
		DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
			// bloody dots at the end of records...
			// we have to do this, mainly for NS and CNAME records
			// Possibly MX records too... hell, let's just do them all
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
