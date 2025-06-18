package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/db/user"
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
		Schema: databaseUserResourceSchema,
	}
}

// readResource is a function to read a user from a stack database.
func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := user.New(conf.Client)

	serverName := fmt.Sprintf("%v", d.Get("server_name"))
	mysqlHost := fmt.Sprintf("%v", d.Get("mysql_host"))
	username := fmt.Sprintf("%v", d.Get("username"))

	d.SetId(fmt.Sprintf("%s/%s/%s", serverName, mysqlHost, username))

	_, err := client.Get(
		ctx,
		user.GetRequest{
			ServerName: serverName,
			MySQLHost:  mysqlHost,
			Username:   username,
		},
	)
	if err != nil {
		return diag.Errorf("error retrieving database user: server %s, host %s, username %s, %s", serverName, mysqlHost, username, err)
	}

	return nil
}

// createResource is a function to create a stack database user.
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}
	client := user.New(conf.Client)

	serverName := fmt.Sprintf("%v", d.Get("server_name"))
	mysqlHost := fmt.Sprintf("%v", d.Get("mysql_host"))
	database := fmt.Sprintf("%v", d.Get("database"))
	username := fmt.Sprintf("%v", d.Get("username"))
	password := fmt.Sprintf("%v", d.Get("password"))

	d.SetId(fmt.Sprintf("%s/%s/%s", serverName, mysqlHost, username))

	response, err := client.Add(
		ctx,
		user.AddRequest{
			ServerName: serverName,
			MySQLHost:  mysqlHost,
			Username:   username,
			Password:   password,
			Database:   database,
		},
	)
	if err != nil {
		return diag.Errorf("error creating database user: server %s, host %s, username %s, %s", serverName, mysqlHost, username, err)
	}

	if err := helper.WaitForAction(conf.Client, job.GetRequest{ID: response.Return.ID, Type: response.Return.Type}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// updateResource is a function to update a stack database user.
func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}
	client := user.New(conf.Client)

	serverName := fmt.Sprintf("%v", d.Get("server_name"))
	mysqlHost := fmt.Sprintf("%v", d.Get("mysql_host"))
	username := fmt.Sprintf("%v", d.Get("username"))
	password := fmt.Sprintf("%v", d.Get("password"))

	response, err := client.Update(
		ctx,
		user.UpdateRequest{
			ServerName: serverName,
			MySQLHost:  mysqlHost,
			Username:   username,
			Password:   password,
		},
	)
	if err != nil {
		return diag.Errorf("error updating database user: server %s, host %s, username %s, %s", serverName, mysqlHost, username, err)
	}

	if err := helper.WaitForAction(conf.Client, job.GetRequest{ID: response.Return.ID, Type: response.Return.Type}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// deleteResource is a function to delete a stack database user.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}
	client := user.New(conf.Client)

	serverName := fmt.Sprintf("%v", d.Get("server_name"))
	mysqlHost := fmt.Sprintf("%v", d.Get("mysql_host"))
	username := fmt.Sprintf("%v", d.Get("username"))

	response, err := client.Delete(
		ctx,
		user.DeleteRequest{
			ServerName: serverName,
			MySQLHost:  mysqlHost,
			Username:   username,
		},
	)
	if err != nil {
		return diag.Errorf("error removing database user: server %s, host %s, username %s, %s", serverName, mysqlHost, username, err)
	}

	if err := helper.WaitForAction(conf.Client, job.GetRequest{ID: response.Return.ID, Type: response.Return.Type}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// importResource is a function to import a stack database user.
func importResource(_ context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	// todo: create parser for database user
	split := strings.Split(d.Id(), "/")
	if len(split) != 3 {
		return nil, fmt.Errorf("invalid id: %s. The ID should be in the format [server_name]/[mysql_host]/[username]", d.Id())
	}

	serverName := split[0]
	mysqlHost := split[1]
	username := split[2]

	d.SetId(fmt.Sprintf("%s/%s/%s", serverName, mysqlHost, username))

	if err := d.Set("server_name", serverName); err != nil {
		return nil, fmt.Errorf("error importing db: server %s, host %s, username %s, %s", serverName, mysqlHost, username, err)
	}

	if err := d.Set("mysql_host", mysqlHost); err != nil {
		return nil, fmt.Errorf("error importing db: server %s, host %s, username %s, %s", serverName, mysqlHost, username, err)
	}

	if err := d.Set("username", username); err != nil {
		return nil, fmt.Errorf("error importing db: server %s, host %s, username %s, %s", serverName, mysqlHost, username, err)
	}

	return []*schema.ResourceData{d}, nil
}
