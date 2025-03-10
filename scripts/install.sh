#!/bin/bash

set -e

# Set default values for environment variables
# DMA_DM_CONFIG The location of the dnsMasq config file the API will manage
DMA_DM_CONFIG="${DMA_DM_CONFIG-/etc/dnsmasq.d/api.conf}"
# DMA_CONFIG The location of the config file for the dnsMasqAPI service. Defaults to /usr/local/etc/dnsMasqAPI/config.yaml
DMA_CONFIG=""
# DMA_GROUP The Linux group that the user running the API service belongs to. Default dnsmasqapi
DMA_GROUP=${DMA_GROUP:-dnsmasqapi}
# DMA_USER The linux user that runs the API service. Default dnsmasqapi
DMA_USER=${DMA_USER:-dnsmasqapi}
# DM_USER The linux user that runs dnsMasq. This is usually dnsmasq
DM_USER=${DM_USER:-dnsmasq}
# ARCH The architecture of this system
ARCH=${ARCH-amd64}
# PLATFORM The platform of this system
PLATFORM=${PLATFORM-linux}

# LOG_PATH - If set, will change the path to the log in the log rotate script. If using the standard config, will also update this.
# LOG_TO_JOURNAL - If set, will set logging to stdout which will goto journals and be accessible via journalctl -xeu dnsMasqAPI.service. Does nothing if using your own config

# SKIP_DB_CREATE - If set, will skip creating /var/lib/dnsMasqAPI/dns.db. Set this if you want to use a different location for the bolt db file
# SKIP_LOG_ROTATE - If set, will skip linking /usr/local/etc/logrotated.conf to /etc/logrotate.d/dnsMasqAPI
# SKIP_RELOAD - If set, will skip reloading the system daemon and enabling the service. You'll need to do these manually
# SKIP_SUDOERS - If set, will skip adding a sudoers file to let the DMA_USER reload dnsmasq. This is auto-skipped if DMA_USER == root (Please don't do this though...)
# SKIP_USER_SETUP - If set, will skip adding user/group and sudoers entry. You'll need to do this manually


# check_or_add_group - check or add a group
check_or_add_group() {
  local group=$1
  if ! getent group "${group}" > /dev/null 2>&1; then
    echo "Group ${group} does not exist. Creating..."
    groupadd "${group}"
  else
    echo "Group ${group} already exists."
  fi
}

# check_or_add_user - check or add a user
check_or_add_user() {
  local user=$1
  local group=$2
  if ! id -u "${user}" > /dev/null 2>&1; then
    echo "User ${user} does not exist. Creating..."
    useradd -g "${group}" -s /sbin/nologin "${user}"
  else
    echo "User ${user} already exists."
  fi
}

# add_user_to_group - add a user to a group
add_user_to_group() {
  local user=$1
  local group=$2
  if id -nG "${user}" | grep -qw "${group}"; then
    echo "User ${user} is already in group ${group}."
  else
    echo "Adding user ${user} to group ${group}."
    usermod -aG "${group}" "${user}"
  fi
}

# if we aren't skipping user setup
if [ -z "${SKIP_USER_SETUP}" ]; then
  echo "Configuring dnsMasqAPI with user ${DMA_USER} and group ${DMA_GROUP}"
  if [ "${DMA_USER}" != "root" ]; then
    # Check or add the group and users for the API
    check_or_add_group "${DMA_GROUP}"
    check_or_add_user "${DMA_USER}" "${DMA_GROUP}"
    add_user_to_group "${DMA_USER}" "syslog"

    if [ -z "${SKIP_SUDOERS}" ] ; then
      echo " - installing sudoers file /etc/sudoers.d/dnsmasqapi"
      # Install sudoers file for DMA_USER to reload dnsMasq
      tee /etc/sudoers.d/dnsmasqapi  > /dev/null <<EOF
# Allow DMA_USER to restart dnsmasq.service without a password
${DMA_USER} ALL=(ALL) NOPASSWD: /bin/systemctl restart dnsmasq.service
${DMA_USER} ALL=(ALL) NOPASSWD: /bin/systemctl status dnsmasq.service
EOF
    fi
  fi
fi

# If not running dnsMasq as root, add that user to the group
if [ "${DM_USER}" != "root" ] ; then
  echo "Configuring dnsMasq user ${DM_USER}"
  add_user_to_group "${DM_USER}" "${DMA_GROUP}"
fi

