package dns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/models"
	"github.com/sitehostnz/gosh/pkg/utils"
	"log"
	"strconv"
)

func updateRecordResource(d *schema.ResourceData, domainRecord *models.DNSRecord) {
	log.Printf("[INFO] updating resource Record: %s", domainRecord)

	d.SetId(domainRecord.ID)
	d.Set("domain", domainRecord.Domain)
	d.Set("name", utils.DeconstructFqdn(domainRecord.Name, domainRecord.Domain))
	d.Set("content", domainRecord.Content)
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
	}

	// sh api returns the name as the fdqn without the trailing dot
	d.Set("fqdn", domainRecord.Name+".")
	d.Set("change_date", domainRecord.ChangeDate)
}
