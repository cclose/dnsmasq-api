#!/bin/bash

set -e

source dist/envvar.sh
if [ -z "${ARCHIVE_NAME}" ] || [ -z "${ARTIFACT_DIR}" ] || [ -z "${ARCHIVE_PATH}" ] || [ -z "${GITHUB_REF}" ] || [ -z "${PLATFORM}" ] || [ -z "${ARCH}" ] ; then
  echo "Missing required variable from dist/envvar. Verify environment is setup correctly!"
  exit 1
fi


# Create directories
echo "Creating dist directories"
mkdir -p dist/bundle/bin
mkdir -p dist/bundle/lib/systemd/system
mkdir -p dist/bundle/etc/dnsMasqAPI
mkdir -p "$ARTIFACT_DIR"

# Move the binary into the dist
echo "Bundling binary"
cp dnsMasqAPI dist/bundle/bin/dnsMasqAPI
chmod +x dist/bundle/bin/dnsMasqAPI

# Move the systemd service into the dist
echo "Bundling Service Unit"
cp systemd/dnsMasqAPI.service dist/bundle/lib/systemd/system/dnsMasqAPI.service

# Move the default config into the dist
echo "Bundling config.yaml"
# This will effectively cat the config file while filtering out the skip line
cp deploy/config.yaml.default dist/bundle/etc/dnsMasqAPI/config.yaml
cp deploy/logrotated.conf dist/bundle/etc/dnsMasqAPI/logrotated.conf

# Create a tar.gz archive
echo "Archiving release as ${ARCHIVE_NAME}.tar.gz"
tar -czf ${ARCHIVE_PATH} -C dist/bundle .

# Add install script
echo "Adding install.sh"
cp scripts/install.sh dist/artifacts/.

# Checksum
echo "Check-summing artifacts"
SUMFILE="${ARTIFACT_DIR}/sha256sum-${PLATFORM}-${ARCH}.txt"
echo -n "" > "${SUMFILE}"
for FILE in $(ls -1 ${ARTIFACT_DIR} | grep -vE '\.sum|\.txt$'); do
  echo " - summing ${FILE}"
  sha256sum "${ARTIFACT_DIR}/${FILE}" | sed "s|$(dirname ${ARTIFACT_DIR}/${FILE})/||g" >> "${SUMFILE}"
done