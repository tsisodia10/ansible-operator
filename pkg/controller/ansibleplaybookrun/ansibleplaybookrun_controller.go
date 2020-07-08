package ansibleplaybookrun

import (
	"context"

	"github.com/ManageIQ/manageiq-pods/manageiq-operator/pkg/helpers/tlstools"
	ansiblev1alpha1 "github.com/ansible-operator/pkg/apis/ansible/v1alpha1"
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
	ap := &ansiblev1.AnsiblePlaybook{}
	apr := &ansiblev1.AnsiblePlaybookRun{}

	err := r.client.Get(context.TODO(), request.NamespacedName, ap)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	ap.Initialize()
	apr.Initialize()

	if e := r.generateSecrets(ap); e != nil {
		return reconcile.Result{}, e
	}

	if e := r.generateOrchestratorResources(apr); e != nil {
		return reconcile.Result{}, e
	}

	return reconcile.Result{}, nil

	//--------------------------------------------------------------------------------------------------------------------------------
	// // Define a new Job object
	// job := AnsiblePlaybookRunJob(apr)

	// // Set AnsiblePlaybookRun instance as the owner and controller
	// if err := controllerutil.SetControllerReference(apr, job, r.scheme); err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Check if this Pod already exists
	// found := &corev1.Job{}
	// err = r.client.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	reqLogger.Info("Creating a new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
	// 	err = r.client.Create(context.TODO(), job)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}

	// 	// Job created successfully - don't requeue
	// 	return reconcile.Result{}, nil
	// } else if err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Job already exists - don't requeue
	// reqLogger.Info("Skip reconcile: Job already exists", "Job.Namespace", found.Namespace, "Job.Name", found.Name)
	// return reconcile.Result{}, nil
	//-------------------------------------------------------------------------------------------------------------------------------
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newAnsiblePlaybookRunJob(cr *ansiblev1alpha1.AnsiblePlaybookRun) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

func (r *ReconcileAnsiblePlaybookRun) generateSecrets(cr *ansiblev1alpha1.AnsiblePlaybook) error {
	appSecret := tools.AppSecret(cr)
	if err := r.createk8sResIfNotExist(cr, appSecret, &corev1.Secret{}); err != nil {
		return err
	}

	tlsSecret, err := ap.TLSSecret(cr)
	if err != nil {
		return err
	}

	if err := r.createk8sResIfNotExist(cr, tlsSecret, &corev1.Secret{}); err != nil {
		return err
	}

	return nil
}

func (r *ReconcileAnsiblePlaybookRun) createk8sResIfNotExist(cr *ansiblev1alpha1.AnsiblePlaybook, res, restype metav1.Object) error {
	reqLogger := logger.WithValues("task: ", "create resource")
	if err := controllerutil.SetControllerReference(cr, res, r.scheme); err != nil {
		return err
	}
	resClient := res.(runtime.Object)
	resTypeClient := restype.(runtime.Object)
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: res.GetName(), Namespace: res.GetNamespace()}, resTypeClient.(runtime.Object))
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating ", "Resource.Namespace", res.GetNamespace(), "Resource.Name", res.GetName())
		if err = r.client.Create(context.TODO(), resClient); err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func TLSSecret(cr *ansiblev1alpha1.AnsiblePlaybook) (*corev1.Secret, error) {
	labels := map[string]string{
		"app": cr.Spec.AppName,
	}

	crt, key, err := tlstools.GenerateCrt("server")
	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"tls.crt": string(crt),
		"tls.key": string(key),
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tlsSecretName(cr),
			Namespace: cr.ObjectMeta.Namespace,
			Labels:    labels,
		},
		StringData: data,
		Type:       "kubernetes.io/tls",
	}
	return secret, nil
}

func tlsSecretName(cr *ansiblev1alpha1.AnsiblePlaybook) string {
	secretName := "tls-secret"
	if cr.Spec.TLSSecret != "" {
		secretName = cr.Spec.TLSSecret
	}

	return secretName
}

func (r *ReconcileAnsiblePlaybookRun) generateOrchestratorResources(cr *ansiblev1alpha1.AnsiblePlaybook) error {
	orchestratorJob, mutateFunc, err := tools.OrchestratorJob(cr, r.scheme)
	if err != nil {
		return err
	}

	if result, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, orchestratorDeployment, mutateFunc); err != nil {
		return err
	} else {
		logger.Info("Job has been reconciled", "component", "orchestrator", "result", result)
	}

	return nil
}
