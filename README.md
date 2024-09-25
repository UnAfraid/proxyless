# proxyless gRPC in kubernetes with istio service mesh POC

This PoC aims to enable load balancing of gRPC services in kubernetes using istio without proxying the requests, but instead leverage the xDS support.

The main benefits of using proxyless load balancing is:
- Lower latency
- Less resources necessary per service in the mesh

## Prerequisites
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
- [istioctl](https://istio.io/latest/docs/setup/additional-setup/download-istio-release/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)


## Setup
### Create kind cluster and install istio
```shell
git clone https://github.com/UnAfraid/proxyless
cd kubernetes/cluster
./install.sh
```

### Install proxyless client and server PoC 
```shell
cd ../workload
./install.sh
```

### Inspect
Make sure all pods are running
```shell
kubectl --context kind-proxyless-poc -n proxyless get pods
```

Inspect logs of the client
```shell
kubectl --context kind-proxyless-poc -n proxyless logs -l app.kubernetes.io/name=proxyless-client
```

Inspect logs of the server
```shell
kubectl --context kind-proxyless-poc -n proxyless logs -l app.kubernetes.io/name=proxyless-server
```



Relevant links:
- https://cloud.google.com/service-mesh/docs/service-routing/proxyless-overview
- https://istio.io/latest/blog/2021/proxyless-grpc/
- https://github.com/grpc/grpc-go/tree/master/examples/features/xds
- https://events.istio.io/istiocon-2022/sessions/proxyless-grpc/
