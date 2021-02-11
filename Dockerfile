FROM golang:1.13 AS builder

# Set some shell options for using pipes and such
SHELL [ "/bin/bash", "-euo", "pipefail", "-c" ]

# Install common CA certificates to blag later
RUN apt-get update \
  && apt-get install --assume-yes --no-install-recommends ca-certificates \
  && apt-get autoremove --assume-yes \
  && rm -rf /root/.cache

# Don't call any C code (the 'scratch' base image used later won't have any libraries to reference)
ENV CGO_ENABLED=0

# Can't use Go modules because of broken vendoring/dependencies in Tyk Gateway v2.9.3 - TODO
ENV GO111MODULE=off

# Precompile the entire Go standard library into a Docker cache layer: useful for other projects too!
# cf. https://www.reddit.com/r/golang/comments/hj4n44/improved_docker_go_module_dependency_cache_for/
RUN go install -ldflags="-buildid= -w" -trimpath -v std

WORKDIR /go/src/github.com/TykTechnologies/mserv

# vvv Can't use Go modules because of broken vendoring/dependencies in Tyk Gateway v2.9.3 - TODO
# # This will save Go dependencies in the Docker cache, until/unless they change
# COPY go.mod go.sum ./

# # Download and precompile all third party libraries
# RUN go mod graph | awk '$1 !~ "@" { print $2 }' | xargs go get -ldflags="-buildid= -w" -trimpath -v
# ^^^ Can't use Go modules because of broken vendoring/dependencies in Tyk Gateway v2.9.3 - TODO

# Add the sources
COPY . .

# Compile!
RUN go build -ldflags="-buildid= -w" -trimpath -v -o /bin/mserv

FROM debian:buster-slim AS runner

# Set some shell options for using pipes and such
SHELL [ "/bin/bash", "-euo", "pipefail", "-c" ]

ENV TYKVERSION 0.1
ENV TYK_MSERV_CONFIG /etc/mserv/mserv.json

LABEL Description="Tyk MServ service docker image" Vendor="Tyk" Version=$TYKVERSION

RUN mkdir -p /opt/mserv/downloads /opt/mserv/plugins

WORKDIR /opt/mserv

# Bring common CA certificates and binary over
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /bin/mserv /opt/mserv/mserv

ENTRYPOINT [ "/opt/mserv/mserv" ]
