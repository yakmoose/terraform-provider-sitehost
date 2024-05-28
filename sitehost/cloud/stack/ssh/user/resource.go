package user

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/ssh/user"
	"github.com/sitehostnz/gosh/pkg/api/job"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
	"strings"
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

// readResource is a function to read a ssh user
func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := user.New(conf.Client)

	serverName := fmt.Sprint(d.Get("server_name"))
	username := fmt.Sprint(d.Get("username"))

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

	d.SetId(fmt.Sprintf("%s@%s", response.Return.Username, response.Return.ServerName))

	sshkeys := []map[string]string{}
	for _, v := range response.Return.SSHKeys {
		sshkeys = append(
			sshkeys,
			map[string]string{
				"id":      v.ID,
				"label":   v.Label,
				"content": v.Content,
			})
	}
	d.Set("ssh_key", sshkeys)

	containers := []map[string]string{}
	for _, v := range response.Return.Containers {
		containers = append(
			containers,
			map[string]string{
				"name": v,
			})
	}
	d.Set("container", containers)

	volumes := []map[string]string{}
	for _, v := range response.Return.Volumes {
		volumes = append(
			containers,
			map[string]string{
				"name": v,
			})
	}
	d.Set("volumes", volumes)
	d.Set("read_only_config", response.Return.ReadOnlyConfig)

	return nil
}

// createResource is a function to create a stack environment.
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	addRequest := user.AddRequest{
		ServerName: fmt.Sprint(d.Get("server_name")),
		Username:   fmt.Sprint(d.Get("username")),
		Password:   fmt.Sprint(d.Get("password")),
	}

	if val := d.Get("container").([]interface{}); val != nil {
		for _, v := range val {
			addRequest.Containers = append(addRequest.Containers, v.(map[string]interface{})["name"].(string))
		}
	}

	if val := d.Get("volume").([]interface{}); val != nil {
		for _, v := range val {
			addRequest.Volumes = append(addRequest.Volumes, v.(map[string]interface{})["name"].(string))
		}
	}

	if val := d.Get("ssh_key").([]interface{}); val != nil {
		for _, v := range val {
			addRequest.SSHKeys = append(addRequest.SSHKeys, v.(map[string]interface{})["id"].(string))
		}
	}

	if v := d.Get("readonly_config"); v != nil {
		addRequest.ReadOnlyConfig = v.(bool)
	}

	client := user.New(conf.Client)
	_, err := client.Add(ctx, addRequest)

	if err != nil {
		return diag.Errorf("error updating ssh user: server %s, username %s, %s", addRequest.ServerName, addRequest.Username, err)
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

	if val := d.Get("container").([]interface{}); val != nil {
		for _, v := range val {
			updateRequest.Containers = append(updateRequest.Containers, v.(map[string]interface{})["name"].(string))
		}
	}

	if val := d.Get("volume").([]interface{}); val != nil {
		for _, v := range val {
			updateRequest.Volumes = append(updateRequest.Volumes, v.(map[string]interface{})["name"].(string))
		}
	}

	if val := d.Get("ssh_key").([]interface{}); val != nil {
		for _, v := range val {
			updateRequest.SSHKeys = append(updateRequest.SSHKeys, v.(map[string]interface{})["id"].(string))
		}
	}

	if v := d.Get("readonly_config"); v != nil {
		updateRequest.ReadOnlyConfig = v.(bool)
	}

	client := user.New(conf.Client)
	response, err := client.Update(ctx, updateRequest)

	if err != nil {
		return diag.Errorf("error updating ssh user: server %s, username %s, %s", updateRequest.ServerName, updateRequest.Username, err)
	}

	if err := helper.WaitForAction(conf.Client, job.GetRequest{JobID: response.Return.JobID, Type: job.SchedulerType}); err != nil {
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
	serverName := fmt.Sprintf("%v", d.Get("server_name"))
	username := fmt.Sprintf("%v", d.Get("username"))

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

	if err := helper.WaitForAction(conf.Client, job.GetRequest{JobID: response.Return.JobID, Type: job.SchedulerType}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func importResource(ctx context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	split := strings.Split(d.Id(), "@")

	if len(split) != 2 {
		return nil, fmt.Errorf("invalid id: %s. The ID should be in the format [username]@[server_name]", d.Id())
	}

	serverName := split[1]
	username := split[0]

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
