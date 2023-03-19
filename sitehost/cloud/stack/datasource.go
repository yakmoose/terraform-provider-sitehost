// Package stack provides the functions to create/get cloud stacks resource via SiteHost API.
package stack

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// DataSource returns a schema with the function to read cloud stack resource.
func DataSource() *schema.Resource {
	recordSchema := stackDataSourceSchema()

	return &schema.Resource{
		ReadContext: readResource,
		Schema:      recordSchema,
	}
}

// readDataSource calls the GoSH client to set the cloud stack schema.
func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := stack.New(conf.Client)

	resp, err := client.Get(ctx, stack.GetRequest{
		ServerName: fmt.Sprintf("%v", d.Get("server_name")),
		Name:       fmt.Sprintf("%v", d.Get("name")),
	})
	if err != nil {
		return diag.Errorf("Error retrieving api info: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error retrieving api info: %s", resp.Msg)
	}

	d.SetId(resp.Stack.Name)

	if err := d.Set("label", resp.Stack.Label); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("server_name", resp.Stack.Server); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("server_label", resp.Stack.ServerLabel); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
