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
	return image.Pt(v.X<<n, v.Y<<n)
}

// I provides a convenient way to convert a fixed-point vector to an integer
// one. It supports positive and negative n so it can be used to invert F and A
// as below
//
//	I(F(v, n), n) == v // true if no overflow
//	I(A(v, n))    == v // true if no overflow and no significant bit reduction
func I(v image.Point, n int) image.Point {
	if n >= 0 {
		m := uint(n)
		round := 1 << (m - 1)
		return image.Pt((v.X+round)>>m, (v.Y+round)>>m)
	}
	m := uint(-n)
	return image.Pt(v.X<<m, v.Y<<m)
}

// A converts an integer vector to a fixed-point one ensuring that both f.X
// and f.Y valeus fit in n bits and at last one of f.X or f.Y has exactly n
// significant bits. A returns the converted vector and the length of its
// fractional part (can be negative if the number of significant bits was
// reduced). For zero v it returns (v, n). Use A to normalize the vector before
// passing it to the Rotate or Polar functions to improve accuracy.
func A(v image.Point, n int) (f image.Point, m int) {
	x, y := v.X, v.Y
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	if x < y {
		x = y
	}
	n -= bits.UintSize - bits.LeadingZeros(uint(x))
	if n >= 0 {
		return image.Pt(v.X<<uint(n), v.Y<<uint(n)), n
	}
	k := uint(-n)
	round := 1 << (k - 1)
	return image.Pt((v.X+round)>>k, (v.Y+round)>>k), n
}
