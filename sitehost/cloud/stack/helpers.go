package stack

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sitehostnz/gosh/pkg/models"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

func extractLabelValueFromList(list []string, label string) (ret string) {
	v := helper.First(list, func(s string) bool { return strings.HasPrefix(s, label+"=") })
	if v != "" {
		ret = strings.TrimPrefix(v, label+"=")
	}
	return ret
}

func extractAliasesFromDockerFile(dockerFile Compose, s models.Stack) []string {
	var aliases []string
	for i := range dockerFile.Services[s.Name].Environment {
		service := dockerFile.Services[s.Name].Environment[i]
		if strings.HasPrefix(service, "VIRTUAL_HOST=") {
			aliases = strings.Split(
				strings.TrimPrefix(service, "VIRTUAL_HOST="),
				",",
			)
			aliases = helper.Filter(aliases, func(l string) bool { return l != s.Label })
			break
		}
	}
	return aliases
}

// ParseStackName parses a stack identifier into its components: server name, project, and service.
// Returns a ParsedStackName struct or an error if the ID format is invalid.
func ParseStackName(id string) (stack *ParsedStackName, err error) {
	split := strings.Split(id, "/")

	switch len(split) {
	case 3:
		return &ParsedStackName{
			ServerName: split[0],
			Project:    split[1],
			Service:    split[2],
		}, nil
	case 2:
		return &ParsedStackName{
			ServerName: split[0],
			Project:    split[1],
			Service:    split[1],
		}, nil
	default:
		re := regexp.MustCompile(`^(?:(?:https://cp.sitehost.nz)?/cloud/manage-container)?/server/(?P<server>[-_.a-zA-Z0-9]+)/stack/(?P<project>[-_.a-zA-Z0-9]+)/?$`)
		matches := re.FindStringSubmatch(id)
		if matches != nil {
			return &ParsedStackName{
				ServerName: matches[1],
				Project:    matches[2],
				Service:    matches[2],
			}, nil
		}
	}

	return nil, fmt.Errorf(
		"invalid id: %s.\n\n"+
			"The ID should be in one of the following formats:\n"+
			"https://cp.sitehost.nz/clound/manage-container/server/[server_name]/stack/[project]\n"+
			"/cloud/manage-container/server/[server_name]/stack/[project]\n"+
			"/server/[server_name]/stack/[project]\n"+
			"[server_name]/[project]/[service]\n"+
			"[server_name]/[project]",
		id)
}
