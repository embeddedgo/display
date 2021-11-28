// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math2d

import "image"

// Rotate rotates the vector v by the given angle.
//
// Rotate can also be used to convert a vector in the polar form (R, theta) to
// the rectangular one as below
//
//	Rotate(image.Pt(R, 0), theta)
//
// Moreover, it can be used to calculate trigonometric functions as below
//
//	Rotate(image.Pt(1<<n, 0), alpha)
//
// The returned vector contains {cos(alpha), sin(alpha)} with n-bit fractional
// part. Since the cosine and sine are calculated simultaneously, the returned
// vector can be also thought of as a fractional representation of tangent and
// cotangent.
//
// The internal calculations require that v.Mul(2) must not overflow.
func Rotate(v image.Point, angle int32) image.Point {
	// make -FullAngle/4 <= angle < FullAngle/4
	if angle >= 0 {
		angle -= RightAngle
		v.X, v.Y = -v.Y, v.X
	} else {
		angle += RightAngle
		v.X, v.Y = v.Y, -v.X
	}
	angle <<= fracTh
	for i, th := range &cordicThs {
		round := 1 << uint(i-1)
		dx := (v.Y + round) >> uint(i)
		dy := (v.X + round) >> uint(i)
		if angle >= 0 {
			v.X -= dx
			v.Y += dy
			angle -= th
		} else {
			v.X += dx
			v.Y -= dy
			angle += th
		}
	}
	const roundK = 1 << (fracK - 1)
	v.X = int((int64(v.X)*cordicK + roundK) >> fracK)
	v.Y = int((int64(v.Y)*cordicK + roundK) >> fracK)
	return v
}
