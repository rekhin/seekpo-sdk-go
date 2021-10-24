package seekpo

import (
	"context"
	"time"
)

type Series struct {
	Measurement Measurement
	Sets        []Set
}

type Measurement = string

type Set struct {
	Field  Field
	Points []Point
}

type Point struct {
	Timestamp time.Time
	Value     Value
}

type SeriesWriter interface {
	WriteSeries(context.Context, Series) error
}

type LastSeriesReader interface {
	ReadLastSeries(context.Context, []Field) (Series, error)
}

type Range struct {
	Start time.Time
	End   time.Time
}

type SeriesReader interface {
	ReadSeries(context.Context, Range, []Measurement, []Field) (Series, error)
}
