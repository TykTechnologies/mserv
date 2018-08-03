FROM debian:jessie-slim

ENV TYKVERSION 0.1
ENV TYK_MSERV_CONFIG /etc/mserv/mserv.json

LABEL Description="Tyk MServ sidecar service docker image" Vendor="Tyk" Version=$TYKVERSION

RUN apt-get update \
 && apt-get upgrade -y \
 && apt-get install -y --no-install-recommends \
            curl ca-certificates apt-transport-https \
 && curl -L https://packagecloud.io/tyk/tyk-pump/gpgkey | apt-key add - \
 && apt-get autoremove -y \
 && rm -rf /root/.cache


RUN mkdir -p /opt/mserv/tmp
RUN mkdir -p /opt/mserv/downloads
RUN mkdir -p /opt/mserv/plugins
COPY mserv /opt/mserv/

RUN chmod +x /opt/mserv/mserv

VOLUME ["/opt/mserv/"]

WORKDIR /opt/mserv

CMD ["/opt/mserv/mserv"]