# SPDX-License-Identifier: Apache-2.0

#########################################################################
##    docker build --no-cache --target certs -t vela-docker:certs .    ##
#########################################################################

FROM alpine as certs

RUN apk add --update --no-cache ca-certificates

##########################################################
##    docker build --no-cache -t vela-docker:local .    ##
##########################################################

FROM docker:24.0-dind

ENV DOCKER_HOST=unix:///var/run/docker.sock

ENV DOCKER_BUILDKIT=1

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-docker /bin/vela-docker

ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "/bin/vela-docker"]