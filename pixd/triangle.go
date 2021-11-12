// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import "image"

func (a *Area) Triangle(p0, p1, p2 image.Point, fill bool) {
	if !fill {
		a.Line(p0, p1)
		a.Line(p1, p2)
		a.Line(p2, p0)
		return
	}
	// order the points by X and calculate the width of the triangle
	if p0.X > p2.X {
		p0, p2 = p2, p0
	}
	if p0.X > p1.X {
		p0, p1 = p1, p0
	} else if p1.X > p2.X {
		p1, p2 = p2, p1
	}
	//width := c.X - a.X

}
