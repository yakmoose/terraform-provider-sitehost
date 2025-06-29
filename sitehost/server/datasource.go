package server

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/server"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// DataSource returns a schema with the function to read Server resource.
func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSource,
		Schema:      serverDataSourceSchema,
	}
}

// readDataSource is a function to read servers (not implemented).
func readDataSource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {

	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := server.New(conf.Client)

	resp, err := client.Get(ctx, server.GetRequest{
		ServerName: fmt.Sprintf("%v", d.Get("name")),
	})

	if err != nil {
		return diag.Errorf("Error retrieving server: %s", err)
	}
	d.SetId(resp.Server.Name)
	
	if !resp.Status {
		return diag.Errorf("Error retrieving server: %s", resp.Msg)
	}

	if err := setServerAttributes(d, resp.Server); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
