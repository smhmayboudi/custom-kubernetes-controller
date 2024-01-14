# Objects In Kubernetes

## Reference

https://kubernetes.io/docs/concepts/overview/working-with-objects/

## Kubernetes Object Management

```shell
kubectl create deployment nginx --image nginx
kubectl create -f nginx.yaml # create a object
kubectl delete -f nginx.yaml -f redis.yaml # delete a object
kubectl replace -f nginx.yaml # update a object
kubectl diff -f configs/ && kubectl apply -f configs/ # make a patch & apply it
```

## Object Names and IDs

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-demo # RFC1123 & RFC1035(start with char), len(name) <= 63
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```

## Labels and Selectors

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: label-demo
  labels:
    environment: production # [a-z0-9A-Z], (-, _, .), len(key) <= 63
    app: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```

> Note: For some API types, such as ReplicaSets, the label selectors of two instances must not overlap within a namespace, or the controller can see that as conflicting instructions and fail to determine how many replicas should be present.

> Caution: For both equality-based and set-based conditions there is no logical OR (||) operator. Ensure your filter statements are structured accordingly.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: cuda-test
spec:
  containers:
    - name: cuda-test
      image: "registry.k8s.io/cuda-vector-add:v0.1"
      resources:
        limits:
          nvidia.com/gpu: 1
  nodeSelector:
    accelerator: nvidia-tesla-p100 # accelerator==nvidia-tesla-p100, in,notin and exists
```
```yaml
labelSelector=environment%3Dproduction,tier%3Dfrontend
```

```shell
kubectl get pods -l environment=production,tier=frontend
```

```yaml
?labelSelector=environment+in+%28production%2Cqa%29%2Ctier+in+%28frontend%29
```

```shell
kubectl get pods -l 'environment in (production),tier in (frontend)'
```

examples: https://github.com/kubernetes/examples/tree/master/guestbook/

```shell
kubectl label pods -l app=nginx tier=fe # select all ngnix and add tier to fe
kubectl get pods -l app=nginx -L tier # (--label-columns) to see it
```

## Namespaces

- default
- kube-node-lease: [Lease](https://kubernetes.io/docs/concepts/architecture/leases/) Objects to send [heartbeats](https://kubernetes.io/docs/concepts/architecture/nodes/#heartbeats)
- kube-public
- kube-system

```shell
kubectl get namespace
```

DNS: <service-name>.<namespace-name>.svc.cluster.local

> By creating namespaces with the same name as [public top-level domains](https://data.iana.org/TLD/tlds-alpha-by-domain.txt), Services in these namespaces can have short DNS names that overlap with public DNS records. Workloads from any namespace performing a DNS lookup without a [trailing dot](https://datatracker.ietf.org/doc/html/rfc1034#page-8) will be redirected to those services, taking precedence over public DNS.
To mitigate this, limit privileges for creating namespaces to trusted users. If required, you could additionally configure third-party security controls, such as [admission webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/), to block creating any namespace with the name of [public TLDs](https://data.iana.org/TLD/tlds-alpha-by-domain.txt).

namespace resources are not themselves in a namespace. And low-level resources, such as [nodes](https://kubernetes.io/docs/concepts/architecture/nodes/) and [persistentVolumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/), are not in any namespace.
```shell
kubectl api-resources --namespaced=true # in a namespace
kubectl api-resources --namespaced=false # not in a namespace
```

`kubernetes.io/metadata.name` is atomic labelling which is set by the control plane.

## Annotations

You can use Kubernetes annotations to attach arbitrary non-identifying metadata to [objects](https://kubernetes.io/docs/concepts/overview/working-with-objects/#kubernetes-objects). Clients such as tools and libraries can retrieve this metadata.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: annotations-demo
  annotations:
    imageregistry: "https://hub.docker.com/" # it shows the imageregistry
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```

# Field Selectors

```shell
kubectl get pods --field-selector metadata.name=my-service
kubectl get pods --field-selector metadata.namespace!=default
kubectl get pods --field-selector status.phase=Pending
```

