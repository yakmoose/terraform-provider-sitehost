package api_info

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/api_info"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// DataSource returns a schema with the function to read Server resource.
func DataSource() *schema.Resource {
	recordSchema := apiInfoDataSourceSchema()

	return &schema.Resource{
		ReadContext: readDataSource,
		Schema:      recordSchema,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := api_info.New(conf.Client)

	apiInfo, err := client.Get(ctx)
	if err != nil {
		return diag.Errorf("Error retrieving api info: %s", err)
	}

	d.SetId("client_id")
	d.Set("client_id", apiInfo.ClientID)
	d.Set("contact_id", apiInfo.ContactID)
	d.Set("roles", apiInfo.Roles)
	d.Set("modules", apiInfo.Modules)

	return nil
}
