---
page_title: "sitehost_ssh_key Resource - terraform-provider-sitehost"
subcategory: ""
description: Provides a SiteHost SSH Key resource. This can be used to create, modify, and delete SSH Keys.

---

# sitehost_ssh_key (Resource)

Provides a SiteHost SSH Key resource. This can be used to create, modify, and delete SSH Keys.

## Example Usage
```hcl
# Create an SSH Key
resource "sitehost_ssh_key" "key" {
    label = "My New SSH Key"
    content = "ssh-rsa AAAAB3..."
}
```

## Schema

### Required

- `label` (String) The SSH Key label.
- `content` (String, Sensitive) The string content of your SSH Public Key.

### Optional

- `custom_image_access` (String) Whether or not the SSH Key will have access to custom images (1: True, 0: False).

### Read-Only

- `id` (String) The ID of the SSH Key.
- `date_added` (String) The timestamp for when the SSH Key was created.
- `date_updated` (String) The timestamp for when the SSH key was last updated. Defaults to the `date_added` value.
