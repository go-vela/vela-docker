# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

#########################################################################
##    docker build --no-cache --target certs -t vela-docker:certs .    ##
#########################################################################

FROM alpine as certs

RUN apk add --update --no-cache ca-certificates

##########################################################
##    docker build --no-cache -t vela-docker:local .    ##
##########################################################

FROM docker:dind

# busybox "ip" is insufficient:
#   [rootlesskit:child ] error: executing [[ip tuntap add name tap0 mode tap] [ip link set tap0 address 02:50:00:00:00:01]]: exit status 1
RUN apk add --no-cache iproute2 fuse-overlayfs

# "/run/user/UID" will be used by default as the value of XDG_RUNTIME_DIR
RUN mkdir /run/user && chmod 1777 /run/user

# create a default user preconfigured for running rootless dockerd
# set subuid and subgid ranges — very high due to occasional docker layers being out-of-range with default
RUN set -eux; \
	adduser -h /home/rootless -g 'Rootless' -D -u 1000 rootless; \
	printf 'dockremap:100000:2000000000\nrootless:2000100000:2000000000' > /etc/subuid; \
	printf 'dockremap:100000:2000000000\nrootless:2000100000:2000000000' > /etc/subgid; \
	chmod 0644 /etc/subuid; \
	chmod 0644 /etc/subgid

# establish docker rootless extras
RUN set -eux; \
	\
	apkArch="$(apk --print-arch)"; \
	case "$apkArch" in \
		'x86_64') \
			url='https://download.docker.com/linux/static/stable/x86_64/docker-rootless-extras-20.10.18.tgz'; \
			;; \
		'aarch64') \
			url='https://download.docker.com/linux/static/stable/aarch64/docker-rootless-extras-20.10.18.tgz'; \
			;; \
		*) echo >&2 "error: unsupported 'rootless.tgz' architecture ($apkArch)"; exit 1 ;; \
	esac; \
	\
	wget -O 'rootless.tgz' "$url"; \
	\
	tar --extract \
		--file rootless.tgz \
		--strip-components 1 \
		--directory /usr/local/bin/ \
		'docker-rootless-extras/rootlesskit' \
		'docker-rootless-extras/rootlesskit-docker-proxy' \
		'docker-rootless-extras/vpnkit' \
	; \
	rm rootless.tgz; \
	\
	rootlesskit --version; \
	vpnkit --version

# pre-create "/var/lib/docker" for our rootless user
RUN set -eux; \
	mkdir -p /home/rootless/.local/share/docker; \
	chown -R rootless:rootless /home/rootless/.local/share/docker
VOLUME /home/rootless/.local/share/docker
USER rootless

ENV DOCKER_HOST=unix:///var/run/user/1000/docker.sock

ENV DOCKER_BUILDKIT=1

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-docker /bin/vela-docker

ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "/bin/vela-docker"]