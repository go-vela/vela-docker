# Copyright (c) 2020 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

# The `clean` target is intended to clean the workspace
# and prepare the local changes for submission.
#
# Usage: `make clean`
.PHONY: clean
clean: tidy vet fmt fix

# The `run` target is intended to build and
# execute the Docker image for the plugin.
#
# Usage: `make run`
.PHONY: run
run: build docker-build docker-run

# The `tidy` target is intended to clean up
# the Go module files (go.mod & go.sum).
#
# Usage: `make tidy`
.PHONY: tidy
tidy:
	@echo
	@echo "### Tidying Go module"
	@go mod tidy

# The `vet` target is intended to inspect the
# Go source code for potential issues.
#
# Usage: `make vet`
.PHONY: vet
vet:
	@echo
	@echo "### Vetting Go code"
	@go vet ./...

# The `fmt` target is intended to format the
# Go source code to meet the language standards.
#
# Usage: `make fmt`
.PHONY: fmt
fmt:
	@echo
	@echo "### Formatting Go Code"
	@go fmt ./...

# The `fix` target is intended to rewrite the
# Go source code using old APIs.
#
# Usage: `make fix`
.PHONY: fix
fix:
	@echo
	@echo "### Fixing Go Code"
	@go fix ./...

# The `test` target is intended to run
# the tests for the Go source code.
#
# Usage: `make test`
.PHONY: test
test:
	@echo
	@echo "### Testing Go Code"
	@go test -race ./...

# The `test-cover` target is intended to run
# the tests for the Go source code and then
# open the test coverage report.
#
# Usage: `make test-cover`
.PHONY: test-cover
test-cover:
	@echo
	@echo "### Creating test coverage report"
	@go test -race -covermode=atomic -coverprofile=coverage.out ./...
	@echo
	@echo "### Opening test coverage report"
	@go tool cover -html=coverage.out

# The `build` target is intended to compile
# the Go source code into a binary.
#
# Usage: `make build`
.PHONY: build
build:
	@echo
	@echo "### Building release/vela-docker binary"
	GOOS=linux CGO_ENABLED=0 \
		go build -a \
		-o release/vela-docker \
		github.com/go-vela/vela-docker/cmd/vela-docker

# The `build-static` target is intended to compile
# the Go source code into a statically linked binary.
#
# Usage: `make build-static`
.PHONY: build-static
build-static:
	@echo
	@echo "### Building static release/vela-docker binary"
	GOOS=linux CGO_ENABLED=0 \
		go build -a \
		-ldflags '-s -w -extldflags "-static"' \
		-o release/vela-docker \
		github.com/go-vela/vela-docker/cmd/vela-docker

# The `check` target is intended to output all
# dependencies from the Go module that need updates.
#
# Usage: `make check`
.PHONY: check
check: check-install
	@echo
	@echo "### Checking dependencies for updates"
	@go list -u -m -json all | go-mod-outdated -update

# The `check-direct` target is intended to output direct
# dependencies from the Go module that need updates.
#
# Usage: `make check-direct`
.PHONY: check-direct
check-direct: check-install
	@echo
	@echo "### Checking direct dependencies for updates"
	@go list -u -m -json all | go-mod-outdated -direct

# The `check-full` target is intended to output
# all dependencies from the Go module.
#
# Usage: `make check-full`
.PHONY: check-full
check-full: check-install
	@echo
	@echo "### Checking all dependencies for updates"
	@go list -u -m -json all | go-mod-outdated

# The `check-install` target is intended to download
# the tool used to check dependencies from the Go module.
#
# Usage: `make check-install`
.PHONY: check-install
check-install:
	@echo
	@echo "### Installing psampaz/go-mod-outdated"
	@go get -u github.com/psampaz/go-mod-outdated

# The `bump-deps` target is intended to upgrade
# non-test dependencies for the Go module.
#
# Usage: `make bump-deps`
.PHONY: bump-deps
bump-deps: check
	@echo
	@echo "### Upgrading dependencies"
	@go get -u ./...

