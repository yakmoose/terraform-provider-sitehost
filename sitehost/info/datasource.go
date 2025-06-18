// Package info provides the function to call api/get_info resource via SiteHost API.
package info

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/info"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// DataSource returns a schema with the function to read api info resource.
func DataSource() *schema.Resource {

	return &schema.Resource{
		ReadContext: readDataSource,
		Schema:      apiInfoDataSourceSchema,
	}
}

// readDataSource is a function to read api info.
func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := info.New(conf.Client)

	resp, err := client.Get(ctx)
	if err != nil {
		return diag.Errorf("Error retrieving api info: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error retrieving api info: %s", resp.Msg)
	}

	d.SetId("client_id")

	if err := d.Set("client_id", resp.Return.ClientID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("contact_id", resp.Return.ContactID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("roles", resp.Return.Roles); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modules", resp.Return.Modules); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
