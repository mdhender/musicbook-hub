#!/bin/bash
############################################################################
# fail on any error
set -e
############################################################################
#
USER=mbooks
HOST=damned.dev
REMOTE_WEB_DIR=/var/www/damned.dev
REMOTE_BIN_DIR=/var/www/damned.dev/bin
SERVICE_NAME=mbooks
############################################################################
#
[ -d build -a -f "bin/deploy.sh" ] || {
  echo error: must run from the root of the repository
  exit 2
}
############################################################################
#
echo "🛠️  removing old builds..."
rm -rf build/*
############################################################################
#
echo "🛠️  Building frontend and backend..."
make build
############################################################################
#
echo "🔁 Stopping systemd service..."
ssh "${HOST}" "systemctl stop   ${SERVICE_NAME}"
sleep 1
ssh "${HOST}" "systemctl status ${SERVICE_NAME} --no-pager ; echo"
############################################################################
#
echo "🚀 Copying backend binary to ${HOST}..."
scp build/musicbook-hub "${HOST}:${REMOTE_BIN_DIR}/"
############################################################################
#
echo "📦 Syncing React dist/ folder to ${HOST}..."
rsync -av --delete web/dist/ "${HOST}:${REMOTE_WEB_DIR}/web/dist/"
############################################################################
#
echo "🔁 Restarting systemd service..."
ssh "${HOST}" "systemctl restart ${SERVICE_NAME}"
sleep 1
ssh "${HOST}" "systemctl status  ${SERVICE_NAME} --no-pager ; echo"
############################################################################
#
echo "✅ Deploy complete!"
############################################################################
#
exit 0
