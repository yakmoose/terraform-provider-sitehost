package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/ssh/user"
	"github.com/sitehostnz/gosh/pkg/api/job"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// Resource returns a schema with the operations for Server resource.
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResource,
		ReadContext:   readResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,

		// assume this is correct here.... wheeeee
		Importer: &schema.ResourceImporter{
			StateContext: importResource,
		},
		Schema: resourceSchema,
	}
}

// readResource is a function to read an ssh user.
func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := user.New(conf.Client)
	serverName := fmt.Sprint(d.Get("server_name"))
	username := fmt.Sprint(d.Get("username"))
	d.SetId(fmt.Sprintf("%s@%s", username, serverName))

	response, err := client.Get(
		ctx,
		user.GetRequest{
			ServerName: serverName,
			Username:   username,
		},
	)
	if err != nil {
		return diag.Errorf("error retrieving ssh user: server %s, username %s, %s", serverName, username, err)
	}

	sshKeys := make([]map[string]string, 0, len(response.Return.SSHKeys))
	for _, v := range response.Return.SSHKeys {
		sshKeys = append(
			sshKeys,
			map[string]string{
				"id":      v.ID,
				"label":   v.Label,
				"content": v.Content,
			})
	}

	if err := d.Set("ssh_key", sshKeys); err != nil {
		return diag.FromErr(err)
	}

	containers := make([]map[string]string, 0, len(response.Return.Containers))
	for _, v := range response.Return.Containers {
		containers = append(
			containers,
			map[string]string{
				"name": v,
			})
	}

	if err := d.Set("container", containers); err != nil {
		return diag.FromErr(err)
	}

	volumes := make([]map[string]string, 0, len(response.Return.Volumes))
	for _, v := range response.Return.Volumes {
		volumes = append(
			volumes,
			map[string]string{
				"name": v,
			})
	}

	if err := d.Set("volume", volumes); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("read_only_config", response.Return.ReadOnlyConfig); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// createResource is a function to create a stack environment.
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	serverName := fmt.Sprint(d.Get("server_name"))
	username := fmt.Sprint(d.Get("username"))
	d.SetId(fmt.Sprintf("%s@%s", username, serverName))

	addRequest := user.AddRequest{
		ServerName: serverName,
		Username:   username,
		Password:   fmt.Sprint(d.Get("password")),
	}

	if val, ok := d.Get("container").([]interface{}); ok && val != nil {
		addRequest.Containers = make([]string, 0, len(val))
		for _, v := range val {
			if v, ok := v.(map[string]interface{})["name"].(string); ok {
				addRequest.Containers = append(addRequest.Containers, v)
			}
		}
	}

	if val, ok := d.Get("volume").([]interface{}); ok && val != nil {
		for _, v := range val {
			addRequest.Volumes = make([]string, 0, len(val))
			if v, ok := v.(map[string]interface{})["name"].(string); ok {
				addRequest.Volumes = append(addRequest.Volumes, v)
			}
		}
	}

	if val, ok := d.Get("ssh_key").(*schema.Set); ok && val != nil {
		addRequest.SSHKeys = make([]string, 0, val.Len())
		for _, v := range val.List() {
			if v, ok := v.(map[string]interface{})["id"].(string); ok {
				addRequest.SSHKeys = append(addRequest.SSHKeys, v)
			}
		}
	}

	if v, ok := d.Get("readonly_config").(bool); ok {
		addRequest.ReadOnlyConfig = v
	}

	if err := d.Set("read_only_config", addRequest.ReadOnlyConfig); err != nil {
		return diag.FromErr(err)
	}

	client := user.New(conf.Client)
	response, err := client.Add(ctx, addRequest)
	if err != nil {
		return diag.Errorf("error creating ssh user: server %s, username %s, %s", addRequest.ServerName, addRequest.Username, err)
	}

	if err := helper.WaitForAction(conf.Client, job.GetRequest{ID: response.Return.ID, Type: response.Return.Type}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// updateResource is a function to update a stack environment.
func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	updateRequest := user.UpdateRequest{
		ServerName: fmt.Sprint(d.Get("server_name")),
		Username:   fmt.Sprint(d.Get("username")),
		Password:   fmt.Sprint(d.Get("password")),
	}

	if val, ok := d.Get("container").([]interface{}); ok && val != nil {
		updateRequest.Containers = make([]string, 0, len(val))
		for _, v := range val {
			if v, ok := v.(map[string]interface{})["name"].(string); ok {
				updateRequest.Containers = append(updateRequest.Containers, v)
			}
		}
	}

	if val, ok := d.Get("volume").([]interface{}); ok && val != nil {
		for _, v := range val {
			updateRequest.Volumes = make([]string, 0, len(val))
			if v, ok := v.(map[string]interface{})["name"].(string); ok {
				updateRequest.Volumes = append(updateRequest.Volumes, v)
			}
		}
	}

	if val, ok := d.Get("ssh_key").(*schema.Set); ok && val != nil {
		updateRequest.SSHKeys = make([]string, 0, val.Len())
		for _, v := range val.List() {
			if v, ok := v.(map[string]interface{})["id"].(string); ok {
				updateRequest.SSHKeys = append(updateRequest.SSHKeys, v)
			}
		}
	}

	if err := d.Set("read_only_config", updateRequest.ReadOnlyConfig); err != nil {
		return diag.FromErr(err)
	}

	client := user.New(conf.Client)
	response, err := client.Update(ctx, updateRequest)
	if err != nil {
		return diag.Errorf("error updating ssh user: server %s, username %s, %s", updateRequest.ServerName, updateRequest.Username, err)
	}

	if err := helper.WaitForAction(conf.Client, job.GetRequest{ID: response.Return.ID, Type: response.Return.Type}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// deleteResource is a function to delete a stack environment.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := user.New(conf.Client)
	serverName := fmt.Sprint(d.Get("server_name"))
	username := fmt.Sprint(d.Get("username"))

	response, err := client.Delete(
		ctx,
		user.DeleteRequest{
			ServerName: serverName,
			Username:   username,
		},
	)
	if err != nil {
		return diag.Errorf("error deleting ssh user: server %s, username %s, %s", serverName, username, err)
	}

	if err := helper.WaitForAction(conf.Client, job.GetRequest{ID: response.Return.ID, Type: response.Return.Type}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func importResource(_ context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	split := strings.Split(d.Id(), "@")

	if len(split) != 2 {
		return nil, fmt.Errorf("invalid id: %s. The ID should be in the format [username]@[server_name]", d.Id())
	}

	serverName := split[1]
	username := split[0]

	d.SetId(fmt.Sprintf("%s@%s", username, serverName))

	err := d.Set("server_name", serverName)
	if err != nil {
		return nil, fmt.Errorf("error importing user: server %s, username %s, %s", serverName, username, err)
	}

	err = d.Set("username", username)
	if err != nil {
		return nil, fmt.Errorf("error importing user: server %s, username %s, %s", serverName, username, err)
	}

	return []*schema.ResourceData{d}, nil
}
