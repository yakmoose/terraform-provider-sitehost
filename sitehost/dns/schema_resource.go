package dns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		ValidateFunc: validation.NoZeroValues,
		Description:  "The subdomain",
	},

	"type": {
		Type:     schema.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice([]string{
			"A",
			"AAAA",
			"CAA",
			"CNAME",
			"MX",
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

	"record": {
		Type:     schema.TypeString,
		Optional: true,
	},

	"change_date": {
		Type:     schema.TypeString,
		Computed: true,
	},
}
