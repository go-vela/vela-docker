## Description

This plugin enables you to build and publish [Docker](https://www.docker.com/) images in a Vela pipeline.

Source Code: https://github.com/go-vela/vela-docker

Registry: https://hub.docker.com/r/target/vela-docker

## Usage

**NOTE: It is not recommended to use `latest` as the tag for the Docker image. Users should use a semantically versioned tag instead.**

Sample of building and publishing an image:

```yaml
steps:
  - name: publish_hello-world
    image: target/vela-docker:latest
    pull: always
    parameters:
      registry: index.docker.io
      tags: [ index.docker.io/octocat/hello-world:latest ]
```

Sample of building an image without publishing:

```diff
steps:
  - name: publish hello world
    image: target/vela-docker:latest
    pull: always
    parameters:
      registry: index.docker.io
      tags: [ index.docker.io/octocat/hello-world:latest ]    
+     dry_run: true
```

Sample of building and publishing an image with custom tags:

```diff
steps:
  - name: publish hello world
    image: target/vela-docker:latest
    pull: always
    parameters:
      registry: index.docker.io
+     tags: 
+       - index.docker.io/octocat/hello-world:latest
+       - index.docker.io/octocat/hello-world:1
+       - index.docker.io/octocat/hello-world:foobar
```

Sample of building and publishing an image with build arguments:

```diff
steps:
  - name: publish hello world
    image: target/vela-docker:latest
    pull: always
    parameters:
      registry: index.docker.io
      tags: [ index.docker.io/octocat/hello-world ]    
+     build_args:
+       - FOO=bar
```

Sample of building and publishing an image with image caching:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-docker:latest
    pull: always
    parameters:
      registry: index.docker.io
      tags: [ index.docker.io/octocat/hello-world:latest ]    
+     cache_from: index.docker.io/octocat/hello-world
```

Sample of building and publishing with custom daemon settings:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-docker:latest
    pull: always
    parameters:
      registry: index.docker.io
      tags: [ index.docker.io/octocat/hello-world:latest ]    
+     daemon: 
+       registry_mirror: mirror.index.docker.io
```

## Secrets

> **NOTE:** Users should refrain from configuring sensitive information in your pipeline in plain text.

### Internal

Users can use [Vela internal secrets](https://go-vela.github.io/docs/concepts/pipeline/secrets/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-docker:latest
    pull: always
+   secrets: [ docker_username, docker_password ]
    parameters:
      registry: index.docker.io
      tags: [ index.docker.io/octocat/hello-world:latest ] 
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
      tags: [ index.docker.io/octocat/hello-world:latest ] 
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
> `key.key` syntax signifies a new yaml object within the definition.

The following parameters are used to configure the image:

| Name                    | Description                                                                                                           | Required | Default           | Environment Variables                                               |
| ----------------------- | --------------------------------------------------------------------------------------------------------------------- | -------- | ----------------- | ------------------------------------------------------------------- |
| `add_hosts`             | enables adding a custom host-to-IP mapping - format (host:ip)                                                         | `false`  | N/A               | `PARAMETER_ADD_HOSTS`<br>`DOCKER_ADD_HOSTS`                         |
| `build_args`            | enables setting build-time variables                                                                                  | `false`  | N/A               | `PARAMETER_BUILD_ARGS`<br>`DOCKER_BUILD_ARGS`                       |
| `cache_from`            | enables setting images to consider as cache sources                                                                   | `false`  | N/A               | `PARAMETER_CACHE_FROM`<br>`DOCKER_CACHE_FROM`                       |
| `cgroup_parent`         | enables setting an optional parent cgroup for the container                                                           | `false`  | N/A               | `PARAMETER_CGROUP_PARENT`<br>`DOCKER_CGROUP_PARENT`                 |
| `compress`              | enables setting compression the build context using gzip                                                              | `false`  | N/A               | `PARAMETER_COMPRESS`<br>`DOCKER_COMPRESS`                           |
| `context`               | enables setting the build context                                                                                     | `true`   | `.`               | `PARAMETER_CONTEXT`<br>`DOCKER_CONTEXT`                             |
| `cpu`                   | enables setting the cpu parameter, see [cpu](#cpu) settings below                                                     | `false`  | N/A               | `PARAMETER_CPU`<br>`DOCKER_CPU`                                     |
| `daemon`                | enables setting the daemon parameter, see [daemon](#daemon) settings below                                            | `false`  | N/A               | `PARAMETER_DAEMON`<br>`DOCKER_DAEMON`                               |
| `disable_content_trust` | enables skipping image verification                                                                                   | `false`  | N/A               | `PARAMETER_DISABLE_CONTENT_TRUST`<br>`DOCKER_DISABLE_CONTENT_TRUST` |
| `dry_run`               | enables building the image without publishing                                                                         | `false`  | `false`           | `PARAMETER_DRY_RUN`<br>`DOCKER_DRY_RUN`                             |  
| `file`                  | enables setting the name of the Dockerfile                                                                            | `false`  | N/A               | `PARAMETER_FILE`<br>`DOCKER_FILE`                                   |
| `force_rm`              | enables setting always remove on intermediate containers                                                              | `false`  | N/A               | `PARAMETER_FORCE_RM`<br>`DOCKER_FORCE_RM`                           |
| `image_id_file`         | enables setting writing the image ID to the file                                                                      | `false`  | N/A               | `PARAMETER_IMAGE_ID_FILE`<br>`DOCKER_IMAGE_ID_FILE`                 |
| `isolation`             | enables container isolation technology                                                                                | `false`  | N/A               | `PARAMETER_ISOLATION`<br>`DOCKER_ISOLATION`                         |
| `labels`                | enables setting metadata for an image                                                                                 | `false`  | N/A               | `PARAMETER_LABELS`<br>`DOCKER_LABELS`                               |
| `log_level`             | set the log level for the plugin                                                                                      | `true`   | `info`            | `PARAMETER_LOG_LEVEL`<br>`DOCKER_LOG_LEVEL`                         |
| `memory`                | enables setting a memory limit                                                                                        | `false`  | N/A               | `PARAMETER_MEMORY`<br>`DOCKER_MEMORY`                               |
| `memory_swaps`          | enables setting a swap limit equal to memory plus swap: '-1' to enable unlimited swap                                 | `false`  | N/A               | `PARAMETER_MEMORY_SWAPS`<br>`DOCKER_MEMORY_SWAPS`                   |
| `network`               | enables setting the networking mode for the RUN instructions during build                                             | `false`  | N/A               | `PARAMETER_NETWORK`<br>`DOCKER_NETWORK`                             |
| `no_cache`              | enables setting not use cache when building the image                                                                 | `false`  | N/A               | `PARAMETER_NO_CACHE`<br>`DOCKER_NO_CACHE`                           |
| `outputs`               | enables setting an output destination - format (type=local,dest=path)                                                 | `false`  | N/A               | `PARAMETER_OUTPUTS`<br>`DOCKER_OUTPUTS`                             |
| `password`              | password for communication with the registry                                                                          | `true`   | N/A               | `PARAMETER_PASSWORD`<br>`DOCKER_PASSWORD`                           |  
| `platform`              | enables setting a platform if server is multi-platform capable                                                        | `false`  | N/A               | `PARAMETER_PLATFORM`<br>`DOCKER_PLATFORM`                           |
| `progress`              | enables setting type of progress output - options (auto\|plain\|tty)                                                  | `false`  | N/A               | `PARAMETER_PROGRESS`<br>`DOCKER_PROGRESS`                           |
| `pull`                  | enables always attempting to pull a newer version of the image                                                        | `false`  | N/A               | `PARAMETER_PULL`<br>`DOCKER_PULL`                                   |
| `quiet`                 | enables suppressing the build output and print image ID on success                                                    | `false`  | N/A               | `PARAMETER_QUIET`<br>`DOCKER_QUIET`                                 |
| `registry`              | Docker registry address to communicate with                                                                           | `true`   | `index.docker.io` | `PARAMETER_REGISTRY`<br>`DOCKER_REGISTRY`                           | 
| `remove`                | enables removing the intermediate containers after a successful build                                                 | `false`  | N/A               | `PARAMETER_REMOVE`<br>`DOCKER_REMOVE`                               |
| `secrets`               | enables setting a secret file to expose to the build - format (id=mysecret,src=/local/secret)                         | `false`  | N/A               | `PARAMETER_SECRETS`<br>`DOCKER_SECRETS`                             |
| `security_opts`         | enables setting security options                                                                                      | `false`  | N/A               | `PARAMETER_SECURITY_OPTS`<br>`DOCKER_SECURITY_OPTS`                 |
| `shm_sizes`             | enables setting the size of /dev/shm                                                                                  | `false`  | N/A               | `PARAMETER_SHM_SIZES`<br>`DOCKER_SHM_SIZES`                         |
| `squash`                | enables setting squash newly built layers into a single new layer                                                     | `false`  | N/A               | `PARAMETER_SQUASH`<br>`DOCKER_SQUASH`                               |
| `ssh_components`        | enables setting an ssh agent socket or keys to expose to the build - format (default\|<id>[=<socket>\|<key>[,<key>]]) | `false`  | N/A               | `PARAMETER_SSH_COMPONENTS`<br>`DOCKER_SSH_COMPONENTS`               |
| `stream`                | enables streaming attaches to server to negotiate build context                                                       | `false`  | N/A               | `PARAMETER_STREAM`<br>`DOCKER_STREAM`                               |
| `tags`                  | enables naming and optionally a tagging - format (name:tag)                                                           | `true`   | N/A               | `PARAMETER_TAGS`<br>`DOCKER_TAGS`                                   |
| `target`                | enables setting the target build stage to build.                                                                      | `false`  | N/A               | `PARAMETER_TARGET`<br>`DOCKER_TARGET`                               |
| `ulimits`               | enables setting ulimit options                                                                                        | `false`  | N/A               | `PARAMETER_ULIMITS`<br>`DOCKER_ULIMITS`                             |
| `username`              | user name for communication with the registry                                                                         | `true`   | N/A               | `PARAMETER_USERNAME`<br>`DOCKER_USERNAME`                           |

### CPU

The following settings are used to configure the `cpu` parameter:

| Name       | Description                                                              | Required | Default |
| ---------- | ------------------------------------------------------------------------ | -------- | ------- |
| `period`   | enables setting limits on the CPU CFS (Completely Fair Scheduler) period | `false`  | N/A     |
| `quota`    | enables setting limit on the CPU CFS (Completely Fair Scheduler) quota   | `false`  | N/A     |
| `shares`   | enables setting CPU shares (relative weight)                             | `false`  | N/A     |
| `set_cpus` | enables setting CPUs in which to allow execution (0-3, 0,1)              | `false`  | N/A     |
| `set_mems` | enables setting MEMs in which to allow execution (0-3, 0,1)              | `false`  | N/A     |

### Daemon

The following settings are used to configure the `daemon` parameter:

| Name                  | Description                                                                  | Required | Default |
| --------------------- | ---------------------------------------------------------------------------- | -------- | ------- |
| `bip`                 | enables specifying a network bridge IP                                       | `false`  | N/A     |
| `dns`                 | enables setting the DNS settings, see [dns](#dns) settings below             | `false`  | N/A     |
| `experimental`        | enables experimental features                                                | `false`  | N/A     |
| `insecure_registries` | enables insecure registry communication                                      | `false`  | N/A     |
| `ipv6`                | enables IPv6 networking                                                      | `false`  | N/A     |
| `mtu`                 | enables setting the containers network MTU                                   | `false`  | N/A     |
| `registry_mirrors`    | enables setting a preferred Docker registry mirror                           | `false`  | N/A     |
| `storage`             | enables setting the storage settings, see [storage](#storage) settings below | `false`  | N/A     |

### DNS

The following settings are used to configure the `dns daemon` setting:

| Name       | Description                                   | Required | Default |
| ---------- | --------------------------------------------- | -------- | ------- |
| `servers`  | enables setting the DNS server to use         | `false`  | N/A     |
| `searches` | enables setting the DNS search domains to use | `false`  | N/A     |

### Storage

The following settings are used to configure the `storage daemon` setting:

| Name      | Description                                             | Required | Default |
| --------- | ------------------------------------------------------- | -------- | ------- |
| `drivers` | enables setting an alternate storage driver             | `false`  | N/A     |
| `opts`    | enables setting options on the alternate storage driver | `false`  | N/A     |

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
      tags: [ index.docker.io/octocat/hello-world:latest ] 
+     log_level: trace

```

Below are a list of common problems and how to solve them:
