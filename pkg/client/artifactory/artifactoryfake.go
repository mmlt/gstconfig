package artifactory

import (
	"context"
	"github.com/mmlt/gstconfig/pkg/gst"
)

type ArtifactoryFake struct {
	ReadTally, WriteTally int
	Data                  gst.Config
}

func (a *ArtifactoryFake) Read(ctx context.Context, name string, data *gst.Config) error {
	a.ReadTally++
	data = &a.Data

	return nil
}

func (a *ArtifactoryFake) Write(ctx context.Context, name string, data *gst.Config) error {
	a.WriteTally++
	a.Data = *data

	return nil
}
