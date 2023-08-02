package grant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSource is the datasource for a cloud database.
func DataSource() *schema.Resource {

	return &schema.Resource{
		ReadContext: readResource,
		Schema:      databaseGrantDataSourceSchema(),
	}
}
