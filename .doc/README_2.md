# Run & Test & Clean

## Preparing

### Kind

Installation, and Test
NOTE: kind-registry:5000

```shell
./.script/kind-with-registry.sh
curl -v http://localhost:5000/v2/_catalog
```

### Minikube

Installation, Running, Start Local Registry, and Test
NOTE: registry

```shell
brew install minikube
minikube start
minikube addons enable registry
curl -v http://localhost:61655/v2/_catalog
```

If you need to change the default port forwarding.

```shell
kubectl port-forward --namespace kube-system service/registry 5000:80
docker run --rm -it --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:host.docker.internal:5000"

```

## Docker

The script provided at [this link](https://kind.sigs.k8s.io/docs/user/local-registry/).

```shell
# Docker Hub
make docker-build docker-push IMG=smhmayboudi/custom-kubernetes-controller:v0.0.1
make deploy IMG=smhmayboudi/custom-kubernetes-controller:v0.0.1

# Local Registry
make docker-build docker-push IMG=localhost:5001/smhmayboudi/custom-kubernetes-controller:v0.0.1
make deploy IMG=kind-registry:5000/smhmayboudi/custom-kubernetes-controller:v0.0.1

kubectl apply -f config/samples/

make undeploy
```

## OLM

```shell
operator-sdk olm install

# Docker Hub
make bundle IMG=smhmayboudi/custom-kubernetes-controller:v0.0.1
make bundle-build bundle-push BUNDLE_IMG=smhmayboudi/custom-kubernetes-controller-bundle:v0.0.1
operator-sdk run bundle smhmayboudi/custom-kubernetes-controller-bundle:v0.0.1

# Local Registry
make bundle IMG=localhost:5001/smhmayboudi/custom-kubernetes-controller:v0.0.1
make bundle-build bundle-push BUNDLE_IMG=localhost:5001/smhmayboudi/custom-kubernetes-controller-bundle:v0.0.1
operator-sdk run bundle kind-registry:5000/smhmayboudi/custom-kubernetes-controller-bundle:v0.0.1

kubectl apply -f config/samples/

operator-sdk cleanup custom-kubernetes-controller
```

## Direct

```shell
# Docker Hub
make deploy IMG="smhmayboudi/custom-kubernetes-controller:v0.0.1"

# Local Registry
make deploy IMG="kind-registry:5000/smhmayboudi/custom-kubernetes-controller:v0.0.1"

kubectl apply -f config/samples/

make undeploy
```

## locally

```shell
make install run

kubectl apply -f config/samples/

make uninstall
```

## Check the Deployments

```shell
kubectl get deployments
kubectl logs deployments/dummy1
kubectl describe deployments/dummy1
```

# Shell

```shell
# If shell completion is not already enabled in your environment you will need
# to enable it. You can execute the following once:
$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ operator-sdk completion zsh > "${fpath[1]}/_operator-sdk"

# You will need to start a new shell for this setup to take effect.
```

# OLM

## OLM Integration Bundle Quickstart

Export environment variables

```shell
export CHANNELS=development
export DEFAULT_CHANNEL=development
export USERNAME=smhmayboudi
export VERSION=0.0.1
export IMAGE_TAG_BASE=docker.io/$USERNAME/custom-kubernetes-controller
export IMG=$IMAGE_TAG_BASE:v$VERSION
# export BUNDLE_IMG=docker.io/$USERNAME/custom-kubernetes-controller-bundle:v$VERSION
# export CATALOG_IMG=docker.io/$USERNAME/custom-kubernetes-controller-catalog:v$VERSION
```

Create a bundle from the root directory of your project

```shell
make bundle
```

Build and push the bundle image

```shell
make bundle-build bundle-push
```

Validate the bundle

```shell
operator-sdk bundle validate
```

Install the bundle with OLM

```shell
operator-sdk run bundle
```

Deploying bundles in production

```shell
make catalog-build catalog-push
```

All

```shell
make bundle-build bundle-push catalog-build catalog-push
```

## OLM and Bundle CLI Overview

OLM installation

```shell
operator-sdk olm install
operator-sdk olm status
operator-sdk olm uninstall
```

Manifests and metadata

```shell
operator-sdk generate kustomize manifests
```

Bundles

```shell
make bundle
make bundle-build
operator-sdk run bundle
```

## Generating Manifests and Metadata



# Extra

```shell
operator-sdk generate kustomize manifests

operator-sdk run bundle localhost:5001/smhmayboudi/custom-kubernetes-controller-bundle:v0.0.1 --skip-tls --skip-tls-verify

operator-sdk bundle validate --list-optional
operator-sdk bundle validate ./bundle --select-optional suite=operatorframework 
operator-sdk bundle validate ./bundle --select-optional suite=operatorframework --optional-values=k8s-version=1.29.0
operator-sdk bundle validate ./bundle --select-optional name=community
operator-sdk bundle validate ./bundle --select-optional name=community --optional-values=image-path=bundle.Dockerfile
operator-sdk bundle validate ./bundle --select-optional name=multiarch
```

```shell
kubectl get pods --all-namespaces
kubectl exec -it <your_pod> -- /bin/bash

kubectl get configmaps
```

```shell
kubectl get catalogsource -n olm
kubectl get catalogsource  operatorhubio-catalog -n olm -o yaml

kubectl get csv -n olm
kubectl get csv packageserver -n olm -o yaml

kubectl get subscriptions -n olm

kubectl get installplans -n olm
kubectl get installplans install-xxxxx -n operators -o yaml

kubectl get operators operator-application.operators -n olm -o yaml
```

```shell
kubebuilder edit --plugins grafana/v1-alpha
kubebuilder create api \
    --group="cache" \
    --image-container-command="memcached,-m=64,-o,modern,-v" \
    --image-container-port="11211" \
    --image="memcached:1.6.23-alpine3.19" \
    --kind="Memcached" \
    --plugins="deploy-image/v1-alpha" \
    --run-as-user="1001" \
    --version="v1alpha1"

operator-sdk create api \
    --group="cache" \
    --image-container-command="memcached,-m=64,-o,modern,-v" \
    --image-container-port="11211" \
    --image="memcached:1.6.23-alpine3.19" \
    --kind="Memcached2" \
    --plugins="deploy-image/v1-alpha" \
    --run-as-user="1001" \
    --version="v2alpha1"
```

# Tools

## Development

- https://github.com/kubernetes-sigs/controller-runtime
- https://github.com/kubernetes-sigs/controller-tools
- https://github.com/kubernetes-sigs/kubebuilder
- https://github.com/kubernetes-sigs/kustomize
- https://github.com/kudobuilder/kuttl
- https://github.com/operator-framework/operator-registry
- https://github.com/operator-framework/operator-sdk


## Management

- https://github.com/ahmetb/kubectx
- https://github.com/danielfoehrKn/kubeswitch
- https://github.com/jonmosco/kube-ps1
- https://github.com/undistro/zora
- https://github.com/wagoodman/dive
