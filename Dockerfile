# SPDX-License-Identifier: Apache-2.0

#########################################################################
##    docker build --no-cache --target certs -t vela-docker:certs .    ##
#########################################################################

FROM alpine:3.24.4@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf as certs

RUN apk add --update --no-cache ca-certificates

##########################################################
##    docker build --no-cache -t vela-docker:local .    ##
##########################################################

FROM docker:29.10-dind@sha256:4bbc7223f81db790097c6dea42fb93e05b8fead4f0431fbe6ccaba6e8473201c

ENV DOCKER_HOST=unix:///var/run/docker.sock

ENV DOCKER_BUILDKIT=1

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-docker /bin/vela-docker

ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "/bin/vela-docker"]