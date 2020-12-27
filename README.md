# Encontainer
Envcontainer, an extremely simple way to create a development environment with docker containers.

## Download
Check latest version on relase page [here](https://github.com/ErickMaria/envcontainer/releases).

## Requirements

Docker version 18.02.0+.

Docker compose 1.27.0+.

Linux based systems
- Ubuntu (64-bit)
- Debian (64-bit)
- CentOS (64-bit)
> Obs: **Windows system have not been yet tested.**


## Usage
### Initialize
Inside your project run command

```bash
envcontainer init
```

This command will create a folder named by ".envcontainer", with a some files like Dockerfile, docker-compose and .variables.

How your can edit Dockerfile into .envcontainer folder to prepare your workspace inside container

### Build

```bash
envcontainer build
```

By default init command already execute build command to build Dockerfile and up container but if you edited Dockerfile make sure run this command to the envcontainer get the latest modifications

### Start

```bash
envcontainer start
```

Using this commmand envcontainer will enter inside container shell.

### Stop

```bash
envcontainer stop
```


the stop command will kill container stopping your envcontainer environment

### Delete

```bash
envcontainer delete
```

will delete container running, image and envcontainer files.

#### for more information about commands run
 
 ```bash
envcontainer help
```
