// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixd

import "image"

// Rectangle draws an empty or filled rectangle with a given diagonal.
func (a *Area) Rectangle(pa, pb image.Point, fill bool) {
	if pa.X > pb.X {
		pb.X, pa.X = pa.X, pb.X
	}
	if pa.Y > pb.Y {
		pb.Y, pa.Y = pa.Y, pb.Y
	}
	var r image.Rectangle
	r.Min = pa
	if fill {
		r.Max.X = pb.X + 1
		r.Max.Y = pb.Y + 1
	} else {
		r.Max.X = pb.X + 1
		r.Max.Y = pa.Y + 1
		a.Fill(r)
		r.Max.X = pa.X + 1
		r.Max.Y = pb.Y + 1
		a.Fill(r)
		r.Max.X = pb.X + 1
		r.Min.Y = pb.Y
		a.Fill(r)
		r.Min.X = pb.X
		r.Min.Y = pa.Y
	}
	a.Fill(r)
}
