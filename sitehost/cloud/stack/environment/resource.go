package environment

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack/environment"
	"github.com/sitehostnz/gosh/pkg/models"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// Resource returns a schema with the operations for Server resource.
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: updateResource,
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

// readResource is a function to read a stack environment.
func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	environmentVariablesResponse, err := client.Get(ctx, environment.GetRequest{ServerName: serverName, Project: project, Service: service})
	if err != nil {
		return diag.Errorf("Error retrieving environment info: %s", err)
	}

	settings := map[string]string{}
	for _, v := range environmentVariablesResponse.EnvironmentVariables {
		settings[strings.ToUpper(v.Name)] = v.Content
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", serverName, project, service))
	d.Set("server_name", serverName)
	d.Set("service", service)
	d.Set("project", project)
	d.Set("settings", settings)

	return nil
}

// updateResource is a function to update a stack environment, there is no create environment outside of when you create a stack, these all work on the assumption that the stack exists.
func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	settings := d.Get("settings").(map[string]interface{})
	environmentVariables := []models.EnvironmentVariable{}
	for k, v := range settings {
		environmentVariables = append(environmentVariables, models.EnvironmentVariable{Name: k, Content: v.(string)})
	}

	client := environment.New(conf.Client)
	job, err := client.Update(
		ctx,
		environment.UpdateRequest{
			ServerName:           serverName,
			Project:              project,
			Service:              service,
			EnvironmentVariables: environmentVariables,
		})

	if nil != err {
		return diag.Errorf("Error updating environment info: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", serverName, project, service))

	if err := helper.WaitForAction(conf.Client, job.Return.JobID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// deleteResource is a function to delete a stack environment.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// how do we delete/destry...
	return nil
}

func importResource(ctx context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	split := strings.Split(d.Id(), "/")
	if len(split) != 3 {
		return nil, fmt.Errorf("invalid id: %s. The ID should be in the format [server_name]/[project]/[service]", d.Id())
	}

	serverName := split[0]
	project := split[1]
	service := split[2]

	err := d.Set("server_name", serverName)
	if err != nil {
		return nil, fmt.Errorf("error importing stack environment: server %s, project %s, service %s, %s", serverName, service, project, err)
	}

	err = d.Set("project", project)
	if err != nil {
		return nil, fmt.Errorf("error importing stack environment: server %s, project %s, service %s, %s", serverName, service, project, err)
	}

	err = d.Set("service", service)
	if err != nil {
		return nil, fmt.Errorf("error importing stack environment: server %s, project %s, service %s, %s", serverName, service, project, err)
	}

	return []*schema.ResourceData{d}, nil
}
