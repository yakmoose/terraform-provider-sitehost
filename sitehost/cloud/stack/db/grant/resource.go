package grant

import (
	"context"
	"fmt"
	"github.com/sitehostnz/gosh/pkg/api/cloud/db/grant"
	"github.com/sitehostnz/gosh/pkg/api/cloud/db/user"
	"github.com/sitehostnz/gosh/pkg/api/job"
	"github.com/sitehostnz/gosh/pkg/models"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"

	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		Schema: databaseGrantResourceSchema(),
	}
}

// readResource is a function to read a user from a stack database.
func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// fun thing, we don't get grants from the grant endpoint, we get them from the user endpoint
	// so we neeed to get the user data and then filter out the grants that we want
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	userClient := user.New(conf.Client)
	// grantClient := grant.New(conf.Client)

	serverName := fmt.Sprintf("%v", d.Get("server_name"))
	mysqlHost := fmt.Sprintf("%v", d.Get("mysql_host"))
	username := fmt.Sprintf("%v", d.Get("username"))
	database := fmt.Sprintf("%v", d.Get("database"))

	response, err := userClient.Get(
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

	if helper.Has(response.User.Grants, func(g models.Grant) bool {
		return g.DBName == database
	}) {
		g := helper.First(response.User.Grants, func(g models.Grant) bool {
			return g.DBName == database
		})
		d.SetId(fmt.Sprintf("%s/%s/%s", response.User.ServerName, response.User.MysqlHost, response.User.Username))
		d.Set("grants", g.Grants)
	}

	return nil
}

// createResource is a function to create a stack database user.
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}
	client := grant.New(conf.Client)

	serverName := fmt.Sprintf("%v", d.Get("server_name"))
	mysqlHost := fmt.Sprintf("%v", d.Get("mysql_host"))
	username := fmt.Sprintf("%v", d.Get("username"))
	database := fmt.Sprintf("%v", d.Get("database"))

	grants := d.Get("grants").([]interface{})
	g := make([]string, len(grants))
	for i, v := range grants {
		g[i] = fmt.Sprint(v)
	}

	response, err := client.Add(
		ctx,
		grant.AddRequest{
			ServerName: serverName,
			MySQLHost:  mysqlHost,
			Username:   username,
			Grants:     g,
			Database:   database,
		},
	)

	if err != nil {
		return diag.Errorf("error retrieving database grants: server %s, host %s, username %s, %s", serverName, mysqlHost, username, err)
	}

	if err := helper.WaitForAction(conf.Client, job.GetRequest{JobID: response.Return.JobID, Type: job.SchedulerType}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", serverName, mysqlHost, username))

	return nil
}

// updateResource is a function to update a stack database user.
func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}
	client := grant.New(conf.Client)

	serverName := fmt.Sprintf("%v", d.Get("server_name"))
	mysqlHost := fmt.Sprintf("%v", d.Get("mysql_host"))
	username := fmt.Sprintf("%v", d.Get("username"))

	response, err := client.Update(
		ctx,
		grant.UpdateRequest{
			ServerName: serverName,
			MySQLHost:  mysqlHost,
			Username:   username,
			//			Grants:     d.Get("grants").([]interface{}),
		},
	)

	if err != nil {
		return diag.Errorf("error updating database grant: server %s, host %s, username %s, %s", serverName, mysqlHost, username, err)
	}

	if err := helper.WaitForAction(conf.Client, job.GetRequest{JobID: response.Return.JobID, Type: job.SchedulerType}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// deleteResource is a function to delete a stack database grant.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}
	client := grant.New(conf.Client)

	serverName := fmt.Sprintf("%v", d.Get("server_name"))
	mysqlHost := fmt.Sprintf("%v", d.Get("mysql_host"))
	username := fmt.Sprintf("%v", d.Get("username"))

	response, err := client.Delete(
		ctx,
		grant.DeleteRequest{
			ServerName: serverName,
			MySQLHost:  mysqlHost,
			Username:   username,
		},
	)
	if err != nil {
		return diag.Errorf("error removing database user: server %s, host %s, username %s, %s", serverName, mysqlHost, username, err)
	}

	if err := helper.WaitForAction(conf.Client, job.GetRequest{JobID: response.Return.JobID, Type: job.SchedulerType}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// importResource is a function to import a stack database user.
func importResource(ctx context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	split := strings.Split(d.Id(), "/")

	if len(split) != 3 {
		return nil, fmt.Errorf("invalid id: %s. The ID should be in the format [server_name]/[mysql_host]/[username]", d.Id())
	}

	serverName := split[0]
	mysqlHost := split[1]
	username := split[2]

	err := d.Set("server_name", serverName)
	if err != nil {
		return nil, fmt.Errorf("error importing db: server %s, host %s, username %s, %s", serverName, mysqlHost, username, err)
	}

	err = d.Set("mysql_host", mysqlHost)
	if err != nil {
		return nil, fmt.Errorf("error importing db: server %s, host %s, username %s, %s", serverName, mysqlHost, username, err)
	}

	err = d.Set("username", username)
	if err != nil {
		return nil, fmt.Errorf("error importing db: server %s, host %s, username %s, %s", serverName, mysqlHost, username, err)
	}

	return []*schema.ResourceData{d}, nil
}
