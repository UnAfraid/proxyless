#!/bin/bash

CONTEXT="kind-proxyless-poc"
NAMESPACE="proxyless"

while [[ $# -gt 0 ]]; do
  case $1 in
    --context)
      CONTEXT="$2"
      shift # past argument
      shift # past value
      ;;
    -|--namespace)
      NAMESPACE="$2"
      shift # past argument
      shift # past value
      ;;
    *)
      echo "Invalid option: $1"
      usage
  esac
done

kubectl --context "$CONTEXT" get namespace | grep -q "^$NAMESPACE " || kubectl --context "$CONTEXT" create namespace "$NAMESPACE"
kubectl --context "$CONTEXT" label namespace "$NAMESPACE" istio-injection=enabled

helm --kube-context "$CONTEXT" upgrade --install -n "$NAMESPACE" proxyless-server charts/server
helm --kube-context "$CONTEXT" upgrade --install -n "$NAMESPACE" proxyless-client charts/client
