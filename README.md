# dnsmasq-api

`dnsmasq-api` is a RESTful API service for managing DNS entries, integrating seamlessly with `dnsmasq` for dynamic DNS management.

## Features

- Manage DNS records via RESTful API
- Retrieve service status and metrics
- Configurable logging
- Systemd service setup

## Prerequisites

- `dnsmasq` installed and configured
- `Go` installed for building the project

## Installation

You can install `dnsmasq-api` using the provided install script:

```
git clone https://github.com/cclose/dnsmasq-api.git
cd dnsmasq-api
sudo bash ./scripts/install.sh
```

or curl just the installer with:

```shell
LATEST_RELEASE=$(curl --silent https://api.github.com/repos/cclose/dnsmasq-api/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
curl -L --silent "https://github.com/cclose/dnsmasq-api/releases/download/${LATEST_RELEASE}/install.sh" -o install.sh
sudo bash install.sh
```

I promise the installer isn't up to anything nefarious... but it does do quite a bit and needs sudo permission, so you 
should probably look through it and make sure you understand what it's doing. The default options should only touch
files and users that don't already exist and are only needed for dnsMasqAPI, but still.

To uninstall, the installer will unpack an uninstaller at `/usr/local/bin/dnsMasqAPI-uninstall.sh`.
Set the environment variable `PRESERVE_SELF` to keep the uninstaller from deleting itself.

### Environment Variables

The install script supports several environment variables:

- `DMA_DM_CONFIG`: Location of the `dnsmasq` config file managed by the API (default: `/etc/dnsmasq.d/api.conf`)
- `DMA_CONFIG`: Location of the `dnsmasq-api` config file (default: `/usr/local/etc/dnsMasqAPI/config.yaml`)
- `DMA_GROUP`: Linux group for the API service user (default: `dnsmasqapi`)
- `DMA_USER`: Linux user for the API service (default: `dnsmasqapi`)
- `DM_USER`: Linux user running `dnsmasq` (default: `dnsmasq`)
- `ARCH`: System architecture (default: `amd64`)
- `PLATFORM`: System platform (default: `linux`)
- `LOG_PATH`: Path for log files (optional)
- `LOG_TO_JOURNAL`: Set to log to stdout (optional)

## Usage

After installation, manage the service using `systemctl`:

```
sudo systemctl start dnsMasqAPI.service
sudo systemctl stop dnsMasqAPI.service
sudo systemctl status dnsMasqAPI.service
```

### API Endpoints

- **DNS Management**
    - `GET /dns`: Retrieve all DNS records
    - `GET /dns/:hostname`: Retrieve a specific DNS record by hostname
    - `POST /dns/:hostname`: Add or update a DNS record
    - `DELETE /dns/:hostname`: Delete a DNS record

- **Service Status and Metrics**
    - `GET /statusz`: Get service status
    - `GET /metricz`: Get service metrics

### Configuration

The configuration file is located at `/usr/local/etc/dnsMasqAPI/config.yaml` by default. Customize this path using the `DMA_CONFIG` environment variable during installation.

### Logging

Specify logging configuration in the `config.yaml` file. Log to a file, stdout, or stderr based on your setup.

## Development

### Building

To build the project, run:

```
make build
```

### Testing

Run tests with:

```
make test
```

### Linting

Lint the code with:

```
make lint
```

### Docker

Build and run the Docker container:

```
make docker
make run
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

# Configuration Notes

Note, the bundled installer `scripts/install.sh` handles all of the below, but I wanted to call it out so you know.

## Sudo Permissions

In order to be able to reload DNSMasq service, the user running the webservice needs
permission to call systemctl. If not running the service as root (Please do not run as root!!)
you need to add the follow entries to your sudoers file, assuming user `dnsmasqapi`:

```
dnsmasqapi ALL=(ALL) NOPASSWD: /bin/systemctl start dnsmasq.service
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

