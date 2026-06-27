#!/bin/sh
set -eu

mkdir -p /app/.ssh /sandbox-keys

if [ ! -f /sandbox-keys/sandbox_id_rsa ]; then
  ssh-keygen -q -t rsa -b 4096 -N '' -f /sandbox-keys/sandbox_id_rsa
fi

if [ ! -f /app/.ssh/.sshman ] || [ ! -s /app/.ssh/.sshman ]; then
  cp /sandbox/sshman.json /app/.ssh/.sshman
fi

exec /usr/local/bin/sshman web --bind 0.0.0.0 --allow-remote --port 8080
