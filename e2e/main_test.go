package e2e

import (
	"context"
	"fmt"
	"github.com/mmlt/gstconfig/controllers"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	//"github.com/go-logr/stdr"
	clusteropsv1 "github.com/mmlt/gstconfig/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	//"log"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	//logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sync"
	"testing"
)

// TestMain instantiates the following vars for usage in tests.
var (
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
)

// Tests use the following config.
var (
	// UseExistingCluster selects what k8s API Server is used when running tests.
	// When true the kubeconfig-current-context api server will be used.
	// When false the envtest apiserver will be used.
	useExistingCluster = false

	// Namespace and name for test resources.
	testNSN = types.NamespacedName{
		Namespace: "default",
		Name:      "gstc",
	}

	// TestCtx is used when invoking test clients.
	testCtx = context.Background()
)

// TestMain sets-up a test API server, runs tests and tears down the API server.
func TestMain(m *testing.M) {
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&zap.Options{
		Development: true,
	})))

	// Setup.
	testEnv = &envtest.Environment{
		BinaryAssetsDirectory: "../testbin/bin",
		UseExistingCluster:    &useExistingCluster,
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
	}

	var err error
	cfg, err = testEnv.Start()
	mustNotErr("starting testenv", err)

	err = corev1.AddToScheme(scheme.Scheme)
	mustNotErr("add to schema", err)
	err = clusteropsv1.AddToScheme(scheme.Scheme)
	mustNotErr("add to schema", err)

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	mustNotErr("creating client", err)

	if !useExistingCluster {
		// to access envtest api server (set alwaysShowLog=true to see this message in time)
		fmt.Printf("kubectl --server=%s\n", cfg.Host)
	}

	// Run.
	r := m.Run()

	// Teardown.
	err = testEnv.Stop()
	mustNotErr("stopping testenv", err)

	os.Exit(r)
}

// TestStartManager starts a Manager with the provided reconciler.
func testStartManager(t *testing.T, ctx context.Context, reconciler *controllers.GSTConfigReconciler) *sync.WaitGroup {
	t.Helper()

	// Setup manager (similar to main.go)

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:         scheme.Scheme,
		LeaderElection: false,
	})
	mustNotErr("new manager", err)

	// Add reconciler to manager.
	reconciler.Client = mgr.GetClient()
	reconciler.Scheme = mgr.GetScheme()
	// watch CR and Secrets
	err = ctrl.NewControllerManagedBy(mgr).
		For(&clusteropsv1.GSTConfig{}).
		Watches(
			&source.Kind{Type: &corev1.Secret{}},
			&handler.EnqueueRequestForObject{}).
		Complete(reconciler)

	mustNotErr("setup with manager", err)

	// Start manager.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = mgr.Start(ctx)
		mustNotErr("starting manager", err)
	}()

	return &wg
}

func mustNotErr(msg string, err error) {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
}
