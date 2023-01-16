// Package server provides the functions to create a Server resource via SiteHost API.
package domain_record

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/domain_record"
	"github.com/sitehostnz/gosh/pkg/models"
	"github.com/sitehostnz/terraform-provider-sitehost/sitehost/helper"
	"log"
	"strconv"
	"strings"
)

// Resource returns a schema with the operations for Server resource.
func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResource,
		ReadContext:   readResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: importResource,
		},
		Schema: resourceSchema,
	}
}

// createResource is a function to create a new Server resource.
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	domain := d.Get("domain").(string)
	name := d.Get("name").(string)

	client := domain_record.New(conf.Client)
	domainRecord, err := client.Create(
		ctx,
		&models.DomainRecord{
			Domain:   domain,
			Type:     d.Get("type").(string),
			Priority: strconv.Itoa(d.Get("priority").(int)),
			Name:     constructFqdn(name, domain),
			Content:  d.Get("record").(string),
		},
	)

	if err != nil {
		return diag.Errorf("Error creating domain_record: %s", err)
	}
	if domainRecord == nil {
		return diag.Errorf("Failed creating domain_record: %s", domainRecord)
	}

	updateRecordResource(d, domainRecord)

	log.Printf("[INFO] Domain Record: %s", d.Id())

	return nil
}

// readResource is a function to read a new Server resource.
func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := domain_record.New(conf.Client)

	domainRecord, err := client.Get(
		ctx,
		domain_record.RecordRequest{DomainName: d.Get("domain").(string), Id: d.Id()},
	)

	if err != nil {
		return diag.Errorf("Error retrieving domain: %s", err)
	}

	if domainRecord == nil {
		return diag.Errorf("Failed retrieving domain_record: %s", err)
	}

	updateRecordResource(d, domainRecord)

	return nil
}

// deleteResource is a function to delete a new Server resource.
func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := domain_record.New(conf.Client)

	_, err := client.Delete(ctx, &models.DomainRecord{Domain: d.Get("domain").(string), Id: d.Id()})

	if err != nil {
		return diag.Errorf("Error deleting domain record: %s", err)
	}

	return nil
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf, ok := meta.(*helper.CombinedConfig)
	if !ok {
		return diag.Errorf("failed to convert meta object")
	}

	client := domain_record.New(conf.Client)

	domainRecord, err := client.Update(
		ctx,
		&models.DomainRecord{
			Id:       d.Id(),
			Domain:   d.Get("domain").(string),
			Name:     d.Get("name").(string),
			Type:     d.Get("type").(string),
			Content:  d.Get("content").(string),
			Priority: strconv.Itoa(d.Get("priority").(int)),
		},
	)

	if err != nil {
		return diag.Errorf("Error updating domain record: %s", err)
	}

	if domainRecord == nil {
		return diag.Errorf("Failed updating domain_record: %s", err)
	}

	updateRecordResource(d, domainRecord)

	return nil
}

func importResource(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")

		d.SetId(s[1])
		d.Set("domain", s[0])
	}

	return []*schema.ResourceData{d}, nil
}
