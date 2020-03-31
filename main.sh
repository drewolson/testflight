#!/usr/bin/env bash

# This file was built to replace the functionality of the Makefile. However, It'd almost be better
# to have no script to handle these commands... or at least a minimal script & provide more learning
# lessons to fill in the gap with commands. This is going to get out of hand to maintain with
# multiple projects because of all the code duplication.

set -e

#
# Prepare the local shell context with any globally needed env vars
#
environment() {
  export SERVICE_NAME=testflight
  export COMPOSE_ENV=${1:-local}

  check_github_access_token
  export AWS_DEFAULT_REGION=us-east-1
  export GIT_VER=$(git rev-parse --short HEAD)
  export COMPOSE_NETWORK=${SERVICE_NAME}_${GIT_VER}

  if [[ ${ON_JENKINS} == true ]]; then
      echo "On Jenkins! Not doing any environment setup"
      return 0
  fi

  aws_credentials
}

aws_credentials() {
  unset AWS_PROFILE
  unset AWS_ACCESS_KEY_ID
  unset AWS_SECRET_ACCESS_KEY
  unset AWS_SESSION_TOKEN
  unset AWS_SECURITY_TOKEN

  if [[ -f ~/.talos/exports/aws.creds ]]; then
      rm ~/.talos/exports/aws.creds
  fi

  if ! type -a aws &>/dev/null; then
      echo "Must install AWS Cli"
      echo "https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html"
      exit 1
  fi

  if ! type -a jq &>/dev/null; then
      echo "Must install jq"
      echo "https://github.com/stedolan/jq/wiki/Installation"
      exit 1
  fi

  session_json=$(aws sts assume-role --role-arn arn:aws:iam::711065218908:role/config --role-session-name local-tmp)

  export AWS_ACCESS_KEY_ID=$(echo "${session_json}"|jq -r '.Credentials.AccessKeyId')
  export AWS_SECRET_ACCESS_KEY=$(echo "${session_json}"|jq -r '.Credentials.SecretAccessKey')
  export AWS_SESSION_TOKEN=$(echo "${session_json}"|jq -r '.Credentials.SessionToken')

  echo "Exported AWS credentials"
}

check_github_access_token() {
  if [[ -z $GITHUB_ACCESS_TOKEN ]]; then
    echo "ERROR: \$GITHUB_ACCESS_TOKEN is not available."
    echo "Please visit the following url for instructions."
    echo "https://github.com/calm/talos/blob/master/docs/github-pat.md#what-to-do-with-this-thing"
    exit 1
  fi
}

cmd_build() {
  environment

  echo "Building container for ${SERVICE_NAME}"
  DOCKER_TAG=$(docker_tag)

  DOCKER_ENV_TAG="$(env_tag)"

  docker build \
      --build-arg GITHUB_ACCESS_TOKEN=${GITHUB_ACCESS_TOKEN} \
      --build-arg SERVICE_NAME=${SERVICE_NAME} \
      -t ${DOCKER_TAG} \
      -t ${DOCKER_ENV_TAG} .
}

# pass in an optional tag or if left empty, git short hash is used
docker_tag() {
  if [[ -z ${SERVICE_NAME} ]]; then
    echo "SERVICE_NAME is required"
    exit 1
  fi

  GIT_VER=$(git rev-parse --short HEAD)
  DOCKER_TAG="${SERVICE_NAME}:${GIT_VER}"

  echo "864879987165.dkr.ecr.us-east-1.amazonaws.com/calm/${DOCKER_TAG}"
}

env_tag() {
  echo "864879987165.dkr.ecr.us-east-1.amazonaws.com/calm/${SERVICE_NAME}:${COMPOSE_ENV}"
}

# runs tests
# add -w or watch to the command for watch mode
cmd_tests() {
  environment
  DOCKER_ENV_TAG=$(env_tag)

  # run tests
  docker run \
    -e "AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION}" \
    -e "AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}" \
    -e "AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}" \
    -e "AWS_SESSION_TOKEN=${AWS_SESSION_TOKEN}" \
    ${DOCKER_ENV_TAG} bash -c "bash ./entrypoints/test.sh ${1}"
}

################################################

#
# Print all available commands
#
help() {
  echo "Usage:
  $ ./main.sh COMMAND

  Where COMMAND can be:
      build               # Build Docker container

      test                # Run tests. Add -w for watch mode
  " 1>&2
  exit 1
}

#
# Command translation from arg to func call
# If you wish to accesas a new function you must add a hook here
#
case "$1" in
  build)
    cmd_build ${@:2};;
  test)
    cmd_tests ${@:2};;
  *)
    help; exit 1
esac
