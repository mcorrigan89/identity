#!/bin/bash

set -e
source ./script/variables.sh

echo "Building tag $DOCKER_TAG"
docker build . -t mcorrigan89/$DOCKER_IMAGE_NAME:$DOCKER_TAG \
  --build-arg FULL_HASH="$FULL_HASH" \
  --build-arg DOCKER_TAG="$DOCKER_TAG"

echo "Pushing ghcr.io/mcorrigan89/$DOCKER_IMAGE_NAME:$DOCKER_TAG"
docker push ghcr.io/mcorrigan89/$DOCKER_IMAGE_NAME:$DOCKER_TAG
