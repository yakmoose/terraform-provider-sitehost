terraform {
  required_providers {
    sitehost = {
      source  = "sitehostnz/sitehost"
      version = "~> 1.1.1"
    }
  }
}

provider "sitehost" {
  client_id    = var.sitehost_client_id
  api_key      = var.sitehost_api_key
  api_endpoint = "https://api.sitehost.nz/1.1/"
}

resource "sitehost_ssh_key" "sitehost_ssh_key_create" {
  label   = var.ssh_key_label
  content = var.ssh_key_data
}