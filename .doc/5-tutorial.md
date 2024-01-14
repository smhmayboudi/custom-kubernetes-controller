# Tutorial

## Reference

https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/

```shell
mkdir -p $HOME/Developer/custom-kubernetes-controller
cd $HOME/Developer/custom-kubernetes-controller
# we'll use a domain of smhmayboudi.github.io
# so all API groups will be <group>.smhmayboudi.github.io
operator-sdk init --domain=smhmayboudi.github.io --repo=github.com/smhmayboudi/custom-kubernetes-controller
```

### MacOS

https://kubebuilder.io/plugins/available-plugins

```shell
mkdir -p $HOME/Developer/custom-kubernetes-controller
cd $HOME/Developer/custom-kubernetes-controller
# we'll use a domain of smhmayboudi.github.io
# so all API groups will be <group>.smhmayboudi.github.io
operator-sdk init --domain=smhmayboudi.github.io --repo=github.com/smhmayboudi/custom-kubernetes-controller --plugins=go/v4
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

```shell
operator-sdk create api --group=cache --version=v1alpha1 --kind=Memcached --resource --controller
```
