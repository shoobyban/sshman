#!/bin/sh
if [ ! -f /etc/ssh/ssh_host_rsa_key  ]; then
    ssh-keygen -f /etc/ssh/ssh_host_rsa_key -N '' -t rsa
    ssh-keygen -f /etc/ssh/ssh_host_dsa_key -N '' -t dsa
fi
/usr/sbin/sshd -D
