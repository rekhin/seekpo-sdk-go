package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/rekhin/seekpo-sdk-go"
	"github.com/rekhin/seekpo-sdk-go/influxdb"
)

func main() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	client := influxdb2.NewClient("http://localhost:8086", "my-super-secret-auth-token")
	writer := influxdb.NewSeriesWriter(client)
	series := seekpo.Series{
		Sets: []seekpo.Set{
			{
				Code: "temperature",
				Points: func() []seekpo.Point {
					var points []seekpo.Point
					for i := 1; i <= 24; i++ {
						point := seekpo.Point{
							Timestamp: time.Date(2021, 10, 22, i, 0, 0, 0, time.UTC),
							Value:     float32(rand.Intn(1500)) / 100,
							Status:    0,
						}
						points = append(points, point)
					}
					return points
				}(),
			},
		},
	}
	if err := writer.WriteSeries(context.Background(), series); err != nil {
		log.Printf("[ERROR] write series failed: %s", err)
	}
}
