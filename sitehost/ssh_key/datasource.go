package sshkey

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sshkey "github.com/sitehostnz/gosh/pkg/api/ssh/key"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// DataSource returns a schema with the function to read Server resource.
func DataSource() *schema.Resource {
	recordSchema := sshKeyDataSourceSchema()

	return &schema.Resource{
		ReadContext: readDataSource,
		Schema:      recordSchema,
	}
}

// readDataSource is a function to read an SSH Key.
func readDataSource(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := sshkey.New(conf.Client)

	resp, err := client.Get(context.Background(), sshkey.GetRequest{
		ID: d.Id(),
	})
	if err != nil {
		return diag.Errorf("Error retrieving SSH Key: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error retrieving SSH Key: %s", resp.Msg)
	}

	if diagErr := setData(resp, d); diagErr != nil {
		return diagErr
	}

	return nil
}
