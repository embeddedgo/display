// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math2d

import (
	"image"
	"math"
	"testing"
)

func toInt(f float64) int {
	if f < 0 {
		f -= 0.5
	} else {
		f += 0.5
	}
	return int(f)
}

func rotate1(v image.Point, angle int32) image.Point {
	a := (math.Pi / (1 << 31)) * float64(angle)
	cosa := math.Cos(a)
	sina := math.Sin(a)
	x := float64(v.X)
	y := float64(v.Y)
	return image.Point{toInt(x*cosa - y*sina), toInt(x*sina + y*cosa)}
}

func TestRotate(t *testing.T) {
	const frac = 2
	v := F(image.Pt(4567, -3456), frac)
	n := 0
	angle := int32(0)
	for {
		v0 := Rotate(v, angle)
		v1 := rotate1(v, angle)
		v0 = I(v0, frac)
		v1 = I(v1, frac)
		dv := v0.Sub(v1)
		if dv.X*dv.X+dv.Y*dv.Y > 2 {
			n++
			t.Error(angle, v0, v1)
		}
		if angle += 512; angle == 0 {
			break
		}
	}
	if n != 0 {
		t.Error("n =", n)
	}
}

func polar1(v image.Point) (r int, angle int32) {
	x, y := float64(v.X), float64(v.Y)
	th := math.Atan2(y, x)
	angle = int32(1<<31*th/math.Pi + 0.5)
	r = int(math.Sqrt(x*x+y*y) + 0.5)
	return
}

func TestPolar(t *testing.T) {
	const (
		fracR    = 8
		imprecTh = 9
	)
	var angle int32
	n := 0
	for {
		v := Rotate(F(image.Pt(1e5, 0), fracR), angle)
		r, th := Polar(v)
		r = (r + 1<<(fracR-1)) >> fracR
		th = (th + 1<<(imprecTh-1)) >> imprecTh << imprecTh
		r1, th1 := polar1(v)
		r1 = (r1 + 1<<(fracR-1)) >> fracR
		th1 = (th1 + 1<<(imprecTh-1)) >> imprecTh << imprecTh
		if r != r1 || th != th1 {
			n++
			t.Error("(", r, ",", th, ") != (", r1, ",", th1, ")")
		}
		if angle += 1 << imprecTh; angle == 0 {
			break
		}
	}
	if n != 0 {
		t.Error("n =", n)
	}
}

func BenchmarkRotateCORDIC(b *testing.B) {
	v := image.Pt(1e6, 1e4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v = Rotate(v, int32(i))
		v = Rotate(v, int32(-i))
	}
	b.StopTimer()
}

func BenchmarkRotateSinCos(b *testing.B) {
	v := image.Pt(1e6, 1e4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v = rotate1(v, int32(i))
		v = rotate1(v, int32(-i))
	}
	b.StopTimer()
}
