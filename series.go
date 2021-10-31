// TODO rename into 'skp'
package seekpo

import (
	"context"
	"time"
)

type Series struct {
	Sets []Set
}

type Set struct {
	Measurement Measurement
	Code        Code
	Unit        string
	Type        Type
	Points      []Point
}

type Measurement = string

type Code = string

type Point struct {
	Timestamp time.Time
	Value     Value
	Status    Status
}

type SeriesWriter interface {
	WriteSeries(context.Context, Series) error
}

type LastSeriesReader interface {
	ReadLastSeries(context.Context, []Measurement, []Code) (Series, error)
}

type Range struct {
	Start time.Time
	End   time.Time
}

type SeriesReader interface {
	// TODO rework with filters (options)
	ReadSeries(context.Context, Range, []Measurement, []Code) (Series, error)
}
