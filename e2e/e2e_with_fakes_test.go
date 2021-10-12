package e2e

import (
	"context"
	"k8s.io/apimachinery/pkg/labels"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api/v1"

	clusteropsv1 "github.com/mmlt/gstconfig/api/v1"
	"github.com/mmlt/gstconfig/controllers"
	"github.com/mmlt/gstconfig/pkg/client/artifactory"
	"github.com/mmlt/gstconfig/pkg/gst"
	"github.com/mmlt/testr"
	"github.com/stretchr/testify/assert"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)

func TestGoodRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	logf.SetLogger(testr.New(t))

	l := map[string]string{
		"clusterops.mmlt.nl/operator": "playgroundenvs",
	}

	ls, err := labels.ValidatedSelectorFromSet(l)
	assert.NoError(t, err)

	artifake := &artifactory.ArtifactoryFake{}
	wg := testStartManager(t, ctx, &controllers.GSTConfigReconciler{
		ClusterLabelSelector: ls,
		Repo:                 artifake,
	})

	t.Run("should_run_without_secrets", func(t *testing.T) {
		testCreateCR(t, testGSTConfigCR(testNSN, &clusteropsv1.GSTConfigSpec{}))
		got := testGetCRWhenConditionReady(t, testNSN)

		// Condition
		assert.Equal(t, 1, len(got.Status.Conditions), "number of Status.Conditions")
		assert.Equal(t, "Synced", got.Status.Conditions[0].Reason)
		// Probe fakes
		assert.Equal(t, 2, artifake.ReadTally)
		assert.Equal(t, gst.Config{}, artifake.Data)
	})

	t.Run("should_run_with_secrets", func(t *testing.T) {
		t.Logf("kubectl -s %s", cfg.Host) // because IDE restrict output to a single t.Run() when debugging.

		testCreateSecret(t, testClusterSecret(t, testNSN.Namespace, "foo", l, &clientcmdapi.Cluster{
			Server: "https://foo.example.com",
		}, nil))
		testCreateSecret(t, testClusterSecret(t, testNSN.Namespace, "bar", l, &clientcmdapi.Cluster{
			Server: "https://bar.example.com",
		}, nil))

		testCreateCR(t, testGSTConfigCR(testNSN, &clusteropsv1.GSTConfigSpec{}))
		got := testGetCRWhenConditionReady(t, testNSN)

		// Condition
		assert.Equal(t, 1, len(got.Status.Conditions), "number of Status.Conditions")
		assert.Equal(t, "Synced", got.Status.Conditions[0].Reason)
		// Probe fakes
		// note: clusters is ordered by Endpoint
		want := gst.Config{
			Clusters: []gst.ClusterEndpoint{
				{Endpoint: "https://bar.example.com", Default: false, Deprecated: false},
				{Endpoint: "https://foo.example.com", Default: false, Deprecated: false},
			},
		}
		assert.Equal(t, want, artifake.Data)
	})

	// TODO Add tests should_detect_secret_create and should_detect_secret_delete
	//  Wait for completion by polling data or updateTime

	// teardown manager
	cancel()
	wg.Wait()
}
