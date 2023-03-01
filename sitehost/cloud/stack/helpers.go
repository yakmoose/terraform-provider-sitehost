package stack

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack/environment"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
	"gopkg.in/yaml.v3"
	"strings"
)

func extractLabelValueFromList(list []string, label string) (ret string) {
	v := helper.First(list, func(s string) bool { return strings.HasPrefix(s, label+"=") })
	if v != "" {
		ret = strings.TrimPrefix(v, label+"=")
	}
	return ret
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

	// unmarshall the docker file
	dockerFile := Compose{}
	yaml.Unmarshal([]byte(stack.DockerFile), &dockerFile)

	// the big assumption here... for now... is that we are going to have only one service?
	// get all the settings from the docker compose that we need...
	// things that exist in the yaml from the server
	// d.Set("enable_ssl", false)
	// 1. virtual hosts
	var virtualHosts []string
	for i := range dockerFile.Services[stack.Name].Environment {
		s := dockerFile.Services[stack.Name].Environment[i]
		if strings.HasPrefix(s, "VIRTUAL_HOST=") {

			virtualHosts = strings.Split(
				strings.TrimPrefix(s, "VIRTUAL_HOST="),
				",",
			)

			virtualHosts = helper.Filter(virtualHosts, func(s string) bool { return s != stack.Label })

			break
		}
	}

	//d.Set("aliases", virtualHosts)
	//d.Set("type", extractLabelValueFromList(dockerFile.Services[stack.Name].Labels, "nz.sitehost.container.type"))
	//
	//v, err := strconv.ParseBool(extractLabelValueFromList(dockerFile.Services[stack.Name].Labels, "nz.sitehost.container.image_update"))
	//if err != nil {
	//	v = false
	//}
	//d.Set("image_update", v)
	//
	//v, err = strconv.ParseBool(extractLabelValueFromList(dockerFile.Services[stack.Name].Labels, "nz.sitehost.container.monitored"))
	//if err != nil {
	//	v = false
	//}
	//d.Set("monitored", v)
	//
	//v, err = strconv.ParseBool(extractLabelValueFromList(dockerFile.Services[stack.Name].Labels, "nz.sitehost.container.backup_disable"))
	//if err != nil {
	//	v = false
	//}
	//d.Set("backup_disable", v)

	return nil
}
