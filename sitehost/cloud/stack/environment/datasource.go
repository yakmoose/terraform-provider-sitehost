package environment

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack/environment"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

func DataSource() *schema.Resource {
	recordSchema := stackDataSourceSchema()
	return &schema.Resource{
		ReadContext: readDataSource,
		Schema:      recordSchema,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	serverName := d.Get("server_name").(string)
	project := d.Get("project").(string)
	service := d.Get("service").(string)

	if service == "" {
		service = project
	}

	client := environment.New(conf.Client)
	environmentVariablesResponse, err := client.Get(ctx, environment.GetRequest{
		ServerName: serverName,
		Project:    project,
		Service:    service,
	})
	if err != nil {
		return diag.Errorf("Error retrieving environment info: %s", err)
	}

	var settings = map[string]string{}
	for _, v := range environmentVariablesResponse.EnvironmentVariables {
		settings[v.Name] = v.Content
	}

	// id is the server/stack/project combo...
	d.SetId(fmt.Sprintf("%s/%s/%s", serverName, project, service))
	d.Set("server_name", serverName)
	d.Set("service", service)
	d.Set("project", project)
	d.Set("settings", settings)

	return nil
}
