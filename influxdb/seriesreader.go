package influxdb

import (
	"context"
	"fmt"
	"log"
	"strconv"
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
	measurements []seekpo.Measurement,
	codes []seekpo.Code,
) (seekpo.Series, error) {
	sets := make([]seekpo.Set, 0, len(codes))
	queryAPI := r.client.QueryAPI(r.orgName)
	query := formatQuery(r.bucketName, dateRange, measurements, codes)
	log.Printf("[INFO] query '%s' is made", query)
	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return seekpo.Series{}, fmt.Errorf("query failed: %s", err)
	}
	var previous int64 = -1
	for result.Next() {
		current := assertInt64(result.Record().ValueByKey("table"))
		if current != previous {
			set := seekpo.Set{
				Measurement: result.Record().Measurement(),
				Code:        result.Record().Field(),
				Unit:        assertString(result.Record().ValueByKey("unit")),
				Type:        assertString(result.Record().ValueByKey("type")),
			}
			sets = append(sets, set)
			previous = current
		}
		point := seekpo.Point{
			Timestamp: result.Record().Time(),
			Value:     result.Record().Value(),
			Status:    parseStatus(assertString(result.Record().ValueByKey("status"))),
		}
		sets[current].Points = append(sets[current].Points, point)
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
	fields []string,
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
	query := strings.Join(append([]string{from, range_}, filters...), "")
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

func assertInt64(v interface{}) int64 {
	i, _ := v.(int64)
	return i
}

func assertString(v interface{}) string {
	s, _ := v.(string)
	return s
}

func parseStatus(s string) seekpo.Status {
	i, err := strconv.ParseUint(s, 16, 32) // TODO panic
	if err != nil {
		log.Printf("[WARNING] parse status failed: %s", err)
	}
	return seekpo.Status(i)
}
