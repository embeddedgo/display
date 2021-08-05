// Copyright 2021 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pixdisp

// RGBA represents a traditional 24-bit color without the alpha channel.
type RGB struct {
	R, G, B uint8
}

func (c RGB) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = 0xffff
	return
}

// RGB565 represents a traditional 16-bit color without the alpha channel.
type RGB565 uint16

func (c RGB565) RGBA() (r, g, b, a uint32) {
	r = uint32(c >> 11)
	r |= r << 13
	g = uint32(c >> 5 & 0x3f)
	g |= g << 12
	b = uint32(c & 0x1f)
	b |= b << 13
	a = 0xffff
	return
}
