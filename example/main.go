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
	bucketName = "PT1H"
)

func main() {
	client := influxdb2.NewClient("http://localhost:8086", "my-super-secret-auth-token")
	ctx := context.Background()
	if err := initDatabase(ctx, client); err != nil {
		log.Printf("[ERROR] init database failed: %s", err)
	}
	if err := writeSeries(ctx, client); err != nil {
		log.Printf("[ERROR] write series failed: %s", err)
	}
	if err := readSeries(ctx, client); err != nil {
		log.Printf("[ERROR] read series failed: %s", err)
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

func writeSeries(ctx context.Context, client influxdb2.Client) error {
	writer := influxdb.NewSeriesWriter(client, orgName, bucketName)
	series := seekpo.Series{
		Measurement: "example",
		Sets: []seekpo.Set{
			{
				Field:  "temperature",
				Points: generatePoints(),
			},
			{
				Field:  "density",
				Points: generatePoints(),
			},
		},
	}
	if err := writer.WriteSeries(ctx, series); err != nil {
		return fmt.Errorf("write series failed: %s", err)
	}
	return nil
}

func generatePoints() []seekpo.Point {
	points := []seekpo.Point{}
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 24; i++ {
		point := seekpo.Point{
			Timestamp: time.Date(2021, 10, 22, i, 0, 0, 0, time.UTC),
			Value:     float32(rand.Intn(1500)) / 100,
			// Status:    0,
		}
		points = append(points, point)
	}
	return points
}

func readSeries(ctx context.Context, client influxdb2.Client) error {
	reader := influxdb.NewSeriesReader(client, orgName, bucketName)
	series, err := reader.ReadSeries(ctx,
		seekpo.Range{
			Start: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2021, 10, 22, 2, 0, 0, 0, time.UTC),
		},
		[]seekpo.Measurement{},
		[]seekpo.Field{"temperature"},
	)
	if err != nil {
		return fmt.Errorf("read series failed: %s", err)
	}
	log.Printf("[INFO] series %+v", series)
	return nil
}
