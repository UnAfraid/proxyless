#!/bin/bash

type kind >/dev/null 2>&1 || { echo >&2 "kind is not installed"; exit 1; }

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

kind delete cluster --name "${CONTEXT#kind-}"
