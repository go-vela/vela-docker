# SPDX-License-Identifier: Apache-2.0

#########################################################################
##    docker build --no-cache --target certs -t vela-docker:certs .    ##
#########################################################################

FROM alpine:3.19.1@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b as certs

RUN apk add --update --no-cache ca-certificates

##########################################################
##    docker build --no-cache -t vela-docker:local .    ##
##########################################################

FROM docker:27.5-dind@sha256:aa3df78ecf320f5fafdce71c659f1629e96e9de0968305fe1de670e0ca9176ce

ENV DOCKER_HOST=unix:///var/run/docker.sock

ENV DOCKER_BUILDKIT=1

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-docker /bin/vela-docker

ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "/bin/vela-docker"]