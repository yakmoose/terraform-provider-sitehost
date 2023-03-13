package api

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func apiInfoDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"client_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "The client id",
		},
		"contact_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,

			Description: "The contact id",
		},
		"roles": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"modules": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
