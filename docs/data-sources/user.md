---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "loriot_user Data Source - loriot"
subcategory: ""
description: |-
  User data source for the current user
---

# loriot_user (Data Source)

User data source for the current user



<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `alerts` (Boolean) Notifcations alerts configuration
- `devices_limit` (Number) Devices limit for this user
- `email` (String) User email address
- `first_name` (String) First name or Forename
- `gateways_limit` (Number) Gateways limit for this user
- `has_credit_card` (Boolean) User has a credit card in their account
- `id` (Number) Internal user ID in Loriot Network Server
- `last_name` (String) Last name or surname
- `level` (Number) Level of the user for admin rights (1 to 100)
- `mcast_devices_limit` (Number) Multicast devices limit by user
- `organization_role` (String) Organization role of the user
- `organization_uuid` (String) Unique identifier of the organization
- `output_limit` (Number) Maximum number of outputs allowed
- `tier` (Number) Tier of the user
