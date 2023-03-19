package environment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack/environment"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// DataSource is the datasource for environments.
func DataSource() *schema.Resource {
	recordSchema := stackDataSourceSchema()
	return &schema.Resource{
		ReadContext: readDataSource,
		Schema:      recordSchema,
	}
}

// readDataSource reads the environment.
func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	serverName := fmt.Sprintf("%v", d.Get("server_name"))
	project := fmt.Sprintf("%v", d.Get("project"))
	service := fmt.Sprintf("%v", d.Get("service"))

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

	settings := map[string]string{}
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
