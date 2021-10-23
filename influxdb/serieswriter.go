package influxdb

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/rekhin/seekpo-sdk-go"
)

type SeriesWriter struct {
	client      influxdb2.Client
	org, bucket string
}

func NewSeriesWriter(client influxdb2.Client, org, bucket string) *SeriesWriter {
	return &SeriesWriter{
		client: client,
		org:    org,
		bucket: bucket,
	}
}

func (w *SeriesWriter) WriteSeries(ctx context.Context, series seekpo.Series) error {
	writeAPI := w.client.WriteAPIBlocking(w.org, w.bucket)
	var points []*write.Point
	for i := range series.Sets {
		for j := range series.Sets[i].Points {
			point := influxdb2.NewPointWithMeasurement("tag").
				AddTag("code", string(series.Sets[i].Code)).
				AddField("value", series.Sets[i].Points[j].Value).
				AddField("status", series.Sets[i].Points[j].Status).
				SetTime(series.Sets[i].Points[j].Timestamp)
			points = append(points, point)
		}
	}
	if err := writeAPI.WritePoint(ctx, points...); err != nil {
		return fmt.Errorf("write points failed: %s", err)
	}
	return nil
}
