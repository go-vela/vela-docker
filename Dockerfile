# SPDX-License-Identifier: Apache-2.0

#########################################################################
##    docker build --no-cache --target certs -t vela-docker:certs .    ##
#########################################################################

FROM alpine:3.21.3@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c as certs

RUN apk add --update --no-cache ca-certificates

##########################################################
##    docker build --no-cache -t vela-docker:local .    ##
##########################################################

FROM docker:28.3-dind@sha256:cf5c79bfb90a1b8ef3947b013fe61b3d66ad790ab4bcf3ee5319e8b88134f553

ENV DOCKER_HOST=unix:///var/run/docker.sock

ENV DOCKER_BUILDKIT=1

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-docker /bin/vela-docker

ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "/bin/vela-docker"]