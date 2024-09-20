#!/bin/bash

docker build -t ghcr.io/unafraid/proxyless/server:latest -f cmd/server/Dockerfile --no-cache .
docker push ghcr.io/unafraid/proxyless/server:latest

docker build -t ghcr.io/unafraid/proxyless/client:latest -f cmd/client/Dockerfile --no-cache .
docker push ghcr.io/unafraid/proxyless/client:latest
