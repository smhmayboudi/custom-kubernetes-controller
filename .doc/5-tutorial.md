# Tutorial

## Reference

https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/

```shell
mkdir -p $HOME/Developer/custom-kubernetes-controller
cd $HOME/Developer/custom-kubernetes-controller
# we'll use a domain of interview.com
# so all API groups will be <group>.interview.com
operator-sdk init --domain interview.com --repo github.com/smhmayboudi/custom-kubernetes-controller
```

### MacOS

https://kubebuilder.io/plugins/available-plugins

```shell
mkdir -p $HOME/Developer/custom-kubernetes-controller
cd $HOME/Developer/custom-kubernetes-controller
# we'll use a domain of interview.com
# so all API groups will be <group>.interview.com
operator-sdk init --domain interview.com --repo github.com/smhmayboudi/custom-kubernetes-controller --plugins go/v4
```

## Manager

https://book.kubebuilder.io/cronjob-tutorial/empty-main.html

https://sdk.operatorframework.io/docs/building-operators/golang/operator-scope/

## Question

https://book.kubebuilder.io/cronjob-tutorial/empty-main.html

```go
    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme: scheme,
        Cache: cache.Options{
            DefaultNamespaces: map[string]cache.Config{
                namespace: {},
            },
        },
        Metrics: server.Options{
            BindAddress: metricsAddr,
        },
        WebhookServer:          webhook.NewServer(webhook.Options{Port: 9443}),
        HealthProbeBindAddress: probeAddr,
        LeaderElection:         enableLeaderElection,
        LeaderElectionID:       "80807133.tutorial.kubebuilder.io",
    })
```

The above example will change the scope of your project to a single Namespace. In this scenario, it is also suggested to restrict the provided authorization to this namespace by replacing the default ClusterRole and ClusterRoleBinding to Role and RoleBinding respectively. For further information see the Kubernetes documentation about Using [RBAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/).


```go
    var namespaces []string // List of Namespaces
    defaultNamespaces := make(map[string]cache.Config)

    for _, ns := range namespaces {
        defaultNamespaces[ns] = cache.Config{}
    }

    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme: scheme,
        Cache: cache.Options{
            DefaultNamespaces: defaultNamespaces,
        },
        Metrics: server.Options{
            BindAddress: metricsAddr,
        },
        WebhookServer:          webhook.NewServer(webhook.Options{Port: 9443}),
        HealthProbeBindAddress: probeAddr,
        LeaderElection:         enableLeaderElection,
        LeaderElectionID:       "80807133.tutorial.kubebuilder.io",
    })
```

Also, it is possible to use the [DefaultNamespaces](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/cache#Options) from cache.Options{} to cache objects in a specific set of namespaces. For further information see [cache.Options{}](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/cache#Options)

## Create a new API and Controller

[Deploy Image plugin](https://book.kubebuilder.io/plugins/deploy-image-plugin-v1-alpha.html)
[Single Group to Multi-Group](https://book.kubebuilder.io/migration/multi-group.html)
[Understanding Kubernetes APIs](https://book.kubebuilder.io/cronjob-tutorial/gvks.html)
[Controller Runtime](https://github.com/kubernetes-sigs/controller-runtime)

```shell
operator-sdk create api --group cache --version v1alpha1 --kind Memcached --resource --controller
# OR
operator-sdk create api --group cache --version v1alpha1 --kind Memcached --plugins "deploy-image/v1-alpha" --image=kind-registry:5000/library/memcached:1.6.23-alpine3.19 --image-container-command="memcached,-m=64,-o,modern,-v" --run-as-user="1001"
```

In general, it’s recommended to have one controller responsible for managing each API created for the project to properly follow the design goals set by [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime).

## Define the API

modify memcached_types

```shell
make generate
```

## Generating CRD manifests

```shell
make manifests
```

## Docker

```shell
./.script/kind-with-registry.sh
```

```shell
docker tag memcached:1.6.23-alpine3.19 localhost:5001/library/memcached:1.6.23-alpine3.19
docker push localhost:5001/library/memcached:1.6.23-alpine3.19
```

```shell
make docker-build docker-push IMG=localhost:5001/library/controller:latest
# OR
make docker-build docker-push
docker tag controller:latest localhost:5001/library/controller:latest
docker push localhost:5001/library/controller:latest
```

## Run the Operator

### Run locally outside the cluster

```shell
make install run MEMCACHED_IMAGE="memcached:1.6.23-alpine3.19"
```

```shell
export
make deploy IMG=kind-registry:5000/library/controller:latest
```

```shell
kubectl -n custom-kubernetes-controller-system get deployment
kubectl -n custom-kubernetes-controller-system describe deployment/custom-kubernetes-controller-controller-manager
kubectl -n custom-kubernetes-controller-system describe pod custom-kubernetes-controller-controller-manager
kubectl -n custom-kubernetes-controller-system logs deployment/custom-kubernetes-controller-controller-manager

kubectl get pods --all-namespaces
kubectl get deployments --all-namespaces
kubectl get services --all-namespaces
```

```shell
kubectl delete namespace custom-kubernetes-controller-system
```
