// Package stack provides the functions to create/get cloud stacks resource via SiteHost API.
package stack

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// DataSource returns a schema with the function to read cloud stack resource.
func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readResource,
		Schema:      stackDataSourceSchema,
	}
}

// ListDataSource is the datasource for listing stacks, with/without a filter.
func ListDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: listResource,
		Schema:      listStackDataSourceSchema,
	}
}

func listResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}
	client := stack.New(conf.Client)

	listResponse, err := client.List(ctx, stack.ListRequest{})
	if err != nil {
		return diag.Errorf("Failed to fetch database list %s", err)
	}

	stacks := []map[string]string{}
	for _, v := range listResponse.Return.Stacks {
		s := map[string]string{
			"name":         v.Name,
			"label":        v.Label,
			"server_id":    v.ServerID,
			"server_name":  v.Server,
			"server_label": v.ServerLabel,
		}
		//
		stacks = append(stacks, s)
	}

	d.SetId("stacks")

	if err := d.Set("stacks", stacks); err != nil {
		return diag.Errorf("Error retrieving stacks: %s", err)
	}

	return nil
}
