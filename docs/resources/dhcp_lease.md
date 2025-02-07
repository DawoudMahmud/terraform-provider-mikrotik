# mikrotik_dhcp_lease (Resource)
Creates a DHCP lease on the mikrotik device.

## Example Usage
```terraform
resource "mikrotik_dhcp_lease" "file_server" {
  address    = "192.168.88.1"
  macaddress = "11:22:33:44:55:66"
  comment    = "file server"
  blocked    = "false"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `address` (String) The IP address of the DHCP lease to be created.
- `macaddress` (String) The MAC addreess of the DHCP lease to be created.

### Optional

- `blocked` (String) Whether to block access for this DHCP client (true|false). Default: `false`.
- `comment` (String) The comment of the DHCP lease to be created.
- `dynamic` (Boolean) Whether the dhcp lease is static or dynamic. Dynamic leases are not guaranteed to continue to be assigned to that specific device. Defaults to false. Default: `false`.
- `hostname` (String) The hostname of the device

### Read-Only

- `id` (String) The ID of this resource.

## Import
Import is supported using the following syntax:
```shell
# The resource ID (*19) is a MikroTik's internal id.
# It can be obtained via CLI:
# [admin@MikroTik] /ip dhcp-server lease> :put [find where address=10.0.1.254]
# *19
terraform import mikrotik_dhcp_lease.file_server '*19'
```
