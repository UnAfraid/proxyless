#!/bin/bash

docker build -t ghcr.io/unafraid/proxyless/client:latest -f cmd/client/Dockerfile .
docker push ghcr.io/unafraid/proxyless/client:latest

docker build -t ghcr.io/unafraid/proxyless/server:latest -f cmd/server/Dockerfile .
docker push ghcr.io/unafraid/proxyless/server:latest
