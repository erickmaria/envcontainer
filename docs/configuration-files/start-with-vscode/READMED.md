
# Envcontainer configuration files:

in this session we show you how start `envcontainer` with vscode


## Configuration File Exemple

```yaml
project:
  name: Envcontainer
  version: 1.0.0
  description: Create a development environment for Envcontainer Application.
container:
  build: |
    FROM ubuntu:22.04

    # INSTALL OPENSSH
    RUN apt-get update && apt-get install sudo openssh-server -y && \
        mkdir /var/run/sshd

    # ADD USER
    RUN useradd -ms /bin/bash envcontainer && \
        echo 'envcontainer:envcontainer' | chpasswd && \
        echo "envcontainer ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers
    USER envcontainer

    # -------------------------------------
    # ADD YOUR DOCKERFILE INSTRUCTION HERE
    # -------------------------------------

    # Start SSH server
    CMD ["sudo", "/usr/sbin/sshd", "-D"]

auto-stop: false

```

Run `envcontainer start --code`