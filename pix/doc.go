// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pix provides set of drawing functions for simple pixmap based
// displays. Such displays, unlike the vector or display-list based ones,
// require data in the form of ready-made array of pixels.
//
// Only one drawing operation is required from the display driver: drawing an
// image on a selected rectangle. This basic drawing primitive is supported in
// hardware by many simple LCD and OLED display controlers like ILI9341,
// ILI948x, ST7796S, HX8357, SSD1351 and can be also easily implemented in
// software for other displays.
package pix
