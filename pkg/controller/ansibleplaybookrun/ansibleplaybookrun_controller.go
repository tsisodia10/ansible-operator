package ansibleplaybookrun

import (
	"context"

	ansiblev1alpha1 "github.com/ansible-operator/pkg/apis/ansible/v1alpha1"
	batch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_ansibleplaybookrun")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new AnsiblePlaybookRun Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileAnsiblePlaybookRun{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("ansibleplaybookrun-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AnsiblePlaybookRun
	err = c.Watch(&source.Kind{Type: &ansiblev1alpha1.AnsiblePlaybookRun{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner AnsiblePlaybookRun
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &ansiblev1alpha1.AnsiblePlaybookRun{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileAnsiblePlaybookRun implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAnsiblePlaybookRun{}

// ReconcileAnsiblePlaybookRun reconciles a AnsiblePlaybookRun object
type ReconcileAnsiblePlaybookRun struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AnsiblePlaybookRun object and makes changes based on the state read
// and what is in the AnsiblePlaybookRun.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAnsiblePlaybookRun) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AnsiblePlaybookRun")

	// Fetch the AnsiblePlaybookRun instance
	ap := &ansiblev1alpha1.AnsiblePlaybook{}
	apr := &ansiblev1alpha1.AnsiblePlaybookRun{}

	err := r.client.Get(context.TODO(), request.NamespacedName, apr)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue

			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		apr.Status.Status = "Unavailable"
		err = r.client.Status().Update(context.TODO(), apr)
		return reconcile.Result{}, err
	}

	err = r.client.Get(context.TODO(), apr.Spec.AnsiblePlaybookRef, ap)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue

			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		apr.Status.Status = "Unavailable"
		err = r.client.Status().Update(context.TODO(), apr)
		return reconcile.Result{}, err
	}

	// Define a new Job object
	job := BuildJobSpec(apr, ap)

	// Set AnsiblePlaybookRun instance as the owner and controller
	if err := controllerutil.SetControllerReference(apr, job, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Job already exists
	found := &batch.Job{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: apr.Name, Namespace: apr.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Job", "Job.Namespace", apr.Namespace, "Job.Name", apr.Name)
		err = r.client.Create(context.TODO(), job)
		if err != nil {
			apr.Status.Status = "Finished"
			reqLogger.Info("Status: ", apr.Status)
			return reconcile.Result{}, nil
		}

		// Job created successfully - don't requeue
		reqLogger.Info("Success: Job created")
		return reconcile.Result{}, nil
	}

	// Job already exists - requeue
	reqLogger.Info("Reconcile requeue: Job already exists", "Job.Namespace", found.Namespace, "Job.Name", found.Name)
	return reconcile.Result{Requeue: true}, nil

	// // Check if this Job already exists
	// found := &batch.Job{}
	// err = r.client.Get(context.TODO(), types.NamespacedName{Name: apr.Name, Namespace: apr.Namespace}, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	reqLogger.Info("Creating a new Job", "Job.Namespace", apr.Namespace, "Job.Name", apr.Name)
	// 	err = r.client.Create(context.TODO(), job)
	// 	if err != nil {
	// 		reqLogger.Info("Failed creating job")
	// 		return reconcile.Result{}, err
	// 	}

	// 	// Job created successfully - don't requeue
	// 	reqLogger.Info("Success: Job created")
	// 	return reconcile.Result{}, nil
	// }

}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func BuildJobSpec(cr *ansiblev1alpha1.AnsiblePlaybookRun, cr1 *ansiblev1alpha1.AnsiblePlaybook) *batch.Job {
	labels := map[string]string{
		"app": cr.Name,
	}

	return &batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-job",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: batch.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "ansible-runner",
							Image:   "quay.io/tsisodia/ansible-runner",
							Command: []string{"sleep", "3600"},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "ANSIBLEPLAYBOOK_NAME",
									Value: cr1.Spec.PlaybookName,
								},
								corev1.EnvVar{
									Name:  "REPOSITORY_TYPE",
									Value: cr1.Spec.RepositoryType,
								},
								corev1.EnvVar{
									Name:  "REPOSITORY_URL",
									Value: cr1.Spec.RepositoryURL,
								},
								corev1.EnvVar{
									Name:  "INVENTORY",
									Value: cr.Spec.Inventory,
								},
								corev1.EnvVar{
									Name:  "HOST_CREDENTIAL",
									Value: cr.Spec.HostCredential,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "extravars-volume",
									MountPath: "/runner/env/extravars",
								},
								corev1.VolumeMount{
									Name:      "password-volume",
									MountPath: "/runner/env/password",
								},
								corev1.VolumeMount{
									Name:      "sshkey-volume",
									MountPath: "/runner/env/ssh_key",
								},
								corev1.VolumeMount{
									Name:      "inventory-volume",
									MountPath: "/runner/inventory/hosts",
								},
								corev1.VolumeMount{
									Name:      "projectvars-volume",
									MountPath: "/runner/project/roles/testrole/vars",
								},
								corev1.VolumeMount{
									Name:      "projectmeta-volume",
									MountPath: "/runner/project/roles/testrole/meta",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: "extravars-volume",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/data",
								},
							},
						},
						corev1.Volume{
							Name: "password-volume",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/data",
								},
							},
						},
						corev1.Volume{
							Name: "sshkey-volume",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/data",
								},
							},
						},
						corev1.Volume{
							Name: "inventory-volume",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/data",
								},
							},
						},
						corev1.Volume{
							Name: "projectvars-volume",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/data",
								},
							},
						},
						corev1.Volume{
							Name: "projectmeta-volume",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/data",
								},
							},
						},
					},
				},
			},
		},
	}
}
