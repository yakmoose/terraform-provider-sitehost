package domain_record

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sitehostnz/gosh/pkg/models"
	"strconv"
	"strings"
)

func constructFqdn(name, domain string) string {
	if name == "@" {
		return domain
	}

	rn := strings.ToLower(name)
	domainSuffix := domain + "."
	if strings.HasSuffix(rn, domainSuffix) {
		rn = strings.TrimSuffix(rn, ".")
	} else {
		rn = strings.Join([]string{name, domain}, ".")
	}
	return rn
}

func deconstructFqdn(name, domain string) string {
	if name == domain {
		return "@"
	}

	rn := strings.ToLower(name)
	rn = strings.TrimSuffix(rn, ".")
	rn = strings.TrimSuffix(rn, "."+domain)

	return rn
}

func updateRecordResource(d *schema.ResourceData, domainRecord *models.DomainRecord) {
	d.SetId(domainRecord.Id)
	d.Set("domain", domainRecord.Domain)
	d.Set("name", deconstructFqdn(domainRecord.Name, domainRecord.Domain))
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

	d.Set("fqdn", constructFqdn(domainRecord.Name, domainRecord.Domain))
	d.Set("change_date", domainRecord.ChangeDate)
}
