package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/rekhin/seekpo-sdk-go"
	"github.com/rekhin/seekpo-sdk-go/influxdb"
)

const (
	orgName    = "my-org"
	bucketName = "my-bucket"
)

func main() {
	client := influxdb2.NewClient("http://localhost:8086", "my-super-secret-auth-token")
	ctx := context.Background()
	if err := initDatabase(ctx, client); err != nil {
		log.Printf("[ERROR] init database failed: %s", err)
	}
	if err := writePoints(ctx, client); err != nil {
		log.Printf("[ERROR] write points failed: %s", err)
	}
	client.Close()
}

func initDatabase(ctx context.Context, client influxdb2.Client) error {
	org, err := client.OrganizationsAPI().FindOrganizationByName(ctx, orgName)
	if err != nil {
		org, err = client.OrganizationsAPI().CreateOrganizationWithName(ctx, orgName)
		if err != nil {
			return fmt.Errorf("create organization with name %q failed: %s", orgName, err)
		}
		log.Printf("[INFO] organization with name %q is created", orgName)
	}
	if _, err := client.BucketsAPI().FindBucketByName(ctx, bucketName); err != nil {
		if _, err := client.BucketsAPI().CreateBucketWithName(ctx, org, bucketName); err != nil {
			return fmt.Errorf("create bucket with name %q failed: %s", bucketName, err)
		}
		log.Printf("[INFO] bucket with name %q is created", bucketName)
	}
	return nil
}

func writePoints(ctx context.Context, client influxdb2.Client) error {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	writer := influxdb.NewSeriesWriter(client, orgName, bucketName)
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
	if err := writer.WriteSeries(ctx, series); err != nil {
		return fmt.Errorf("write series failed: %s", err)
	}
	return nil
}
