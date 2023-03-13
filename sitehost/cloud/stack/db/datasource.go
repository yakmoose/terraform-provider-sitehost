package db

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	recordSchema := databaseDataSourceSchema()

	return &schema.Resource{
		ReadContext: readResource,
		Schema:      recordSchema,
	}
}
