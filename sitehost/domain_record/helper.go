package domain_record

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/api/domain_record"
	"github.com/sitehostnz/gosh/pkg/models"
	"log"
	"strconv"
)

func updateRecordResource(d *schema.ResourceData, domainRecord *models.DomainRecord) {
	log.Printf("[INFO] updating resource Record: %s", domainRecord)

	d.SetId(domainRecord.ID)
	d.Set("domain", domainRecord.Domain)
	d.Set("name", domain_record.DeconstructFqdn(domainRecord.Name, domainRecord.Domain))
	d.Set("record", domainRecord.Content)
	d.Set("type", domainRecord.Type)

	ttl, e := strconv.Atoi(domainRecord.TTL)
	if e == nil {
		d.Set("ttl", ttl)
	} else {
		d.Set("ttl", 3600)
	}

	priority, e := strconv.Atoi(domainRecord.Priority)
	if e == nil {
		d.Set("priority", priority)
	} else {
		d.Set("priority", 0)
	}

	d.Set("fqdn", domain_record.ConstructFqdn(domainRecord.Name, domainRecord.Domain))
	d.Set("change_date", domainRecord.ChangeDate)
}
