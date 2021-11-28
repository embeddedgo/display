// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math2d

import "image"

// Polar returns the polar form of v.  The internal calculations require that
// v.Mul(2) must not overflow.
func Polar(v image.Point) (r int, angle int32) {
	if v.X < 0 {
		v.X = -v.X
		v.Y = -v.Y
		angle = -FullAngle / 2
	}
	var theta int32
	for i, th := range &cordicThs {
		round := 1 << uint(i-1)
		dx := (v.Y + round) >> uint(i)
		dy := (v.X + round) >> uint(i)
		if v.Y <= 0 {
			v.X -= dx
			v.Y += dy
			theta -= th
		} else {
			v.X += dx
			v.Y -= dy
			theta += th
		}
	}
	const (
		roundK  = 1 << (fracK - 1)
		roundTh = 1 << (fracTh - 1)
	)
	angle += (theta + roundTh) >> fracTh
	r = int((int64(v.X)*cordicK + roundK) >> fracK)
	return
}
