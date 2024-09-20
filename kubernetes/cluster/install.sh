#!/bin/bash

type kind >/dev/null 2>&1 || { echo >&2 "kind is not installed"; exit 1; }
type istioctl >/dev/null 2>&1 || { echo >&2 "istioctl is not installed"; exit 1; }

CONTEXT="kind-proxyless-poc"

while [[ $# -gt 0 ]]; do
  case $1 in
    --context)
      CONTEXT="$2"
      shift # past argument
      shift # past value
      ;;
    *)
      echo "Invalid option: $1"
      usage
  esac
done

kind create cluster --name "${CONTEXT#kind-}" --config manifests/cluster.yaml

echo "Waiting 30 seconds..."
sleep 30

echo "Installing istio operator..."
istioctl --context "$CONTEXT" operator init

echo "Configuring istio operator..."
kubectl --context "$CONTEXT" apply -f manifests/istio-operator.yaml
