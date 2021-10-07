package e2e

import (
	"context"
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

	artifake := &artifactory.ArtifactoryFake{}
	wg := testStartManager(t, ctx, &controllers.GSTConfigReconciler{
		Repo: artifake,
	})

	t.Run("should_complete_with_secrets", func(t *testing.T) {
		testCreateCR(t, testGSTConfigCR(testNSN, &clusteropsv1.GSTConfigSpec{}))
		got := testGetCRWhenConditionReady(t, testNSN)

		// Condition
		assert.Equal(t, 1, len(got.Status.Conditions), "number of Status.Conditions")
		assert.Equal(t, "Synced", got.Status.Conditions[0].Reason)
		// Probe fakes
		assert.Equal(t, 2, artifake.ReadTally)
		assert.Equal(t, gst.Config{}, artifake.Data)
	})

	t.Run("should_complete_with_secrets", func(t *testing.T) {
		testCreateCR(t, testGSTConfigCR(testNSN, &clusteropsv1.GSTConfigSpec{}))
		got := testGetCRWhenConditionReady(t, testNSN)

		// Condition
		assert.Equal(t, 1, len(got.Status.Conditions), "number of Status.Conditions")
		assert.Equal(t, "Synced", got.Status.Conditions[0].Reason)
		// Probe fakes
		assert.Equal(t, 2, artifake.ReadTally)
		assert.Equal(t, gst.Config{}, artifake.Data)
	})

	// teardown manager
	cancel()
	wg.Wait()
}
