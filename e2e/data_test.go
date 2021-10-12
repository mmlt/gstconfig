package e2e

import (
	clusteropsv1 "github.com/mmlt/gstconfig/api/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api/v1"
	"testing"
)

// This file contains data used by multiple tests.

//
func testGSTConfigCR(nn types.NamespacedName, spec *clusteropsv1.GSTConfigSpec) *clusteropsv1.GSTConfig {
	return &clusteropsv1.GSTConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nn.Name,
			Namespace: nn.Namespace,
		},
		Spec: *spec,
	}
}

// TestClusterSecret return a Secret with all the data needed to access a k8s APIServer.
func testClusterSecret(t *testing.T, namespace, name string, labels map[string]string, cluster *clientcmdapi.Cluster, authInfo *clientcmdapi.AuthInfo) *corev1.Secret {
	t.Helper()

	cjson, err := json.Marshal(cluster)
	assert.NoError(t, err)

	aijson, err := json.Marshal(authInfo)
	assert.NoError(t, err)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: map[string][]byte{
			"cluster":  cjson,
			"authInfo": aijson,
		},
	}
	return secret
}
