package gst

// Config represent the Getting Started Tool cluster configuration.
type Config struct {
	Clusters []ClusterEndpoint `json:"clusters"`
}

// ClusterEndpoint represents a single cluster.
type ClusterEndpoint struct {
	// Endpoint is the URL of the APIServer
	Endpoint string `json:"endpoint"`
	// Default sets the kubeconfig current context to this cluster.
	Default bool `json:"default"`
	// Deprecated indicates the cluster is no longer used.
	Deprecated bool `json:"deprecated"`
}
