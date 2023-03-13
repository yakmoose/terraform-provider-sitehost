package dns

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/dns"
	"github.com/sitehostnz/gosh/pkg/models"
	"github.com/sitehostnz/gosh/pkg/utils"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
)

// RecordResource returns a schema with the operations for Server resource.
func RecordResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: dnsRecordCreateResource,
		ReadContext:   dnsRecordReadResource,
		UpdateContext: dnsRecordUpdateResource,
		DeleteContext: dnsRecordDeleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: importResource,
		},
		Schema: dnsRecordResourceSchema,
	}
}

// nameClean is a dirty little helper...
func nameClean(name, domain string) string {
	return strings.TrimSuffix(
		utils.ConstructFqdn(
			name,
			domain,
		),
		".",
	)
}

// dnsRecordCreateResource is a function to create a new Server resource.
func dnsRecordCreateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	domain := d.Get("domain").(string)
	name := nameClean(d.Get("name").(string), domain)
	rr := d.Get("type").(string)
	content := d.Get("content").(string)
	priority := strconv.Itoa(d.Get("priority").(int))

	client := dns.New(conf.Client)
	_, err := client.AddRecord(
		ctx,
		dns.AddRecordRequest{
			Domain:   domain,
			Type:     rr,
			Priority: priority,
			Name:     name,
			Content:  content,
		},
	)
	if err != nil {
		return diag.Errorf("Error creating domain_record: %s", err)
	}

	domainRecord, err := client.GetRecordWithRecord(ctx, models.DNSRecord{Name: name, Domain: domain, Type: rr, Priority: priority, Content: content})
	if domainRecord == nil || err != nil {
		return diag.Errorf("Could not find new record: %s", err)
	}
	// need to read the record back, and granb an id.
	// or we need to start generating an id based on name/type/content
	updateRecordResource(d, domainRecord)

	log.Printf("[INFO] Domain Record: %s", d.Id())

	return nil
}

// dnsRecordReadResource is a function to read a new Server resource.
func dnsRecordReadResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)

	domainRecord, err := client.GetRecord(
		ctx,
		dns.RecordRequest{DomainName: d.Get("domain").(string), ID: d.Id()},
	)

	log.Printf("[INFO] Reading Domain Record: %s", domainRecord)

	if err != nil {
		return diag.Errorf("Error retrieving domain: %s", err)
	}

	if domainRecord == nil {
		return diag.Errorf("Failed retrieving domain_record: %s", err)
	}

	updateRecordResource(d, domainRecord)

	return nil
}

// dnsRecordDeleteResource is a function to delete a new Server resource.
func dnsRecordDeleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)

	domain := d.Get("domain").(string)
	id := d.Id()

	_, err := client.DeleteRecord(ctx, dns.DeleteRecordRequest{Domain: domain, RecordID: id})
	if err != nil {
		return diag.Errorf("Error deleting domain record: %s", err)
	}

	return nil
}

// dnsRecordUpdateResource updates the record.
func dnsRecordUpdateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := dns.New(conf.Client)

	id := d.Id()
	domain := d.Get("domain").(string)
	name := nameClean(d.Get("name").(string), domain)

	_, err := client.UpdateRecord(
		ctx,
		dns.UpdateRecordRequest{
			RecordID: id,
			Domain:   domain,
			Name:     name,
			Type:     d.Get("type").(string),
			Content:  d.Get("content").(string),
			Priority: strconv.Itoa(d.Get("priority").(int)),
		},
	)
	if err != nil {
		return diag.Errorf("Error updating domain record: %s", err)
	}

	// do we need to push back the updates
	// updateRecordResource(d, domainRecord)

	return nil
}

func importResource(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), "/") {
		s := strings.Split(d.Id(), "/")
		d.SetId(s[1])
		d.Set("domain", s[0])
	}

	return []*schema.ResourceData{d}, nil
}
