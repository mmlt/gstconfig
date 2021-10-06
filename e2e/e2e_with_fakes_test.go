package e2e

import (
	"context"
	clusteropsv1 "github.com/mmlt/gstconfig/api/v1"
	"github.com/mmlt/testr"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"testing"
)

func TestGoodRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	logf.SetLogger(testr.New(t))

	wg := testManagerWithFakeClients(t, ctx)

	t.Run("should_run_all_steps", func(t *testing.T) {
		testCreateCR(t, testGSTConfigCR(testNSN, &clusteropsv1.GSTConfigSpec{}))

		//TODO enable tesy
		/*		got := testGetCRWhenConditionReady(t, testNSN)

				// Condition
				assert.Equal(t, 1, len(got.Status.Conditions), "number of Status.Conditions")
				assert.Equal(t, clusteropsv1.ReasonReady, got.Status.Conditions[0].Reason)
		*/
	})

	// teardown manager
	cancel()
	wg.Wait()
}
