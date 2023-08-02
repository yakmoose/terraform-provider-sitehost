package db

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/db"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// DataSource is the datasource for a cloud database.
func DataSource() *schema.Resource {

	return &schema.Resource{
		ReadContext: readResource,
		Schema:      databaseDataSourceSchema(),
	}
}

// ListDataSource is the datasource for listing databases, with a filter.
func ListDataSource() *schema.Resource {

	return &schema.Resource{
		ReadContext: listResource,
		Schema:      listDatabaseDataSourceSchema(),
	}
}

func listResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}
	client := db.New(conf.Client)

	listResponse, err := client.List(ctx, db.ListOptions{})

	if err != nil {
		return diag.Errorf("Failed to fetch database list %s", err)
	}

	databases := []map[string]string{}
	for _, v := range listResponse.Return.Databases {
		d := map[string]string{
			"name":             v.DBName,
			"server_id":        v.ServerID,
			"server_name":      v.ServerName,
			"server_label":     v.ServerLabel,
			"mysql_host":       v.MySQLHost,
			"backup_container": v.Container,

			// I've intentionally left out the grants here
		}

		databases = append(databases, d)
	}

	d.SetId("databases")

	if err := d.Set("databases", databases); err != nil {
		return diag.Errorf("Error retrieving api info: %s", err)
	}

	return nil
}
