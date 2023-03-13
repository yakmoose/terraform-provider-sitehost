package dns

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/dns"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// ZoneResource returns a schema with the operations for Server resource.
func ZoneResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: dnsZoneCreateResource,
		ReadContext:   dnsZoneReadResource,
		DeleteContext: dnsZoneDeleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: dnsZoneResourceSchema,
	}
}

// createResource is a function to create a new Server resource.
func dnsZoneCreateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	domainName := d.Get("name").(string)

	client := dns.New(conf.Client)
	_, err := client.CreateZone(ctx, dns.CreateZoneRequest{DomainName: domainName})
	if err != nil {
		return diag.Errorf("Error creating domain: %s", err)
	}

	d.SetId(domainName)

	log.Printf("[INFO] Domain Name: %s", d.Id())

	return nil
}

// dnsZoneReadResource is a function to read a new Server resource.
func dnsZoneReadResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	domainResponse, err := client.GetZone(ctx, dns.GetZoneRequest{DomainName: d.Id()})
	if err != nil {
		return diag.Errorf("Error retrieving domain: %s", err)
	}

	if domainResponse.Return == nil {
		return diag.Errorf("Failed retrieving domain: %s", err)
	}

	// domain := domainResponse.Return
	d.SetId(d.Id())
	d.Set("name", d.Id())

	log.Printf("[INFO] Domain Name: %s", d.Id())

	return nil
}

// dnsZoneDeleteResource is a function to delete a new Server resource.
func dnsZoneDeleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	_, err := client.DeleteZone(ctx, dns.DeleteZoneRequest{DomainName: d.Id()})
	if err != nil {
		return diag.Errorf("Error deleting server: %s", err)
	}

	return nil
}
