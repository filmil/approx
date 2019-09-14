// Package approx contains code for computing with approximate numbers.
//
// An approximate number is a value (float64) plus, or minus some uncertainty.
// Approximate numbers are what you normally get as result of any real-world
// measurement.  This package allows you to use approximate numbers and use
// "regular" mathematical operations to compute with them.
//
// Why is this useful?
//
// Approximate numbers come out of any sort of physical measurement.  No real
// life measurement ever yields a single number.  Though we sometimes choose to
// ignore measurement error, that error is always present.  The question
// becomes, suppose we do *not* want to disregard the error, what happens then?
// For example, measuring one side of a kitchen table with a tape measure would
// yield the result:
//
//     width = (50±0.5)cm
//
// The 0.5cm error comes from the fact that a tape measure has the smallest
// division of 1 centimeter.  Since the divisions are large enough that we can
// estimate if we are off more than one half of the division, we can say that
// we are confident in not making more than 0.5cm of a measurement error.
//
// Suppose now that we measure the length of the table too:
//
//     length = (100±0.5)cm
//
// Since we are measuring with the same tape measure, the outcome in terms of
// error is similar: we're making another error of about 0.5cm.
//
// Now, what is the perimeter of the table?  It is:
//
//     perimeter = 2 * (width + length)
//
// But, since the original width and length that we computed are approximate
// numbers, we will also have some error in the computation of the perimeter.
//
// Since we could have overshot our measurement for all values at one extreme,
// or undershot at another, our perimeter falls in the interval:
//
//     perimeter = (300±2)cm
//
// Note here that individual measurement errors added up.  Now, if we wanted to
// compute the difference between length and width, we'd get:
//
//     length - width = (50±1)cm
//
// What happened here?  We see that even though the data points were
// subtracted, the errors were *added* together.  This is because again the
// errors could have conspired to make our measurement less accurate, and we
// have to account for that.
//
// This library has a few functions that make working with approximate numbers
// easy.  You can load up some approximate numbers like so:
//
//     import "github.com/filmil/approx"
//     width, _ := approx.Parse("50±0.5")
//     length, _ := approx.Parse("50±0.5")
//     perimeter := approx.Add(
//         approx.Add(width, length),
//         approx.Add(width, length),
//     )
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

// String implements Stringer.
//
// This implementation prints the most basic version of the number.  If you want
// more specific formatting, use Value() and Delta() to extract the components
// from the number, and format them at will.
func (f Float64) String() string {
	return fmt.Sprintf("%v±%v", f.val, f.delta)
}

// Value returns the value at the center of f's interval.
func (f Float64) Value() float64 {
	return f.val
}

// Delta returns the delta around the interval.  delta is nonnegative.
func (f Float64) Delta() float64 {
	return f.delta
}

// Min returns the minimal extreme value for f.
func (f Float64) Min() float64 {
	return f.val - f.delta
}

// Max returns the maximal extreme value for f.
func (f Float64) Max() float64 {
	return f.val + f.delta
}

// RelDelta returns the relative error of f.
func (f Float64) RelDelta() float64 {
	return math.Abs(f.delta / f.val)
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

// Mul computes a scalar product of f with a number c.
func (f Float64) Mul(c float64) Float64 {
	return New(c*f.val, math.Abs(c*f.delta))
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

// Lt returns true if f is definitely less than t.
func (f Float64) Lt(t Float64) bool {
	return f.val+f.delta < t.val-t.delta
}

// Le returns true if f is definitely either less than, or equal to t.
func (f Float64) Le(t Float64) bool {
	return f.val+f.delta <= t.val-t.delta
}

// Gt returns true if f is definitely greater than t.
func (f Float64) Gt(t Float64) bool {
	return t.Le(f)
}

// Ge returns true if f is definitely either greather than, or equal to t.
func (f Float64) Ge(t Float64) bool {
	return t.Lt(f)
}

// Overlap returns true if t and f may overlap.
func Overlap(f, t Float64) bool {
	return !f.Le(t) && !t.Le(f)
}

// eqFunc is a dirty trick which compares function based on their address in
// memory.
func eqFunc(f1, f2 func(float64) float64) bool {
	return fmt.Sprintf("%p", f1) == fmt.Sprintf("%p", f2)
}

// applyLog computes approximate value for a natural logarithm.
//
// Based on first-order Taylor expansion around x:
//   ln(x+dx) = ln(x) + 1/x * dx
func (f Float64) applyLog() Float64 {
	delta := f.delta / f.val
	val := math.Log(f.val)
	return New(val, delta)
}

// applyExp computes approximate value for e^x.
//
// Based on first-order Taylor expansion around x:
//   e^(x+dx) = e^x + e^x*dx
func (f Float64) applyExp() Float64 {
	val := math.Exp(f.val)
	delta := val * f.delta
	return New(val, delta)
}

// Apply applies the function fx to f.
//
// Based on first order Taylor expansion of fx around f:
// f := x + dx
// fx(f) = x + fx'(f) * dx.
//
// For well known functions, the computation is exact.  For user-defined
// functions, the computation is via computing numeric derivative around the
// centerpoint of f, for which 'eps' is the interval to compute numeric
// derivative on.
func (f Float64) Apply(fx func(float64) float64, eps float64) Float64 {
	// Special-case some interesting functions.
	if eqFunc(fx, math.Log) {
		return f.applyLog()
	}
	if eqFunc(fx, math.Exp) {
		return f.applyExp()
	}

	// Central difference numeric derivative computation.
	fmin := fx(f.val - eps)
	fmax := fx(f.val + eps)
	d := 2 * eps
	dfx := (fmax - fmin) / d
	return New(fx(f.val), math.Abs(dfx*f.delta))
}
