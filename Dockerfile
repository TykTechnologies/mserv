FROM golang:1.20.11 AS builder

# Set some shell options for using pipes and such
SHELL [ "/bin/bash", "-euo", "pipefail", "-c" ]

# Install common CA certificates to blag later
RUN apt-get update \
  && apt-get install --assume-yes --no-install-recommends ca-certificates \
  && apt-get autoremove --assume-yes \
  && rm -rf /root/.cache

# Don't call any C code (the 'scratch' base image used later won't have any libraries to reference)
ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/TykTechnologies/mserv

COPY . .

RUN go build -ldflags="-buildid= -w" -trimpath -v -o /bin/mserv
RUN mkdir -p /opt/mserv/downloads /opt/mserv/plugins

FROM gcr.io/distroless/base:nonroot AS runner
USER 65532

ENV TYK_MSERV_CONFIG /etc/mserv/mserv.json

LABEL Description="Tyk MServ service docker image" Vendor="Tyk" Version=$TYKVERSION

WORKDIR /opt/mserv

# Bring common CA certificates and binary over.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /bin/mserv /opt/mserv/mserv
COPY --from=builder /opt/mserv/downloads /opt/mserv/downloads
COPY --from=builder /opt/mserv/plugins /opt/mserv/plugins

ENTRYPOINT [ "/opt/mserv/mserv" ]
