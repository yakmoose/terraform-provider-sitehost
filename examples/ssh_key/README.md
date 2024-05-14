# SiteHost SSH Key

This example adds an SSH Key to your account, for use when creating new servers.

> Note: To run this example, first configure your SiteHost provider as described in <https://kb.sitehost.nz/developers/terraform-provider>.

## Prerequisites

Your `Client ID` and `API Key` values can be generated<https://cp.sitehost.nz/api/list-keys> by clicking "Add API Key". You can find out how to do this in our [knowledge base](https://kb.sitehost.nz/developers/api).

You will need to export your SiteHost Client ID and API Key as an environment variable:

```sh
export TF_VAR_sitehost_client_id="Put your SiteHost Client ID here"
export TF_VAR_sitehost_api_key="Put your SiteHost API key here"
export TF_VAR_ssh_key_label="Put the SSH Key label here"
export TF_VAR_ssh_key_data="Put the contents of your SSH Public Key here"
```

## Run this example

From the `examples/ssh_key` directory.

```sh
terraform init
terraform apply
```

> The SSH Key should be added 

## Destroy the Resources

Clean up by removing all the resources that were created in one command:

```sh
terraform destroy
```
