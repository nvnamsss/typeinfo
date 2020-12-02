package examples

import "math"

//Math contains mathematic methods
type Math struct {
}

// Sin returns the sine of the radian argument x.
//
// Special cases are:
//	Sin(±0) = ±0
//	Sin(±Inf) = NaN
//	Sin(NaN) = NaN
func (Math) Sin(rad float64) float64 {
	return math.Sin(rad)
}

// Cos returns the cosine of the radian argument x.
//
// Special cases are:
//	Cos(±Inf) = NaN
//	Cos(NaN) = NaN
func (Math) Cos(rad float64) float64 {
	return math.Cos(rad)
}

//Clamps a value between a minimum and maximum value.
func (Math) ClampInt(value, min, max int64) int64 {
	if value < min {
		value = min
	} else if value > max {
		value = max
	}

	return value
}

//Clamps a value between a minimum and maximum value.
func (Math) ClampFloat(value, min, max float64) float64 {
	if value < min {
		value = min
	} else if value > max {
		value = max
	}

	return value
}
