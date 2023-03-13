// Package server provides the functions to create a Server resource via SiteHost API.
package server

import (
	"context"
	"fmt"
	"github.com/sitehostnz/gosh/pkg/api/job"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/server"
	"github.com/sitehostnz/gosh/pkg/models"
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

// createResource is a function to create a new server.
func createResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := server.New(conf.Client)

	keys, ok := d.Get("ssh_keys").([]any)
	if !ok {
		return diag.Errorf("failed to convert ssh keys object")
	}

	sshKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		sshKeys = append(sshKeys, fmt.Sprint(key))
	}

	opts := server.CreateRequest{
		Label:       fmt.Sprint(d.Get("label")),
		Location:    fmt.Sprint(d.Get("location")),
		ProductCode: fmt.Sprint(d.Get("product_code")),
		Image:       fmt.Sprint(d.Get("image")),
		Params: server.ParamsOptions{
			SSHKeys: sshKeys,
		},
	}

	res, err := client.Create(ctx, opts)
	if err != nil {
		return diag.Errorf("Error creating server: %s", err)
	}

	if !res.Status {
		return diag.Errorf("Error creating server: %s", res.Msg)
	}

	// Set data
	d.SetId(res.Return.Name)
	if err := d.Set("name", res.Return.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("password", res.Return.Password); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ips", res.Return.Ips); err != nil {
		return diag.FromErr(err)
	}

	// wait for "Completed" status
	if err := helper.WaitForAction(conf.Client, job.GetRequest{JobID: res.Return.JobID, Type: job.DaemonType}); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Server Name: %s", d.Id())

	return nil
}

// readResource is a function to read a new server.
func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := server.New(conf.Client)

	resp, err := client.Get(context.Background(), server.GetRequest{
		ServerName: d.Id(),
	})
	if err != nil {
		return diag.Errorf("Error retrieving server: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error retrieving server: %s", resp.Msg)
	}

	if err := setServerAttributes(d, resp.Server); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// updateResource is a function to update a server.
func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := server.New(conf.Client)

	if d.HasChange("product_code") {
		return upgradePlan(conf, client, d)
	}

	if d.HasChange("label") {
		return updateLabel(client, d)
	}

	return readResource(ctx, d, meta)
}

// upgradePlan is a function to upgrade and commit a server to the next plan.
func upgradePlan(conf *helper.CombinedConfig, client *server.Client, d *schema.ResourceData) diag.Diagnostics {
	res, err := client.Upgrade(context.Background(), server.UpgradeRequest{
		Name: d.Id(),
		Plan: fmt.Sprint(d.Get("product_code")),
	})
	if err != nil {
		return diag.Errorf("Error upgrading server: %s", err)
	}

	if !res.Status {
		return diag.Errorf("Error upgrading server: %s", res.Msg)
	}

	resp, err := client.CommitDiskChanges(context.Background(), server.CommitDiskChangesRequest{
		ServerName: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if !res.Status {
		return diag.Errorf("Error upgrading server: %s", res.Msg)
	}

	if err := helper.WaitForAction(conf.Client, job.GetRequest{JobID: resp.Return.JobID, Type: job.DaemonType}); err != nil {

		return diag.FromErr(err)
	}

	return nil
}

// updateLabel is a function to update a label of a server.
func updateLabel(client *server.Client, d *schema.ResourceData) diag.Diagnostics {
	res, err := client.Update(context.Background(), server.UpdateRequest{
		Name:  d.Id(),
		Label: fmt.Sprint(d.Get("label")),
	})
	if err != nil {
		return diag.Errorf("Error updating server: %s", err)
	}

	if !res.Status {
		return diag.Errorf("Error updating server: %s", res.Msg)
	}

	return nil
}

// deleteResource is a function to delete a server.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := server.New(conf.Client)

	resp, err := client.Delete(context.Background(), server.DeleteRequest{
		Name: d.Id(),
	})
	if err != nil {
		return diag.Errorf("Error deleting server: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error deleting server: %s", resp.Msg)
	}

	if err := helper.WaitForAction(conf.Client, job.GetRequest{JobID: resp.Return.JobID, Type: job.DaemonType}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// setServerAttributes is a function to set data to a server.
func setServerAttributes(d *schema.ResourceData, server models.Server) error {
	return d.Set("name", server.Name)
}
