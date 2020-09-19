#!/bin/bash

#######################
## DEFAULT VARAIBLES ##
#######################

DOCKER_COMPOSE_PATH=compose
DOCKER_COMPOSE=$DOCKER_COMPOSE_PATH/docker-compose.yaml

DEAFULT_PROJECT_NAME=environment
DEFAULT_DOCKER_COMPOSE_ENV_FILE_PATH=$DOCKER_COMPOSE_PATH/env
DEFAULT_DOCKER_COMPOSE_ENV_FILE=$DEFAULT_DOCKER_COMPOSE_ENV_FILE_PATH/.env

##############
## FUNCIONS ##
##############

options(){
  # flags validation

  while test $# -gt 0;
  do
    case "$1" in
      --project)
          shift
          PROJECT_NAME=$1
          shift
          ;;
      -p)
          shift
          PROJECT_NAME=$1
          shift
          ;;
      --listener)
          shift
          PORT_LISTENER=$1
          shift
          ;;
      -l)
          shift
          PORT_LISTENER=$1
          shift
          ;;
      --envfile)
          shift
          ENV_FILE=$1
          shift
          ;;
      -e)
          shift
          ENV_FILE=$1
          shift
          ;;
      *)
      echo "$1 is not a recognized flag!"
      exit;
      ;;
    esac
  done

}

docker_compose_port_tag_create(){
  DOCKER_COMPOSE_PORT_TAG=`cat <<EOF
ports:
      - "$PORT_LISTENER"
EOF
`
  echo $DOCKER_COMPOSE_PORT_TAG
}

default_options(){
  # TEST && TRUE || FALSE

  [ -z "$PROJECT_NAME" ] && PROJECT_NAME=$DEAFULT_PROJECT_NAME
  [ -z "$PORT_LISTENER" ] && echo "" || docker_compose_port_tag_create
  [[ -z "$ENV_FILE" ]] && ENV_FILE=$DEFAULT_DOCKER_COMPOSE_ENV_FILE; env_file_generate
}

env_file_generate(){
  mkdir -p $DEFAULT_DOCKER_COMPOSE_ENV_FILE_PATH
  touch $DEFAULT_DOCKER_COMPOSE_ENV_FILE
}

dockerfile_create() {
  touch Dockerfile
}

docker_compose_create(){
  
  mkdir -p $DOCKER_COMPOSE_PATH
  
  # env_file_generate

  cat <<EOF > $DOCKER_COMPOSE
version: "3.6"

services:
  environment:
    build:
      dockerfile: Dockerfile
      context: ../
    working_dir: /home/${PROJECT_NAME}
    env_file:
      - $ENV_FILE
    $DOCKER_COMPOSE_PORT_TAG
    volumes:
      - type: bind
        source: ../../
        target: /home/${PROJECT_NAME}
    stdin_open: true
    tty: true
EOF
  
  sed -i '/^[[:space:]]*$/d' $DOCKER_COMPOSE

}

###############
## BOOTSTRAP ##
############### 

while test $# -gt 0;
do
  case "$1" in
    init)
      shift
        echo "init"
        options $@
        default_options
        docker_compose_create
        dockerfile_create
      shift
      ;;
    build)
      shift
      echo "build"
      shift
      ;;
    up)
      shift
      echo "up"
      shift
      ;;
    *)
    # echo "$1 is not a recognized flag!"
    exit;
    ;;
  esac
done