# The `bump-deps-full` target is intended to upgrade
# all dependencies for the Go module.
#
# Usage: `make bump-deps-full`
.PHONY: bump-deps-full
bump-deps-full: check
	@echo
	@echo "### Upgrading all dependencies"
	@go get -t -u ./...

# The `docker-build` target is intended to build
# the Docker image for the plugin.
#
# Usage: `make docker-build`
.PHONY: docker-build
docker-build:
	@echo
	@echo "### Building vela-docker:local image"
	@docker build --no-cache -t vela-docker:local .

# The `docker-test` target is intended to execute
# the Docker image for the plugin with test variables.
#
# Usage: `make docker-test`
.PHONY: docker-test
docker-test:
	@echo
	@echo "### Testing vela-docker:local image"
	@docker run --rm \
		-e BUILD_COMMIT \
		-e BUILD_EVENT \
		-e BUILD_TAG \
		-e DOCKER_USERNAME \
		-e DOCKER_PASSWORD \
		-e PARAMETER_ADD_HOSTS=host.company.com \
		-e PARAMETER_BUILD_ARGS=FOO=BAR \
		-e PARAMETER_CACHE_FROM=index.docker.in/target/vela-docker \
		-e PARAMETER_CGROUP_PARENT=parent \
		-e PARAMETER_COMPRESS=true \
		-e PARAMETER_CONTEXT=. \
		-e PARAMETER_CPU='{"period": 1, "quota": 1, "shares": 1, "set_cpus": "(0-3, 0,1)", "set_mems": "(0-3, 0,1)"}' \
		-e PARAMETER_DISABLE_CONTENT_TRUST=true \
		-e PARAMETER_FILE=Dockerfile.other \
		-e PARAMETER_FORCE_RM=true \
		-e PARAMETER_IMAGE_ID_FILE=path/to/file \
		-e PARAMETER_ISOLATION=hyperv \
		-e PARAMETER_LABELS=build.number=1 \
		-e PARAMETER_MEMORY=1 \	
		-e PARAMETER_MEMORY_SWAPS=1 \
		-e PARAMETER_NETWORK=default \
		-e PARAMETER_NO_CACHE=true \
		-e PARAMETER_OUTPUTS='type=local,dest=path' \
		-e PARAMETER_PLATFORM=linux \
		-e PARAMETER_PROGRESS=plain \
		-e PARAMETER_PULL=true \
		-e PARAMETER_QUIET=true \
		-e PARAMETER_REMOVE=true \
		-e PARAMETER_SECRETS='id=mysecret,src=/local/secret' \
		-e PARAMETER_SECURITY_OPTS=seccomp \
		-e PARAMETER_SHM_SIZES=1 \
		-e PARAMETER_SQUASH=true \
		-e PARAMETER_SSH_COMPONENTS='default|<id>[=<socket>|<key>[,<key>]]' \
		-e PARAMETER_STREAM=true \
		-e PARAMETER_TAGS=index.docker.io/target/vela-docker:latest \
		-e PARAMETER_TARGET=build \
		-e PARAMETER_ULIMITS=1 \
		-v $(shell pwd):/workspace \
		vela-docker:local

# The `docker-run` target is intended to execute
# the Docker image for the plugin.
#
# Usage: `make docker-run`
.PHONY: docker-run
docker-run:
	@echo
	@echo "### Executing vela-docker:local image"
	@docker run --rm --privileged --workdir /workspace \
		-e BUILD_COMMIT \
		-e BUILD_EVENT \
		-e BUILD_TAG \
		-e REGISTRY_NAME \
		-e REGISTRY_PASSWORD \
		-e REGISTRY_USERNAME \
		-e PARAMETER_CACHE_FROM \
		-e PARAMETER_FILE \
		-e PARAMETER_PLATFORM \
		-e PARAMETER_PROGRESS \
		-e PARAMETER_TAGS \
		-v $(shell pwd):/workspace \
		vela-docker:local
