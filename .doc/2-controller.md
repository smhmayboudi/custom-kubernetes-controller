# Controller

## Reference

https://kubernetes.io/docs/concepts/architecture/controller/

In robotics and automation, a control loop is a non-terminating loop that regulates the state of a system.

LOOP
    Current State => Job Controller (via Calling API Server) => Desired State

> Note: There can be several controllers that create or update the same kind of object. Behind the scenes, Kubernetes controllers make sure that they only pay attention to the resources linked to their controlling resource.
For example, you can have Deployments and Jobs; these both create Pods. The Job controller does not delete the Pods that your Deployment created, because there is information (labels) the controllers can use to tell those Pods apart.

Kubernetes comes with a set of built-in controllers that run inside the kube-controller-manager. These built-in controllers provide important core behaviors.

The Deployment controller and Job controller are examples of controllers that come as part of Kubernetes itself ("built-in" controllers). Kubernetes lets you run a resilient control plane, so that if any of the built-in controllers were to fail, another part of the control plane will take over the work.

The most common way to deploy an operator is to add the Custom Resource Definition and its associated Controller to your cluster. The Controller will normally run outside of the control plane, much as you would run any containerized application. For example, you can run the controller in your cluster as a Deployment.

> Note: This section links to third party projects that provide functionality required by Kubernetes. The Kubernetes project authors aren't responsible for these projects, which are listed alphabetically. To add a project to this list, read the [content guide](https://kubernetes.io/docs/contribute/style/content-guide/#third-party-content) before submitting a change. [More information](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/#third-party-content-disclaimer).

> Note: [Operator SDK](https://github.com/operator-framework/operator-sdk/blob/v1.33.0/) uses the [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder/tree/v3.12.0) plugin feature to include non-Go operators e.g. operator-sdk's Ansible and Helm-based language Operators. To learn more see [how to create your own plugins](https://book.kubebuilder.io/plugins/creating-plugins.html).
