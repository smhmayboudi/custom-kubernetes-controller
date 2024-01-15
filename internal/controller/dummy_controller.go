/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	interviewcomv1alpha1 "github.com/smhmayboudi/custom-kubernetes-controller/api/v1alpha1"
)

// DummyReconciler reconciles a Dummy object
type DummyReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=interview.com,resources=dummies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=interview.com,resources=dummies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=interview.com,resources=dummies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Dummy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *DummyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	const (
		dummyFinalizer   = "interview.com/finalizer"
		podStatusPending = "Pending"
		podStatusRunning = "Running"
	)

	// TODO(user): your logic here

	dummy := &interviewcomv1alpha1.Dummy{}

	err := r.Get(ctx, req.NamespacedName, dummy)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("dummy resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get dummy")
		return ctrl.Result{}, err
	}

	if dummy.Status.PodStatus == "" {
		dummy.Status.PodStatus = podStatusPending

		if err := r.Status().Update(ctx, dummy); err != nil {
			log.Error(err, "Failed to update Dummy status 1")
			return ctrl.Result{}, err
		}

		if err := r.Get(ctx, req.NamespacedName, dummy); err != nil {
			log.Error(err, "Failed to re-fetch memcached")
			return ctrl.Result{}, err
		}
	}

	log.Info(fmt.Sprintf("Name=%s Namespace=%s dummy.Spec=%v dummy.Status=%v", req.Name, req.Namespace, dummy.Spec, dummy.Status))

	// Let's add a finalizer. Then, we can define some operations which should
	// occurs before the custom resource to be deleted.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers
	if !controllerutil.ContainsFinalizer(dummy, dummyFinalizer) {
		log.Info("Adding Finalizer for Dummy")
		if ok := controllerutil.AddFinalizer(dummy, dummyFinalizer); !ok {
			log.Error(err, "Failed to add finalizer into the custom resource")
			return ctrl.Result{Requeue: true}, nil
		}

		if err = r.Update(ctx, dummy); err != nil {
			log.Error(err, "Failed to update custom resource to add finalizer")
			return ctrl.Result{}, err
		}
	}

	// Check if the Dummy instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isDummyMarkedToBeDeleted := dummy.GetDeletionTimestamp() != nil
	if isDummyMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(dummy, dummyFinalizer) {
			log.Info("Performing Finalizer Operations for Dummy before delete CR")

			// // Let's add here an status "Downgrade" to define that this resource begin its process to be terminated.
			// meta.SetStatusCondition(&dummy.Status.Conditions, metav1.Condition{Type: typeDegradedDummy,
			// 	Status: metav1.ConditionUnknown, Reason: "Finalizing",
			// 	Message: fmt.Sprintf("Performing finalizer operations for the custom resource: %s ", dummy.Name)})

			// dummy.Status.PodStatus = podStatusPending

			// if err := r.Status().Update(ctx, dummy); err != nil {
			// 	log.Error(err, "Failed to update Dummy status 2")
			// 	return ctrl.Result{}, err
			// }

			// Perform all operations required before remove the finalizer and allow
			// the Kubernetes API to remove the custom resource.
			r.doFinalizerOperationsForDummy(dummy)

			// TODO(user): If you add operations to the doFinalizerOperationsForDummy method
			// then you need to ensure that all worked fine before deleting and updating the Downgrade status
			// otherwise, you should requeue here.

			// Re-fetch the dummy Custom Resource before update the status
			// so that we have the latest state of the resource on the cluster and we will avoid
			// raise the issue "the object has been modified, please apply
			// your changes to the latest version and try again" which would re-trigger the reconciliation
			if err := r.Get(ctx, req.NamespacedName, dummy); err != nil {
				log.Error(err, "Failed to re-fetch dummy")
				return ctrl.Result{}, err
			}

			// meta.SetStatusCondition(&dummy.Status.Conditions, metav1.Condition{Type: typeDegradedDummy,
			// 	Status: metav1.ConditionTrue, Reason: "Finalizing",
			// 	Message: fmt.Sprintf("Finalizer operations for custom resource %s name were successfully accomplished", dummy.Name)})

			// dummy.Status.PodStatus = podStatusPending

			// if err := r.Status().Update(ctx, dummy); err != nil {
			// 	log.Error(err, "Failed to update Dummy status 3")
			// 	return ctrl.Result{}, err
			// }

			log.Info("Removing Finalizer for Dummy after successfully perform the operations")
			if ok := controllerutil.RemoveFinalizer(dummy, dummyFinalizer); !ok {
				log.Error(err, "Failed to remove finalizer for Dummy")
				return ctrl.Result{Requeue: true}, nil
			}

			if err := r.Update(ctx, dummy); err != nil {
				log.Error(err, "Failed to remove finalizer for Dummy")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Check if the deployment already exists, if not create a new one or nginx.
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: dummy.Name, Namespace: dummy.Namespace}, found)
	if err != nil && apierrors.IsNotFound(err) {
		// Define a new deployment
		dep, err := r.deploymentForNginx(dummy)
		if err != nil {
			log.Error(err, "Failed to define new Deployment resource for Dummy")

			// // The following implementation will update the status
			// meta.SetStatusCondition(&dummy.Status.Conditions, metav1.Condition{Type: typeAvailableDummy,
			// 	Status: metav1.ConditionFalse, Reason: "Reconciling",
			// 	Message: fmt.Sprintf("Failed to create Deployment for the custom resource (%s): (%s)", interviewcomv1alpha1.Name, err)})

			dummy.Status.PodStatus = podStatusPending

			if err := r.Status().Update(ctx, dummy); err != nil {
				log.Error(err, "Failed to update Dummy status 4")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, err
		}

		log.Info(
			"Creating a new Deployment",
			"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		if err = r.Create(ctx, dep); err != nil {
			log.Error(err, "Failed to create new Deployment",
				"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}

		dummy.Status.PodStatus = podStatusRunning

		if err := r.Status().Update(ctx, dummy); err != nil {
			log.Error(err, "Failed to update Dummy status 5")
			return ctrl.Result{}, err
		}

		// Deployment created successfully
		// We will requeue the reconciliation so that we can ensure the state
		// and move forward for the next operations
		return ctrl.Result{RequeueAfter: time.Second}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		// Let's return the error for the reconciliation be re-trigged again
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// doFinalizerOperationsForDummy will perform the required operations before delete the CR.
func (r *DummyReconciler) doFinalizerOperationsForDummy(cr *interviewcomv1alpha1.Dummy) {
	// TODO(user): Add the cleanup steps that the operator
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.

	// Note: It is not recommended to use finalizers with the purpose of delete resources which are
	// created and managed in the reconciliation. These ones, such as the Deployment created on this reconcile,
	// are defined as depended of the custom resource. See that we use the method ctrl.SetControllerReference.
	// to set the ownerRef which means that the Deployment will be deleted by the Kubernetes API.
	// More info: https://kubernetes.io/docs/tasks/administer-cluster/use-cascading-deletion/

	// The following implementation will raise an event
	r.Recorder.Event(cr, "Warning", "Deleting",
		fmt.Sprintf("Custom Resource %s is being deleted from the namespace %s",
			cr.Name,
			cr.Namespace))
}

// deploymentForNginx returns a Dummy Deployment object
func (r *DummyReconciler) deploymentForNginx(
	dummy *interviewcomv1alpha1.Dummy) (*appsv1.Deployment, error) {
	ls := labelsForDummy(dummy.Name)
	replicas := int32(1)

	// Get the Operand image
	image, err := imageForNginx()
	if err != nil {
		return nil, err
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dummy.Name,
			Namespace: dummy.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					// TODO(user): Uncomment the following code to configure the nodeAffinity expression
					// according to the platforms which are supported by your solution. It is considered
					// best practice to support multiple architectures. build your manager image using the
					// makefile target docker-buildx. Also, you can use docker manifest inspect <image>
					// to check what are the platforms supported.
					// More info: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#node-affinity
					//Affinity: &corev1.Affinity{
					//	NodeAffinity: &corev1.NodeAffinity{
					//		RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
					//			NodeSelectorTerms: []corev1.NodeSelectorTerm{
					//				{
					//					MatchExpressions: []corev1.NodeSelectorRequirement{
					//						{
					//							Key:      "kubernetes.io/arch",
					//							Operator: "In",
					//							Values:   []string{"amd64", "arm64", "ppc64le", "s390x"},
					//						},
					//						{
					//							Key:      "kubernetes.io/os",
					//							Operator: "In",
					//							Values:   []string{"linux"},
					//						},
					//					},
					//				},
					//			},
					//		},
					//	},
					//},
					// SecurityContext: &corev1.PodSecurityContext{
					// 	RunAsNonRoot: &[]bool{true}[0],
					// 	// IMPORTANT: seccomProfile was introduced with Kubernetes 1.19
					// 	// If you are looking for to produce solutions to be supported
					// 	// on lower versions you must remove this option.
					// 	SeccompProfile: &corev1.SeccompProfile{
					// 		Type: corev1.SeccompProfileTypeRuntimeDefault,
					// 	},
					// },
					Containers: []corev1.Container{{
						Image:           image,
						Name:            "ngnix",
						ImagePullPolicy: corev1.PullIfNotPresent,
						// Ensure restrictive context for the container
						// More info: https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted
						// SecurityContext: &corev1.SecurityContext{
						// 	// WARNING: Ensure that the image used defines an UserID in the Dockerfile
						// 	// otherwise the Pod will not run and will fail with "container has runAsNonRoot and image has non-numeric user"".
						// 	// If you want your workloads admitted in namespaces enforced with the restricted mode in OpenShift/OKD vendors
						// 	// then, you MUST ensure that the Dockerfile defines a User ID OR you MUST leave the "RunAsNonRoot" and
						// 	// "RunAsUser" fields empty.
						// 	RunAsNonRoot: &[]bool{true}[0],
						// 	// The dummy image does not use a non-zero numeric user as the default user.
						// 	// Due to RunAsNonRoot field being set to true, we need to force the user in the
						// 	// container to a non-zero numeric user. We do this using the RunAsUser field.
						// 	// However, if you are looking to provide solution for K8s vendors like OpenShift
						// 	// be aware that you cannot run under its restricted-v2 SCC if you set this value.
						// 	RunAsUser:                &[]int64{1001}[0],
						// 	AllowPrivilegeEscalation: &[]bool{false}[0],
						// 	Capabilities: &corev1.Capabilities{
						// 		Drop: []corev1.Capability{
						// 			"ALL",
						// 		},
						// 	},
						// },
						Ports: []corev1.ContainerPort{{
							ContainerPort: int32(80),
							Name:          "ngnix",
						}},
						// Command: []string{},
					}},
				},
			},
		},
	}

	// Set the ownerRef for the Deployment
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/
	if err := ctrl.SetControllerReference(dummy, dep, r.Scheme); err != nil {
		return nil, err
	}
	return dep, nil
}

// labelsForDummy returns the labels for selecting the resources
// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
func labelsForDummy(name string) map[string]string {
	var imageTag string
	image, err := imageForNginx()
	if err == nil {
		imageTag = strings.Split(image, ":")[1]
	}
	return map[string]string{
		"app.kubernetes.io/name":       "Dummy",
		"app.kubernetes.io/instance":   name,
		"app.kubernetes.io/version":    imageTag,
		"app.kubernetes.io/part-of":    "dummy-operator",
		"app.kubernetes.io/created-by": "controller-manager",
	}
}

// imageForNginx gets the Operand image which is managed by this controller
// from the NGINX_IMAGE environment variable defined in the config/manager/manager.yaml
func imageForNginx() (string, error) {
	var imageEnvVar = "NGINX_IMAGE"
	image, found := os.LookupEnv(imageEnvVar)
	if !found {
		return "", fmt.Errorf("Unable to find %s environment variable with the image", imageEnvVar)
	}
	return image, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DummyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&interviewcomv1alpha1.Dummy{}).
		Complete(r)
}
