
# Envcontainer configuration files:

in this session you will find exemples of `envcontainer` configuration files.


## Minimal Exemple

```yaml
project:
  name: Envcontainer
  version: 0.0.1
  description: Envcontainer, an extremely simple way to create a development environment with docker containers.
container:
  build: |
    FROM ubuntu:latest
```

## Full Exemple


```yaml
project:
  name: Envcontainer
  version: 0.0.1
  description: Envcontainer, an extremely simple way to create a development environment with docker containers.
container:
  # your host user
  user: user
  # exemple: containerPort:HostPort
  ports:
    - 8081
    - 8080:8080
    - 8090:8010
  build: |
    FROM ubuntu:latest
# 
always-update: false
# Stop the container when leaving it
auto-stop: false
# mount container volumes
# support types 'volume' and 'bind' 
# ex: source:target:(volume|bind)
# default volume type is 'volume'
# you can omit the source and volume type in the mounts session declaration, below are some examples
# [IMPORTANT] if you are using bind volume type with default prefix path. Be careful with this folder, Don't share this before checking if you have sensitive files
mounts:
  # {prefix-name}-{volume-path}:{target}:{default-volume}
  # ex: envcontainer-tmp:/tmp/:volume
  - /tmp/
  # {prefix-name}-{volume-path}:{target}:{declared-volume}
  # ex: envcontainer-tmp:/tmp/:volume
  - /tmp/:volume
  # {prefix-path}-{bind-path}:{target}:{declared-volume}
  # ex: ${CURRENT_PROJECT_PATH}/.envcontainer/{CURRENT_BIND_PATH}:{target}:bind
  # [notice] path on host will be created if it does not exist
  - /tmp/:bind
  # {prefix-name}-{volume-name}:{target}:{declared-volume}
  # ex: envcontainer-data-tmp:/tmp/:volume
  - data-tmp:/tmp:volume
  # {prefix-name}-{volume-path}:{target}:{declared-volume}
  # {source}:{target}:{declared-volume}
  - /home/erick/.envcontainer/tmp:/tmp/test/ad:bind
```