#!/usr/bin/env bash

TAG=v0.10.6

set -euox pipefail

docker build --no-cache --pull -t kubermatic/kubernetes-test-binaries:${TAG} .
docker push kubermatic/kubernetes-test-binaries:${TAG}
