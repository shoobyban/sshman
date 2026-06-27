#!/bin/sh
set -eu

mkdir -p /root/.ssh
chmod 700 /root/.ssh
touch /root/.ssh/authorized_keys
chmod 600 /root/.ssh/authorized_keys

KEY_FILE="${SANDBOX_PUBKEY_PATH:-/sandbox-keys/sandbox_id_rsa.pub}"
attempt=0
while [ ! -f "$KEY_FILE" ] && [ "$attempt" -lt 60 ]; do
  attempt=$((attempt + 1))
  sleep 1
done

if [ -f "$KEY_FILE" ]; then
  cp "$KEY_FILE" /root/.ssh/authorized_keys
  chmod 600 /root/.ssh/authorized_keys
fi

if [ ! -f /etc/ssh/ssh_host_rsa_key ]; then
  ssh-keygen -f /etc/ssh/ssh_host_rsa_key -N '' -t rsa
fi

if [ ! -f /etc/ssh/ssh_host_ed25519_key ]; then
  ssh-keygen -f /etc/ssh/ssh_host_ed25519_key -N '' -t ed25519
fi

exec /usr/sbin/sshd -D
