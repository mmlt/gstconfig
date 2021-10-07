package artifactory

import (
	"context"
	"github.com/mmlt/gstconfig/pkg/gst"
)

type Artifactory struct {
	URL  string
	Path string
}

func (a *Artifactory) Read(ctx context.Context, name string, data *gst.Config) error {
	panic("implement me")
}

func (a *Artifactory) Write(ctx context.Context, name string, data *gst.Config) error {
	panic("implement me")
}
