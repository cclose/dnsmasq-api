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
# DM_USER The linux user that runs dnsMasq. This is usually root
DM_USER=${DM_USER:-root}
# ARCH The architecture of this system
ARCH=${ARCH-amd64}
# PLATFORM The platform of this system
PLATFORM=${PLATFORM-linux}

# SKIP_RELOAD - If set, will skip reloading the system daemon and enabling the service. You'll need to do these manually
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

    echo " - installing sudoers file /etc/sudoers.d/dnsmasqapi"
    # Install sudoers file for DMA_USER to reload dnsMasq
    tee /etc/sudoers.d/dnsmasqapi <<EOF
# Allow DMA_USER to reload dnsmasq.service without a password
${DMA_USER} ALL=(ALL) NOPASSWD: /bin/systemctl reload dnsmasq.service
EOF

  fi
fi

# If not running dnsMasq as root, add that user to the group
if [ "${DM_USER}" != "root" ] ; then
  echo "Configuring dnsMasq user ${DM_USER}"
  add_user_to_group "${DM_USER}" "${DMA_GROUP}"
}

# Calculate archive download name
REPO_USER="cclose"
REPO_NAME="dnsmasq-api"
LATEST_RELEASE=$(curl --silent "https://api.github.com/repos/$REPO_USER/$REPO_NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
ARCHIVE_NAME="dnsMasqAPI-${LATEST_RELEASE}-${PLATFORM}-${ARCH}.tar.gz"
ARCHIVE_URL="https://github.com/$REPO_USER/$REPO_NAME/releases/download/$LATEST_RELEASE/${ARCHIVE_NAME}"

# Download and Install archive
echo "Downloading release $LATEST_RELEASE from $ARCHIVE_URL..."
curl -L $ARCHIVE_URL -o "/var/cache/${ARCHIVE_NAME}" || { echo "Failed to download ${ARCHIVE_URL}"; exit 1; }
echo "Extracting to /usr/local"
tar -xzf "/var/cache/${ARCHIVE_NAME}" -C /usr/local || { echo "Failed to extract ${ARCHIVE_NAME}"; exit 1; }

# Set the user and group int he Systemd unit
sed -i "s/nobody/${DMA_USER}/" /usr/local/lib/systemd/system/dnsMasqAPI.service
sed -i "s/nogroup/${DMA_GROUP}/" /usr/local/lib/systemd/system/dnsMasqAPI.service

# If they want a different config, remove the default and update the systemd unit
if [ "${DMA_CONFIG}" != "" ] ; then
  sed -i "s|/usr/local/etc/dnsMasqAPI/config.yaml|${DMA_CONFIG}|" /usr/local/lib/systemd/system/dnsMasqAPI.service
  rm -f /usr/local/etc/dnsMasqAPI/config.yaml
fi

# Touch /etc/dnsmasq.conf.d/api.conf
if [ ! -f "${DMA_DM_CONFIG}" ] ; then
  echo "Seeding config ${DMA_DM_CONFIG}"
  touch "${DMA_DM_CONFIG}"
fi

if [ -z "${SKIP_USER_SETUP}" ] ; then
  if [ "${DMA_USER}" != "root" ] ; then
    echo " - setting permissions"
    chown "${DMA_USER}:${DMA_GROUP}" "${DMA_DM_CONFIG}"
    chmod 550 "$DMA_DM_CONFIG"
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