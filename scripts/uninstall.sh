#!/bin/bash

set -e

# Set default values for environment variables
DMA_DM_CONFIG="${DMA_DM_CONFIG-/etc/dnsmasq.d/api.conf}"
DMA_CONFIG="${DMA_CONFIG-/usr/local/etc/dnsMasqAPI/config.yaml}"
DMA_GROUP="${DMA_GROUP:-dnsmasqapi}"
DMA_USER="${DMA_USER:-dnsmasqapi}"
DM_USER="${DM_USER:-root}"
ARCH="${ARCH-amd64}"
PLATFORM="${PLATFORM-linux}"

# SKIP_USER_REMOVE - Don't remove the DM_USER and DM_GROUP form the linux system
# SKIP_DAEMON_RELOAD - Don't reload the system daemon
# PRESERVE_DM_CONFIG - Don't remove the dnsMasq config that was managed by the service
# PRESERVE_LOGS - Don't remove the log files
# PRESERVE_SELF - Don't self-delete this script

if [ -z "${PRESERVE_LOGS}" ]; then
  echo "TODO: Removing Logs"
  # TODO detect and remove logs
fi

# Remove logrotate configuration
if [ -L /etc/logrotate.d/dnsMasqAPI ]; then
  echo "Removing logrotate configuration..."
  rm /etc/logrotate.d/dnsMasqAPI
fi

# Remove systemd service and unit file
echo "Stopping and disabling systemd service..."
systemctl stop dnsMasqAPI.service || true
systemctl disable dnsMasqAPI.service || true
rm -f /usr/local/lib/systemd/system/dnsMasqAPI.service
if [ -z "${SKIP_DAEMON_RELOAD}" ]; then
  systemctl daemon-reload
fi

# Remove installation files and directories
echo "Removing installed files and directories..."
rm -f /usr/local/bin/dnsMasqAPI
rm -rf /usr/local/etc/dnsMasqAPI
rm -rf /var/lib/dnsMasqAPI

# Remove config files if they exist
if [ -z "${PRESERVE_DM_CONFIG}" ]; then
  if [ -f "${DMA_DM_CONFIG}" ]; then
    echo "Removing dnsMasq config file ${DMA_DM_CONFIG}..."
    rm -f "${DMA_DM_CONFIG}"
  fi
fi

if [ -f "${DMA_CONFIG}" ]; then
  if [ -f "${DMA_CONFIG}" ]; then
    echo "Removing config file ${DMA_CONFIG}..."
    rm -f "${DMA_CONFIG}"
  fi
fi

# Optionally remove user and group
if [ -z "${SKIP_USER_REMOVE}" ]; then
  if [ "${DMA_USER}" != "root" ]; then
    echo "Removing user ${DMA_USER}..."
    userdel "${DMA_USER}" || true
  fi

  if [ -n "${DMA_GROUP}" ]; then
    echo "Removing group ${DMA_GROUP}..."
    groupdel "${DMA_GROUP}" || true
  fi
fi

# Schedule the script to delete itself after a short delay
if [ -z "${PRESERVE_SELF}" ] ; then
  SELF=$(realpath "$0")
  {
      echo "sleep 1; rm -f '$SELF'"
  } | at now
  echo "Scheduled self-deletion of the script."
fi
echo "Uninstallation complete!"
