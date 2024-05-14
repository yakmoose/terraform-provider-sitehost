---
page_title: "SiteHost Provider"
subcategory: ""
description:

---

# SiteHost Provider

The SiteHost (SH) provider is used to interact with the
resources supported by SiteHost. The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.


## Example Usage
```hcl
terraform {
  required_providers {
    sitehost = {
      source = "sitehostnz/sitehost"
      version = "~> 1.0"
    }
  }
}

# Configure the SiteHost Provider
provider "sitehost" {
	client_id = ""
	api_key = ""
	api_endpoint = "https://api.sitehost.nz/1.1/"
}

# Create an SSH Key
resource "sitehost_ssh_key" "key" {
    label = "My New SSH Key"
    content = "ssh-rsa AAAAB3..."
}

# Create a web server with the SSH Key we just created
resource "sitehost_server" "web" {
	label = "webserver"
	location = "AKLCITY"
	product_code = "XENLIT"
	image = "ubuntu-xenial.amd64"
	ssh_keys = [sitehost_ssh_key.key.content]
}
```

## Schema

### Required

- `api_key` (String, Sensitive) api authentication key
- `client_id` (String) client identifier

### Optional

- `api_endpoint` (String) url prefix of the api server
