package approx

import (
	"fmt"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var opts []cmp.Option = []cmp.Option{
	cmp.AllowUnexported(Float64{}),
}

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected Float64
		err      error
	}{
		{
			input:    "0",
			expected: Float64{0, 0},
		},
		{
			input:    "1",
			expected: Float64{1, 0},
		},
		{
			// 4.2±0.3
			input:    "-1.23",
			expected: Float64{-1.23, 0},
		},
		{
			input:    "4.2±0.3",
			expected: Float64{4.2, 0.3},
		},
		{
			input:    "4.2±-0.3",
			expected: Float64{4.2, 0.3},
		},
		{
			input: "4.2±--0.3",
			err:   fmt.Errorf("could not parse as delta float: [4.2 --0.3]"),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			actual, err := Parse(test.input)
			if err != nil {
				if test.err == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !cmp.Equal(err.Error(), test.err.Error(), opts...) {
					t.Errorf("different error: expected: %v, actual: %v, diff: %v",
						test.err, err, cmp.Diff(test.err.Error(), err.Error(), opts...))
				}
			}
			if !cmp.Equal(actual, test.expected, opts...) {
				t.Errorf("expected: %v, actual: %v", test.expected, actual)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	tests := []struct {
		v, d     float64
		expected Float64
	}{
		{
			expected: Float64{0, 0},
		},
		{
			v:        10.0,
			d:        -1.0,
			expected: Float64{10, 1},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("v=%v;d=%v", test.v, test.d), func(t *testing.T) {
			actual := New(test.v, test.d)
			if !cmp.Equal(actual, test.expected, opts...) {
				t.Errorf("expected: %v, actual: %v", test.expected, actual)
			}
		})
	}
}

func TestNewMinMax(t *testing.T) {
	t.Parallel()
	tests := []struct {
		min, max float64
		expected Float64
	}{
		{
			expected: Float64{0, 0},
		},
		{
			min:      0,
			max:      10,
			expected: Float64{5, 5},
		},
		{
			min:      1,
			max:      10,
			expected: Float64{5.5, 4.5},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("min=%v;max=%v", test.min, test.max), func(t *testing.T) {
			actual, err := NewMinMax(test.min, test.max)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !cmp.Equal(actual, test.expected, opts...) {
				t.Errorf("expected: %v, actual: %v", test.expected, actual)
			}
		})
	}
}

func must(v Float64, err error) Float64 {
	if err != nil {
		panic(fmt.Sprintf("unexpected error while parsing: %v", err))
	}
	return v
}

func TestOps(t *testing.T) {
	t.Parallel()
	tests := []struct {
		op1, op2 Float64
		sum      Float64
		sub      Float64
		product  Float64
		quotient Float64
	}{
		{
			op1:      New(1, 2),
			op2:      New(3, 4),
			sum:      New(4, 6),
			sub:      New(-2, 6),
			product:  New(3, 10),
			quotient: must(Parse("0.3333333333333333±1.111111111111111")),
		},
		{
			op1:      New(1, 2),
			op2:      New(3, 4),
			sum:      New(4, 6),
			sub:      New(-2, 6),
			product:  New(3, 10),
			quotient: must(Parse("0.3333333333333333±1.111111111111111")),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf(fmt.Sprintf("(%v;%v)", test.op1, test.op2)), func(t *testing.T) {
			sum := Add(test.op1, test.op2)
			if !cmp.Equal(sum, test.sum, opts...) {
				t.Errorf("sum: expected: %v, actual: %v", test.sum, sum)
			}
			sub := Sub(test.op1, test.op2)
			if !cmp.Equal(sub, test.sub, opts...) {
				t.Errorf("sub: expected: %v, actual: %v", test.sub, sub)
			}
			product := Mul(test.op1, test.op2)
			if !cmp.Equal(product, test.product, opts...) {
				t.Errorf("product: expected: %v, actual: %v", test.product, product)
			}
			quotient := Div(test.op1, test.op2)
			if !cmp.Equal(quotient, test.quotient, opts...) {
				t.Errorf("quotient: expected: %v, actual: %v", test.quotient, quotient)
			}
		})
	}
}

func TestRelOps(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		op1, op2           Float64
		lt, le, gt, ge, ov bool
		rd                 float64
	}{
		{
			op1: New(0, 1),
			op2: New(3, 1),

			lt: true,
			le: true,
			gt: false,
			ge: false,
			ov: false,
			rd: math.Inf(1),
		},
		{
			op1: New(1, 1),
			op2: New(2, 1),

			lt: false,
			le: false,
			gt: false,
			ge: false,
			ov: true,
			rd: 1,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf(fmt.Sprintf("(%v;%v)", test.op1, test.op2)), func(t *testing.T) {
			lt := test.op1.Lt(test.op2)
			le := test.op1.Le(test.op2)
			gt := test.op1.Gt(test.op2)
			ge := test.op1.Ge(test.op2)
			ov := Overlap(test.op1, test.op2)
			rd := test.op1.RelDelta()
			if lt != test.lt || le != test.le || gt != test.gt ||
				ge != test.ge || ov != test.ov || rd != test.rd {
				t.Errorf(
					"was : (lt=%v,le=%v,gt=%v,ge=%v,ov=%v,rd=%v)\nwant: (lt=%v,le=%v,gt=%v,ge=%v,ov=%v,rd=%v)",
					lt, le, gt, ge, ov, rd,
					test.lt, test.le, test.gt, test.ge, test.ov, test.rd)
			}
		})
	}
}

func TestApply(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    Float64
		f        func(float64) float64
		expected Float64
	}{
		{
			name:     "log",
			input:    New(1, 0.1),
			f:        math.Log,
			expected: New(0, 0.1),
		},
		{
			name:     "exp",
			input:    New(1, 0.1),
			f:        math.Exp,
			expected: must(Parse("2.718281828459045±0.27182818284590454")),
		},
		{
			name:  "x^2",
			input: New(10, 0.1),
			f: func(x float64) float64 {
				return x * x
			},
			expected: must(Parse("100±1.9999999999988916")),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			actual := test.input.Apply(test.f, 1e-3)
			if !cmp.Equal(actual, test.expected, opts...) {
				t.Errorf("was : %v\nwant: %v", actual, test.expected)
			}
		})
	}
}
