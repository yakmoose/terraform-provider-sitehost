package user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSource is the datasource for a cloud database.
func DataSource() *schema.Resource {
	recordSchema := databaseUserDataSourceSchema()

	return &schema.Resource{
		ReadContext: readResource,
		Schema:      recordSchema,
	}
}
