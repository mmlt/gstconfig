package e2e

import (
	clusteropsv1 "github.com/mmlt/gstconfig/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// This file contains data used by multiple tests.

func testGSTConfigCR(nn types.NamespacedName, spec *clusteropsv1.GSTConfigSpec) *clusteropsv1.GSTConfig {
	return &clusteropsv1.GSTConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nn.Name,
			Namespace: nn.Namespace,
		},
		Spec: *spec,
	}
}
