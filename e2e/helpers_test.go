package e2e

import (
	clusteropsv1 "github.com/mmlt/gstconfig/api/v1"
	"github.com/stretchr/testify/assert"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"testing"
	"time"
)

func testCreateCR(t *testing.T, cr *clusteropsv1.GSTConfig) {
	t.Helper()

	nsn := types.NamespacedName{
		Namespace: cr.Namespace,
		Name:      cr.Name,
	}
	testDeleteCR(t, nsn)
	err := k8sClient.Create(testCtx, cr)
	assert.NoError(t, err)
}

func testDeleteCR(t *testing.T, nsn types.NamespacedName) {
	t.Helper()

	obj := &clusteropsv1.GSTConfig{}
	err := k8sClient.Get(testCtx, nsn, obj)
	if apierrors.IsNotFound(err) {
		return
	}
	assert.NoError(t, err)
	err = k8sClient.Delete(testCtx, obj)
	assert.NoError(t, err)
}

func testGetCR(t *testing.T, nsn types.NamespacedName) *clusteropsv1.GSTConfig {
	t.Helper()

	obj := &clusteropsv1.GSTConfig{}
	err := k8sClient.Get(testCtx, nsn, obj)
	assert.NoError(t, err)
	return obj
}

func testGetCRWhenConditionReady(t *testing.T, nsn types.NamespacedName) *clusteropsv1.GSTConfig {
	t.Helper()

	obj := &clusteropsv1.GSTConfig{}
	err := wait.Poll(time.Second, 10*time.Minute, func() (done bool, err error) {
		err = k8sClient.Get(testCtx, nsn, obj)
		if err != nil {
			return false, err
		}
		//TODO enable when CR has status.conditions field
		for _, c := range obj.Status.Conditions {
			if c.Type == "Ready" {
				return c.Status == metav1.ConditionTrue, nil
			}
		}
		return false, nil
	})
	assert.NoError(t, err)
	return obj
}
