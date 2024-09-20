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

helm uninstall -n "$NAMESPACE" proxyless-client
helm uninstall -n "$NAMESPACE" proxyless-server

if [[ "$STRICT" == "true" ]]; then
  kubectl -n "$NAMESPACE" delete -f manifests/destination_rules.yaml
  kubectl -n "$NAMESPACE" delete -f manifests/peer_authorization.yaml
fi

kubectl delete namespace "$NAMESPACE"
