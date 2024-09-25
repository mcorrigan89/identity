#!/bin/bash

echo "generating variables... ";

# different names this service goes by
DOCKER_IMAGE_NAME=identity
FULL_HASH=`git rev-parse HEAD`
SHORT_HASH=`git rev-parse HEAD | cut -c1-8`
DOCKER_TAG=release-identity-$SHORT_HASH
