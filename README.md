# approx: calculations with approximate numbes, implemented in go

[![Go Report Card](https://goreportcard.com/badge/github.com/filmil/approx)](https://goreportcard.com/report/github.com/filmil/approx)
[![Build Status](https://travis-ci.com/filmil/approx.svg?branch=master)](https://travis-ci.com/filmil/approx)
[![Coverage Status](https://coveralls.io/repos/github/filmil/approx/badge.svg?branch=master)](https://coveralls.io/github/filmil/approx?branch=master)
[![GoDoc](https://godoc.org/github.com/filmil/approx?status.svg)](https://godoc.org/github.com/filmil/approx)

Package approx contains code for computing with approximate numbers.

An approximate number is a value (float64) plus, or minus some uncertainty.
Approximate numbers are what you normally get as result of any real-world
measurement.  This package allows you to use approximate numbers and use
"regular" mathematical operations to compute with them.

Why is this useful?

Approximate numbers come out of any sort of physical measurement.  No real
life measurement ever yields a single number.  Though we sometimes choose to
ignore measurement error, that error is always present.  The question
becomes, suppose we do *not* want to disregard the error, what happens then?
For example, measuring one side of a kitchen table with a tape measure would
yield the result:

    width = (50±0.5)cm

The 0.5cm error comes from the fact that a tape measure has the smallest
division of 1 centimeter.  Since the divisions are large enough that we can
estimate if we are off more than one half of the division, we can say that
we are confident in not making more than 0.5cm of a measurement error.

Suppose now that we measure the length of the table too:

    length = (100±0.5)cm

Since we are measuring with the same tape measure, the outcome in terms of
error is similar: we're making another error of about 0.5cm.

Now, what is the perimeter of the table?  It is:

    perimeter = 2 * (width + length)

But, since the original width and length that we computed are approximate
numbers, we will also have some error in the computation of the perimeter.

Since we could have overshot our measurement for all values at one extreme,
or undershot at another, our perimeter falls in the interval:

    perimeter = (300±2)cm

Note here that individual measurement errors added up.  Now, if we wanted to
compute the difference between length and width, we'd get:

    length - width = (50±1)cm

What happened here?  We see that even though the data points were
subtracted, the errors were *added* together.  This is because again the
errors could have conspired to make our measurement less accurate, and we
have to account for that.

This library has a few functions that make working with approximate numbers
easy.  You can load up some approximate numbers like so:

    import "github.com/filmil/approx"
    width, _ := approx.Parse("50±0.5")
    length, _ := approx.Parse("50±0.5")
    perimeter := approx.Add(
        approx.Add(width, length),
        approx.Add(width, length),
    )
