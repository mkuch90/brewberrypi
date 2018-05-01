package brewery

import (
	"errors"

	temp "github.com/brewduino/temperature"
)

var ErrInvalidCurrentVals = errors.New("could not retrive")

type Target struct {
	Min float64
	Max float64
}

type Comparison string

var GreaterThan = Comparison("GREATER_THAN")
var LessThan = Comparison("LESS_THAN")
var OnTarget = Comparison("ON_TARGETS")

type ComparisonResult struct {
	Key temp.TemperatureSensorName
	Cmp Comparison
}
