#!/bin/bash

NAMESPACE="proxyless"
STRICT="false"

while [[ $# -gt 0 ]]; do
  case $1 in
    -n|--namespace)
      NAMESPACE="$2"
      shift # past argument
      shift # past value
      ;;
    -s|--strict)
      STRICT="true"
      shift # past argument
      ;;
    *)
      echo "Invalid option: $1"
      usage
  esac
done

kubectl get namespace | grep -q "^$NAMESPACE " || kubectl create namespace "$NAMESPACE"
kubectl label namespace "$NAMESPACE" istio-injection=enabled

if [[ "$STRICT" == "true" ]]; then
  kubectl -n "$NAMESPACE" apply -f manifests/destination_rules.yaml
  kubectl -n "$NAMESPACE" apply -f manifests/peer_authorization.yaml
fi

helm upgrade --install -n "$NAMESPACE" proxyless-server charts/server
helm upgrade --install -n "$NAMESPACE" proxyless-client charts/client
