// Package sshkey provides the functions to create a SSH Key resource via SiteHost API.
package sshkey

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sshkey "github.com/sitehostnz/gosh/pkg/api/ssh/key"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// Resource returns a schema with the operations for server.
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResource,
		ReadContext:   readResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceSchema,
	}
}

// createResource is a function to create a new SSH Key.
func createResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := sshkey.New(conf.Client)

	opts := sshkey.CreateRequest{
		Label:             fmt.Sprint(d.Get("label")),
		Content:           fmt.Sprint(d.Get("content")),
		CustomImageAccess: fmt.Sprint(d.Get("custom_image_access")),
	}

	res, err := client.Create(ctx, opts)
	if err != nil {
		return diag.Errorf("Error creating ssh key: %s", err)
	}

	if !res.Status {
		return diag.Errorf("Error creating ssh key: %s", res.Msg)
	}

	getOpts := sshkey.GetRequest{
		ID: res.Return.KeyID,
	}

	getRes, err := client.Get(ctx, getOpts)
	if err != nil {
		return diag.Errorf("Error getting ssh key: %s", err)
	}

	if !getRes.Status {
		return diag.Errorf("Error getting ssh key: %s", getRes.Msg)
	}

	if diagErr := setData(getRes, d); diagErr != nil {
		return diagErr
	}

	log.Printf("[INFO] SSH Key: %s", d.Id())

	return nil
}

// setData is a function to set the ResourceData values based on a getResponse.
func setData(res sshkey.GetResponse, d *schema.ResourceData) diag.Diagnostics {
	// Set data
	d.SetId(res.Return.ID)

	if err := d.Set("label", res.Return.Label); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("content", res.Return.Content); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("date_added", res.Return.DateAdded); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("date_updated", res.Return.DateUpdated); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// readResource is a function to read a new server.
func readResource(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

// updateResource is a function to update a server.
func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := sshkey.New(conf.Client)

	if d.HasChange("label") || d.HasChange("content") {
		updateKey(client, d)
	}

	return readResource(ctx, d, meta)
}

// updateKey is a function to update an SSH Key.
func updateKey(client *sshkey.Client, d *schema.ResourceData) diag.Diagnostics {
	res, err := client.Update(context.Background(), sshkey.UpdateRequest{
		ID:                d.Id(),
		Label:             fmt.Sprint(d.Get("label")),
		Content:           fmt.Sprint(d.Get("content")),
		CustomImageAccess: fmt.Sprint(d.Get("custom_image_access")),
	})
	if err != nil {
		return diag.Errorf("Error updating SSH Key: %s", err)
	}

	if !res.Status {
		return diag.Errorf("Error updating SSH Key: %s", res.Msg)
	}

	return nil
}

// deleteResource is a function to delete an SSH Key.
func deleteResource(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := sshkey.New(conf.Client)

	resp, err := client.Delete(context.Background(), sshkey.DeleteRequest{
		ID: d.Id(),
	})
	if err != nil {
		return diag.Errorf("Error deleting SSH Key: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error deleting SSH Key: %s", resp.Msg)
	}

	return nil
}
