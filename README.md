# Encontainer
Envcontainer, an extremely simple way to create a development environment with docker containers.

## Download
Check latest version on relase page [here](https://github.com/ErickMaria/envcontainer/releases).

## Requirements

- Docker version 18.02.0+.

Linux based systems:
- Ubuntu (64-bit)
- Debian (64-bit)
- CentOS (64-bit)
> [!NOTE] 
> Obs: **Windows system have not been yet tested.**

## Installation

### Download binary
Download a  [binary release](https://github.com/erickmaria/envcontainer/releases)

#### or use installation script
```bash
curl -fsSL https://raw.githubusercontent.com/erickmaria/envcontainer/refs/heads/main/scripts/install.sh | bash
```

## Quick Start

Using Envcontainer is a three-step process:

1. Define your app's container environment with a `.envcontainer.yaml` file.
    > [!TIP]
    > configuration files exemples [here](docs/configuration-files/READMED.md).
2. Run `envcontainer build`
3. Lastly, run `envcontainer start` and Envcontainer will start and enter your container.

A Envcontainer file looks like this:

```yaml
project:
  name: <YOUR_PROJECT_NAME> # Envcontainer
  version: <YOUR_PROJECT_VERSION> # 1.0.0
  description: <YOUR_PROJECT_DESCRIPTION> # Create a development environment for Envcontainer Application.
container:
  # write Dockerfile to build container
  build: |
    FROM ubuntu:latest
auto_stop: false

```
For more information about envcontainer, run `envcontainer help` 
 
 ```bash
UEnvcontainer helps you create, run, and manage reproducible development environments backed by Docker containers. Use the subcommands to build images, start/stop containers, and list projects.

Usage:
  envcontainer [command]

Available Commands:
  build       Build a Docker image from the envcontainer configuration in the current directory
  down        Stop and remove containers created by envcontainer in the current directory
  help        Help about any command
  init        Create a starter .envcontainer.yaml configuration file
  list        List discovered envcontainer projects and their status
  run         Run a one-off container from an image without saving configuration
  up          Start a container from envcontainer config and bind it to the current directory
  version     Show envcontainer CLI version

Flags:
  -h, --help   help for envcontainer

Use "envcontainer [command] --help" for more information about a command.
```

> [!NOTE] 
> `devcontainer` commands do not support all features to manage your containers, in this case you can use `docker` cli commands if you need.