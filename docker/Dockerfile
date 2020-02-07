FROM debian:buster-slim

ENV TYKVERSION 0.1
ENV TYK_MSERV_CONFIG /etc/mserv/mserv.json

LABEL Description="Tyk MServ service docker image" Vendor="Tyk" Version=$TYKVERSION

RUN apt-get update \
 && apt-get upgrade -y \
 && apt-get install -y --no-install-recommends \
            ca-certificates \
 && apt-get autoremove -y \
 && rm -rf /root/.cache


RUN mkdir -p /opt/mserv/downloads
RUN mkdir -p /opt/mserv/plugins
COPY mserv /opt/mserv/

RUN chmod +x /opt/mserv/mserv

WORKDIR /opt/mserv

ENTRYPOINT ["/opt/mserv/mserv"]
