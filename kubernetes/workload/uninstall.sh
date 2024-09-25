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

helm --kube-context "$CONTEXT" uninstall -n "$NAMESPACE" proxyless-client-1x
helm --kube-context "$CONTEXT" uninstall -n "$NAMESPACE" proxyless-client-2x
helm --kube-context "$CONTEXT" uninstall -n "$NAMESPACE" proxyless-client-4x

helm --kube-context "$CONTEXT" uninstall -n "$NAMESPACE" proxyless-server-1
helm --kube-context "$CONTEXT" uninstall -n "$NAMESPACE" proxyless-server-2
helm --kube-context "$CONTEXT" uninstall -n "$NAMESPACE" proxyless-server-3
helm --kube-context "$CONTEXT" uninstall -n "$NAMESPACE" proxyless-server-4

kubectl --context "$CONTEXT" delete namespace "$NAMESPACE"
