# dnsmasq-api
Go API for DNS Masq

## Sudo Permissions

In order to be able to reload DNSMasq service, the user running the webservice needs
permission to call systemctl. If not running the service as root (Please do not run as root!!)
you need to add the follow entries to your sudoers file, assuming user `dnsmasqapi`:

```
dnsmasqapi ALL=(ALL) NOPASSWD: /bin/systemctl reload dnsmasq.service
dnsmasqapi ALL=(ALL) NOPASSWD: /bin/systemctl status dnsmasq.service
```

## Configuring DNSMasq

It is recommended to avoid using the main configuration file (`/etc/dnsmasq.conf`) for the 
DNSMasq settings managed by this API. Instead, use a configuration file in the confdir
(`/etc/dnsmasq.d/`), such as `/etc/dnsmasq.d/api.conf`.

### Permissions for Configuration File

To ensure both DNSMasq and the web service user can access and modify the configuration 
file securely, DNSMasq and the API should use users that belong to a common group and
the configuration file should belong to this group.

This is not a concern if the API runs as root, but you also should not run an API as root,
even in a container. That's a great way to make your infrastructure vulnerable to container
break attacks.

