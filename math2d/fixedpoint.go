// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math2d

import (
	"image"
	"math/bits"
)

// F returns image.Point{v.X<<n, v.Y<<n}. It provides a convenient way to
// convert an integer vector to a fixed-point one. If you want to ensure a
// specific number of significant bits use A instead.
func F(v image.Point, n int) image.Point {
	return image.Point{v.X << n, v.Y << n}
}

// I provides a convenient way to convert a fixed-point vector to an integer
// one. It supports positive and negative n so it can be used to invert F and A
// as below
//
//	I(F(v, n), n) == v // if no overflow
//	I(A(v, n))    == v // if no overflow and no reduction of significant bits
func I(v image.Point, n int) image.Point {
	if n >= 0 {
		m := uint(n)
		round := 1 << (m - 1)
		return image.Point{(v.X + round) >> m, (v.Y + round) >> m}
	}
	m := uint(-n)
	return image.Point{v.X << m, v.Y << m}
}

// A converts an integer vector to a n bit fixed-point one ensuring that all n
// bits are used by the value stored in X or Y of the returned vector. It
// returns the converted vector and the length of fractional part (can be
// negative if the number of significant bits was reduced). Use A to normalize
// the vector before passing it to the Rotate or Polar function to improve
// accuracy.
func A(v image.Point, n int) (image.Point, int) {
	x := v.X
	if x < 0 {
		x = -x
	}
	y := v.Y
	if y < 0 {
		y = -y
	}
	if x < y {
		x = y
	}
	n -= bits.UintSize - bits.LeadingZeros(uint(x))
	if n >= 0 {
		return image.Point{v.X << uint(n), v.Y << uint(n)}, n
	}
	m := uint(-n)
	round := 1 << (m - 1)
	return image.Point{(v.X + round) >> m, (v.Y + round) >> m}, n
}
