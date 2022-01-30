FROM alpine

RUN apk add --update openssh
COPY sshd_config /etc/ssh/sshd_config
RUN passwd -d root
RUN mkdir -p /root/.ssh
RUN chmod 700 /root/.ssh
COPY id_rsa.pub /root/.ssh/authorized_keys
RUN chmod 600 /root/.ssh/authorized_keys
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
