FROM alpine:3.22

RUN apk add --no-cache openssh
COPY tests/docker/sshd_config /etc/ssh/sshd_config
COPY sandbox/host-entrypoint.sh /entrypoint.sh

RUN passwd -d root && mkdir -p /root/.ssh && chmod 700 /root/.ssh && chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
