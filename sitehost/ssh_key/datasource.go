package sshkey

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sshkey "github.com/sitehostnz/gosh/pkg/api/ssh/key"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// DataSource returns a schema with the function to read Server resource.
func DataSource() *schema.Resource {
	recordSchema := sshKeyDataSourceSchema

	return &schema.Resource{
		ReadContext: readDataSource,
		Schema:      recordSchema,
	}
}

// ListDataSource is the datasource for listing stacks, with a filter.
func ListDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: listDataSource,
		Schema:      listSSHKeysDataSourceSchema,
	}
}

// readDataSource is a function to read an SSH Key.
func readDataSource(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := sshkey.New(conf.Client)

	resp, err := client.Get(context.Background(), sshkey.GetRequest{
		ID: d.Id(),
	})
	if err != nil {
		return diag.Errorf("Error retrieving SSH Key: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error retrieving SSH Key: %s", resp.Msg)
	}

	if diagErr := setData(resp, d); diagErr != nil {
		return diagErr
	}

	return nil
}

// listDataSource is a function to read an SSH Key.
func listDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}
	client := sshkey.New(conf.Client)

	listResponse, err := client.List(ctx)
	if err != nil {
		return diag.Errorf("Failed to fetch ssh users list %s", err)
	}

	sshKeys := []map[string]string{}
	for _, v := range listResponse.Return.SSHKeys {
		k := map[string]string{
			"id":           v.ID,
			"label":        v.Label,
			"content":      v.Content,
			"date_added":   v.DateAdded,
			"date_updated": v.DateUpdated,
			// "custom_image_access": v.CustomImageAccess,

			// I've intentionally left out the grants here
		}

		sshKeys = append(sshKeys, k)
	}

	d.SetId("ssh_keys")

	if err := d.Set("ssh_keys", sshKeys); err != nil {
		return diag.Errorf("Error retrieving ssh keys info: %s", err)
	}

	return nil
}
