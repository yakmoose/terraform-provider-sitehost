package stack

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack/environment"
	"github.com/sitehostnz/gosh/pkg/models"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
	"gopkg.in/yaml.v3"
	"log"
	"strconv"
	"strings"
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

	serverName := d.Get("server_name").(string)
	name := d.Get("name").(string)

	stackClient := stack.New(conf.Client)
	stackResponse, err := stackClient.Get(ctx, stack.GetRequest{ServerName: serverName, Name: name})
	if err != nil {
		return diag.Errorf("Error retrieving stack info: server %s, stack %s, %s", serverName, stackResponse.Stack, err)
	}

	stack := stackResponse.Stack

	d.SetId(fmt.Sprintf("%s/%s", serverName, name))
	d.Set("server_ip_address", stack.IPAddress)
	d.Set("server_label", stack.Server)
	d.Set("docker_file", stack.DockerFile)
	d.Set("label", stack.Label)

	// this is a nested fun time...
	// d.Set("", stack.Containers)
	// do we need to expose the containers in terraform? we have no real way of getting them
	// or changing them... other than the docker file?

	environmentClient := environment.New(conf.Client)
	environmentVariablesResponse, err := environmentClient.Get(ctx, environment.GetRequest{ServerName: serverName, Project: name, Service: name})
	if err != nil {
		return diag.Errorf("Error retrieving environment info: server %s, stack %s, %s", serverName, stack, err)
	}
	var settings = map[string]string{}
	for _, v := range environmentVariablesResponse.EnvironmentVariables {
		settings[v.Name] = v.Content
	}

	if len(settings) > 0 {
		d.Set("settings", settings)
	}

	// unmarshall the docker file so we can get bits out of it.
	dockerFile := Compose{}
	yaml.Unmarshal([]byte(stack.DockerFile), &dockerFile)

	// set the docker file here, for fun/read it back...
	d.Set("docker_file", stack.DockerFile)

	// the big assumption here... for now... is that we are going to have only one service?
	// get all the settings from the docker compose that we need...
	// things that exist in the yaml from the server

	// 1. virtual hosts
	var aliases []string
	for i := range dockerFile.Services[stack.Name].Environment {
		s := dockerFile.Services[stack.Name].Environment[i]
		if strings.HasPrefix(s, "VIRTUAL_HOST=") {

			aliases = strings.Split(
				strings.TrimPrefix(s, "VIRTUAL_HOST="),
				",",
			)
			aliases = helper.Filter(aliases, func(s string) bool { return s != stack.Label })
			break
		}
	}
	d.Set("aliases", aliases)

	// what happens when these things go missing?
	// and what do we need?
	// container type
	// does this have any impacts on things
	d.Set("type", extractLabelValueFromList(dockerFile.Services[stack.Name].Labels, "nz.sitehost.container.type"))

	v, err := strconv.ParseBool(extractLabelValueFromList(dockerFile.Services[stack.Name].Labels, "nz.sitehost.container.image_update"))
	if err != nil {
		v = false
	}
	d.Set("image_update", v)

	// is the stack monitored
	v, err = strconv.ParseBool(extractLabelValueFromList(dockerFile.Services[stack.Name].Labels, "nz.sitehost.container.monitored"))
	if err != nil {
		v = false
	}
	d.Set("monitored", v)

	// should we disable the backup
	v, err = strconv.ParseBool(extractLabelValueFromList(dockerFile.Services[stack.Name].Labels, "nz.sitehost.container.backup_disable"))
	if err != nil {
		v = false
	}
	d.Set("backup_disable", v)

	// store the containers, if we need this level of detail
	// d.Set("containers", stack.Containers)

	// now find the first container, that is the one, with the things... and stuff...
	// we are mainly interested in the ssl enabled value.
	// and the assumption I am making is that the first one that is true in the stack wins.

	// likely we want to remove this since we can only turn it on, but
	// not turn it off with an update
	enableSSL := false
	for _, container := range stack.Containers {
		if container.SslEnabled {
			enableSSL = true
			break
		}
	}
	d.Set("enable_ssl", enableSSL)

	return nil
}

// createResource is a function to create a stack environment.
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	stackClient := stack.New(conf.Client)

	//1. get a stack container id
	stackNameResponse, err := stackClient.GenerateName(ctx)
	if err != nil {
		return diag.Errorf("Failed to generate stack name: %s", err)
	}
	serverName := d.Get("server_name").(string)
	name := stackNameResponse.Return.Name

	//2. we need to congfiure the docker file
	//3. we need to rollllll out the variables

	settings := d.Get("settings").(map[string]string)
	environmentVariables := make([]models.EnvironmentVariable, len(settings))
	for environmentVariableName, content := range settings {
		environmentVariables = append(environmentVariables, models.EnvironmentVariable{Name: environmentVariableName, Content: content})
	}

	addRequest := stack.AddRequest{
		ServerName:           serverName,
		Name:                 name,
		Label:                d.Get("label").(string),
		EnableSSL:            0,
		DockerCompose:        "",
		EnvironmentVariables: nil,
	}
	log.Printf("[INFO] stack.addRequest: %s", addRequest)

	return diag.Errorf("giving up")
	// set the id once we're actually happy...
	// 	d.SetId(fmt.Sprintf("%s/%s", serverName, name))
	// fire the request
	// stackAddResponse, err := stackClient.Add(ctx, )
	// return nil
}

// updateResource is a function to update a stack environment.
func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//conf, ok := meta.(*helper.CombinedConfig)
	//if !ok {
	//	return diag.Errorf("failed to convert meta object")
	//}

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
