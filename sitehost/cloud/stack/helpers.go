package stack

import (
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
	"strings"
)

func extractLabelValueFromList(list []string, label string) (ret string) {
	v := helper.First(list, func(s string) bool { return strings.HasPrefix(s, label+"=") })
	if v != "" {
		ret = strings.TrimPrefix(v, label+"=")
	}
	return ret
}
