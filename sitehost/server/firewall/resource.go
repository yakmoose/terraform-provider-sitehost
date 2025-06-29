// Package firewall provides the functions to create a Firewall resource via SiteHost API.
package firewall

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/server/firewall"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// Resource returns a schema with the operations for server.
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: updateResource,
		ReadContext:   readResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceSchema,
	}
}

// readResource is a function to read the firewall of a server.
func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := firewall.New(conf.Client)

	resp, err := client.Get(ctx, firewall.GetRequest{
		ServerName: d.Get("server").(string),
	})
	if err != nil {
		return diag.Errorf("Error reading server: %s", err)
	}

	groups := make([]string, len(resp.Return))
	for i, group := range resp.Return {
		groups[i] = group.Group
	}
	if err := d.Set("groups", groups); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// updateResource is a function to update the firewall of a server.
func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := firewall.New(conf.Client)

	serverName, ok := d.Get("server").(string)
	if !ok {
		return diag.Errorf("failed to convert server name to string")
	}
	list, ok := d.Get("groups").([]interface{})
	if !ok {
		return diag.Errorf("failed to convert groups to []interface{}")
	}
	groups := make([]string, 0)
	for _, v := range list {
		groups = append(groups, fmt.Sprint(v))
	}

	response, diags := updateFirewallGroups(ctx, client, serverName, groups)
	if diags != nil {
		return diags
	}

	if err := helper.WaitForJob(conf.Client, response.Return.Job); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serverName)
	return nil
}

// deleteResource is a function to clear the firewall of a server.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := firewall.New(conf.Client)

	serverName, ok := d.Get("server").(string)
	if !ok {
		return diag.Errorf("failed to convert server name to string")
	}

	response, diags := updateFirewallGroups(ctx, client, serverName, []string{})
	if diags != nil {
		return diags
	}

	if err := helper.WaitForJob(conf.Client, response.Return.Job); err != nil {
		return diag.FromErr(err)
	}

	// Clear the ID when the resource is deleted
	d.SetId("")
	return nil
}

// updateFirewallGroups is a function to update the firewall security groups of a server.
func updateFirewallGroups(ctx context.Context, client *firewall.Client, serverName string, groups []string) (firewall.UpdateResponse, diag.Diagnostics) {
	res, err := client.Update(ctx, firewall.UpdateRequest{
		ServerName:     serverName,
		SecurityGroups: groups,
	})
	if err != nil {
		return res, diag.Errorf("Error updating server: %s", err)
	}

	if !res.Status {
		return res, diag.Errorf("Error updating server: %s", res.Msg)
	}

	return res, nil
}
