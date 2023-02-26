package helper

import (
	"strings"
)

func ConstructFqdn(name, domain string) string {
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

func DeconstructFqdn(name, domain string) string {
	if name == domain {
		return "@"
	}

	rn := strings.ToLower(name)
	rn = strings.TrimSuffix(rn, ".")
	rn = strings.TrimSuffix(rn, "."+domain)

	return rn
}
