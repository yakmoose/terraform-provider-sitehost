// Package dns provides the functions to create/get dns zones/records resource via SiteHost API.
package dns

import (
	"context"
	"fmt"
	"github.com/sitehostnz/gosh/pkg/utils"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/dns"
	"github.com/sitehostnz/gosh/pkg/models"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// ZoneResource returns a schema with the operations for DNS Zone resource.
func ZoneResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createZoneResource,
		ReadContext:   readZoneResource,
		DeleteContext: deleteZoneResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceZoneSchema,
	}
}

// createZoneResource is a function to create a new DNS Zone.
func createZoneResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	domain := fmt.Sprintf("%v", d.Get("name"))

	// The response don't have the domain name, so we need to get it from the request.
	resp, err := client.CreateZone(ctx, dns.CreateZoneRequest{DomainName: domain})
	if err != nil {
		return diag.Errorf("Error creating domain: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error creating domain: %s", resp.Msg)
	}

	d.SetId(domain)

	log.Printf("[INFO] Domain Name: %s", d.Id())

	return nil
}

// readZoneResource is a function to read a DNS Zone.
func readZoneResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	response, err := client.GetZone(ctx, dns.GetZoneRequest{DomainName: d.Id()})
	if err != nil {
		return diag.Errorf("Error retrieving domain: %s", err)
	}

	if !response.Status {
		return diag.Errorf("Error retrieving domain: %s", response.Msg)
	}

	// iterate over the domains to find the one we are looking for.
	for _, domain := range response.Return {
		if domain.Name == d.Id() {
			if err := d.Set("name", domain.Name); err != nil {
				return diag.FromErr(err)
			}
			return nil
		}
	}

	return diag.Errorf("Error finding the domain")
}

// deleteZoneResource is a function to delete a DNS Zone.
func deleteZoneResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	resp, err := client.DeleteZone(ctx, dns.DeleteZoneRequest{DomainName: d.Id()})
	if err != nil {
		return diag.Errorf("Error deleting server: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error deleting server: %s", resp.Msg)
	}

	return nil
}

// RecordResource returns a schema with the operations for DNS Record resource.
func RecordResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createRecordResource,
		ReadContext:   readRecordResource,
		UpdateContext: updateRecordResource,
		DeleteContext: deleteRecordResource,
		Importer: &schema.ResourceImporter{
			StateContext: importRecordResource,
		},
		Schema: resourceRecordSchema,
	}
}

// createRecordResource is a function to create a new DNS Record.
func createRecordResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	domain := fmt.Sprintf("%v", d.Get("domain"))
	name := utils.ConstructFqdn(fmt.Sprintf("%v", d.Get("name")), domain)

	domainRecord := models.DNSRecord{
		Name:     strings.TrimSuffix(name, "."),
		Domain:   domain,
		Content:  fmt.Sprintf("%v", d.Get("record")),
		Type:     fmt.Sprintf("%v", d.Get("type")),
		TTL:      fmt.Sprintf("%v", d.Get("ttl")),
		Priority: fmt.Sprintf("%v", d.Get("priority")),
	}

	client := dns.New(conf.Client)
	resp, err := client.AddRecord(ctx, dns.AddRecordRequest{
		Domain:   domainRecord.Domain,
		Type:     domainRecord.Type,
		Name:     domainRecord.Name,
		Content:  domainRecord.Content,
		Priority: domainRecord.Priority,
	})

	if err != nil {
		return diag.Errorf("Error creating DNS record: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error creating DNS record: %s", resp.Msg)
	}

	d.SetId(resp.Return.ID)

	log.Printf("[INFO] Domain Record: %s", d.Id())

	return nil
}

// readRecordResource is a function to read a DNS Record.
func readRecordResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)

	domain := fmt.Sprintf("%v", d.Get("domain"))
	record, err := client.GetRecord(ctx, dns.RecordRequest{
		ID:         d.Id(),
		DomainName: domain,
	})

	if err != nil {
		return diag.Errorf("Error retrieving DNS record: %s", err)
	}

	//if record == nil {
	//	log.Printf("[WARN] Record (%s,%s) not found, removing from state", d.Id(), domain)
	//	d.SetId("")
	//	return nil
	//}

	if err := setRecordAttributes(d, record); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// deleteRecordResource is a function to delete a DNS Record.
func deleteRecordResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	resp, err := client.DeleteRecord(ctx, dns.DeleteRecordRequest{
		Domain:   fmt.Sprintf("%v", d.Get("domain")),
		RecordID: d.Id(),
	})
	if err != nil {
		return diag.Errorf("Error deleting DNS record: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error deleting DNS record: %s", resp.Msg)
	}

	return nil
}

// updateRecordResource is a function to update a DNS Record.
func updateRecordResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)
	resp, err := client.UpdateRecord(
		ctx,
		dns.UpdateRecordRequest{
			Domain:   fmt.Sprintf("%v", d.Get("domain")),
			RecordID: d.Id(),
			Type:     fmt.Sprintf("%v", d.Get("type")),
			Name:     fmt.Sprintf("%v", d.Get("name")),
			Content:  fmt.Sprintf("%v", d.Get("content")),
			Priority: fmt.Sprintf("%v", d.Get("priority")),
		},
	)
	if err != nil {
		return diag.Errorf("Error updating DNS record: %s", err)
	}

	if !resp.Status {
		return diag.Errorf("Error updating DNS record: %s", resp.Msg)
	}

	record, err := client.GetRecord(ctx, dns.RecordRequest{
		ID:         d.Id(),
		DomainName: fmt.Sprintf("%v", d.Get("domain")),
	})

	if err != nil {
		return diag.Errorf("Error creating DNS record: %s", err)
	}

	if err := setRecordAttributes(d, record); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// importRecordResource is a function to import a DNS Record.
func importRecordResource(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")

		d.SetId(s[1])
		if err := d.Set("domain", s[0]); err != nil {
			return nil, err
		}
	}

	return []*schema.ResourceData{d}, nil
}

// setRecordAttributes is a function to set the attributes of a DNS Record.
func setRecordAttributes(d *schema.ResourceData, record models.DNSRecord) error {
	d.SetId(record.ID)

	if err := d.Set("domain", record.Domain); err != nil {
		return err
	}

	if err := d.Set("name", record.Name); err != nil {
// 	} else {
		return err
	}

	if err := d.Set("type", record.Type); err != nil {
		return err
	}

	priority, err := strconv.Atoi(record.Priority)
	if err != nil {
		priority = 0
	}

	if err := d.Set("priority", priority); err != nil {
		return err
	}

	if err := d.Set("record", record.Content); err != nil {
		return err
	}

	return d.Set("change_date", record.ChangeDate)
}
