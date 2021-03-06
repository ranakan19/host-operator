package hostoperatorconfig

import (
	"context"

	toolchainv1alpha1 "github.com/codeready-toolchain/api/pkg/apis/toolchain/v1alpha1"
	"github.com/codeready-toolchain/host-operator/pkg/configuration"

	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_hostoperatorconfig")

// Add creates a new HostOperatorConfig Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, _ *configuration.Config) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileHostOperatorConfig{client: mgr.GetClient()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("hostoperatorconfig-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource HostOperatorConfig
	return c.Watch(&source.Kind{Type: &toolchainv1alpha1.HostOperatorConfig{}}, &handler.EnqueueRequestForObject{})
}

// blank assignment to verify that ReconcileHostOperatorConfig implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileHostOperatorConfig{}

// ReconcileHostOperatorConfig reconciles a HostOperatorConfig object
type ReconcileHostOperatorConfig struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
}

// Reconcile reads that state of the cluster for a HostOperatorConfig object and makes changes based on the state read
// and what is in the HostOperatorConfig.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileHostOperatorConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling HostOperatorConfig")

	// Fetch the HostOperatorConfig instance
	hostOperatorConfig := &toolchainv1alpha1.HostOperatorConfig{}
	err := r.client.Get(context.TODO(), request.NamespacedName, hostOperatorConfig)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Error(err, "it looks like the HostOperatorConfig resource with the name 'config' was removed - the cache will use the latest version of the resource")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	updateConfig(hostOperatorConfig)
	return reconcile.Result{}, nil
}
