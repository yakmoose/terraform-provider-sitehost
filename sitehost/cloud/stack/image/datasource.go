// Package image stack provides the functions to create/get cloud stacks resource via SiteHost API.
package image

//
//import (
//	"context"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//	"github.com/sitehostnz/gosh/pkg/api/cloud/image"
//	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
//)
//
//// DataSource returns a schema with the function to read cloud stack resource.
//func DataSource() *schema.Resource {
//	return &schema.Resource{
//		ReadContext: readDataSource,
//		Schema:      imagesDatasourceSchema(),
//	}
//}
//
//// readDataSource calls the GoSH client to set the cloud stack schema.
//func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
//	conf, ok := meta.(*helper.CombinedConfig)
//	if !ok {
//		return diag.Errorf("failed to convert meta object")
//	}
//
//	client := image.New(conf.Client)
//
//	resp, err := client.List(ctx)
//
//	// now put them in the things....
//
//}
