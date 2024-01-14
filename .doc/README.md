# References

1. https://kubernetes.io/docs/concepts/overview/working-with-objects/
2. https://kubernetes.io/docs/concepts/architecture/controller/
3. https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
4. https://sdk.operatorframework.io/docs/installation/
5. https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/

# Run & Test & Clean

## Localy

```shell
make install run

kubectl apply -f config/samples/

make uninstall
```

## Docker

```shell
./.script/kind-with-registry.sh

# Docker Hub
make docker-build docker-push IMG=smhmayboudi/custom-kubernetes-controller:latest
make deploy IMG=smhmayboudi/custom-kubernetes-controller:latest

# Local Registry
make docker-build docker-push IMG=localhost:5001/smhmayboudi/custom-kubernetes-controller:latest
make deploy IMG=http://kind-registry:5000/smhmayboudi/custom-kubernetes-controller:latest

kubectl apply -f config/samples/

make undeploy
```

## OLM

```shell
operator-sdk olm install
make bundle IMG=smhmayboudi/custom-kubernetes-controller:v0.0.1
make bundle-build bundle-push BUNDLE_IMG=smhmayboudi/custom-kubernetes-controller-bundle:v0.0.1
operator-sdk run bundle smhmayboudi/custom-kubernetes-controller-bundle:v0.0.1

kubectl apply -f config/samples/

operator-sdk cleanup custom-kubernetes-controller
```
