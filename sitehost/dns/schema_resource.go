package dns

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sitehostnz/gosh/pkg/utils"
)

// resourceZoneSchema is the schema with values for a DNS zone resource.
var resourceZoneSchema = map[string]*schema.Schema{
	"name": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.NoZeroValues,
		Description:  "The domain name",
	},
}

// resourceRecordSchema is the schema with values for a DNS record resource.
var resourceRecordSchema = map[string]*schema.Schema{
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
			domain := fmt.Sprintf("%v", d.Get("domain"))

			oldValue = utils.ConstructFqdn(fmt.Sprintf("%v.", oldValue), domain)
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
			"TXT",
			"SRV",
			"NS", // added this back, as creating a zone does not appear to set the DNS records
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
		DiffSuppressFunc: func(k, oldRecord, newRecord string, d *schema.ResourceData) bool {
			// bloody dots at the end of records...
			// we have to do this, mainly for NS and CNAME records
			// Possibly MX records too... hell, let's just do them all
			return strings.TrimSuffix(oldRecord, ".") == strings.TrimSuffix(newRecord, ".")
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