# Calculate archive download name
REPO_USER="cclose"
REPO_NAME="dnsmasq-api"
LATEST_RELEASE=$(curl --silent "https://api.github.com/repos/$REPO_USER/$REPO_NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
ARCHIVE_NAME="dnsMasqAPI-${LATEST_RELEASE}-${PLATFORM}-${ARCH}.tar.gz"
ARCHIVE_URL="https://github.com/$REPO_USER/$REPO_NAME/releases/download/$LATEST_RELEASE/${ARCHIVE_NAME}"

# Download and Install archive
echo "Downloading release $LATEST_RELEASE from $ARCHIVE_URL..."
curl -L --silent $ARCHIVE_URL -o "/var/cache/${ARCHIVE_NAME}" || { echo "Failed to download ${ARCHIVE_URL}"; exit 1; }
echo "Extracting to /usr/local"
tar -xzf "/var/cache/${ARCHIVE_NAME}" -C /usr/local || { echo "Failed to extract ${ARCHIVE_NAME}"; exit 1; }

# Set the user and group in the Systemd unit
sed -i "s/nobody/${DMA_USER}/" /usr/local/lib/systemd/system/dnsMasqAPI.service
sed -i "s/nogroup/${DMA_GROUP}/" /usr/local/lib/systemd/system/dnsMasqAPI.service

# If they want a different config, remove the default and update the systemd unit
if [ -n "${DMA_CONFIG}" ] ; then
  echo "Changing Systemd Unit Config to ${DMA_CONFIG}"
  sed -i "s|/usr/local/etc/dnsMasqAPI/config.yaml|${DMA_CONFIG}|" /usr/local/lib/systemd/system/dnsMasqAPI.service
  rm -f /usr/local/etc/dnsMasqAPI/config.yaml
# If using standard config but custom log path
else
  if [ -n "${LOG_TO_JOURNAL}" ] ; then
    echo "Setting Logging to Journald (STDOUT)"
    sed -i "/log:/d" /usr/local/etc/dnsMasqAPI/config.yaml
    sed -i "/dnsMasqAPI.log/d" /usr/local/etc/dnsMasqAPI/config.yaml
  elif [ -n "${LOG_PATH}" ] ; then
    echo "Setting log path to ${LOG_PATH}"
    sed -i "s|/var/log/dnsMasqAPI.log|${LOG_PATH}|" /usr/local/etc/dnsMasqAPI/config.yaml
  fi
fi

# Set user and group in the logrotate config (even if we're skipping it)
sed -i "s/nobody/${DMA_USER}/" /usr/local/etc/dnsMasqAPI/logrotated.conf
sed -i "s/nogroup/${DMA_GROUP}/" /usr/local/etc/dnsMasqAPI/logrotated.conf
if [ -z "${SKIP_LOG_ROTATE}" ] && [ -z "${LOG_TO_JOURNAL}" ] ; then
  echo "Deploying LogRotate"
  # Adjust the path if using a custom path
  if [ ! -z "${LOG_PATH}" ] ; then
    sed -i "s|/var/log/dnsMasqAPI.log|${LOG_PATH}|" /usr/local/etc/dnsMasqAPI/logrotated.conf
  fi
  # link the logrotate config
  ln -s /usr/local/etc/dnsMasqAPI/logrotated.conf /etc/logrotate.d/dnsMasqAPI
fi

# Touch /etc/dnsmasq.conf.d/api.conf
if [ ! -f "${DMA_DM_CONFIG}" ] ; then
  echo "Seeding config ${DMA_DM_CONFIG}"
  touch "${DMA_DM_CONFIG}"
fi

# Create the /var/lib directory for default db file and pid file
echo "Creating /var/lib/dnsMasqAPI"
mkdir -p /var/lib/dnsMasqAPI
if [ -z "${SKIP_USER_SETUP}" ] ; then
  if [ "${DMA_USER}" != "root" ] ; then
    echo "Chowning /var/lib/dnsMasqAPI to ${DMA_USER}:${DMA_GROUP}"
    chown -R "${DMA_USER}:${DMA_GROUP}" /var/lib/dnsMasqAPI
    chmod -R 770 /var/lib/dnsMasqAPI
  fi
fi

if [ -z "${SKIP_USER_SETUP}" ] ; then
  if [ "${DMA_USER}" != "root" ] ; then
    echo " - setting permissions on dnsMasq config"
    chown "${DMA_USER}:${DMA_GROUP}" "${DMA_DM_CONFIG}"
    chmod 770 "$DMA_DM_CONFIG"
  fi
fi


# Enable and Start Service
if [ -z "${SKIP_SYSTEM}" ]; then
  echo "Configuring systemd"
  systemctl daemon-reload
  systemctl enable dnsMasqAPI.service
  systemctl start dnsMasqAPI.service
fi

echo "Installation complete!"