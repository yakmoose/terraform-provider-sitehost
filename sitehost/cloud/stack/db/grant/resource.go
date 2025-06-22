package grant

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/db/grant"
	"github.com/sitehostnz/gosh/pkg/api/cloud/db/user"
	"github.com/sitehostnz/gosh/pkg/models"
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
		Schema: databaseGrantResourceSchema,
	}
}

// readResource is a function to read a user from a stack database.
func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// fun thing, we don't get grants from the grant endpoint, we get them from the user endpoint
	// we need to get the user data and then filter out the grants that we want
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	userClient := user.New(conf.Client)

	serverName := fmt.Sprint(d.Get("server_name"))
	mysqlHost := fmt.Sprint(d.Get("mysql_host"))
	username := fmt.Sprint(d.Get("username"))
	database := fmt.Sprint(d.Get("database"))
	d.SetId(fmt.Sprintf("%s/%s/%s/%s", serverName, mysqlHost, database, username))

	response, err := userClient.Get(
		ctx,
		user.GetRequest{
			ServerName: serverName,
			MySQLHost:  mysqlHost,
			Username:   username,
		},
	)
	if err != nil {
		return diag.Errorf("error retrieving database user: server %s, host %s, database %s, username %s, %s", serverName, mysqlHost, database, username, err)
	}

	if helper.Has(response.User.Grants, func(g models.Grant) bool {
		return g.DBName == database
	}) {
		g := helper.First(response.User.Grants, func(g models.Grant) bool {
			return g.DBName == database
		})

		// what do the grants look like/get stored as just a list... or an empty list... mmmm
		if err := d.Set("grants", g.Grants); err != nil {
			return diag.FromErr(err)
		}
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

	serverName := fmt.Sprint(d.Get("server_name"))
	mysqlHost := fmt.Sprint(d.Get("mysql_host"))
	username := fmt.Sprint(d.Get("username"))
	database := fmt.Sprint(d.Get("database"))
	d.SetId(fmt.Sprintf("%s/%s/%s/%s", serverName, mysqlHost, database, username))

	var g []string
	if grants, ok := d.Get("grants").([]interface{}); ok {
		g = make([]string, len(grants))
		for i, v := range grants {
			g[i] = fmt.Sprint(v)
		}
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
		return diag.Errorf("error retrieving database grants: server %s, host %s, database %s, username %s, %s", serverName, mysqlHost, database, username, err)
	}

	if err := helper.WaitForJob(conf.Client, response.Return.Job); err != nil {
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
	client := grant.New(conf.Client)

	serverName := fmt.Sprint(d.Get("server_name"))
	mysqlHost := fmt.Sprint(d.Get("mysql_host"))
	username := fmt.Sprint(d.Get("username"))

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

	if err := helper.WaitForJob(conf.Client, response.Return.Job); err != nil {
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

	serverName := fmt.Sprint(d.Get("server_name"))
	mysqlHost := fmt.Sprint(d.Get("mysql_host"))
	username := fmt.Sprint(d.Get("username"))

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

	if err := helper.WaitForJob(conf.Client, response.Return.Job); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// importResource is a function to import a stack database user.
func importResource(_ context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	// todo: create parser/helper for database grants
	split := strings.Split(d.Id(), "/")
	if len(split) != 4 {
		return nil, fmt.Errorf("invalid id: %s. The ID should be in the format [server_name]/[mysql_host]/[database]/[username]", d.Id())
	}

	serverName := split[0]
	mysqlHost := split[1]
	database := split[2]
	username := split[3]

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", serverName, mysqlHost, database, username))

	if err := d.Set("server_name", serverName); err != nil {
		return nil, fmt.Errorf("error importing db: server %s, host %s, username %s, database %s, %s", serverName, mysqlHost, username, database, err)
	}

	if err := d.Set("mysql_host", mysqlHost); err != nil {
		return nil, fmt.Errorf("error importing db: server %s, host %s, username %s, database %s, %s", serverName, mysqlHost, username, database, err)
	}

	if err := d.Set("username", username); err != nil {
		return nil, fmt.Errorf("error importing db: server %s, host %s, username %s, database %s, %s", serverName, mysqlHost, username, database, err)
	}

	if err := d.Set("database", database); err != nil {
		return nil, fmt.Errorf("error importing db: server %s, host %s, username %s, database %s, %s", serverName, mysqlHost, username, database, err)
	}

	return []*schema.ResourceData{d}, nil
}
