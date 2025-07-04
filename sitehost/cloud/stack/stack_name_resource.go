package stack

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/cloud/stack"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// NameResource is a simple helper for creating container names independently of the stack.
func NameResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createStackNameResource,
		ReadContext:   readStackNameResource,
		DeleteContext: deleteStackNameResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The cloud stack name",
			},
		},
	}
}

func readStackNameResource(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// this is a no-op.  we don't need to read anything from the API
	// we just need to keep the number.
	// if this changes, it will result in the stack relying on this being removed
	if err := d.Set("name", d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func createStackNameResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := stack.New(conf.Client)
	response, err := client.GenerateName(ctx)
	if err != nil {
		return diag.Errorf("Error generating stack name: %s", err)
	}

	d.SetId(response.Return.Name)
	if err := d.Set("name", response.Return.Name); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteStackNameResource(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// we don't need to delete anything from the API, only need to remove the item from the state
	return nil
}
