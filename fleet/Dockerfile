FROM progrium/busybox
RUN opkg-install bash
RUN mkdir -p /etc/ssl && mkdir -p /etc/ssl/certs
ADD certs /etc/ssl/certs/
ADD config.json config.json
ADD provider/cloud-config-agent.yaml cloud-config-agent.yaml
ADD provider/cloud-config-init.yaml cloud-config-init.yaml
ADD fleet fleet
CMD ["./fleet"]