package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/db"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// Resource returns a schema with the operations for Server resource.
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResource,
		ReadContext:   readResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,

		// assume this is correct here....
		Importer: &schema.ResourceImporter{
			StateContext: importResource,
		},
		Schema: databaseResourceSchema,
	}
}

// readResource is a function to read a stack environment.
func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}
	client := db.New(conf.Client)

	serverName := fmt.Sprint(d.Get("server_name"))
	mysqlHost := fmt.Sprint(d.Get("mysql_host"))
	database := fmt.Sprint(d.Get("name"))

	d.SetId(fmt.Sprintf("%s/%s/%s", serverName, mysqlHost, database))

	response, err := client.Get(
		ctx,
		db.GetRequest{
			ServerName: serverName,
			MySQLHost:  mysqlHost,
			Database:   database,
		},
	)
	if err != nil {
		return diag.Errorf("error retrieving stack: server %s, name %s, database %s, %s", serverName, mysqlHost, database, err)
	}

	if err := d.Set("backup_container", response.Database.Container); err != nil {
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
	client := db.New(conf.Client)

	serverName := fmt.Sprint(d.Get("server_name"))
	mysqlHost := fmt.Sprint(d.Get("mysql_host"))
	database := fmt.Sprint(d.Get("name"))
	container := fmt.Sprint(d.Get("backup_container"))

	d.SetId(fmt.Sprintf("%s/%s/%s", serverName, mysqlHost, database))

	response, err := client.Add(
		ctx,
		db.AddRequest{
			ServerName: serverName,
			MySQLHost:  mysqlHost,
			Database:   database,
			Container:  container,
		},
	)
	if err != nil {
		return diag.Errorf("error retrieving db: server %s, name %s, database %s, %s", serverName, mysqlHost, database, err)
	}

	if err := helper.WaitForJob(conf.Client, response.Return.Job); err != nil {
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
	client := db.New(conf.Client)

	serverName := fmt.Sprint(d.Get("server_name"))
	mysqlHost := fmt.Sprint(d.Get("mysql_host"))
	database := fmt.Sprint(d.Get("name"))
	container := fmt.Sprint(d.Get("backup_container"))

	_, err := client.Update(
		ctx,
		db.UpdateRequest{
			ServerName: serverName,
			MySQLHost:  mysqlHost,
			Database:   database,
			Container:  container,
		},
	)
	if err != nil {
		return diag.Errorf("error updating db: server %s, name %s, database %s, %s", serverName, mysqlHost, database, err)
	}

	return nil
}

// deleteResource is a function to delete a stack environment.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}
	client := db.New(conf.Client)

	serverName := fmt.Sprint(d.Get("server_name"))
	mysqlHost := fmt.Sprint(d.Get("mysql_host"))
	database := fmt.Sprint(d.Get("name"))

	response, err := client.Delete(
		ctx,
		db.DeleteRequest{
			ServerName: serverName,
			MySQLHost:  mysqlHost,
			Database:   database,
		},
	)
	if err != nil {
		return diag.Errorf("error removing db: server %s, name %s, database %s, %s", serverName, mysqlHost, database, err)
	}

	if err := helper.WaitForJob(conf.Client, response.Return.Job); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func importResource(_ context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	// todo: make database name parser helper
	split := strings.Split(d.Id(), "/")
	if len(split) != 3 {
		return nil, fmt.Errorf("invalid id: %s. The ID should be in the format [server_name]/[mysql_host]/[database]", d.Id())
	}

	serverName := split[0]
	mysqlHost := split[1]
	database := split[2]

	d.SetId(fmt.Sprintf("%s/%s/%s", serverName, mysqlHost, database))

	if err := d.Set("server_name", serverName); err != nil {
		return nil, fmt.Errorf("error importing db: server %s, name %s, database %s, %s", serverName, mysqlHost, database, err)
	}

	if err := d.Set("mysql_host", mysqlHost); err != nil {
		return nil, fmt.Errorf("error importing db: server %s, name %s, database %s, %s", serverName, mysqlHost, database, err)
	}

	if err := d.Set("name", database); err != nil {
		return nil, fmt.Errorf("error importing db: server %s, name %s, database %s, %s", serverName, mysqlHost, database, err)
	}

	return []*schema.ResourceData{d}, nil
}
