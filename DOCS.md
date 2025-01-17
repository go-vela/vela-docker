## Description

This plugin enables you to build and publish [Docker](https://www.docker.com/) images in a Vela pipeline.

Source Code: https://github.com/go-vela/vela-docker

Registry: https://hub.docker.com/r/target/vela-docker

## Usage

> **NOTE:**
>
> Users should refrain from using latest as the tag for the Docker image.
>
> It is recommended to use a semantically versioned tag instead.

Samples of building and publishing an image:

```yaml
steps:
  - name: publish_hello-world
    image: target/vela-docker:latest
    pull: always
    parameters:
      registry: index.docker.io
      repo: octocat/hello-world
      tags: [ latest ]
```

```yaml
steps:
  - name: publish_hello-world
    image: target/vela-docker:latest
    pull: always
    parameters:
      registry: index.docker.io
      tags: [ index.docker.io/octocat/hello-world:latest ]
```

> **NOTE:** The two above samples are functionally equivalent.

Sample of building an image without publishing:

```diff
steps:
  - name: publish hello world
    image: target/vela-docker:latest
    pull: always
    parameters:
+     dry_run: true
      registry: index.docker.io
      repo: octocat/hello-world
      tags: [ latest ]
```

Sample of building and publishing an image with custom tags:

```diff
steps:
  - name: publish hello world
    image: target/vela-docker:latest
    pull: always
    parameters:
      registry: index.docker.io
      repo: octocat/hello-world
-     tags: [ latest ]
+     tags: 
+       - latest
+       - octocat/hello-world:1
+       - index.docker.io/octocat/hello-world:foobar
```

Sample of building and publishing an image with build arguments:

```diff
steps:
  - name: publish hello world
    image: target/vela-docker:latest
    pull: always
    parameters:
+     build_args:
+       - FOO=bar
      registry: index.docker.io
      repo: octocat/hello-world
      tags: [ latest ]
```

Sample of building and publishing an image with image caching:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-docker:latest
    pull: always
    parameters:
+     cache_from: index.docker.io/octocat/hello-world
      registry: index.docker.io
      repo: octocat/hello-world
      tags: [ latest ]
```

Sample of building and publishing with custom daemon settings:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-docker:latest
    pull: always
    parameters:
+     daemon: 
+       registry_mirrors: mirror.index.docker.io
      registry: index.docker.io
      repo: octocat/hello-world
      tags: [ latest ]
```

## Secrets

> **NOTE:** Users should refrain from configuring sensitive information in your pipeline in plain text.

### Internal

Users can use [Vela internal secrets](https://go-vela.github.io/docs/tour/secrets/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-docker:latest
    pull: always
+   secrets: [ docker_username, docker_password ]
    parameters:
      registry: index.docker.io
      repo: octocat/hello-world
      tags: [ latest ]
-     username: octocat
-     password: superSecretPassword
```

> This example will add the secrets to the `publish_hello-world` step as environment variables:
>
> * `DOCKER_USERNAME=<value>`
> * `DOCKER_PASSWORD=<value>`

### External

The plugin accepts the following files for authentication:

| Parameter  | Volume Configuration                                                          |
| ---------- | ----------------------------------------------------------------------------- |
| `password` | `/vela/parameters/docker/password`, `/vela/secrets/docker/password` |
| `username` | `/vela/parameters/docker/username`, `/vela/secrets/docker/username` |

Users can use [Vela external secrets](https://go-vela.github.io/docs/concepts/pipeline/secrets/origin/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-docker:latest
    pull: always
    parameters:
      registry: index.docker.io
      repo: octocat/hello-world
      tags: [ latest ]
-     username: octocat
-     password: superSecretPassword
```

> This example will read the secret values in the volume stored at `/vela/secrets/`
## Parameters

> **NOTE:**
>
> The plugin supports reading all parameters via environment variables or files.
>
> Any values set from a file take precedence over values set from the environment.
>
> By default [build kit](https://docs.docker.com/develop/develop-images/build_enhancements/) is on; it can be turned off by setting `DOCKER_BUILDKIT=0` in the environment.
>
> The `key.key` syntax signifies a new yaml object within the definition.

The following parameters are used to configure the image:

| Name                    | Description                                                                                                                       | Required | Default           | Environment Variables                                               |
| ----------------------- | --------------------------------------------------------------------------------------------------------------------------------- | -------- | ----------------- | ------------------------------------------------------------------- |
| `add_hosts`             | set a custom host-to-IP mapping - format (host:ip)                                                                                | `false`  | N/A               | `PARAMETER_ADD_HOSTS`<br/>`DOCKER_ADD_HOSTS`                         |
| `build_args`            | set variables to pass to the image at build-time                                                                                  | `false`  | N/A               | `PARAMETER_BUILD_ARGS`<br/>`DOCKER_BUILD_ARGS`                       |
| `cache_from`            | set of images to consider as cache sources                                                                                        | `false`  | N/A               | `PARAMETER_CACHE_FROM`<br/>`DOCKER_CACHE_FROM`                       |
| `cgroup_parent`         | set a parent cgroup for the container                                                                                             | `false`  | N/A               | `PARAMETER_CGROUP_PARENT`<br/>`DOCKER_CGROUP_PARENT`                 |
| `compress`              | enable compressing the build context using gzip                                                                                   | `false`  | `false`           | `PARAMETER_COMPRESS`<br/>`DOCKER_COMPRESS`                           |
| `context`               | set of files and/or directory to build the image from                                                                             | `true`   | `.`               | `PARAMETER_CONTEXT`<br/>`DOCKER_CONTEXT`                             |
| `cpu`                   | set the cpu parameter, see [cpu](#cpu) settings below                                                                             | `false`  | N/A               | `PARAMETER_CPU`<br/>`DOCKER_CPU`                                     |
| `daemon`                | set the daemon parameter, see [daemon](#daemon) settings below                                                                    | `false`  | N/A               | `PARAMETER_DAEMON`<br/>`DOCKER_DAEMON`                               |
| `disable_content_trust` | enable skipping verification of the image                                                                                         | `false`  | `true`            | `PARAMETER_DISABLE_CONTENT_TRUST`<br/>`DOCKER_DISABLE_CONTENT_TRUST` |
| `dry_run`               | enable building the image without publishing                                                                                      | `false`  | `false`           | `PARAMETER_DRY_RUN`<br/>`DOCKER_DRY_RUN`                             |
| `file`                  | set the name of the Dockerfile                                                                                                    | `false`  | N/A               | `PARAMETER_FILE`<br/>`DOCKER_FILE`                                   |
| `force_rm`              | enable always removing the intermediate containers after a successful build                                                       | `false`  | `false`           | `PARAMETER_FORCE_RM`<br/>`DOCKER_FORCE_RM`                           |
| `image_id_file`         | set the file to write the image ID to                                                                                             | `false`  | N/A               | `PARAMETER_IMAGE_ID_FILE`<br/>`DOCKER_IMAGE_ID_FILE`                 |
| `isolation`             | set container isolation technology                                                                                                | `false`  | N/A               | `PARAMETER_ISOLATION`<br/>`DOCKER_ISOLATION`                         |
| `labels`                | set metadata for an image                                                                                                         | `false`  | N/A               | `PARAMETER_LABELS`<br/>`DOCKER_LABELS`                               |
| `log_level`             | set the log level for the plugin                                                                                                  | `true`   | `info`            | `PARAMETER_LOG_LEVEL`<br/>`DOCKER_LOG_LEVEL`                         |
| `memory`                | set memory limit                                                                                                                  | `false`  | N/A               | `PARAMETER_MEMORY`<br/>`DOCKER_MEMORY`                               |
| `memory_swaps`          | set the swap limit equal to memory plus swap: '-1' to enable unlimited swap                                                       | `false`  | N/A               | `PARAMETER_MEMORY_SWAPS`<br/>`DOCKER_MEMORY_SWAPS`                   |
| `network`               | set the networking mode for the RUN instructions during build                                                                     | `false`  | N/A               | `PARAMETER_NETWORK`<br/>`DOCKER_NETWORK`                             |
| `no_cache`              | disable caching when building the image                                                                                           | `false`  | `false`           | `PARAMETER_NO_CACHE`<br/>`DOCKER_NO_CACHE`                           |
| `output`                | set the output destination - format (type=local,dest=path)                                                                        | `false`  | N/A               | `PARAMETER_OUTPUTS`<br/>`DOCKER_OUTPUTS`                             |
| `password`              | set password for communication with the registry                                                                                  | `true`   | N/A               | `PARAMETER_PASSWORD`<br/>`DOCKER_PASSWORD`                           |
| `platform`              | set a platform if server is multi-platform capable                                                                                | `false`  | N/A               | `PARAMETER_PLATFORM`<br/>`DOCKER_PLATFORM`                           |
| `progress`              | set type of progress output - options (auto\|plain\|tty)                                                                          | `false`  | N/A               | `PARAMETER_PROGRESS`<br/>`DOCKER_PROGRESS`                           |
| `pull`                  | enable always attempting to pull a newer version of the image                                                                     | `false`  | `false`           | `PARAMETER_PULL`<br/>`DOCKER_PULL`                                   |
| `quiet`                 | enable suppressing the build output and print image ID on success                                                                 | `false`  | `false`           | `PARAMETER_QUIET`<br/>`DOCKER_QUIET`                                 |
| `registry`              | set Docker registry address to communicate with                                                                                   | `true`   | `index.docker.io` | `PARAMETER_REGISTRY`<br/>`DOCKER_REGISTRY`                           |
| `remove`                | enable removing the intermediate containers after a successful build                                                              | `false`  | `true`            | `PARAMETER_REMOVE`<br/>`DOCKER_REMOVE`                               |
| `repo`                  | set Docker repository for the image                                                                                               | `false`  | N/A               | `PARAMETER_REPO`<br/>`DOCKER_REPO`                                   |
| `secret`                | set secret file to expose to the build (only if BuildKit enabled) - format (id=mysecret,src=/local/secret)                        | `false`  | N/A               | `PARAMETER_SECRETS`<br/>`DOCKER_SECRETS`                             |
| `security_opts`         | set options for security                                                                                                          | `false`  | N/A               | `PARAMETER_SECURITY_OPTS`<br/>`DOCKER_SECURITY_OPTS`                 |
| `shm_sizes`             | set the size of /dev/shm                                                                                                          | `false`  | N/A               | `PARAMETER_SHM_SIZES`<br/>`DOCKER_SHM_SIZES`                         |
| `squash`                | enable squashing newly built layers into a single new layer                                                                       | `false`  | `false`           | `PARAMETER_SQUASH`<br/>`DOCKER_SQUASH`                               |
| `ssh_components`        | set SSH agent socket or keys to expose to the build (only if BuildKit enabled) - format `(default\|<id>[=<socket>\|<key>[,<key>]])` | `false`  | N/A               | `PARAMETER_SSH_COMPONENTS`<br/>`DOCKER_SSH_COMPONENTS`               |
| `stream`                | enable stream attaching to the server to negotiate build context                                                                  | `false`  | `false`           | `PARAMETER_STREAM`<br/>`DOCKER_STREAM`                               |
| `tags`                  | set the tags for the Docker image - format (name:tag)                                                                             | `true`   | N/A               | `PARAMETER_TAGS`<br/>`DOCKER_TAGS`                                   |
| `target`                | set the target build stage to build                                                                                               | `false`  | N/A               | `PARAMETER_TARGET`<br/>`DOCKER_TARGET`                               |
| `ulimits`               | set options for ulimits                                                                                                           | `false`  | N/A               | `PARAMETER_ULIMITS`<br/>`DOCKER_ULIMITS`                             |
| `username`              | set user name for communication with the registry                                                                                 | `true`   | N/A               | `PARAMETER_USERNAME`<br/>`DOCKER_USERNAME`                           |

### CPU

The following settings are used to configure the `cpu` parameter:

| Name       | Description                                                 | Required | Default |
| ---------- | ----------------------------------------------------------- | -------- | ------- |
| `period`   | set limit on the CPU CFS (Completely Fair Scheduler) period | `false`  | N/A     |
| `quota`    | set limit on the CPU CFS (Completely Fair Scheduler) quota  | `false`  | N/A     |
| `shares`   | set CPU shares (relative weight)                            | `false`  | N/A     |
| `set_cpus` | set CPUs in which to allow execution (0-3, 0,1)             | `false`  | N/A     |
| `set_mems` | set MEMs in which to allow execution (0-3, 0,1)             | `false`  | N/A     |

### Daemon

The following settings are used to configure the `daemon` parameter:

| Name                  | Description                                                      | Required | Default |
| --------------------- | ---------------------------------------------------------------- | -------- | ------- |
| `bip`                 | set a network bridge IP                                          | `false`  | N/A     |
| `dns`                 | set the DNS settings, see [dns](#dns) settings below             | `false`  | N/A     |
| `experimental`        | enable experimental features                                     | `false`  | N/A     |
| `insecure_registries` | set the insecure Docker registries                               | `false`  | N/A     |
| `ipv6`                | enable IPv6 networking                                           | `false`  | N/A     |
| `mtu`                 | set the network MTU for the contain                              | `false`  | N/A     |
| `registry_mirrors`    | set the Docker registry mirrors                                  | `false`  | N/A     |
| `storage`             | set the storage settings, see [storage](#storage) settings below | `false`  | N/A     |

### DNS

The following settings are used to configure the `dns daemon` setting:

| Name       | Description                | Required | Default |
| ---------- | -------------------------- | -------- | ------- |
| `servers`  | set the DNS nameservers    | `false`  | N/A     |
| `searches` | set the DNS search domains | `false`  | N/A     |

### Storage

The following settings are used to configure the `storage daemon` setting:

| Name     | Description                            | Required | Default |
| -------- | -------------------------------------- | -------- | ------- |
| `driver` | set the storage driver for the daemon  | `false`  | N/A     |
| `opts`   | set the storage options for the daemon | `false`  | N/A     |

## Template

COMING SOON!

## Troubleshooting

You can start troubleshooting this plugin by tuning the level of logs being displayed:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-docker:latest
    pull: always
    parameters:
      registry: index.docker.io
      repo: octocat/hello-world
      tags: [ latest ]
+     log_level: trace

```

Below are a list of common problems and how to solve them:
