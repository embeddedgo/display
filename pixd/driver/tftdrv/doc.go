// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tftdrv defines the Display Controller Interface (DCI) and provides
// generic drivers to the graphics controllers that can be found in a wide
// variety of LCD and OLED displays. The supported controllers are generally
// similar to each other both in terms of used DCI and the command set they
// provide. The de facto standard in this area was set in the early 2000s by
// Philips PCF8833 and Epson S1D15G00 controllers used in the first mobile
// phones with color display like Nokia 6100 or Siemens S65.
//
// The subpackages provides drivers to the specific controllers.
package tftdrv
