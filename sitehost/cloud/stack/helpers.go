package stack

import (
	"strings"

	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

func extractLabelValueFromList(list []string, label string) (ret string) {
	v := helper.First(list, func(s string) bool { return strings.HasPrefix(s, label+"=") })
	if v != "" {
		ret = strings.TrimPrefix(v, label+"=")
	}
	return ret
}
