<a href="https://terraform.io">
    <img src=".github/terraform_logo.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# Terraform Provider for SiteHost

### Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.22.2 (to build the provider plugin)

### Example Usage

```
terraform {
    required_providers {
        sitehost = {
            source = "sitehostnz/sitehost"
            version = "~> 1.0"
        }
    }
}

provider "sitehost" {
    client_id = "****"
    api_key = "***"
    api_endpoint = "https://api.sitehost.nz/1.1/"
}

resource "sitehost_server" "web" {
    label = "webserver"
    location = "AKLCITY"
    product_code = "XENLIT"
    image = "ubuntu-jammy-pvh.amd64"
}
```

You can [read more about the steps to create a SiteHost API KEY in our knowledge base](http://kb.sitehost.nz/developers/api).

### Locations - Product Codes - Images

Check our KB to see the complete list of [locations](http://kb.sitehost.nz/developers/api/locations), [product codes](http://kb.sitehost.nz/developers/api/product-codes) and [images](http://kb.sitehost.nz/developers/api/images).


### Building The Provider (For developers)

Clone repository to: `$GOPATH/src/github.com/sitehostnz/
terraform-provider-sitehost`

```bash
$ git clone git@github.com:sitehostnz/terraform-provider-sitehost
```

Enter the provider directory and build the provider

```bash
$ make install
```

## Contributing
If you're interested in contributing to our project:
- Start by reading our [style guide](https://github.com/sitehostnz/go-style-guide/blob/master/style.md).
- Explore our [issues](https://github.com/sitehostnz/terraform-provider-sitehost/issues).
- Or send us feature PRs.

## License
SiteHost Terraform project is distributed under [MIT](./LICENSE.md).
