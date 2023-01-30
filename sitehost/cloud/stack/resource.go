package stack

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack/environment"
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

		//assume this is correct here.... wheeeee
		Importer: &schema.ResourceImporter{
			StateContext: importResource,
		},
		Schema: resourceSchema,
	}
}

// createResource is a function to create a stack environment.
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//conf, ok := meta.(*helper.CombinedConfig)
	//if !ok {
	//	return diag.Errorf("failed to convert meta object")
	//}
	//
	//client := domain.New(conf.Client)
	//domain, err := client.Create(ctx, &models.Domain{Name: d.Get("name").(string)})
	//
	//if err != nil {
	//	return diag.Errorf("Error creating domain: %s", err)
	//}
	//if domain == nil {
	//	return diag.Errorf("Failed retrieving domain: %s", err)
	//}
	//d.SetId(domain.Name)
	//log.Printf("[INFO] Domain Name: %s", d.Id())
	//creating takes a bunch of stuff, and returns a job
	// wait for the job and stuf it back
	// there are a bunch of thigns that we get that are computed

	return nil
}

// readResource is a function to read a stack environment.
func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	serverName := d.Get("server").(string)
	name := d.Get("name").(string)

	stackClient := stack.New(conf.Client)
	stack, err := stackClient.Get(ctx, stack.GetRequest{ServerName: serverName, Name: name})
	if err != nil {
		return diag.Errorf("Error retrieving stack info: server %s, stack %s, %s", serverName, stack, err)
	}

	environmentClient := environment.New(conf.Client)
	environment, err := environmentClient.Get(ctx, environment.GetRequest{ServerName: serverName, Project: name, Service: name})
	if err != nil {
		return diag.Errorf("Error retrieving environment info: server %s, stack %s, %s", serverName, stack, err)
	}

	var settings = map[string]string{}
	for _, v := range *environment {
		settings[v.Name] = v.Content
	}

	d.SetId(fmt.Sprintf("%s/%s", serverName, name))
	d.Set("server_name", serverName)

	if len(settings) > 0 {
		d.Set("settings", settings)
	}

	// what things do we need to set here???
	// all the computed thigns that we may or may not have...

	return nil
}

// updateResource is a function to update a stack environment.
func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//conf, ok := meta.(*helper.CombinedConfig)
	//if !ok {
	//	return diag.Errorf("failed to convert meta object")
	//}
	//
	//client := domain.New(conf.Client)
	//domain, err := client.Create(ctx, &models.Domain{Name: d.Get("name").(string)})

	return nil
}

// deleteResource is a function to delete a stack environment.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//conf, ok := meta.(*helper.CombinedConfig)
	//if !ok {
	//	return diag.Errorf("failed to convert meta object")
	//}
	//
	//client := domain.New(conf.Client)
	//domain, err := client.Create(ctx, &models.Domain{Name: d.Get("name").(string)})

	return nil
}

func importResource(ctx context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	split := strings.Split(d.Id(), "/")

	if len(split) != 2 {
		return nil, fmt.Errorf("invalid id: %s. The ID should be in the format [server_name]/[stack]", d.Id())
	}

	serverName := split[0]
	name := split[1]

	err := d.Set("server_name", serverName)
	if err != nil {
		return nil, fmt.Errorf("error importing stack: server %s, name %s", serverName, name, err)
	}

	//
	err = d.Set("name", name)
	if err != nil {
		return nil, fmt.Errorf("error importing stack environment: server %s, name %s, service %s, %s", serverName, name, err)
	}

	return []*schema.ResourceData{d}, nil
}
