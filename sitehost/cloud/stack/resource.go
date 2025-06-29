package stack

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
	"gopkg.in/yaml.v3"
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
		Schema: resourceSchema,
	}
}

// readResource is a function to read a stack environment.
func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	serverName := fmt.Sprintf("%v", d.Get("server_name"))
	name := fmt.Sprintf("%v", d.Get("name"))
	d.SetId(fmt.Sprintf("%s/%s", serverName, name))

	stackClient := stack.New(conf.Client)
	stackResponse, err := stackClient.Get(ctx, stack.GetRequest{ServerName: serverName, Name: name})
	if err != nil {
		return diag.Errorf("Error retrieving stack info: server %s, stack %s, %s", serverName, name, err)
	}

	s := stackResponse.Stack
	if err := d.Set("server_ip_address", s.IPAddress); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("server_label", s.Server); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("docker_file", s.DockerFile); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("label", s.Label); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("server_name", serverName); err != nil {
		return diag.FromErr(err)
	}

	// unmarshall the docker file so we can get bits out of it.
	dockerFile := Compose{}
	if err := yaml.Unmarshal([]byte(s.DockerFile), &dockerFile); err != nil {
		return diag.FromErr(err)
	}

	// set the docker file here, for fun/read it back...
	if err := d.Set("docker_file", s.DockerFile); err != nil {
		return diag.FromErr(err)
	}

	// the big assumption here... for now... is that we are going to have only one service?
	// get all the settings from the docker compose that we need...
	// things that exist in the yaml from the server

	// 1. virtual hosts
	aliases := extractAliasesFromDockerFile(dockerFile, s)
	if err := d.Set("aliases", aliases); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("type", extractLabelValueFromList(dockerFile.Services[s.Name].Labels, "nz.sitehost.container.type")); err != nil {
		return diag.FromErr(err)
	}

	v, err := strconv.ParseBool(extractLabelValueFromList(dockerFile.Services[s.Name].Labels, "nz.sitehost.container.image_update"))
	if err != nil {
		v = false
	}

	if err := d.Set("image_update", v); err != nil {
		return diag.FromErr(err)
	}

	// is the stack monitored
	v, err = strconv.ParseBool(extractLabelValueFromList(dockerFile.Services[s.Name].Labels, "nz.sitehost.container.monitored"))
	if err != nil {
		v = false
	}

	if err := d.Set("monitored", v); err != nil {
		return diag.FromErr(err)
	}

	// should we disable the backup?
	v, err = strconv.ParseBool(extractLabelValueFromList(dockerFile.Services[s.Name].Labels, "nz.sitehost.container.backup_disable"))
	if err != nil {
		v = false
	}

	if err := d.Set("backup_disable", v); err != nil {
		return diag.FromErr(err)
	}

	// now find the first container that is the one, with the things... and stuff...
	// we are mainly interested in the ssl enabled value.
	// and the assumption I am making is that the first one that is true in the stack wins.

	// likely we want to remove this since we can only turn it on, but
	// not turn it off with an update
	enableSSL := false
	for _, container := range s.Containers {
		if container.SslEnabled {
			enableSSL = true
			break
		}
	}
	if err := d.Set("enable_ssl", enableSSL); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// createResource is a function to create a stack environment.
func createResource(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Errorf("giving up")
}

// updateResource is a function to update a stack environment.
func updateResource(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Errorf("giving up")
}

// deleteResource is a function to delete a stack environment.
func deleteResource(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Errorf("giving up")
}

func importResource(_ context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	parsedStackName, err := ParseStackName(d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", parsedStackName.ServerName, parsedStackName.Project, parsedStackName.Service))

	return []*schema.ResourceData{d}, nil
}
