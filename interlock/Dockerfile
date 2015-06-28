FROM debian:jessie
RUN apt-get update
RUN apt-get install -y --no-install-recommends \
    haproxy \
    nginx-full \
    dnsmasq \
    ca-certificates
ADD interlock /usr/local/bin/interlock
EXPOSE 53/udp 80 443
ENTRYPOINT ["/usr/local/bin/interlock"]
