package seekpo

import (
	"context"
	"time"
)

type Series struct {
	Sets []Set
}

type Set struct {
	Code   Code
	Points []Point
}

type Point struct {
	Timestamp time.Time
	Value     Value
	Status    Status
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
	ReadLastSeries(context.Context, []Code) (Series, error)
}

type Range struct {
	Start time.Time
	End   time.Time
}

type SeriesReader interface {
	ReadSeries(context.Context, []Code, Range) (Series, error)
}
