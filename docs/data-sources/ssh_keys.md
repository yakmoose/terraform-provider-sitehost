---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sitehost_ssh_keys Data Source - terraform-provider-sitehost"
subcategory: ""
description: |-
  
---

# sitehost_ssh_keys (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `ssh_keys` (Set of Object) The list of ssh keys (see [below for nested schema](#nestedatt--ssh_keys))

<a id="nestedatt--ssh_keys"></a>
### Nested Schema for `ssh_keys`

Read-Only:

- `content` (String)
- `custom_image_access` (Boolean)
- `date_added` (String)
- `date_updated` (String)
- `id` (String)
- `label` (String)


