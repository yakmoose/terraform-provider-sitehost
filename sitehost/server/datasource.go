package server

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSource returns a schema with the function to read Server resource.
func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSource,
		Schema:      serverDataSourceSchema,
	}
}

// readDataSource is a function to read servers (not implemented).
func readDataSource(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
	fmt.Println("not implemented")
	return nil
}
