package influxdb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/rekhin/seekpo-sdk-go"
)

type SeriesReader struct {
	client     influxdb2.Client
	orgName    string
	bucketName string
}

var _ seekpo.SeriesReader = new(SeriesReader)

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
	ids []uuid.UUID,
) (seekpo.Series, error) {
	sets := make([]seekpo.Set, len(ids))
	queryAPI := r.client.QueryAPI(r.orgName)
	query := formatQuery(r.bucketName, dateRange, measurement, ids)
	log.Printf("[INFO] query '%s' is made", query)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return seekpo.Series{}, fmt.Errorf("query failed: %s", err)
	}
	i := -1
	for result.Next() {
		if result.TableChanged() {
			i++
			sets[i].ID = uuid.MustParse(result.Record().Field()) // TODO panic
		}
		point := seekpo.Point{
			Timestamp: result.Record().Time(),
			Value:     result.Record().Value(),
			// Status:    seekpo.Status(result.Record().Field()),
		}
		sets[i].Points = append(sets[i].Points, point)
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
	ids []uuid.UUID,
) string {
	from := fmt.Sprintf(`from(bucket: "%s")`, bucketName)
	range_ := fmt.Sprintf(`|> range(start: %s, stop: %s)`,
		dateRange.Start.Format(time.RFC3339Nano),
		dateRange.End.Format(time.RFC3339Nano),
	)
	fields := make([]string, len(ids))
	for i := range ids {
		fields[i] = ids[i].String()
	}
	filters := []string{
		formatFilter("_measurement", measurements),
		formatFilter("_field", fields),
	}
	query := fmt.Sprintf(`%s %s %s`, from, range_, strings.Join(filters, " "))
	return query
}

func formatFilter(s string, ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	conditions := formatConditions(s, ss)
	filter := fmt.Sprintf(`|> filter(fn: (r) => %s)`, strings.Join(conditions, " or "))
	return filter
}

func formatConditions(s string, ss []string) []string {
	conditions := make([]string, len(ss))
	for i := range ss {
		conditions[i] = fmt.Sprintf(`r.%s == "%s"`, s, ss[i])
	}
	return conditions
}
