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
Usage: envcontainer COMMAND --FLAGS

Commands
build:                  build a image using envcontainer configuration in the current directory
down:                   remove all envcontainer configuration running in the current directory
    --name:                     container name
    --get-closer:               will get the closest configuration file and remove all envcontainer
help:                   Run build COMMAND' for more information on a command. See: 'build help'
ls:                     list envcontainers
run:                    execute an .envcontainer on the current directory without saving it locally
    --name:                     container name
    --image:                    envcontainer image
    --command:                  execute command inside container
up:                     run the envcontainer configuration to start the container and link it to the current directory
    --get-closer:               will get the closest configuration file to run a new container
    --auto-stop:                terminal shell that must be used
    --code:                     open with vscode
    --port:                     ssh port that vscode will use connect
version:                show envcontainer version
```

> [!NOTE] 
> `devcontainer` commands do not support all features to manage your containers, in this case you can use `docker` cli commands if you need.