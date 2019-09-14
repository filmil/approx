// Package approx contains code for computing with approximate numbers.
package approx

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// Float64 represents a floating point number with a degree of
// uncertainty.
//
// Every Float64 has an exact value and a delta about it.
type Float64 struct {
	val, delta float64
}

func (f Float64) String() string {
	return fmt.Sprintf("%v±%v", f.val, f.delta)
}

// Parse parses an uncertain number from a string.
//
// Example:
//     approx.Parse("4.2±0.3") -> {4.2, 0.3}
func Parse(s string) (Float64, error) {
	// First strip all spaces from the thing.
	s = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
	splitstr := strings.Split(s, "±")
	if len(splitstr) == 1 {
		// Exact value
	}
	switch len(splitstr) {
	case 1: // Exact
		val, err := strconv.ParseFloat(splitstr[0], 64)
		if err != nil {
			return Float64{}, fmt.Errorf("could not parse as exact float: %v", splitstr)
		}
		return Float64{val: val, delta: 0.0}, nil
	case 2: // Inexact
		val, err := strconv.ParseFloat(splitstr[0], 64)
		if err != nil {
			return Float64{}, fmt.Errorf("could not parse as exact float: %v", splitstr)
		}
		delta, err := strconv.ParseFloat(splitstr[1], 64)
		if err != nil {
			return Float64{}, fmt.Errorf("could not parse as delta float: %v", splitstr)
		}
		return Float64{val: val, delta: math.Abs(delta)}, nil

	default: // Everything else
		return Float64{}, fmt.Errorf("could not parse as approximate number: %v", splitstr)
	}
	return Float64{}, nil
}

// New constructs a new Float64 from exact float components.
//
// The recorded delta is always nonnegative, so
//   New(10,1) == New(10,-1)
func New(val, delta float64) Float64 {
	return Float64{val: val, delta: math.Abs(delta)}
}

// NewMinMax constructs a new Float64 from a minimum and maximum interval boundaries.  min *must*
// be less than or equal to max.
func NewMinMax(min, max float64) (Float64, error) {
	if max < min {
		return Float64{}, fmt.Errorf("min must be less or equal to max: min:%v, max:%v", min, max)
	}
	val := (min + max) / 2
	delta := math.Abs((max - min) / 2)
	return New(val, delta), nil
}

// Add computes a sum of two approximate numbers a and b.
func Add(a, b Float64) Float64 {
	return New(a.val+b.val, a.delta+b.delta)
}

// Sub computes a diference when subtracting a from b.
func Sub(a, b Float64) Float64 {
	return New(a.val-b.val, a.delta+b.delta)
}

// Mul computes a multplication of a and b.
func Mul(a, b Float64) Float64 {
	relA := math.Abs(a.delta / a.val)
	relB := math.Abs(b.delta / b.val)
	rel := relA + relB
	val := a.val * b.val
	delta := math.Abs(val * rel)
	return New(val, delta)
}

// Div computes a quotient of a and b. Zeroes cause infinities, as expected.
func Div(a, b Float64) Float64 {
	relA := math.Abs(a.delta / a.val)
	relB := math.Abs(b.delta / b.val)
	rel := relA + relB
	val := a.val / b.val
	delta := math.Abs(val * rel)
	return New(val, delta)
}

// LessThan returns true if f is definitely smaller than t.
func (f Float64) LessThan(t Float64) bool {
	return f.val+f.delta < t.val-t.delta
}
