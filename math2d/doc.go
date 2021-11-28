// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package math2d provides a set of fixed-point 2D operations that work with
// image.Point vectors and int32 angles.
//
// The angle unit is equal to 1/4294967296 of the full angle. The int32 type can
// represent angles from -FullAngle/2 to FullAngle/2 - 1 (inclusive).
//
// The fixed-point CORDIC algorithm used by this package is at least 10 times
// faster than using math.Cos and math.Sin on the ISA that does not support
// cos/sin in hardware. If the ISA provides hardware cos/sin the CORDIC is about
// 2 times slower than using math.Cos/math.Sin but it does not use FPU context
// which could be an advantage in some cases.
//
// The accuracy of the fixed-point CORDIC algorithm depends on many factors but
// is usually much worse than the accuracy of any algorithm based on cos/sin and
// floating-point calculations (see "Numerical Accuracy and Hardware Tradeoffs
// for CORDIC Arithmetic for Special Purpose Processors", K. Kota, J. Cavallaro,
// http://scholarship.rice.edu/bitstream/handle/1911/20039/Kot1993Jul1NumericalA.PDF)
package math2d

// Angle constants. FullAngle does not fit in int32 but can be used to calculate
// smaller angle constants.
const (
	FullAngle        = 1 << 32       // 2π rad = 360°
	RightAngle int32 = FullAngle / 4 // π/2 rad = 90°
)
