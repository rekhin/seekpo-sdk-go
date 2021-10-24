package influxdb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/rekhin/seekpo-sdk-go"
)

type SeriesReader struct {
	client     influxdb2.Client
	orgName    string
	bucketName string
}

func NewSeriesReader(client influxdb2.Client, org, bucket string) *SeriesReader {
	return &SeriesReader{
		client:     client,
		orgName:    org,
		bucketName: bucket,
	}
}

func (r *SeriesReader) ReadSeries(
	ctx context.Context,
	dateRange seekpo.Range,
	measurement []seekpo.Measurement,
	fields []seekpo.Field,
) (seekpo.Series, error) {
	sets := make([]seekpo.Set, len(fields))
	queryAPI := r.client.QueryAPI(r.orgName)
	query := formatQuery(r.bucketName, dateRange, measurement, fields)
	log.Printf("[INFO] query '%s' is made", query)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return seekpo.Series{}, fmt.Errorf("query failed: %s", err)
	}
	i := -1
	for result.Next() {
		if result.TableChanged() {
			i++
			sets[i].Field = result.Record().Field()
			// fmt.Printf("table: %s\n", result.TableMetadata().String())
		}
		point := seekpo.Point{
			Timestamp: result.Record().Time(),
			Value:     result.Record().Value,
			// Status:    seekpo.Status(result.Record().Field()),
		}
		sets[i].Points = append(sets[i].Points, point)
		// fmt.Printf("row: %s\n", result.Record().String())
	}
	if result.Err() != nil {
		return seekpo.Series{}, fmt.Errorf("next failed: %s", result.Err())
	}
	return seekpo.Series{Sets: sets}, nil
}

func formatQuery(
	bucketName string,
	dateRange seekpo.Range,
	measurements []seekpo.Measurement,
	fields []seekpo.Field,
) string {
	from := fmt.Sprintf(`from(bucket: "%s")`, bucketName)
	range_ := fmt.Sprintf(`|> range(start: %s, stop: %s)`,
		dateRange.Start.Format(time.RFC3339Nano),
		dateRange.End.Format(time.RFC3339Nano),
	)
	filters := []string{
		formatFilter("_measurement", measurements),
		formatFilter("_field", fields),
	}
	query := fmt.Sprintf(`%s %s %s`, from, range_, strings.Join(filters, " "))
	return query
}

func formatFilter(tag string, items []string) string {
	conditions := formatConditions(tag, items)
	filter := fmt.Sprintf(`|> filter(fn: (r) => %s)`, strings.Join(conditions, " or "))
	return filter
}

func formatConditions(tag string, items []seekpo.Field) []string {
	conditions := make([]string, len(items))
	for i := range items {
		conditions[i] = fmt.Sprintf(`r.%s == "%s"`, tag, items[i])
	}
	return conditions
}
