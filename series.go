package seekpo

import (
	"context"
	"time"
)

type Series struct {
	Measurement Measurement
	Sets        []Set
}

type Measurement int

//go:generate stringer -type=Measurement -trimprefix=Measurement
const (
	MeasurementValue Measurement = iota
	MeasurementStatus
)

type Set struct {
	Field  Field
	Points []Point
}

type Point struct {
	Timestamp time.Time
	Value     Value
	// Status    Status
}

// type BoolPoint struct {
// 	Timestamp time.Time
// 	Bool      bool
// 	Status    Status
// }

type Value interface{}

type Status uint32

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
	ReadSeries(context.Context, []Field, Range) (Series, error)
}
