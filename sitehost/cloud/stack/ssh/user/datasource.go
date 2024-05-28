// Package user provides the functions to create/get cloud users resource via SiteHost API.
package user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSource returns a schema with the function to read cloud stack resource.
func DataSource() *schema.Resource {
	recordSchema := userDataSourceSchema()

	return &schema.Resource{
		ReadContext: readResource,
		Schema:      recordSchema,
	}
}
