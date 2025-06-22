package environment

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack/environment"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/cloud/stack"
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

	serverName := fmt.Sprint(d.Get("server_name"))
	project := fmt.Sprint(d.Get("project"))
	service := fmt.Sprint(d.Get("service"))

	if service == "" {
		service = project
		if d.Set("service", service) != nil {
			return nil
		}
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", serverName, project, service))

	client := environment.New(conf.Client)
	environmentVariablesResponse, err := client.Get(
		ctx,
		environment.GetRequest{ServerName: serverName, Project: project, Service: service},
	)

	if err != nil {
		return diag.Errorf("Error retrieving environment info: %s", err)
	}

	settings := map[string]string{}
	for _, v := range environmentVariablesResponse.EnvironmentVariables {
		settings[strings.ToUpper(v.Name)] = v.Content
	}

	if err := d.Set("settings", settings); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// updateResource is a function to update a stack environment. There is no create environment outside creating a stack, these all work on the assumption that the stack exists.
func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	serverName := fmt.Sprint(d.Get("server_name"))
	project := fmt.Sprint(d.Get("project"))
	service := fmt.Sprint(d.Get("service"))

	if service == "" {
		service = project
		if d.Set("service", service) != nil {
			return nil
		}
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", serverName, project, service))

	if !d.HasChange("settings") {
		return nil
	}

	ov, nv := d.GetChange("settings")
	environmentVariables := createEnvironmentVariableChangeSet(ov, nv)

	// if we have changes... then we need to push em...
	// what happens if we have an empty list...
	client := environment.New(conf.Client)

	if len(environmentVariables) > 0 {
		response, err := client.Update(
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

		if err := helper.WaitForJob(conf.Client, response.Return.Job); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

// deleteResource is a function to delete a stack environment.
func deleteResource(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// the environment doesn't go away, we just clear it out...
	return nil
}

func importResource(_ context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	parsedStackName, err := stack.ParseStackName(d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", parsedStackName.ServerName, parsedStackName.Project, parsedStackName.Service))

	if err := d.Set("server_name", parsedStackName.ServerName); err != nil {
		return nil, fmt.Errorf("error importing stack environment: server %s, project %s, service %s, %s", parsedStackName.ServerName, parsedStackName.Service, parsedStackName.Project, err)
	}

	if err := d.Set("project", parsedStackName.Project); err != nil {
		return nil, fmt.Errorf("error importing stack environment: server %s, project %s, service %s, %s", parsedStackName.ServerName, parsedStackName.Service, parsedStackName.Project, err)
	}

	if err := d.Set("service", parsedStackName.Service); err != nil {
		return nil, fmt.Errorf("error importing stack environment: server %s, project %s, service %s, %s", parsedStackName.ServerName, parsedStackName.Service, parsedStackName.Project, err)
	}

	return []*schema.ResourceData{d}, nil
}