```shell
kubectl get pods,statefulsets,services --all-namespaces --field-selector=status.phase!=Running,spec.restartPolicy=Always
```

# Finalizers

Finalizers are namespaced keys that tell Kubernetes to wait until specific conditions are met before it fully deletes resources marked for deletion. Finalizers alert controllers to clean up resources the deleted object owned.

When you tell Kubernetes to delete an object that has finalizers specified for it, the Kubernetes API marks the object for deletion by populating `.metadata.deletionTimestamp`, and returns a 202 status code (HTTP "Accepted"). The target object remains in a terminating state while the control plane, or other components, take the actions defined by the finalizers. After these actions are complete, the controller removes the relevant finalizers from the target object. When the `metadata.finalizers` field is empty, Kubernetes considers the deletion complete and deletes the object.

> When you DELETE an object, Kubernetes adds the deletion timestamp for that object and then immediately starts to restrict changes to the .metadata.finalizers field for the object that is now pending deletion. You can remove existing finalizers (deleting an entry from the finalizers list) but you cannot add a new finalizer. You also cannot modify the deletionTimestamp for an object once it is set.
After the deletion is requested, you can not resurrect this object. The only way is to delete it and make a new similar object.

> Note: In cases where objects are stuck in a deleting state, avoid manually removing finalizers to allow deletion to continue. Finalizers are usually added to resources for a reason, so forcefully removing them can lead to issues in your cluster. This should only be done when the purpose of the finalizer is understood and is accomplished in another way (for example, manually cleaning up some dependent object).

# Owners and Dependents

A valid owner reference (`metadata.ownerReferences`) consists of the object name and a UID within the same namespace as the dependent object.
`ownerReferences.blockOwnerDeletion`

> Note: Cross-namespace owner references are disallowed by design. Namespaced dependents can specify cluster-scoped or namespaced owners. A namespaced owner must exist in the same namespace as the dependent. If it does not, the owner reference is treated as absent, and the dependent is subject to deletion once all owners are verified absent.
Cluster-scoped dependents can only specify cluster-scoped owners. In v1.20+, if a cluster-scoped dependent specifies a namespaced kind as an owner, it is treated as having an unresolvable owner reference, and is not able to be garbage collected.
In v1.20+, if the garbage collector detects an invalid cross-namespace ownerReference, or a cluster-scoped dependent with an ownerReference referencing a namespaced kind, a warning Event with a reason of OwnerRefInvalidNamespace and an involvedObject of the invalid dependent is reported. You can check for that kind of Event by running `kubectl get events -A --field-selector=reason=OwnerRefInvalidNamespace`.

# Recommended Labels

Shared labels and annotations share a common prefix: `app.kubernetes.io`. Labels without a prefix are private to users. The shared prefix ensures that shared labels do not interfere with custom user labels.

```yaml
# This is an excerpt StatefulSet object
apiVersion: apps/v1
kind: StatefulSet
metadata:
  labels:
    app.kubernetes.io/name: mysql
    app.kubernetes.io/instance: mysql-abcxzy # every instance of an application must have a unique name.
    app.kubernetes.io/version: "5.7.21"
    app.kubernetes.io/component: database
    app.kubernetes.io/part-of: wordpress
    app.kubernetes.io/managed-by: helm
```

```yaml
# This is an excerpt Deployment object, to oversee the pods running the application itself
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app.kubernetes.io/name: wordpress
    app.kubernetes.io/instance: wordpress-abcxzy # every instance of an application must have a unique name.
    app.kubernetes.io/version: "4.9.4"
    app.kubernetes.io/component: server
    app.kubernetes.io/part-of: wordpress
    app.kubernetes.io/managed-by: helm
```

```yaml
# This is an excerpt Service object, to expose the application
apiVersion: apps/v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/name: wordpress
    app.kubernetes.io/instance: wordpress-abcxzy # every instance of an application must have a unique name.
    app.kubernetes.io/version: "4.9.4"
    app.kubernetes.io/component: server
    app.kubernetes.io/part-of: wordpress
    app.kubernetes.io/managed-by: helm
```
