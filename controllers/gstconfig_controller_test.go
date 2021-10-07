package controllers

import (
	"github.com/mmlt/gstconfig/pkg/gst"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api/v1"
	"testing"
)

func Test_mapSecretsToClusterConfig(t *testing.T) {
	tests := []struct {
		it      string
		in      *corev1.SecretList
		want    *gst.Config
		wantErr bool
	}{
		{
			it:   "should_handle_empty_input",
			in:   &corev1.SecretList{},
			want: &gst.Config{},
		},
		{
			it: "should_handle_normal_input",
			in: &corev1.SecretList{
				Items: []corev1.Secret{
					{
						StringData: map[string]string{
							"cluster": mustJSON2String(&clientcmdapi.Cluster{
								Server: "https://foo.example.com",
							}),
						},
					},
				},
			},
			want: &gst.Config{
				Clusters: []gst.ClusterEndpoint{
					{
						Endpoint: "https://foo.example.com",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			got, err := mapSecretsToClusterConfig(tt.in)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// MustJSON2String is a test helper that converts objects to a JSON string.
func mustJSON2String(in interface{}) string {
	out, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return string(out)
}
