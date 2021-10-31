package influxdb

import (
	"context"
	"fmt"
	"strconv"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/rekhin/seekpo-sdk-go"
)

type SeriesWriter struct {
	client     influxdb2.Client
	orgName    string
	bucketName string
}

var _ seekpo.SeriesWriter = new(SeriesWriter)

func NewSeriesWriter(client influxdb2.Client, org, bucket string) *SeriesWriter {
	return &SeriesWriter{
		client:     client,
		orgName:    org,
		bucketName: bucket,
	}
}

// TODO Сделать запись массивов и объектов помимо примитивных типов
func (w *SeriesWriter) WriteSeries(ctx context.Context, series seekpo.Series) error {
	writeAPI := w.client.WriteAPIBlocking(w.orgName, w.bucketName)
	points := []*write.Point{}
	for i := range series.Sets {
		for j := range series.Sets[i].Points {
			point := influxdb2.NewPointWithMeasurement(series.Sets[i].Measurement).
				SetTime(series.Sets[i].Points[j].Timestamp).
				AddField(series.Sets[i].Code, series.Sets[i].Points[j].Value).
				AddTag("status", strconv.FormatUint(uint64(series.Sets[i].Points[j].Status), statusBase)).
				AddTag("unit", series.Sets[i].Unit).
				AddTag("type", series.Sets[i].Type.String())
			points = append(points, point)
		}
	}
	if err := writeAPI.WritePoint(ctx, points...); err != nil {
		return fmt.Errorf("write points failed: %s", err)
	}
	return nil
}
