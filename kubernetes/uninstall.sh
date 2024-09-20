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

helm --kube-context "$CONTEXT" uninstall -n "$NAMESPACE" proxyless-client
helm --kube-context "$CONTEXT" uninstall -n "$NAMESPACE" proxyless-server

kubectl --context "$CONTEXT" delete namespace "$NAMESPACE"
